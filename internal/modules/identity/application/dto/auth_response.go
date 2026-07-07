package dto

import "time"

type AuthUserResponse struct {
	ID string `json:"id"`

	Email string `json:"email"`

	Status string `json:"status"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`

	RefreshToken string `json:"refreshToken"`

	TokenType string `json:"tokenType"`

	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`

	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

type LoginResponse struct {
	User AuthUserResponse `json:"user"`

	Tokens TokenResponse `json:"tokens"`
}

type RefreshTokenResponse struct {
	Tokens TokenResponse `json:"tokens"`
}

type CurrentUserResponse struct {
	ID string `json:"id"`

	Email string `json:"email"`

	Status string `json:"status"`

	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}
