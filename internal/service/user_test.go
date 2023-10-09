package service

import (
	"context"
	"errors"
	reflect "reflect"
	"testing"
	"time"

	"github.com/elgntt/segmentation-service/internal/model"
	gomock "github.com/golang/mock/gomock"
)

func TestUserService_GetActiveUserSegments(t *testing.T) {
	userId := 100
	tests := []struct {
		name           string
		userRepoBehave func(repository *MockUserRepo)
		userId         int
		want           []string
		wantErr        bool
	}{
		{
			name:   "success",
			userId: userId,
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetActiveUserSegments(gomock.Any(), userId).Return([]string{"AVITO_TECH", "AVITO_DISCOUNT_30"}, nil)
			},

			want:    []string{"AVITO_TECH", "AVITO_DISCOUNT_30"},
			wantErr: false,
		},
		{
			name:   "no segments to user",
			userId: userId,
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetActiveUserSegments(gomock.Any(), userId).Return([]string{}, nil)
			},

			want:    []string{},
			wantErr: false,
		},
		{
			name:   "error from accessing the GetActiveUserSegments() repository",
			userId: userId,
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().GetActiveUserSegments(gomock.Any(), userId).Return(nil, errors.New("sql error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockSegmentRepo := NewMockSegmentRepo(ctrl)
			mockUserRepo := NewMockUserRepo(ctrl)
			mockHistoryRepo := NewMockHistoryRepo(ctrl)

			if tt.userRepoBehave != nil {
				tt.userRepoBehave(mockUserRepo)
			}

			s := &UserService{
				userRepo:    mockUserRepo,
				segmentRepo: mockSegmentRepo,
				historyRepo: mockHistoryRepo,
			}
			got, err := s.GetActiveUserSegments(context.Background(), tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetActiveUserSegments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetActiveUserSegments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UserSegmentAction(t *testing.T) {
	userId := 100
	segmentsToAdd := []string{"AVITO_TECH", "AVITO_DISCOUNT_30", "AVITO_DISCOUNT_11"}
	segmentsToRemove := []string{"AVITO_DISCOUNT_5", "AVITO_DISCOUNT_12"}
	notExistsSegments := []string{"RANDOM", "TEST", "SEGMENT"}
	allSegments := append(segmentsToAdd, segmentsToRemove...)
	expirationTime := time.Now().Add(10 * time.Hour)

	repoError := "repo error"
	tests := []struct {
		name              string
		userSegments      model.UserSegmentAction
		segmentRepoBehave func(repository *MockSegmentRepo)
		historyRepoBehave func(repository *MockHistoryRepo)
		userRepoBehave    func(repository *MockUserRepo)
		wantErr           bool
	}{
		{
			name: "success",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), allSegments).Return(allSegments, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().AddUserToMultipleSegments(gomock.Any(), &expirationTime, segmentsToAdd, userId).Return(segmentsToAdd, nil)
				repository.EXPECT().RemoveUserFromMultipleSegments(gomock.Any(), segmentsToRemove, userId).Return(segmentsToRemove, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToAdd,
					Operation:   addOperationStr,
				}).Return(nil)
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToRemove,
					Operation:   removeOperationStr,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "remove segments empty",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: []string{},
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), segmentsToAdd).Return(segmentsToAdd, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().AddUserToMultipleSegments(gomock.Any(), &expirationTime, segmentsToAdd, userId).Return(segmentsToAdd, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToAdd,
					Operation:   addOperationStr,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "add segments empty",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    []string{},
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), segmentsToRemove).Return(segmentsToRemove, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().RemoveUserFromMultipleSegments(gomock.Any(), segmentsToRemove, userId).Return(segmentsToRemove, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToRemove,
					Operation:   removeOperationStr,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error from GetSegmentsBySlug()",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), allSegments).Return(allSegments, errors.New(repoError))
			},
			wantErr: true,
		},
		{
			name: "error AddUserToMultipleSegments()",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), allSegments).Return(allSegments, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().AddUserToMultipleSegments(gomock.Any(), &expirationTime, segmentsToAdd, userId).Return(nil, errors.New(repoError))
			},
			wantErr: true,
		},
		{
			name: "error from RemoveUserFromMultipleSegments()",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), allSegments).Return(allSegments, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().AddUserToMultipleSegments(gomock.Any(), &expirationTime, segmentsToAdd, userId).Return(segmentsToAdd, nil)
				repository.EXPECT().RemoveUserFromMultipleSegments(gomock.Any(), segmentsToRemove, userId).Return(segmentsToRemove, errors.New(repoError))
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToAdd,
					Operation:   addOperationStr,
				}).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "error from RemoveUserFromMultipleSegments()",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), allSegments).Return(allSegments, nil)
			},
			userRepoBehave: func(repository *MockUserRepo) {
				repository.EXPECT().AddUserToMultipleSegments(gomock.Any(), &expirationTime, segmentsToAdd, userId).Return(segmentsToAdd, nil)
				repository.EXPECT().RemoveUserFromMultipleSegments(gomock.Any(), segmentsToRemove, userId).Return(segmentsToRemove, nil)
			},
			historyRepoBehave: func(repository *MockHistoryRepo) {
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToAdd,
					Operation:   addOperationStr,
				}).Return(nil)
				repository.EXPECT().RecordUserMultipleSegmentsToHistory(gomock.Any(), model.HistoryDataMultipleSegments{
					UserId:      userId,
					SegmentSlug: segmentsToRemove,
					Operation:   removeOperationStr,
				}).Return(errors.New(repoError))
			},
			wantErr: true,
		},
		{
			name: "not exists segments to add",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    notExistsSegments,
				SegmentsSlugsToRemove: segmentsToRemove,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), append(notExistsSegments, segmentsToRemove...)).Return(nil, errors.New(repoError))
			},
			wantErr: true,
		},
		{
			name: "not exists segments to remove",
			userSegments: model.UserSegmentAction{
				UserID:                userId,
				SegmentsSlugsToAdd:    segmentsToAdd,
				SegmentsSlugsToRemove: notExistsSegments,
				SegmentExpirationTime: &expirationTime,
			},
			segmentRepoBehave: func(repository *MockSegmentRepo) {
				repository.EXPECT().GetSegmentsBySlug(gomock.Any(), append(segmentsToAdd, notExistsSegments...)).Return(nil, errors.New(repoError))
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

			s := &UserService{
				userRepo:    mockUserRepo,
				segmentRepo: mockSegmentRepo,
				historyRepo: mockHistoryRepo,
			}
			if err := s.UserSegmentAction(context.Background(), tt.userSegments); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UserSegmentAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findAbsenceInSecondSlice(t *testing.T) {
	longer := []string{"a", "b", "c"}
	smaller := []string{"a", "b"}
	tests := []struct {
		name   string
		first  []string
		second []string
		want   []string
	}{
		{
			name:   "match",
			first:  smaller,
			second: smaller,
			want:   nil,
		},
		{
			name:   "absent",
			first:  longer,
			second: smaller,
			want:   []string{"c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAbsenceInSecondSlice(tt.first, tt.second); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findNonMatchingInArrays() = %v, want %v", got, tt.want)
			}
		})
	}
}
