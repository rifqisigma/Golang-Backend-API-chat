package handler

import (
	"chat/internal/usecase"
	"chat/utils"
	"chat/utils/middleware"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatUC usecase.ChatUsecase
}

func NewChatHandler(chatUC usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{chatUC}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	var input struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Message == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	response, err := h.chatUC.CreateChat(input.Message, uint(roomId), claims.UserID)
	if err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, &response)

}

func (h *ChatHandler) UpdateChat(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	var input struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if input.Message == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.chatUC.UpdateChat(uint(roomId), claims.UserID, input.Message); err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"room_id": roomId,
		"user_id": claims.UserID,
		"message": input.Message,
	})

}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	if err := h.chatUC.DeleteChat(uint(roomId), claims.UserID); err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "succeed delete",
	})

}
func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	result, err := h.chatUC.GetAllChatByRoomId(uint(roomId))
	if err != nil {
		switch err {
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, result)

}
