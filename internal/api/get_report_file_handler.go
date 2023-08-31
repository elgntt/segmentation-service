package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	response "github.com/elgntt/avito-internship-2023/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type parameters struct {
	Month  int
	Year   int
	UserId int
}

type responseUrl struct {
	URL string `json:"url"`
}

// GetReportFile
// @Summary GetReportFile
// @Tags History
// @Description Allows you to get a link to a csv file with the user's history for the transferred month-year
// @Produce application/json
// @Param 	month query int true "actual month"
// @Param 	year query int true "actual year"
// @Param 	userId query int true "actual userId"
// @Success 200 {object} api.responseUrl
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /history/file [get]
func (h *handler) GetReportFile(c *gin.Context) {
	ctx := context.Background()

	params, err := parseParameters(c.Query("month"), c.Query("year"), c.Query("userId"))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	filePath, err := h.historyService.GenerateCSVFile(ctx, params.Month, params.Year, params.UserId)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, responseUrl{
		URL: filePath,
	})
}

func parseParameters(monthQuery, yearQuery string, userIdQuery string) (parameters, error) {
	if yearQuery == "" {
		return parameters{}, app_err.NewBusinessError(ErrInvalidYearParameter)
	}
	if monthQuery == "" {
		return parameters{}, app_err.NewBusinessError(ErrInvalidMonthParameter)
	}
	if userIdQuery == "" {
		return parameters{}, app_err.NewBusinessError(ErrInvalidUserIdParameter)
	}

	var err error
	params := parameters{}

	params.Year, err = strconv.Atoi(yearQuery)
	if err != nil {
		return parameters{}, app_err.NewBusinessError(ErrInvalidYearParameter)
	}

	params.Month, err = strconv.Atoi(monthQuery)
	if err != nil {
		return parameters{}, app_err.NewBusinessError(ErrInvalidMonthParameter)
	}

	params.UserId, err = strconv.Atoi(userIdQuery)
	if err != nil {
		return parameters{}, app_err.NewBusinessError(ErrInvalidUserIdParameter)
	}

	return params, nil
}
