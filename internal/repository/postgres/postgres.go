package repository

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
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

func (r *repo) DeleteSegment(ctx context.Context, slug string) error {
	_, err := r.pool.Exec(ctx,
		` DELETE FROM segments 
		  WHERE slug = $1`, slug)

	return err
}

func (r *repo) AddUserToSegment(ctx context.Context, expirationTime *time.Time, segmentToAdd, userId int) error {
	_, err := r.pool.Exec(ctx,
		` INSERT INTO users_segments (user_id, segment_id, expiration_time)
		  VALUES ($1, $2, $3) ON CONFLICT (user_id, segment_id) DO UPDATE
		  	SET expiration_time = excluded.expiration_time`, userId, segmentToAdd, expirationTime)

	return err
}

func (r *repo) RemoveUserFromSegment(ctx context.Context, segmentFromRemove, userId int) error {
	_, err := r.pool.Exec(ctx,
		` DELETE FROM users_segments 
		  WHERE user_id = $1 
		  AND segment_id = $2`, userId, segmentFromRemove)

	return err
}

func (r *repo) GetActiveUserSegmentsIDs(ctx context.Context, userId int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		` SELECT segment_id 
		  FROM users_segments
		  WHERE user_id = $1
		  AND (expiration_time IS NULL OR expiration_time >= CURRENT_TIMESTAMP)`, userId)
	if err != nil {
		return nil, err
	}

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

func (r *repo) GetIdBySlugs(ctx context.Context, slugs []string) ([]int, error) {
	slugArray := &pgtype.TextArray{}
	if err := slugArray.Set(slugs); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		` SELECT id 
		  FROM segments 
		  WHERE slug = ANY($1)`, slugArray)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
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
