package api

import (
	"context"
	"net/http"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) UserSegmentAction(c *gin.Context) {
	ctx := context.Background()
	request := model.UserSegmentAction{}
	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}

	if err := validateTime(request.SegmentExpirationTime); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	err := h.service.UserSegmentAction(ctx, request)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func validateTime(expirationTime *time.Time) error {
	if expirationTime == nil {
		return nil
	}

	transmittedTime, err := time.Parse(time.RFC3339, expirationTime.Format(time.RFC3339))
	if err != nil {
		return err
	}

	if transmittedTime.Before(time.Now()) {
		return app_err.NewBusinessError("invalid expiration time argument")
	}

	return nil
}
