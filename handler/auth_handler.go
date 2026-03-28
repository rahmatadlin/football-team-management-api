package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"football-team-management-api/service"
	"football-team-management-api/utils/response"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	token, admin, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "login successful", gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   nil,
		"admin": gin.H{
			"id":    admin.ID,
			"email": admin.Email,
		},
	})
}
