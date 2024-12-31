package users

import (
	"date-app/models"
	usersModel "date-app/models/users"
	"date-app/service/users"
	"date-app/utils/jwt"
	"encoding/json"
	"net/http"
)

type Users interface {
	UserRegister(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetRandomUser(w http.ResponseWriter, r *http.Request)
	UpdateIsPremium(w http.ResponseWriter, r *http.Request)
	UserLike(w http.ResponseWriter, r *http.Request)
}

type usersHandler struct {
	usersService users.Users
}

func NewUsersHandler(usersService users.Users) Users {
	return &usersHandler{
		usersService: usersService,
	}
}

func (u *usersHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var res models.GeneralResponse
	var req usersModel.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "error bind request: " + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.Username == "" || req.Password == "" {
		res.Code = http.StatusBadRequest
		res.Message = "username and password must be filled"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	err = u.usersService.Create(&req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Code = http.StatusCreated
	res.Message = "user registered"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (u *usersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var res models.GeneralResponse
	var req usersModel.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "error bind request: " + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.Username == "" || req.Password == "" {
		res.Code = http.StatusBadRequest
		res.Message = "username and password must be filled"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	data, err := u.usersService.Login(&req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Code = http.StatusOK
	res.Message = "login success"
	res.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (u *usersHandler) GetRandomUser(w http.ResponseWriter, r *http.Request) {
	var res models.GeneralResponse
	w.Header().Set("Content-Type", "application/json")
	claims, err := jwt.GetClaims(r)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	data, err := u.usersService.GetRandomUser(claims.UserID)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Code = http.StatusOK
	res.Message = "get random user success"
	res.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (u *usersHandler) UpdateIsPremium(w http.ResponseWriter, r *http.Request) {
	var res models.GeneralResponse
	var req usersModel.UpdateIsPremiumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "error bind request: " + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.IsPremium == nil || req.UserID == 0 {
		res.Code = http.StatusBadRequest
		res.Message = "is_premium and user_id must be filled"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	err = u.usersService.UpdateIsPremium(&req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Code = http.StatusCreated
	res.Message = "update success"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (u *usersHandler) UserLike(w http.ResponseWriter, r *http.Request) {
	var res models.GeneralResponse
	var req usersModel.UserLikeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Message = "error bind request: " + err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if req.UserID == 0 {
		res.Code = http.StatusBadRequest
		res.Message = "user id must be filled"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	claims, err := jwt.GetClaims(r)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	data, err := u.usersService.UserLike(claims.UserID, &req)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Code = http.StatusCreated
	res.Message = "user like success"
	res.Data = data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
