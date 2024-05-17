package postgresql

import (
	"Booking/establishment-service-booking/internal/entity"
	"Booking/establishment-service-booking/internal/pkg/otlp"
	"Booking/establishment-service-booking/internal/pkg/postgres"
	"context"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	reviewTableName = "review_table"
	reviewServiceName = "reviewService"
	reviewSpanRepoPrefix = "reviewRepo"
)

type reviewRepo struct {
	reviewTableName string
	db              *postgres.PostgresDB
}

func NewReviewRepo(db *postgres.PostgresDB) *reviewRepo {
	return &reviewRepo{
		reviewTableName: reviewTableName,
		db:              db,
	}
}

func (r *reviewRepo) ReviewSelectQueryPrefix() squirrel.SelectBuilder {
	return r.db.Sq.Builder.Select(
		"review_id",
		"establishment_id",
		"user_id",
		"rating",
		"comment",
		"created_at",
		"updated_at",
	).From(r.reviewTableName)
}

// create a new review to an establishment
func (r *reviewRepo) CreateReview(ctx context.Context, review *entity.Review) (*entity.Review, error) {

	ctx, span := otlp.Start(ctx, reviewServiceName, reviewSpanRepoPrefix+"Create")
	defer span.End()

	data := map[string]interface{}{
		"review_id":        review.ReviewId,
		"establishment_id": review.EstablishmentId,
		"user_id":          review.UserId,
		"rating":           review.Rating,
		"comment":          review.Comment,
		"created_at":       time.Now().Local(),
		"updated_at":       time.Now().Local(),
	}

	query, args, err := r.db.Sq.Builder.Insert(r.reviewTableName).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var respReview entity.Review

	queryBuilder := r.ReviewSelectQueryPrefix().Where(r.db.Sq.Equal("review_id", review.ReviewId)).Where(r.db.Sq.Equal("deleted_at", nil))

	query, args, err = queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, query, args...).Scan(
		&respReview.ReviewId,
		&respReview.EstablishmentId,
		&respReview.UserId,
		&respReview.Rating,
		&respReview.Comment,
		&respReview.CreatedAt,
		&respReview.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &respReview, nil
}

// list reviews by establishment_id
func (r *reviewRepo) ListReviews(ctx context.Context, establishment_id string) ([]*entity.Review, uint64, error) {
	
	ctx, span := otlp.Start(ctx, reviewServiceName, reviewSpanRepoPrefix+"List")
	defer span.End()
	
	var reviews []*entity.Review

	queryBuilder := r.ReviewSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(r.db.Sq.Equal("establishment_id", establishment_id)).Where(r.db.Sq.Equal("deleted_at", nil))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var review entity.Review

		if err := rows.Scan(
			&review.ReviewId,
			&review.EstablishmentId,
			&review.UserId,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		reviews = append(reviews, &review)
	}

	var count uint64

	queryC := `SELECT COUNT(*) FROM review_table WHERE deleted_at is NULL`

	if err := r.db.QueryRow(ctx, queryC).Scan(&count); err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}

// delete review softly by review_id
func (r *reviewRepo) DeleteReview(ctx context.Context, review_id string) error {
	
	ctx, span := otlp.Start(ctx, reviewServiceName, reviewSpanRepoPrefix+"Delete")
	defer span.End()
	
	// Build the SQL query
	sqlStr, args, err := r.db.Sq.Builder.Update(r.reviewTableName).
		Set("deleted_at", time.Now().Local()).
		Where(r.db.Sq.Equal("review_id", review_id)).
		ToSql()
	if err != nil {
		return err
	}

	// Execute the SQL query
	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return err
	}

	return nil
}
