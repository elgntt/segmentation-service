package api

import (
	"context"
	"net/http"

	"github.com/elgntt/segmentation-service/internal/pkg/app_err"
	response "github.com/elgntt/segmentation-service/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type DeleteSegmentRequest struct {
	SegmentSlug string `json:"slug"`
}

// DeleteSegment
// @Summary DeleteSegment
// @Tags Segment
// @Description Delete segment
// @Produce application/json
// @Param input body api.DeleteSegmentRequest true "segment info"
// @Success 200
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /segment [delete]
func (h *handler) DeleteSegment(c *gin.Context) {
	ctx := context.Background()
	request := DeleteSegmentRequest{}

	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid request body"))
		return
	}

	if err := validateSegmentSlug(request.SegmentSlug); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	err := h.segmentService.DeleteSegment(ctx, request.SegmentSlug)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}
