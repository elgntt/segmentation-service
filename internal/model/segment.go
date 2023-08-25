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
