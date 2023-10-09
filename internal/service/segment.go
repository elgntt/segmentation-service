package service

import (
	"context"

	"github.com/elgntt/segmentation-service/internal/model"
	"github.com/elgntt/segmentation-service/internal/pkg/app_err"
)

type SegmentService struct {
	segmentRepo SegmentRepo
	historyRepo HistoryRepo
	userRepo    UserRepo
}

func NewSegmentService(segmentRepo SegmentRepo, historyRepo HistoryRepo, userRepo UserRepo) *SegmentService {
	return &SegmentService{
		segmentRepo: segmentRepo,
		historyRepo: historyRepo,
		userRepo:    userRepo,
	}
}
func (s *SegmentService) CreateSegment(ctx context.Context, segmentData model.AddSegment) error {
	addedSegmentId, err := s.segmentRepo.CreateSegment(ctx, segmentData.SegmentSlug)
	if err != nil {
		return err
	}

	if segmentData.AutoJoinPercent == 0 {
		return nil
	}

	usersIDs, err := s.userRepo.GetPercentUsers(ctx, segmentData.AutoJoinPercent)
	if err != nil {
		return err
	}

	if usersIDs == nil {
		return nil
	}

	return s.AddMultipleUsersToSegment(ctx, addedSegmentId, segmentData.SegmentSlug, usersIDs)
}

func (s *SegmentService) DeleteSegment(ctx context.Context, segmentSlug string) error {
	removedSegmentId, err := s.segmentRepo.DeleteSegment(ctx, segmentSlug)
	if err != nil {
		return err
	}

	if removedSegmentId == nil {
		return app_err.NewBusinessError("segment does not exists")
	}

	usersIDs, err := s.segmentRepo.RemoveUsersFromDeletedSegment(ctx, *removedSegmentId)
	if err != nil {
		return err
	}

	if usersIDs == nil {
		return nil
	}

	return s.RecordMultipleUsersToHistory(ctx, segmentSlug, removeOperationStr, usersIDs)
}

func (s *SegmentService) AddMultipleUsersToSegment(ctx context.Context, segmentId int, segmentSlug string, usersIDs []int) error {
	err := s.segmentRepo.AddMultipleUsersToSegment(ctx, segmentId, usersIDs)
	if err != nil {
		return err
	}

	return s.RecordMultipleUsersToHistory(ctx, segmentSlug, addOperationStr, usersIDs)
}

func (s *SegmentService) RecordMultipleUsersToHistory(ctx context.Context, segmentSlug, operation string, usersIDs []int) error {
	historyData := model.HistoryDataMultipleUsers{
		UsersIDs:    usersIDs,
		SegmentSlug: segmentSlug,
		Operation:   operation,
	}

	return s.historyRepo.RecordMultipleUsersToHistory(ctx, historyData)
}
