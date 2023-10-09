package service

import (
	"context"
	"errors"
	"testing"

	"github.com/elgntt/segmentation-service/internal/model"
	"github.com/golang/mock/gomock"
)

func TestHistoryService_DeleteExpiredUserSegments(t *testing.T) {
	repoError := "error from repo"
	userId := 100
	segmentSlugs := []string{"AVITO_TECH", "AVITO_DISCOUNT_11"}
	userSegments := []model.UsersSegments{
		{
			UserId:       100,
			SegmentSlugs: segmentSlugs,
		},
	}
	historyData := model.HistoryDataMultipleSegments{
		UserId:      userId,
		SegmentSlug: []string{"AVITO_TECH", "AVITO_DISCOUNT_11"},
		Operation:   removeOperationStr,
	}
	tests := []struct {
		name              string
		historyRepoBehave func(repository *MockHistoryRepo)
		segmentRepoBehave func(repository *MockSegmentRepo)
		wantErr           bool
	}{
		{
			name: "success",
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().DeleteExpiredUserSegments(gomock.Any()).Return(userSegments, nil)
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), historyData).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error from DeleteExpiredUserSegments()",
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().DeleteExpiredUserSegments(gomock.Any()).Return(nil, errors.New(repoError))
			},
			wantErr: true,
		},
		{
			name: "there are no deleted user segments",
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().DeleteExpiredUserSegments(gomock.Any()).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "error from RecordUserMultipleUserSegments()",
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().DeleteExpiredUserSegments(gomock.Any()).Return(userSegments, nil)
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), historyData).Return(errors.New(repoError))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockSegmentRepo := NewMockSegmentRepo(ctrl)
			mockHistoryRepo := NewMockHistoryRepo(ctrl)
			if tt.segmentRepoBehave != nil {
				tt.segmentRepoBehave(mockSegmentRepo)
			}
			if tt.historyRepoBehave != nil {
				tt.historyRepoBehave(mockHistoryRepo)
			}
			s := &HistoryService{
				historyRepo: mockHistoryRepo,
				segmentRepo: mockSegmentRepo,
			}
			if err := s.DeleteExpiredUserSegments(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("HistoryService.DeleteExpiredUserSegments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
