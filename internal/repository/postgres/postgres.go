package repository

import (
	"context"
	"errors"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/model"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repo {
	return &repo{
		pool: pool,
	}
}

func (r *repo) CreateSegment(ctx context.Context, slug string) (int, error) {
	row := r.pool.QueryRow(ctx,
		` INSERT INTO segments (slug)
		  VALUES ($1)
		  RETURNING id`, slug)

	var segmentId int

	err := row.Scan(&segmentId)
	if err != nil {
		return 0, err
	}

	return segmentId, nil
}

func (r *repo) DeleteSegment(ctx context.Context, slug string) (*int, error) {
	var removedSegmentId int
	err := r.pool.QueryRow(ctx,
		` DELETE FROM segments 
		  WHERE slug = $1
		  RETURNING id`, slug).Scan(&removedSegmentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &removedSegmentId, err
}

func (r *repo) RemoveUsersFromDeletedSegment(ctx context.Context, sigmentId int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		` DELETE FROM users_segments
		  WHERE segment_id = $1
		  RETURNING user_id`, sigmentId)

	if err != nil {
		return nil, err
	}
	var usersIDs []int
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		usersIDs = append(usersIDs, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usersIDs, nil
}

func (r *repo) AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsIDsToAdd []int, userId int) error {
	query := `
		INSERT INTO users_segments (user_id, segment_id, expiration_time)
		VALUES ($1, $2, $3) ON CONFLICT (user_id, segment_id) DO UPDATE
		SET expiration_time = excluded.expiration_time`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, segmentId := range segmentsIDsToAdd {
		_, err := tx.Exec(ctx, query, userId, segmentId, expirationTime)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repo) AddMultipleUsersToSegment(ctx context.Context, segmentId int, usersIDs []int) error {
	query := `
	INSERT INTO users_segments (user_id, segment_id)
	VALUES ($1, $2)`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, userId := range usersIDs {
		_, err := tx.Exec(ctx, query, userId, segmentId)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repo) RemoveUserFromMultipleSegments(ctx context.Context, segmentsIDsToRemove []int, userId int) error {
	query := `
		  DELETE FROM users_segments 
			WHERE user_id = $1 
			AND segment_id = $2`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, segmentId := range segmentsIDsToRemove {
		_, err := tx.Exec(ctx, query, userId, segmentId)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repo) GetActiveUserSegmentsIDs(ctx context.Context, userId int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		` SELECT segment_id 
		  FROM users_segments
		  WHERE user_id = $1
		  AND (expiration_time IS NULL OR expiration_time > CURRENT_TIMESTAMP)`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSegmentsIDs []int

	for rows.Next() {
		var segmentId int
		if err := rows.Scan(&segmentId); err != nil {
			return nil, err
		}

		userSegmentsIDs = append(userSegmentsIDs, segmentId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegmentsIDs, nil
}

func (r *repo) GetSlugsByIDs(ctx context.Context, segmentsIDs []int) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		` SELECT slug
		  FROM segments
		  WHERE id = ANY($1)`, segmentsIDs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	segmentsSlug := []string{}
	for rows.Next() {
		var segment string
		if err := rows.Scan(&segment); err != nil {
			return nil, err
		}

		segmentsSlug = append(segmentsSlug, segment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segmentsSlug, nil
}

func (r *repo) GetIdsBySlugs(ctx context.Context, slugs []string) ([]int, []string, error) {
	slugArray := &pgtype.TextArray{}
	if err := slugArray.Set(slugs); err != nil {
		return nil, nil, err
	}

	rows, err := r.pool.Query(ctx,
		` SELECT id, slug 
		  FROM segments 
		  WHERE slug = ANY($1)`, slugArray)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	segmentsSlugsIds := make([]int, 0)
	segmentsSlugs := make([]string, 0)
	for rows.Next() {
		var id int
		var slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return nil, nil, err
		}
		segmentsSlugsIds = append(segmentsSlugsIds, id)
		segmentsSlugs = append(segmentsSlugs, slug)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return segmentsSlugsIds, segmentsSlugs, nil
}

func (r *repo) GetAllUsers(ctx context.Context) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT user_id 
		 FROM users_segments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func (r repo) DeleteExpiredUserSegments(ctx context.Context) (map[int][]int, error) {
	rows, err := r.pool.Query(ctx,
		` DELETE FROM users_segments
		  WHERE expiration_time IS NOT NULL 
		  AND expiration_time <= CURRENT_TIMESTAMP
		  RETURNING user_id, segment_id`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	usersSegments := make(map[int][]int)
	for rows.Next() {
		var userID int
		var segmentID int
		if err := rows.Scan(&userID, &segmentID); err != nil {
			return nil, err
		}
		usersSegments[userID] = append(usersSegments[userID], segmentID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usersSegments, nil
}

func (r *repo) AddUserMultipleSegmentsToHistory(ctx context.Context, historyData model.HistoryDataMultipleSegments) error {
	if len(historyData.SegmentSlug) == 0 {
		return nil
	}

	query := `
		INSERT INTO user_segment_history (user_id, segment_slug, operation)
		VALUES ($1, $2, $3)`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, segmentSlug := range historyData.SegmentSlug {
		_, err := tx.Exec(ctx, query, historyData.UserId, segmentSlug, historyData.Operation)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repo) AddMultipleUsersToHistory(ctx context.Context, historyData model.HistoryDataMultipleUsers) error {
	if len(historyData.SegmentSlug) == 0 {
		return nil
	}

	query := `
		INSERT INTO user_segment_history (user_id, segment_slug, operation)
		VALUES ($1, $2, $3)`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, userId := range historyData.UsersIDs {
		_, err := tx.Exec(ctx, query, userId, historyData.SegmentSlug, historyData.Operation)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repo) GetHistory(ctx context.Context, month, year, userId int) ([]model.History, error) {
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
		if err := rows.Scan(&historyRow.UserId, &historyRow.SegmentSlug, &historyRow.Operation, &historyRow.OperationTime); err != nil {
			return nil, err
		}
		history = append(history, historyRow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
