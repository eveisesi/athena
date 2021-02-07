package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberAssetsRepository struct {
	db    *sqlx.DB
	table string
}

func NewMemberAssetRepository(db *sql.DB) athena.MemberAssetsRepository {
	return &memberAssetsRepository{
		db:    sqlx.NewDb(db, "mysql"),
		table: "member_assets",
	}
}

func (r *memberAssetsRepository) MemberAsset(ctx context.Context, memberID, itemID uint) (*athena.MemberAsset, error) {

	query, args, err := sq.Select(
		"member_id", "item_id", "type_id",
		"location_id", "location_flag", "location_type",
		"quantity", "is_blueprint_copy", "is_singleton",
		"created_at", "updated_at",
	).From(r.table).Where(sq.Eq{"member_id": memberID, "item_id": itemID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var asset = new(athena.MemberAsset)
	err = r.db.GetContext(ctx, asset, query, args...)

	return asset, err

}

func (r *memberAssetsRepository) MemberAssets(ctx context.Context, id uint, operators ...*athena.Operator) ([]*athena.MemberAsset, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "item_id", "type_id",
		"location_id", "location_flag", "location_type",
		"quantity", "is_blueprint_copy", "is_singleton",
		"created_at", "updated_at",
	).From(r.table), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var assets = make([]*athena.MemberAsset, 0)
	err = r.db.SelectContext(ctx, &assets, query, args...)

	return assets, err

}

func (r *memberAssetsRepository) CreateMemberAssets(ctx context.Context, memberID uint, assets []*athena.MemberAsset) ([]*athena.MemberAsset, error) {

	i := sq.Insert(r.table).Columns(
		"member_id", "item_id", "type_id",
		"location_id", "location_flag", "location_type",
		"quantity", "is_blueprint_copy", "is_singleton",
		"created_at", "updated_at",
	)
	for _, asset := range assets {
		i = i.Values(
			memberID, asset.ItemID, asset.TypeID,
			asset.LocationID, asset.LocationFlag, asset.LocationType,
			asset.Quantity, asset.IsBlueprintCopy, asset.IsSingleton,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberAssets(ctx, memberID)

}

func (r *memberAssetsRepository) UpdateMemberAssets(ctx context.Context, memberID, itemID uint, asset *athena.MemberAsset) (*athena.MemberAsset, error) {

	query, args, err := sq.Update(r.table).
		Set("type_id", asset.TypeID).
		Set("location_id", asset.LocationID).
		Set("location_flag", asset.LocationFlag).
		Set("location_type", asset.LocationType).
		Set("quantity", asset.Quantity).
		Set("is_blueprint_copy", asset.IsBlueprintCopy).
		Set("is_singleton", asset.IsSingleton).
		Set("updated_at", asset.UpdatedAt).
		Where(sq.Eq{"member_id": memberID, "item_ID": itemID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberAsset(ctx, memberID, itemID)

}

func (r *memberAssetsRepository) DeleteMemberAssets(ctx context.Context, memberID uint, assets []*athena.MemberAsset) (bool, error) {

	for _, asset := range assets {
		query, args, err := sq.Delete(r.table).Where(sq.Eq{"member_id": memberID, "item_id": asset.ItemID}).ToSql()
		if err != nil {
			return false, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return false, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
		}

	}
	return true, nil

}

func (r *memberAssetsRepository) DeleteMemberAssetAll(ctx context.Context, memberID uint) (bool, error) {

	query, args, err := sq.Delete(r.table).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return true, nil

}
