//go:generate mockgen -source=$GOFILE -destination=mocks_test.go -package=$GOPACKAGE
package service

import (
	"context"
	"time"

	"github.com/elgntt/segmentation-service/internal/model"
)

type SegmentRepo interface {
	CreateSegment(ctx context.Context, slug string) (int, error)
	DeleteSegment(ctx context.Context, slug string) (*int, error)
	AddMultipleUsersToSegment(ctx context.Context, segmentId int, usersIDs []int) error

	GetSegmentsBySlug(ctx context.Context, slugs []string) ([]string, error)

	RemoveUsersFromDeletedSegment(ctx context.Context, sigmentId int) ([]int, error)
}

type HistoryRepo interface {
	RecordUserMultipleSegmentsToHistory(ctx context.Context, historyData model.HistoryDataMultipleSegments) error
	RecordMultipleUsersToHistory(ctx context.Context, historyData model.HistoryDataMultipleUsers) error
	DeleteExpiredUserSegments(ctx context.Context) ([]model.UsersSegments, error)
	GetHistory(ctx context.Context, month, year, userId int) ([]model.History, error)
}

type UserRepo interface {
	GetActiveUserSegments(ctx context.Context, userId int) ([]string, error)
	RemoveUserFromMultipleSegments(ctx context.Context, segmentsSlugsToRemove []string, userId int) ([]string, error)
	AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsSlugs []string, userId int) ([]string, error)
	GetPercentUsers(ctx context.Context, usersPercent int) ([]int, error)
}

const (
	addOperationStr    = "adding"
	removeOperationStr = "removal"

	csvFilesDir = "assets/csv_reports/"
)

const (
	ErrNoDataAvailable = "no data available"
)
