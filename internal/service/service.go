package service

import (
	"context"
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	"github.com/google/uuid"
)

type repository interface {
	CreateSegment(ctx context.Context, slug string) (int, error)
	DeleteSegment(ctx context.Context, slug string) (*int, error)
	GetActiveUserSegmentsIDs(ctx context.Context, userId int) ([]int, error)
	RemoveUserFromMultipleSegments(ctx context.Context, segmentsIDsToRemove []int, userId int) error
	AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsIDsToAdd []int, userId int) error
	AddMultipleUsersToSegment(ctx context.Context, segmentId int, usersIDs []int) error

	RemoveUsersFromDeletedSegment(ctx context.Context, sigmentId int) ([]int, error)
	GetAllUsers(ctx context.Context) ([]int, error)
	GetIdsBySlugs(ctx context.Context, slugs []string) ([]int, []string, error)
	GetSlugsByIDs(ctx context.Context, segmentsIDs []int) ([]string, error)

	AddUserMultipleSegmentsToHistory(ctx context.Context, historyData model.HistoryDataMultipleSegments) error
	AddMultipleUsersToHistory(ctx context.Context, historyData model.HistoryDataMultipleUsers) error

	DeleteExpiredUserSegments(ctx context.Context) (map[int][]int, error)

	GetHistory(ctx context.Context, month, year, userId int) ([]model.History, error)
}
type Service struct {
	repository
}

const (
	addOperationStr    = "adding"
	removeOperationStr = "removal"

	csvFilesDir = "./assets/csv_reports/"
)

const (
	ErrNoDataAvailable = "no data available"
)

func New(repo repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) CreateSegment(ctx context.Context, segmentData model.AddSegment) error {
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

	err = s.AddMultipleUsersToSegment(ctx, addedSegmentId, segmentData.Slug, usersIDs[:numUsersToAdd])
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteSegment(ctx context.Context, slug string) error {
	removedSegmentId, err := s.repository.DeleteSegment(ctx, slug)
	if err != nil {
		return err
	}
	if removedSegmentId == nil {
		return nil
	}

	usersIDs, err := s.repository.RemoveUsersFromDeletedSegment(ctx, *removedSegmentId)
	if err != nil {
		return err
	}

	if usersIDs == nil {
		return nil
	}

	historyData := model.HistoryDataMultipleUsers{
		UsersIDs:    usersIDs,
		SegmentSlug: slug,
		Operation:   removeOperationStr,
	}

	return s.repository.AddMultipleUsersToHistory(ctx, historyData)
}

func (s *Service) GetActiveUserSegments(ctx context.Context, userId int) ([]string, error) {
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

func (s *Service) UserSegmentAction(ctx context.Context, userSegment model.UserSegmentAction) error {
	segmentsIdToAdd, segmentsSlugsToAdd, err := s.repository.GetIdsBySlugs(ctx, userSegment.SegmentsToAdd)
	if err != nil {
		return err
	}

	err = s.AddUserToMultipleSegments(ctx, userSegment.SegmentExpirationTime, segmentsSlugsToAdd, segmentsIdToAdd, userSegment.UserID)
	if err != nil {
		return err
	}

	segmentsIdsToRemove, segmentsSlugsToRemove, err := s.repository.GetIdsBySlugs(ctx, userSegment.SegmentsToRemove)
	if err != nil {
		return err
	}

	return s.RemoveUserFromMultipleSegments(ctx, segmentsIdsToRemove, segmentsSlugsToRemove, userSegment.UserID)
}

func (s *Service) AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsSlugs []string, segmentsIDs []int, userId int) error {
	if segmentsSlugs == nil || segmentsIDs == nil {
		return nil
	}

	err := s.repository.AddUserToMultipleSegments(ctx, expirationTime, segmentsIDs, userId)
	if err != nil {
		return err
	}

	historyData := model.HistoryDataMultipleSegments{
		UserId:      userId,
		SegmentSlug: segmentsSlugs,
		Operation:   addOperationStr,
	}
	err = s.repository.AddUserMultipleSegmentsToHistory(ctx, historyData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveUserFromMultipleSegments(ctx context.Context, segmentsIDsToRemove []int, segmentsSlugs []string, userId int) error {
	if segmentsIDsToRemove == nil || segmentsSlugs == nil {
		return nil
	}
	err := s.repository.RemoveUserFromMultipleSegments(ctx, segmentsIDsToRemove, userId)
	if err != nil {
		return err
	}

	historyData := model.HistoryDataMultipleSegments{
		UserId:      userId,
		SegmentSlug: segmentsSlugs,
		Operation:   removeOperationStr,
	}

	err = s.repository.AddUserMultipleSegmentsToHistory(ctx, historyData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddMultipleUsersToSegment(ctx context.Context, segmentId int, segmentSlug string, usersIDs []int) error {
	if usersIDs == nil {
		return nil
	}

	err := s.repository.AddMultipleUsersToSegment(ctx, segmentId, usersIDs)
	if err != nil {
		return err
	}

	historyData := model.HistoryDataMultipleUsers{
		UsersIDs:    usersIDs,
		SegmentSlug: segmentSlug,
		Operation:   addOperationStr,
	}

	err = s.repository.AddMultipleUsersToHistory(ctx, historyData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteExpiredUserSegments(ctx context.Context) error {
	usersSegments, err := s.repository.DeleteExpiredUserSegments(ctx)
	if err != nil {
		return err
	}

	for userId, userSegmentsIDs := range usersSegments {
		segmentsSlug, err := s.repository.GetSlugsByIDs(ctx, userSegmentsIDs)
		if err != nil {
			return err
		}

		historyData := model.HistoryDataMultipleSegments{
			UserId:      userId,
			SegmentSlug: segmentsSlug,
			Operation:   removeOperationStr,
		}
		err = s.repository.AddUserMultipleSegmentsToHistory(ctx, historyData)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Service) GenerateCSVFile(ctx context.Context, month, year, userId int) (string, error) {
	history, err := s.repository.GetHistory(ctx, month, year, userId)
	if err != nil {
		return "", err
	}

	if history == nil {
		return "", app_err.NewBusinessError(ErrNoDataAvailable)
	}

	filePath := csvFilesDir + uuid.NewString() + ".csv"
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)

	csvWriter.Comma = ';'

	for _, historyRow := range history {
		if err := csvWriter.Write([]string{
			strconv.Itoa(historyRow.UserId),
			historyRow.SegmentSlug,
			historyRow.Operation,
			historyRow.OperationTime.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return "", err
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", err
	}

	return filePath, nil
}
