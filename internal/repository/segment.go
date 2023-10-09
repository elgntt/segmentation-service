package repository

import (
	"context"
	"errors"
	"github.com/elgntt/segmentation-service/internal/pkg/app_err"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentRepo struct {
	pool *pgxpool.Pool
}

func NewSegmentRepo(pool *pgxpool.Pool) *SegmentRepo {
	return &SegmentRepo{
		pool: pool,
	}
}

func (r *SegmentRepo) CreateSegment(ctx context.Context, slug string) (int, error) {
	row := r.pool.QueryRow(ctx,
		` INSERT INTO segments (slug)
		  VALUES ($1)
		  RETURNING id`, slug)

	var segmentId int

	err := row.Scan(&segmentId)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return 0, app_err.NewBusinessError("segment slug already exists")
			}
		}
		return 0, err
	}

	return segmentId, nil
}

func (r *SegmentRepo) DeleteSegment(ctx context.Context, slug string) (*int, error) {
	var removedSegmentId *int
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

	return removedSegmentId, err
}

func (r *SegmentRepo) RemoveUsersFromDeletedSegment(ctx context.Context, sigmentId int) ([]int, error) {
	rows, err := r.pool.Query(ctx,
		` DELETE FROM users_segments
		  WHERE segment_id = $1
		  RETURNING user_id`, sigmentId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *SegmentRepo) AddMultipleUsersToSegment(ctx context.Context, segmentId int, usersIDs []int) error {
	query := `
		INSERT INTO users_segments (user_id, segment_id)
		VALUES ($1, $2)`

	for _, userId := range usersIDs {
		_, err := r.pool.Exec(ctx, query, userId, segmentId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SegmentRepo) GetSegmentsBySlug(ctx context.Context, slugs []string) ([]string, error) {
	slugArray := &pgtype.TextArray{}
	if err := slugArray.Set(slugs); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		` SELECT slug 
		  FROM segments 
		  WHERE slug = ANY($1)`, slugArray)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	segments := make([]string, 0, len(slugs))
	for rows.Next() {
		var segment string
		if err := rows.Scan(&segment); err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}
