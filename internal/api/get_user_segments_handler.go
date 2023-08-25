package api

import (
	"context"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	"net/http"
	"strconv"

	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetUserSegments(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid userId parameter"))
		return
	}

	ctx := context.Background()
	userSegments, err := h.service.GetActiveUserSegments(ctx, userId)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userSegments": userSegments,
	})
}
