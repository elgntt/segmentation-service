package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	"golang.org/x/exp/slices"
)

type UserService struct {
	userRepo    UserRepo
	segmentRepo SegmentRepo
	historyRepo HistoryRepo
}

func NewUserService(userRepo UserRepo, segmentRepo SegmentRepo, historyRepo HistoryRepo) *UserService {
	return &UserService{
		userRepo:    userRepo,
		segmentRepo: segmentRepo,
		historyRepo: historyRepo,
	}
}

func (s *UserService) GetActiveUserSegments(ctx context.Context, userId int) ([]string, error) {
	userSegments, err := s.userRepo.GetActiveUserSegments(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userSegments, nil
}

func (s *UserService) UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error {
	err := s.validateSegments(ctx, userSegment)
	if err != nil {
		return err
	}

	if len(userSegment.SegmentsSlugsToAdd) != 0 {
		err = s.AddUserToMultipleSegments(ctx, userSegment.SegmentExpirationTime, userSegment.SegmentsSlugsToAdd, userSegment.UserID)
		if err != nil {
			return err
		}
	}

	if len(userSegment.SegmentsSlugsToRemove) != 0 {
		return s.RemoveUserFromMultipleSegments(ctx, userSegment.SegmentsSlugsToRemove, userSegment.UserID)
	}

	return nil
}

func (s *UserService) AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsSlugs []string, userId int) error {
	addedSlugs, err := s.userRepo.AddUserToMultipleSegments(ctx, expirationTime, segmentsSlugs, userId)
	if err != nil {
		return err
	}

	if addedSlugs == nil {
		return nil
	}

	return s.RecordUserMultipleSegmentsToHistory(ctx, addedSlugs, addOperationStr, userId)
}

func (s *UserService) RemoveUserFromMultipleSegments(ctx context.Context, segmentsSlugs []string, userId int) error {
	deletedSegmentsSlugs, err := s.userRepo.RemoveUserFromMultipleSegments(ctx, segmentsSlugs, userId)
	if err != nil {
		return err
	}
	if deletedSegmentsSlugs == nil {
		return nil
	}

	return s.RecordUserMultipleSegmentsToHistory(ctx, deletedSegmentsSlugs, removeOperationStr, userId)
}

func (s *UserService) RecordUserMultipleSegmentsToHistory(ctx context.Context, segmentsSlugs []string, operation string, userId int) error {
	historyData := model.HistoryDataMultipleSegments{
		UserId:      userId,
		SegmentSlug: segmentsSlugs,
		Operation:   operation,
	}

	return s.historyRepo.RecordUserMultipleSegmentsToHistory(ctx, historyData)
}

func findAbsenceInSecondSlice(first, second []string) []string {
	var hash = make(map[string]bool, len(second))
	for _, elem := range second {
		hash[elem] = true
	}
	var nonMatchingStrings []string
	for _, elem := range first {
		if _, ok := hash[elem]; !ok {
			nonMatchingStrings = append(nonMatchingStrings, elem)
		}
	}

	return nonMatchingStrings
}

func (s *UserService) validateSegments(ctx context.Context, userSegment model.UserSegmentAction) error {
	totalUserSegments := append(userSegment.SegmentsSlugsToAdd, userSegment.SegmentsSlugsToRemove...)
	segments, err := s.segmentRepo.GetSegmentsBySlug(ctx, totalUserSegments)
	if err != nil {
		return err
	}

	missingSegments := findAbsenceInSecondSlice(totalUserSegments, segments)
	if missingSegments != nil {
		return app_err.NewBusinessError(fmt.Sprintf("these segments do not exist: [%s]", (strings.Join(slices.Compact(missingSegments), ", "))))
	}

	return nil
}
