package api

import (
	"context"
	"net/http"

	"github.com/elgntt/segmentation-service/internal/model"
	"github.com/elgntt/segmentation-service/internal/pkg/app_err"
	response "github.com/elgntt/segmentation-service/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

// CreateSegment
// @Summary CreateSegment
// @Tags Segment
// @Description Create segment
// @Produce application/json
// @Param input body model.AddSegment true "segment info"
// @Success 201
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /segment [post]
func (h *handler) CreateSegment(c *gin.Context) {
	ctx := context.Background()
	request := model.AddSegment{}

	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}

	if err := validateReqData(request); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	err := h.segmentService.CreateSegment(ctx, request)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func validateReqData(segmentData model.AddSegment) error {
	if segmentData.SegmentSlug == "" {
		return app_err.NewBusinessError(ErrInvalidSegmentSlug)
	}
	if segmentData.AutoJoinPercent < 0 || segmentData.AutoJoinPercent > 100 {
		return app_err.NewBusinessError(ErrInvalidAutoJoinPercent)
	}

	return nil
}
