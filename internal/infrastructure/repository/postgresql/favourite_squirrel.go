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
	favouriteTableName = "favourite_table"
	favouriteServiceName    = "favouriteService"
	favouriteSpanRepoPrefix = "favouriteRepo"
)

type favouriteRepo struct {
	favouriteTableName string
	db                 *postgres.PostgresDB
}

func NewFavouriteRepo(db *postgres.PostgresDB) *favouriteRepo {
	return &favouriteRepo{
		favouriteTableName: favouriteTableName,
		db:                 db,
	}
}

func (f *favouriteRepo) FavouriteSelectQueryPrefix() squirrel.SelectBuilder {
	return f.db.Sq.Builder.Select(
		"favourite_id",
		"establishment_id",
		"user_id",
		"created_at",
		"updated_at",
	).From(f.favouriteTableName)
}

func (f *favouriteRepo) AddToFavourites(ctx context.Context, favourite *entity.Favourite) (*entity.Favourite, error) {

	ctx, span := otlp.Start(ctx, favouriteServiceName, favouriteSpanRepoPrefix+"Create")
	defer span.End()

	data := map[string]interface{}{
		"favourite_id":     favourite.FavouriteId,
		"establishment_id": favourite.EstablishmentId,
		"user_id":          favourite.UserId,
		"created_at":       time.Now().Local(),
		"updated_at":       time.Now().Local(),
	}

	query, args, err := f.db.Sq.Builder.Insert(favouriteTableName).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = f.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var respFavourite entity.Favourite

	queryBuilder := f.FavouriteSelectQueryPrefix().Where(f.db.Sq.Equal("favourite_id", favourite.FavouriteId)).Where(f.db.Sq.Equal("deleted_at", nil))

	query, args, err = queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	if err := f.db.QueryRow(ctx, query, args...).Scan(
		&respFavourite.FavouriteId,
		&respFavourite.EstablishmentId,
		&respFavourite.UserId,
		&respFavourite.CreatedAt,
		&respFavourite.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &respFavourite, nil
}

func (f *favouriteRepo) RemoveFromFavourites(ctx context.Context, favourite_id string) error {
	
	ctx, span := otlp.Start(ctx, favouriteServiceName, favouriteSpanRepoPrefix+"Delete")
	defer span.End()
	
	// Build the SQL query
	sqlStr, args, err := f.db.Sq.Builder.Update(f.favouriteTableName).
		Set("deleted_at", time.Now().Local()).
		Where(f.db.Sq.Equal("favourite_id", favourite_id)).
		ToSql()
	if err != nil {
		return err
	}

	// Execute the SQL query
	commandTag, err := f.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return err
	}

	return nil
}

func (f *favouriteRepo) ListFavouritesByUserId(ctx context.Context, user_id string) ([]*entity.Favourite, error) {
	
	ctx, span := otlp.Start(ctx, favouriteServiceName, favouriteSpanRepoPrefix+"List")
	defer span.End()
	
	var favourites []*entity.Favourite

	queryBuilder := f.FavouriteSelectQueryPrefix()

	queryBuilder = queryBuilder.Where(f.db.Sq.Equal("user_id", user_id)).Where(f.db.Sq.Equal("deleted_at", nil))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := f.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var favourite entity.Favourite

		if err := rows.Scan(
			&favourite.FavouriteId,
			&favourite.EstablishmentId,
			&favourite.UserId,
			&favourite.CreatedAt,
			&favourite.UpdatedAt,
		); err != nil {
			return nil, err
		}
		favourites = append(favourites, &favourite)
	}

	return favourites, nil
}
