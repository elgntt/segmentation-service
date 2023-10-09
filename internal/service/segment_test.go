package service

import (
	"context"
	"errors"
	"testing"

	"github.com/elgntt/segmentation-service/internal/model"
	gomock "github.com/golang/mock/gomock"
)

func TestSegmentService_CreateSegment(t *testing.T) {
	autoJoinPercent := 80
	segmentSlug := "test"
	tests := []struct {
		name              string
		segmentData       model.AddSegment
		segmentRepoBehave func(repository *MockSegmentRepo)
		historyRepoBehave func(repository *MockHistoryRepo)
		userRepoBehave    func(repository *MockUserRepo)
		wantErr           bool
	}{
		{
			name: "success",
			segmentData: model.AddSegment{
				SegmentSlug:     segmentSlug,
				AutoJoinPercent: autoJoinPercent,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().CreateSegment(gomock.Any(), segmentSlug).Return(1, nil)
				repository.EXPECT().AddMultipleUsersToSegment(gomock.Any(), 1, []int{1}).Return(nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetPercentUsers(gomock.Any(), autoJoinPercent).Return([]int{1}, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordMultipleUsersToHistory(gomock.Any(), model.HistoryDataMultipleUsers{
					UsersIDs:    []int{1},
					SegmentSlug: segmentSlug,
					Operation:   addOperationStr,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "0 percent to add users",
			segmentData: model.AddSegment{
				SegmentSlug:     segmentSlug,
				AutoJoinPercent: 0,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().CreateSegment(gomock.Any(), segmentSlug).Return(1, nil)
			},
			wantErr: false,
		},
		{
			name: "percentage of users not received",
			segmentData: model.AddSegment{
				SegmentSlug:     segmentSlug,
				AutoJoinPercent: autoJoinPercent,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().CreateSegment(gomock.Any(), segmentSlug).Return(1, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetPercentUsers(gomock.Any(), autoJoinPercent).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "error from repository to CreateSegment func",
			segmentData: model.AddSegment{
				SegmentSlug:     segmentSlug,
				AutoJoinPercent: autoJoinPercent,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().CreateSegment(gomock.Any(), segmentSlug).Return(0, errors.New("sql error"))
			},
			wantErr: true,
		},
		{
			name: "error from repository to GetProcentUsers func",
			segmentData: model.AddSegment{
				SegmentSlug:     segmentSlug,
				AutoJoinPercent: autoJoinPercent,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().CreateSegment(gomock.Any(), segmentSlug).Return(1, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetPercentUsers(gomock.Any(), autoJoinPercent).Return(nil, errors.New("sql error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockSegmentRepo := NewMockSegmentRepo(ctrl)
			mockUserRepo := NewMockUserRepo(ctrl)
			mockHistoryRepo := NewMockHistoryRepo(ctrl)
			if tt.segmentRepoBehave != nil {
				tt.segmentRepoBehave(mockSegmentRepo)
			}
			if tt.historyRepoBehave != nil {
				tt.historyRepoBehave(mockHistoryRepo)
			}
			if tt.userRepoBehave != nil {
				tt.userRepoBehave(mockUserRepo)
			}

			s := &SegmentService{
				segmentRepo: mockSegmentRepo,
				historyRepo: mockHistoryRepo,
				userRepo:    mockUserRepo,
			}

			err := s.CreateSegment(context.Background(), tt.segmentData)
			if (err != nil) != tt.wantErr {
				t.Errorf("SegmentService.CreateSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegmentService_DeleteSegment(t *testing.T) {
	deletedSegmentId := 1
	segmentSlug := "test"
	tests := []struct {
		name              string
		segmentSlug       string
		segmentRepoBehave func(repository *MockSegmentRepo)
		historyRepoBehave func(repository *MockHistoryRepo)
		wantErr           bool
	}{
		{
			name:        "success",
			segmentSlug: segmentSlug,
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().DeleteSegment(gomock.Any(), segmentSlug).Return(&deletedSegmentId, nil)
				repository.EXPECT().RemoveUsersFromDeletedSegment(gomock.Any(), deletedSegmentId).Return([]int{1}, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordMultipleUsersToHistory(gomock.Any(), model.HistoryDataMultipleUsers{
					UsersIDs:    []int{1},
					SegmentSlug: segmentSlug,
					Operation:   removeOperationStr,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "nil id deleted",
			segmentSlug: segmentSlug,
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().DeleteSegment(context.Background(), segmentSlug).Return(nil, errors.New("nothing was deleted"))
			},
			wantErr: true,
		},
		{
			name:        "none of the users have a segment deleted",
			segmentSlug: segmentSlug,
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().DeleteSegment(context.Background(), segmentSlug).Return(&deletedSegmentId, nil)
				repository.EXPECT().RemoveUsersFromDeletedSegment(gomock.Any(), deletedSegmentId)
			},
			wantErr: false,
		},
		{
			name:        "error from accessing the DeleteSegment() repository",
			segmentSlug: segmentSlug,
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().DeleteSegment(context.Background(), segmentSlug).Return(&deletedSegmentId, errors.New("sql error"))
			},
			wantErr: true,
		},
		{
			name:        "error from accessing the RemoveUsersFromDeletedSegment() repository",
			segmentSlug: segmentSlug,
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().DeleteSegment(context.Background(), segmentSlug).Return(&deletedSegmentId, nil)
				repository.EXPECT().RemoveUsersFromDeletedSegment(gomock.Any(), deletedSegmentId).Return(nil, errors.New("sql error"))
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

			s := &SegmentService{
				segmentRepo: mockSegmentRepo,
				historyRepo: mockHistoryRepo,
			}

			err := s.DeleteSegment(context.Background(), tt.segmentSlug)
			if (err != nil) != tt.wantErr {
				t.Errorf("SegmentService.DeleteSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
