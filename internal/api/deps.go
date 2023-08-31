package api

import (
	"context"

	"github.com/elgntt/avito-internship-2023/internal/model"
)

// todo make private
type userService interface {
	GetActiveUserSegments(ctx context.Context, userId int) ([]string, error)
	UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error
}

type segmentService interface {
	CreateSegment(ctx context.Context, segmentData model.AddSegment) error
	DeleteSegment(ctx context.Context, slug string) error
}

type historyService interface {
	GenerateCSVFile(ctx context.Context, month, year, userId int) (string, error)
}
