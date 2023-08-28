package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetReportFile(c *gin.Context) {
	ctx := context.Background()

	month, year, userId, err := parseParameters(c.Query("month"), c.Query("year"), c.Query("userId"))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	filePath, err := h.service.GenerateCSVFile(ctx, month, year, userId)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": "http://localhost:8080/" + filePath,
	})
}

func parseParameters(monthQuery, yearQuery string, userIdQuery string) (int, int, int, error) {
	if yearQuery == "" {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidYearParameter)
	}
	if monthQuery == "" {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidMonthParameter)
	}
	if userIdQuery == "" {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidUserIdParameter)
	}

	year, err := strconv.Atoi(yearQuery)
	if err != nil {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidYearParameter)
	}

	month, err := strconv.Atoi(monthQuery)
	if err != nil {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidMonthParameter)
	}

	userId, err := strconv.Atoi(userIdQuery)
	if err != nil {
		return 0, 0, 0, app_err.NewBusinessError(ErrInvalidUserIdParameter)
	}

	return month, year, userId, nil
}
