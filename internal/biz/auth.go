package biz

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"user-center/internal/conf"
	"user-center/internal/email"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrAccountLocked     = errors.New("account locked, please try again later")
	ErrInvalidCode       = errors.New("invalid verification code")
	ErrCodeExpired       = errors.New("verification code expired")
)

const (
	MaxPasswordErrors = 5
	LockDuration      = 5 * time.Minute
)

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

type AuthUseCase struct {
	repo       UserRepo
	jwtSecret  string
	expireTime int64
	JwtSecret  string
	ExpireTime int64
	smtpCfg    *conf.SMTP
}

func NewAuthUseCase(repo UserRepo, jwtSecret string, expireTime int64, smtpCfg *conf.SMTP) *AuthUseCase {
	return &AuthUseCase{
		repo:       repo,
		jwtSecret:  jwtSecret,
		expireTime: expireTime,
		JwtSecret:  jwtSecret,
		ExpireTime: expireTime,
		smtpCfg:    smtpCfg,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, phone, password, nickname string) (*User, error) {
	// Check if user exists by email
	existingUser, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, *User, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return "", nil, ErrAccountLocked
	}

	// Check password with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Increment password error count
		user.PasswordErrors++
		if user.PasswordErrors >= MaxPasswordErrors {
			lockedUntil := time.Now().Add(LockDuration)
			user.LockedUntil = &lockedUntil
		}
		// Update user in database
		uc.repo.UpdateUser(ctx, user)
		return "", nil, ErrInvalidPassword
	}

	// Password correct, reset error count
	if user.PasswordErrors > 0 || user.LockedUntil != nil {
		user.PasswordErrors = 0
		user.LockedUntil = nil
		uc.repo.UpdateUser(ctx, user)
	}

	// Generate JWT token
	token, err := uc.generateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

func (uc *AuthUseCase) generateToken(userID int64) (string, error) {
	now := time.Now()
	expire := now.Add(time.Duration(uc.expireTime) * time.Second)

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     expire.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.JwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func (uc *AuthUseCase) DeleteAccount(ctx context.Context, userID int64) error {
	return uc.repo.DeleteUser(ctx, userID)
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.repo.UpdateUser(ctx, user)
}

const (
	verifyCodeExpire = 5 * time.Minute
	verifyCodeLen    = 6
)

func generateVerifyCode() (string, error) {
	code := ""
	for i := 0; i < verifyCodeLen; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code, nil
}

func (uc *AuthUseCase) SendVerifyCode(ctx context.Context, emailAddr string) (string, error) {
	code, err := generateVerifyCode()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("verify_code:%s", emailAddr)
	err = uc.repo.SetCache(ctx, key, code, verifyCodeExpire)
	if err != nil {
		return "", err
	}

	if uc.smtpCfg != nil {
		err = email.SendVerifyCode(uc.smtpCfg, emailAddr, code)
		if err != nil {
			return "", fmt.Errorf("failed to send email: %v", err)
		}
	}

	return code, nil
}

func (uc *AuthUseCase) RegisterWithCode(ctx context.Context, email, password, code, nickname string) (*User, error) {
	key := fmt.Sprintf("verify_code:%s", email)
	storedCode, err := uc.repo.GetCache(ctx, key)
	if err != nil {
		return nil, ErrCodeExpired
	}
	if storedCode != code {
		return nil, ErrInvalidCode
	}

	uc.repo.DeleteCache(ctx, key)

	return uc.Register(ctx, email, "", password, nickname)
}

func (uc *AuthUseCase) LoginWithCode(ctx context.Context, email, code string) (string, *User, error) {
	key := fmt.Sprintf("verify_code:%s", email)
	storedCode, err := uc.repo.GetCache(ctx, key)
	if err != nil {
		return "", nil, ErrCodeExpired
	}
	if storedCode != code {
		return "", nil, ErrInvalidCode
	}

	uc.repo.DeleteCache(ctx, key)

	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}

	token, err := uc.generateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
