package api

import (
	"context"
	"net/http"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) DeleteSegment(c *gin.Context) {
	ctx := context.Background()
	request := struct {
		Slug string `json:"slug"`
	}{}
	
	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}

	err := h.service.DeleteSegment(ctx, request.Slug)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
