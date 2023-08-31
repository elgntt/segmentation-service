package api

import (
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/elgntt/avito-internship-2023/docs"
)

const (
	ErrInvalidYearParameter   = `invalid "year" parameter`
	ErrInvalidMonthParameter  = `invalid "month" parameter`
	ErrInvalidSegmentSlug     = `invalid segment slug`
	ErrInvalidUserIdParameter = `invalid "userId" parameter`
	ErrInvalidAutoJoinProcent = `invalid "autoJoinProcent" value`
	ErrInvalidUserId          = "invalid userId"
)

type handler struct {
	userService
	historyService
	segmentService
}

func NewHandler(userService userService, historyService historyService, segmentService segmentService) *handler {
	return &handler{
		userService,
		historyService,
		segmentService,
	}
}

func New(us userService, hs historyService, ss segmentService) *gin.Engine {
	h := NewHandler(us, hs, ss)

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
