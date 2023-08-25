package api

import (
	"context"
	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateSegment(c *gin.Context) {
	ctx := context.Background()
	segment := model.Segment{}

	if err := c.BindJSON(&segment); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}
	err := h.service.CreateSegment(ctx, segment.Slug)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}
