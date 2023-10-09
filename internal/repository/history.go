package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/elgntt/segmentation-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryRepo struct {
	pool *pgxpool.Pool
}

func NewHistoryRepo(pool *pgxpool.Pool) *HistoryRepo {
	return &HistoryRepo{
		pool: pool,
	}
}

func (r *HistoryRepo) DeleteExpiredUserSegments(ctx context.Context) ([]model.UsersSegments, error) {
	rows, err := r.pool.Query(ctx,
		` WITH deleted_segments AS (
				DELETE FROM users_segments
				WHERE expiration_time IS NOT NULL
				AND expiration_time <= CURRENT_TIMESTAMP
				RETURNING user_id, segment_id
			)
			SELECT d.user_id,
				   array_agg(s.slug) AS segment_slugs
			FROM deleted_segments d
			JOIN segments s ON d.segment_id = s.id
			GROUP BY d.user_id`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	usersSegments := []model.UsersSegments{}
	for rows.Next() {
		userSegmentsTemp := model.UsersSegments{}
		if err := rows.Scan(&userSegmentsTemp.UserId, &userSegmentsTemp.SegmentSlugs); err != nil {
			return nil, fmt.Errorf("error zdes:%w", err)
		}
		usersSegments = append(usersSegments, userSegmentsTemp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usersSegments, nil
}

func (r *HistoryRepo) RecordUserMultipleSegmentsToHistory(ctx context.Context, historyData model.HistoryDataMultipleSegments) (err error) {
	query := `
		INSERT INTO user_segment_history (user_id, segment_slug, operation)
		VALUES ($1, $2, $3)`

	for _, segmentSlug := range historyData.SegmentSlug {
		_, err := r.pool.Exec(ctx, query, historyData.UserId, segmentSlug, historyData.Operation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *HistoryRepo) RecordMultipleUsersToHistory(ctx context.Context, historyData model.HistoryDataMultipleUsers) error {
	query := `
		INSERT INTO user_segment_history (user_id, segment_slug, operation)
		VALUES ($1, $2, $3)`

	for _, userId := range historyData.UsersIDs {
		_, err := r.pool.Exec(ctx, query, userId, historyData.SegmentSlug, historyData.Operation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *HistoryRepo) GetHistory(ctx context.Context, month, year, userId int) ([]model.History, error) {
	rows, err := r.pool.Query(ctx,
		` SELECT user_id, segment_slug, operation, operation_time
			FROM user_segment_history
			WHERE DATE_TRUNC('month', operation_time) = $1::date
			AND user_id = $2
			ORDER BY operation_time`, time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var history []model.History
	for rows.Next() {
		var historyRow model.History
		if err := rows.Scan(&historyRow.UserID, &historyRow.SegmentSlug, &historyRow.Operation, &historyRow.OperationTime); err != nil {
			return nil, err
		}
		history = append(history, historyRow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
