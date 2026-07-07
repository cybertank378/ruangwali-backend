package dto

type LoginRequest struct {
	Email string `json:"email"`

	Password string `json:"password"`

	UserAgent string `json:"-"`

	IPAddress string `json:"-"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`

	UserAgent string `json:"-"`

	IPAddress string `json:"-"`
}
