package handler

import (
	"chat/internal/usecase"
	"chat/utils"
	"chat/utils/middleware"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userUC usecase.UserUsecase
}

func NewUserHandler(userUC usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body requst")
		return
	}

	err := h.userUC.Register(input.Username, input.Email, input.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"massage": "link verification has been send on your email",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userUC.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !user.IsVerified {
		http.Error(w, "belum verified", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWTLogin(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed generate token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "berhasil login",
		"token":   token,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.userUC.UpdateUser(claims.Email, input.Username, input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"email":    input.Email,
		"username": input.Username,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	fmt.Println(claims.Email)
	if err := h.userUC.DeleteUser(claims.Email); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user berhasil dihapus",
	})
}

func (h *UserHandler) Verifikasi(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "jwt error")
		return
	}
	if err := h.userUC.ValidateUser(claims.Email); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed verified user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"massage": "succed verified user",
	})
}

func (h *UserHandler) ResendLinkVerif(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.userUC.ResendLinkVerif(input.Email); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to send link")
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"massage": "succed send link verification",
	})
}
