package biz

import (
	"context"
	"time"

	"user-center/internal/entity"
)

type User struct {
	ID             int64
	Email          string
	Phone          string
	Password       string
	Nickname       string
	PasswordErrors int
	LockedUntil    *time.Time
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
	SetCache(ctx context.Context, key string, value interface{}, expire time.Duration) error
	GetCache(ctx context.Context, key string) (string, error)
	DeleteCache(ctx context.Context, key string) error
}

func UserFromEntity(e *entity.User) *User {
	if e == nil {
		return nil
	}
	return &User{
		ID:             e.ID,
		Email:          e.Email,
		Phone:          e.Phone,
		Password:       e.Password,
		Nickname:       e.Nickname,
		PasswordErrors: e.PasswordErrors,
		LockedUntil:    e.LockedUntil,
	}
}

func UserToEntity(u *User) *entity.User {
	if u == nil {
		return nil
	}
	return &entity.User{
		ID:             u.ID,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Nickname:       u.Nickname,
		PasswordErrors: u.PasswordErrors,
		LockedUntil:    u.LockedUntil,
	}
}
