package route

import (
	"date-app/config"
	userHandler "date-app/handler/users"
	userRepo "date-app/repository/users/mysql"
	userSvc "date-app/service/users"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRoute() *mux.Router {
	userRepository := userRepo.NewUsersRepository(config.DB)
	userService := userSvc.NewUsersService(userRepository)
	userHandlers := userHandler.NewUsersHandler(userService)

	router := mux.NewRouter()
	router.HandleFunc("/user", userHandlers.UserRegister).Methods(http.MethodPost)
	router.HandleFunc("/user/login", userHandlers.Login).Methods(http.MethodPost)
	return router
}
