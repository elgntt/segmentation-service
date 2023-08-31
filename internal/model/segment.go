package model

import "time"

type AddSegment struct {
	SegmentSlug     string `json:"slug"`
	AutoJoinPercent int    `json:"autoJoinPercent"`
}

type UserSegmentAction struct {
	UserID                int        `json:"userId"`
	SegmentsSlugsToAdd    []string   `json:"segmentsToAdd"`
	SegmentsSlugsToRemove []string   `json:"segmentsToRemove"`
	SegmentExpirationTime *time.Time `json:"expirationTime,omitempty"`
}

type UsersSegments struct {
	UserId       int
	SegmentSlugs []string
}

type History struct {
	UserID        int
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
