package route

import (
	"date-app/config"
	userHandler "date-app/handler/users"
	userRepo "date-app/repository/users/mysql"
	userRedisRepo "date-app/repository/users/redis"
	userSvc "date-app/service/users"
	"date-app/utils/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRoute() *mux.Router {
	userRepository := userRepo.NewUsersRepository(config.DB)
	userRedisRepository := userRedisRepo.NewUsersRepository(config.RDB)
	userService := userSvc.NewUsersService(userRepository, userRedisRepository)
	userHandlers := userHandler.NewUsersHandler(userService)

	router := mux.NewRouter()
	router.Handle("/user", middleware.ApiKeyMiddleware(http.HandlerFunc(userHandlers.UserRegister))).Methods(http.MethodPost)
	router.HandleFunc("/user/login", userHandlers.Login).Methods(http.MethodPost)
	router.Handle("/user/view", middleware.JWTMiddleware(http.HandlerFunc(userHandlers.GetRandomUser))).Methods(http.MethodGet)
	router.Handle("/user/premium", middleware.ApiKeyMiddleware(http.HandlerFunc(userHandlers.UpdateIsPremium))).Methods(http.MethodPut)
	router.Handle("/user/like", middleware.JWTMiddleware(http.HandlerFunc(userHandlers.UserLike))).Methods(http.MethodPost)

	return router
}
