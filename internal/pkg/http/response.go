package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func WriteErrorResponse(c *gin.Context, err error) {
	var bErr app_err.BusinessError

	if errors.As(err, &bErr) {
		errorResponse := ErrorResponse{
			Status: "error",
			ErrorMessage: ErrorMessage{
				Message: bErr.Error(),
			},
		}

		c.JSON(http.StatusBadRequest, errorResponse)

	} else {
		errorResponse := ErrorResponse{
			Status: "error",
			ErrorMessage: ErrorMessage{
				Message: "Internal server error",
			},
		}

		log.Println(err)

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}
