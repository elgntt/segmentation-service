package api

import (
	"context"
	"github.com/elgntt/avito-internship-2023/internal/model"

	"github.com/gin-gonic/gin"
)

type service interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error
	GetActiveUserSegments(ctx context.Context, userId int) ([]string, error)
}

type handler struct {
	service
}

func New(service service) *gin.Engine {
	h := handler{
		service: service,
	}

	r := gin.New()

	r.POST("/segment/create", h.CreateSegment)
	r.POST("/user/segment/action", h.UserSegmentAction)
	r.DELETE("/segment/delete", h.DeleteSegment)
	r.GET("/user/segment/getAllActive", h.GetUserSegments)

	return r
}
