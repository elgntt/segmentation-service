package repository

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) AddUserToMultipleSegments(ctx context.Context, expirationTime *time.Time, segmentsSlugs []string, userId int) ([]string, error) {
	query := `
		INSERT INTO users_segments (user_id, segment_id, expiration_time)
		VALUES ($1, (
			SELECT id FROM segments WHERE slug = $2
		), $3)
		ON CONFLICT (user_id, segment_id) DO NOTHING`

	addedSlugs := make([]string, 0, len(segmentsSlugs))
	for _, segmentSlug := range segmentsSlugs {
		result, err := r.pool.Exec(ctx, query, userId, segmentSlug, expirationTime)
		if err != nil {
			return nil, err
		}
		if result.RowsAffected() == 1 {
			addedSlugs = append(addedSlugs, segmentSlug)
		}
	}

	return addedSlugs, nil
}

func (r *UserRepo) RemoveUserFromMultipleSegments(ctx context.Context, segmentsSlugsToRemove []string, userId int) ([]string, error) {
	slugArray := &pgtype.TextArray{}
	if err := slugArray.Set(segmentsSlugsToRemove); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`DELETE FROM users_segments
			WHERE user_id = $1
	  		AND segment_id IN (SELECT id FROM segments WHERE slug = ANY($2))
			RETURNING (SELECT slug FROM segments WHERE id = segment_id)`, userId, slugArray)
	if err != nil {
		return nil, err
	}

	deletedSegmentsSlugs := make([]string, 0, len(segmentsSlugsToRemove))

	for rows.Next() {
		var deletedSegment string
		if err := rows.Scan(&deletedSegment); err != nil {
			return nil, err
		}
		deletedSegmentsSlugs = append(deletedSegmentsSlugs, deletedSegment)
	}

	return deletedSegmentsSlugs, nil
}

func (r *UserRepo) GetPercentUsers(ctx context.Context, usersPercent int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		`WITH users AS (
			SELECT 
			  DISTINCT user_id 
			FROM 
			  users_segments
		  ) 
		  SELECT 
			  users.user_id 
		  FROM 
			  users 
		  ORDER BY 
			  RANDOM() 
		  LIMIT 
			  (SELECT COUNT(*) FROM USERS) * $1/100`, usersPercent)
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

func (r *UserRepo) GetActiveUserSegments(ctx context.Context, userId int) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		` SELECT segments.slug
			FROM users_segments us
			JOIN segments  ON us.segment_id = segments.id
			WHERE us.user_id = $1
			AND (us.expiration_time IS NULL OR us.expiration_time > CURRENT_TIMESTAMP)`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userSegmentSlugs := []string{}

	for rows.Next() {
		var segmentSlug string
		if err := rows.Scan(&segmentSlug); err != nil {
			return nil, err
		}

		userSegmentSlugs = append(userSegmentSlugs, segmentSlug)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userSegmentSlugs, nil
}
