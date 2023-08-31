package service

import (
	"context"
	"encoding/csv"
	"github.com/elgntt/avito-internship-2023/internal/config"
	"os"
	"strconv"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/elgntt/avito-internship-2023/internal/pkg/app_err"
	"github.com/google/uuid"
)

type HistoryService struct {
	historyRepo HistoryRepo
	segmentRepo SegmentRepo
}

func NewHistoryService(historyRepo HistoryRepo, segmentRepo SegmentRepo) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		segmentRepo: segmentRepo,
	}
}

func (s *HistoryService) DeleteExpiredUserSegments(ctx context.Context) error {
	usersSegments, err := s.historyRepo.DeleteExpiredUserSegments(ctx)
	if err != nil {
		return err
	}

	for _, userSegments := range usersSegments {
		err = s.RecordUserMultipleSegmentsToHistory(ctx, userSegments.SegmentSlugs, removeOperationStr, userSegments.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HistoryService) GenerateCSVFile(ctx context.Context, month, year, userId int) (string, error) {
	history, err := s.historyRepo.GetHistory(ctx, month, year, userId)
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
			strconv.Itoa(historyRow.UserID),
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

	serverEndpoint := config.GetServerConfig().ServerEndpoint

	return serverEndpoint + filePath, nil
}

func (s *HistoryService) RecordUserMultipleSegmentsToHistory(ctx context.Context, segmentsSlugs []string, operation string, userId int) error {
	historyData := model.HistoryDataMultipleSegments{
		UserId:      userId,
		SegmentSlug: segmentsSlugs,
		Operation:   operation,
	}

	return s.historyRepo.RecordUserMultipleSegmentsToHistory(ctx, historyData)
}
