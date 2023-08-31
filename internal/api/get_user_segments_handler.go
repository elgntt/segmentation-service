package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"

	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type UserSegmentsResponse struct {
	UserId   int      `json:"userId"`
	Segments []string `json:"segments"`
}

// GetUserSegments GetReportSegments
// @Summary GetUserSegments
// @Tags User
// @Description Allows you to get data on segments of some user
// @Produce application/json
// @Param 	userId query int true "actual userId"
// @Success 200 {object} api.UserSegmentsResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /user/segment/getAllActive [get]
func (h *handler) GetUserSegments(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		response.WriteErrorResponse(c, app_err.NewBusinessError("invalid userId parameter"))
		return
	}
	if userId < 1 {
		response.WriteErrorResponse(c, app_err.NewBusinessError(ErrInvalidUserId))
	}
	ctx := context.Background()
	userSegments, err := h.userService.GetActiveUserSegments(ctx, userId)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, UserSegmentsResponse{
		UserId:   userId,
		Segments: userSegments,
	})
}
