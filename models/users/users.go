package users

type UserData struct {
	ID         int
	Username   string
	Password   string
	Fullname   string
	Gender     string
	Location   string
	Education  string
	Occupation string
	Bio        string
	IsPremium  int
}

type CreateUserRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Fullname   string `json:"fullname"`
	Gender     string `json:"gender"`
	Location   string `json:"location"`
	Education  string `json:"education"`
	Occupation string `json:"occupation"`
	Bio        string `json:"bio"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

type UserResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	Gender     string `json:"gender"`
	Location   string `json:"location"`
	Education  string `json:"education"`
	Occupation string `json:"occupation"`
	Bio        string `json:"bio"`
}
