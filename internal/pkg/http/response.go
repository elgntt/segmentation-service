package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status       string `json:"status"`
	errorMessage `json:"error"`
}

type errorMessage struct {
	Message string `json:"message"`
}

func WriteErrorResponse(c *gin.Context, err error) {
	var bErr app_err.BusinessError

	if errors.As(err, &bErr) {
		errorResponse := errorResponse{
			Status: "error",
			errorMessage: errorMessage{
				Message: bErr.Error(),
			},
		}

		c.JSON(http.StatusBadRequest, errorResponse)

	} else {
		errorResponse := errorResponse{
			Status: "error",
			errorMessage: errorMessage{
				Message: "Internal server error",
			},
		}

		log.Println(err)

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}
