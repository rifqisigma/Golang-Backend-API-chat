package router

import (
	"chat/internal/handler"
	"chat/utils/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(userHandler *handler.UserHandler, roomChatHandler *handler.RoomChatHandler, chatHandler *handler.ChatHandler) *mux.Router {
	r := mux.NewRouter()

	//auth
	r.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/verification", userHandler.Verifikasi).Methods(http.MethodGet)
	r.HandleFunc("/resendlink/verif", userHandler.ResendLinkVerif).Methods(http.MethodPost)

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(middleware.JWTAuthMiddleware)

	//user
	userRouter.HandleFunc("/delete", userHandler.DeleteUser).Methods(http.MethodDelete)
	userRouter.HandleFunc("/update", userHandler.UpdateUser).Methods(http.MethodPut)

	chatRouter := r.PathPrefix("/chat").Subrouter()
	chatRouter.Use(middleware.JWTAuthMiddleware)

	chatRouter.HandleFunc("/getroom", roomChatHandler.GetRoomChatByUserId).Methods(http.MethodGet)
	chatRouter.HandleFunc("/getroom/{id}", roomChatHandler.GetRoomChatById).Methods(http.MethodGet)
	chatRouter.HandleFunc("/createroom", roomChatHandler.CreateRoom).Methods(http.MethodPost)
	chatRouter.HandleFunc("/deleteroom/{roomId}", roomChatHandler.DeleteRoom).Methods(http.MethodDelete)
	chatRouter.HandleFunc("/updateroom/{roomId}", roomChatHandler.UpdateRoom).Methods(http.MethodPut)

	//room member
	chatRouter.HandleFunc("/getmember/{id}", roomChatHandler.GetRoomMemberById).Methods(http.MethodGet)
	chatRouter.HandleFunc("/addmember/{roomId}", roomChatHandler.AddMembers).Methods(http.MethodPost)
	chatRouter.HandleFunc("/deletemember/{roomId}", roomChatHandler.DeleteMembersByAdmin).Methods(http.MethodPost)
	chatRouter.HandleFunc("/leaveroom/{roomId}", roomChatHandler.LeaveRoom).Methods(http.MethodDelete)

	//chat
	chatRouter.HandleFunc("/createchat/{roomId}", chatHandler.CreateChat).Methods(http.MethodPost)
	chatRouter.HandleFunc("/updatechat/{roomId}", chatHandler.UpdateChat).Methods(http.MethodPut)
	chatRouter.HandleFunc("/deletechat/{roomId}", chatHandler.DeleteChat).Methods(http.MethodDelete)
	chatRouter.HandleFunc("/getchat/{roomId}", chatHandler.GetChat).Methods(http.MethodGet)

	return r
}
