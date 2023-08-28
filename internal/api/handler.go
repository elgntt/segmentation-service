package api

import (
	"context"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/elgntt/avito-internship-2023/docs"
)

type service interface {
	CreateSegment(ctx context.Context, segmentData model.AddSegment) error
	DeleteSegment(ctx context.Context, slug string) error
	UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error
	GetActiveUserSegments(ctx context.Context, userId int) ([]string, error)
	GenerateCSVFile(ctx context.Context, month, year, userId int) (string, error)
}

const (
	ErrInvalidYearParameter   = `invalid "year" parameter`
	ErrInvalidMonthParameter  = `invalid "month" parameter`
	ErrInvalidSegmentSlug     = `invalid segment slug`
	ErrInvalidUserIdParameter = `invalid "userId" parameter`
)

type handler struct {
	service
}

func New(service service) *gin.Engine {
	h := handler{
		service: service,
	}

	r := gin.New()

	r.Static("/assets/csv_reports", "./assets/csv_reports")

	r.POST("/segment/create", h.CreateSegment)
	r.POST("/user/segment/action", h.UserSegmentAction)
	r.DELETE("/segment/delete", h.DeleteSegment)
	r.GET("/user/segment/getAllActive", h.GetUserSegments)
	r.GET("/history/file", h.GetReportFile)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func validateSegmentSlug(slug string) error {
	if slug == "" {
		return app_err.NewBusinessError(ErrInvalidSegmentSlug)
	}

	return nil
}
