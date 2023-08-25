package model

import "time"

type Segment struct {
	Slug string `json:"slug"`
}

type UserSegmentAction struct {
	UserID                int        `json:"userId"`
	SegmentsToAdd         []string   `json:"segmentsToAdd"`
	SegmentsToRemove      []string   `json:"segmentsToRemove"`
	SegmentExpirationTime *time.Time `json:"expirationTime,omitempty"`
}
