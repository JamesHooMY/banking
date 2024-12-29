package user

import (
	"github.com/shopspring/decimal"
)

type CreateUserReq struct {
	Name     string `json:"name" binding:"required,min=3,max=20,alphanumunicode"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type User struct {
	ID      uint            `json:"id"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Balance decimal.Decimal `json:"balance"`
}

type CreateUserResp struct {
	Data *User `json:"data"`
}

type GetUserResp struct {
	Data *User `json:"data"`
}

type GetUsersResp struct {
	Data []*User `json:"data"`
}

type APIKey struct {
	Key    string `json:"key"`
	Secret string `json:"secret,omitempty"`
	UserID uint   `json:"userId"`
}

type CreateAPIKeyResp struct {
	Data *APIKey `json:"data"`
}
