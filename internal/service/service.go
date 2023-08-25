package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/model"
)

type repository interface {
	CreateSegment(ctx context.Context, slug string) (int, error)
	DeleteSegment(ctx context.Context, slug string) error
	GetActiveUserSegmentsIDs(ctx context.Context, userId int) ([]int, error)
	RemoveUserFromSegment(ctx context.Context, segmentFromRemove, userId int) error
	AddUserToSegment(ctx context.Context, expirationTime *time.Time, segmentToAdd, userId int) error

	GetAllUsers(ctx context.Context) ([]int, error)
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

func (s *service) CreateSegment(ctx context.Context, segmentData model.AddSegment) error {
	addedSegmentId, err := s.repository.CreateSegment(ctx, segmentData.Slug)
	if err != nil {
		return err
	}

	usersIDs, err := s.repository.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	// Shuffle user IDs in random order
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(usersIDs), func(i, j int) { usersIDs[i], usersIDs[j] = usersIDs[j], usersIDs[i] })

	numUsersToAdd := segmentData.AutoJoinProcent * len(usersIDs) / 100
	for _, val := range usersIDs[:numUsersToAdd] {
		s.repository.AddUserToSegment(ctx, nil, addedSegmentId, val)
	}

	return nil
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
		err := s.repository.AddUserToSegment(ctx, userSegment.SegmentExpirationTime, val, userSegment.UserID)
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
