package model

import "time"

type AddSegment struct {
	Slug            string `json:"slug"`
	AutoJoinProcent int    `json:"autoJoinProcent"`
}

type UserSegmentAction struct {
	UserID                int        `json:"userId"`
	SegmentsToAdd         []string   `json:"segmentsToAdd"`
	SegmentsToRemove      []string   `json:"segmentsToRemove"`
	SegmentExpirationTime *time.Time `json:"expirationTime,omitempty"`
}

type History struct {
	UserId        int
	SegmentSlug   string
	Operation     string
	OperationTime time.Time
}
type HistoryDataMultipleSegments struct {
	UserId      int
	SegmentSlug []string
	Operation   string
}

type HistoryDataMultipleUsers struct {
	UsersIDs    []int
	SegmentSlug string
	Operation   string
}
