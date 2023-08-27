package api

import (
	"context"
	"net/http"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateSegment(c *gin.Context) {
	ctx := context.Background()
	request := model.AddSegment{}

	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}

	if err := validateSegmentSlug(request.Slug); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	err := h.service.CreateSegment(ctx, request)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}
