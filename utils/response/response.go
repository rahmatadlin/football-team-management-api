package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"football-team-management-api/utils/apperror"
)

type envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, envelope{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	code := apperror.HTTPStatus(err)
	msg := apperror.Message(err)
	if code == http.StatusInternalServerError {
		msg = "internal server error"
	}
	c.JSON(code, envelope{
		Status:  "error",
		Message: msg,
		Data:    nil,
	})
}

func ErrorMessage(c *gin.Context, code int, message string) {
	c.JSON(code, envelope{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}
