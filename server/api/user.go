package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/janto-pee/fintech-platform.git/db/sqlc"
	"github.com/janto-pee/fintech-platform.git/util"
)

type CreateUserRequest struct {
	Username       string `json:"Username" binding:"required,alphanum"`
	FullName       string `json:"FullName" binding:"required"`
	HashedPassword string `json:"HashedPasword" binding:"required,min=6"`
	Email          string `json:"Email" binding:"required,email"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	ctx.JSON(http.StatusOK, req)
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
