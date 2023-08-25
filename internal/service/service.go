package service

import (
	"context"
	"github.com/elgntt/avito-internship-2023/internal/model"
)

type repository interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	GetActiveUserSegmentsIDs(ctx context.Context, userId int) ([]int, error)
	RemoveUserFromSegment(ctx context.Context, segmentFromRemove, userId int) error
	AddUserToSegment(ctx context.Context, segmentToAdd, userId int) error

	GetIdBySlugs(ctx context.Context, slugs []string) ([]int, error)
	GetSlugsByIDs(ctx context.Context, segmentsIDs []int) ([]string, error)
}
type service struct {
	repository
}

func New(repo repository) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) CreateSegment(ctx context.Context, slug string) error {
	return s.repository.CreateSegment(ctx, slug)
}

func (s *service) DeleteSegment(ctx context.Context, slug string) error {
	return s.repository.DeleteSegment(ctx, slug)
}

func (s *service) GetActiveUserSegments(ctx context.Context, userId int) ([]string, error) {
	userSegmentsIDs, err := s.repository.GetActiveUserSegmentsIDs(ctx, userId)
	if err != nil {
		return nil, err
	}

	segmentSlugs, err := s.repository.GetSlugsByIDs(ctx, userSegmentsIDs)
	if err != nil {
		return nil, err
	}

	return segmentSlugs, nil
}

func (s *service) UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error {
	segmentsIdToAdd, err := s.repository.GetIdBySlugs(ctx, userSegment.SegmentsToAdd)
	if err != nil {
		return err
	}
	for _, val := range segmentsIdToAdd {
		err := s.repository.AddUserToSegment(ctx, val, userSegment.UserID)
		if err != nil {
			return err
		}
	}

	segmentsIdsToRemove, err := s.repository.GetIdBySlugs(ctx, userSegment.SegmentsToRemove)
	if err != nil {
		return err
	}
	for _, val := range segmentsIdsToRemove {
		err := s.repository.RemoveUserFromSegment(ctx, val, userSegment.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}
