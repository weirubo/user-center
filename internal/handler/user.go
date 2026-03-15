package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"user-center/internal/biz"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	authUC    *biz.AuthUseCase
	jwtSecret string
}

func NewUserHandler(authUC *biz.AuthUseCase) *UserHandler {
	return &UserHandler{authUC: authUC, jwtSecret: authUC.JwtSecret}
}

func (h *UserHandler) getUserIDFromRequest(r *http.Request) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, nil
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil {
		return 0, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int64(claims["user_id"].(float64)), nil
	}

	return 0, nil
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := h.authUC.Register(r.Context(), req.Email, req.Phone, req.Password, req.Nickname)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"message": "register success",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	token, user, err := h.authUC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"nickname": user.Nickname,
	})
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}

	user, err := h.authUC.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"phone":    user.Phone,
		"nickname": user.Nickname,
	})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromRequest(r)
	if err != nil || userID == 0 {
		http.Error(w, "unauthorized", 401)
		return
	}

	err = h.authUC.DeleteAccount(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "account deleted successfully",
	})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromRequest(r)
	if err != nil || userID == 0 {
		http.Error(w, "unauthorized", 401)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.authUC.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "password changed successfully",
	})
}

func (h *UserHandler) SendVerifyCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err := h.authUC.SendVerifyCode(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "verification code sent",
	})
}

func (h *UserHandler) RegisterWithCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
		Nickname string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := h.authUC.RegisterWithCode(r.Context(), req.Email, req.Password, req.Code, req.Nickname)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"message": "register success",
	})
}

func (h *UserHandler) LoginWithCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	token, user, err := h.authUC.LoginWithCode(r.Context(), req.Email, req.Code)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"id":       user.ID,
		"nickname": user.Nickname,
	})
}
