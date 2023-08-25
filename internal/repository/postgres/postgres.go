package repository

import (
	"context"

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

func (r *repo) CreateSegment(ctx context.Context, slug string) error {
	_, err := r.pool.Exec(ctx,
		` INSERT INTO segments (slug)
		  VALUES ($1)`, slug)

	return err
}

func (r *repo) DeleteSegment(ctx context.Context, slug string) error {
	_, err := r.pool.Exec(ctx,
		` DELETE FROM segments 
		  WHERE slug = $1`, slug)

	return err
}

func (r *repo) AddUserToSegment(ctx context.Context, segmentToAdd, userId int) error {
	_, err := r.pool.Exec(ctx,
		` INSERT INTO users_segments (user_id, segment_id)
		  VALUES ($1, $2) ON CONFLICT (user_id, segment_id) DO NOTHING`, userId, segmentToAdd)

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
		  WHERE user_id = $1`, userId)
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
	var segmentsSlug []string
	
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
