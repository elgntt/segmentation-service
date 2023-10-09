package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/elgntt/segmentation-service/internal/pkg/app_err"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func WriteErrorResponse(c *gin.Context, err error) {
	var bErr app_err.BusinessError

	if errors.As(err, &bErr) {
		errorResponse := ErrorResponse{
			ErrorMessage: ErrorMessage{
				Message: bErr.Error(),
			},
		}

		c.JSON(http.StatusBadRequest, errorResponse)

	} else {
		errorResponse := ErrorResponse{
			ErrorMessage: ErrorMessage{
				Message: "Internal server error",
			},
		}

		log.Println(err)

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}
