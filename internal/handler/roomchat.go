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

type RoomChatHandler struct {
	RoomChatUC usecase.RoomChatUseCase
}

func NewRoomChatUserHandler(roomChatUC usecase.RoomChatUseCase) *RoomChatHandler {
	return &RoomChatHandler{roomChatUC}
}

func (h *RoomChatHandler) GetRoomChatByUserId(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	result, err := h.RoomChatUC.GetGroupsByUserID(claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)

}

func (h *RoomChatHandler) GetRoomChatById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["id"])
	result, err := h.RoomChatUC.GetRoomChatByID(uint(roomId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

func (h *RoomChatHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var input struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.RoomChatUC.CreateRoom(claims.UserID, input.Name, input.Desc); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"name":            input.Name,
		"desc":            input.Desc,
		"creator (admin)": claims.UserID,
	})

}

func (u *RoomChatHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	if err := u.RoomChatUC.DeleteRoom(uint(roomId), claims.UserID); err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		case utils.ErrInternal:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succeed delete room",
	})

}

func (h *RoomChatHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	var input struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid body requset")
		return
	}

	if err := h.RoomChatUC.UpdateRoom(uint(roomId), claims.UserID, input.Desc, input.Name); err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		case utils.ErrInternal:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"name": input.Name,
		"desc": input.Desc,
	})
}

// room member
func (h *RoomChatHandler) GetRoomMemberById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["id"])
	room, err := h.RoomChatUC.GetRoomMember(uint(roomId))
	if err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrInternal:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, room)
}

func (h *RoomChatHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])
	var input struct {
		TargetIDS []uint `json:"target_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(input.TargetIDS) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := h.RoomChatUC.Addmembers(uint(roomId), claims.UserID, input.TargetIDS); err != nil {
		switch err {
		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		case utils.ErrInternal:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "succeed add members",
		"members": input.TargetIDS,
	})
}

func (h *RoomChatHandler) DeleteMembersByAdmin(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	var input struct {
		TargetIDS []uint `json:"target_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(input.TargetIDS) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid body request")
		return
	}

	if err := h.RoomChatUC.DeleteMembersByAdmin(uint(roomId), claims.UserID, input.TargetIDS); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "succeed kick members",
		"members": input.TargetIDS,
	})
}

func (h *RoomChatHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	params := mux.Vars(r)
	roomId, _ := strconv.Atoi(params["roomId"])

	if err := h.RoomChatUC.LeaveRoom(uint(roomId), claims.UserID, claims.UserID); err != nil {
		switch err {

		case utils.ErrRoomNotFound:
			utils.WriteError(w, http.StatusNotFound, err.Error())
		case utils.ErrUnauthorized:
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
		case utils.ErrInternal:
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "succeed leave room",
	})
}
