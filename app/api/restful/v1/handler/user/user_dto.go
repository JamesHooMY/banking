package user

import (
	"github.com/shopspring/decimal"
)

type CreateUserReq struct {
	Name string `json:"name" binding:"required,min=3,max=20,alphanumunicode"`
}

type User struct {
	ID      uint            `json:"id"`
	Name    string          `json:"name"`
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
