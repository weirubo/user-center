package service

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	v1 "user-center/api/user/v1"
	"user-center/internal/biz"
)

var ErrUnauthorized = &biz.AuthError{Message: "unauthorized"}
var ErrUserNotFound = &biz.AuthError{Message: "user not found"}

type UserService struct {
	v1.UnimplementedUserServer
	
	authUC *biz.AuthUseCase
}

func NewUserService(authUC *biz.AuthUseCase) *UserService {
	return &UserService{
		authUC: authUC,
	}
}

func (s *UserService) getUserIDFromContext(ctx context.Context) (int64, error) {
	// Try to get from context (set by middleware)
	if userID, ok := ctx.Value("user_id").(int64); ok {
		return userID, nil
	}
	
	// Try to get from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authHeader := md.Get("authorization")
		if len(authHeader) > 0 {
			token := strings.TrimPrefix(authHeader[0], "Bearer ")
			if token != authHeader[0] {
				userID, err := s.authUC.ValidateToken(token)
				if err == nil {
					return userID, nil
				}
			}
		}
	}
	
	return 0, ErrUnauthorized
}

func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	user, err := s.authUC.Register(ctx, req.Email, req.Phone, req.Password, req.Nickname)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterReply{
		Id:      user.ID,
		Message: "register success",
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	token, user, err := s.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token:    token,
		Id:       user.ID,
		Nickname: user.Nickname,
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoReply, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	user, err := s.authUC.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &v1.GetUserInfoReply{
		Id:       user.ID,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.Nickname,
	}, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, req *v1.DeleteAccountRequest) (*v1.DeleteAccountReply, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	err = s.authUC.DeleteAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteAccountReply{
		Message: "account deleted successfully",
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordReply, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	err = s.authUC.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &v1.ChangePasswordReply{
		Message: "password changed successfully",
	}, nil
}
