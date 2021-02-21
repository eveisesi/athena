package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberContractRepository struct {
	db *sqlx.DB
	contracts,
	items,
	bids string
}

func NewMemberContractRepository(db *sql.DB) athena.MemberContractRepository {
	return &memberContractRepository{
		db:        sqlx.NewDb(db, "mysql"),
		contracts: "member_contracts",
		items:     "member_contract_items",
		bids:      "member_contract_bids",
	}
}

func (r *memberContractRepository) MemberContract(ctx context.Context, memberID, contractID uint) (*athena.MemberContract, error) {

	contracts, err := r.MemberContracts(ctx, memberID, athena.NewEqualOperator("contract_id", contractID))
	if err != nil {
		return nil, err
	}

	if len(contracts) == 0 {
		return nil, nil
	}

	return contracts[0], nil

}

func (r *memberContractRepository) MemberContracts(ctx context.Context, memberID uint, operators ...*athena.Operator) ([]*athena.MemberContract, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "contract_id", "acceptor_id", "assignee_id",
		"availability", "buyout", "collateral", "date_accepted",
		"date_completed", "date_expired", "date_issued", "days_to_complete",
		"end_location_id", "for_corporation", "issuer_corporation_id", "issuer_id",
		"price", "reward", "start_location_id", "status", "title",
		"type", "volume", "created_at", "updated_at",
	).From(r.contracts), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	var contracts = make([]*athena.MemberContract, 0)
	err = r.db.SelectContext(ctx, &contracts, query, args...)

	return contracts, err

}

func (r *memberContractRepository) CreateContracts(ctx context.Context, memberID uint, contracts []*athena.MemberContract) ([]*athena.MemberContract, error) {

	i := sq.Insert(r.contracts).Columns(
		"member_id", "contract_id", "acceptor_id", "acceptor_type", "assignee_id", "assignee_type",
		"availability", "buyout", "collateral", "date_accepted",
		"date_completed", "date_expired", "date_issued", "days_to_complete",
		"end_location_id", "end_location_type", "for_corporation", "issuer_corporation_id", "issuer_id",
		"price", "reward", "start_location_id", "start_location_type", "status", "title",
		"type", "volume", "created_at", "updated_at",
	)

	contractIDs := make([]uint, len(contracts))
	for j, contract := range contracts {
		contractIDs[j] = contract.ContractID
		i = i.Values(
			memberID, contract.ContractID, contract.AcceptorID, contract.AcceptorType, contract.AssigneeID, contract.AssigneeType,
			contract.Availability, contract.Buyout, contract.Collateral, contract.DateAccepted,
			contract.DateCompleted, contract.DateExpired, contract.DateIssued, contract.DaysToComplete,
			contract.EndLocationID, contract.EndLocationType, contract.ForCorporation, contract.IssuerCorporationID, contract.IssuerID,
			contract.Price, contract.Reward, contract.StartLocationID, contract.StartLocationType,
			contract.Status, contract.Title, contract.Type,
			contract.Volume, sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert records: %w", err)
	}

	return r.MemberContracts(ctx, memberID, athena.NewInOperator("contract_id", contractIDs))

}

func (r *memberContractRepository) UpdateContract(ctx context.Context, memberID uint, contract *athena.MemberContract) (*athena.MemberContract, error) {

	query, args, err := sq.Update(r.contracts).
		Set("acceptor_id", contract.AcceptorID).
		Set("assignee_id", contract.AssigneeID).
		Set("availability", contract.Availability).
		Set("buyout", contract.Buyout).
		Set("collateral", contract.Collateral).
		Set("date_accepted", contract.DateAccepted).
		Set("date_completed", contract.DateCompleted).
		Set("date_expired", contract.DateExpired).
		Set("date_issued", contract.DateIssued).
		Set("days_to_complete", contract.DaysToComplete).
		Set("end_location_id", contract.EndLocationID).
		Set("for_corporation", contract.ForCorporation).
		Set("issuer_corporation_id", contract.IssuerCorporationID).
		Set("issuer_id", contract.IssuerID).
		Set("price", contract.Price).
		Set("reward", contract.Reward).
		Set("start_location_id", contract.StartLocationID).
		Set("status", contract.Status).
		Set("title", contract.Title).
		Set("type", contract.Type).
		Set("volume", contract.Volume).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID, "contract_id": contract.ContractID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert records: %w", err)
	}

	return r.MemberContract(ctx, memberID, contract.ContractID)

}

func (r *memberContractRepository) MemberContractItems(ctx context.Context, memberID, contractID uint, operators ...*athena.Operator) ([]*athena.MemberContractItem, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "contract_id", "record_id",
		"type_id", "quantity", "raw_quantity",
		"is_included", "is_singleton",
		"created_at", "updated_at",
	).From(r.items), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	var items = make([]*athena.MemberContractItem, 0)
	err = r.db.SelectContext(ctx, &items, query, args...)

	return items, err

}

func (r *memberContractRepository) CreateMemberContractItems(ctx context.Context, memberID, contractID uint, items []*athena.MemberContractItem) ([]*athena.MemberContractItem, error) {

	i := sq.Insert(r.items).Columns(
		"member_id", "contract_id", "record_id",
		"type_id", "quantity", "raw_quantity",
		"is_included", "is_singleton",
		"created_at", "updated_at",
	)

	for _, item := range items {
		i = i.Values(
			memberID, contractID,
			item.RecordID, item.TypeID, item.Quantity,
			item.RawQuantity, item.IsIncluded, item.IsSingleton,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert records: %w", err)
	}

	return r.MemberContractItems(ctx, memberID, contractID)

}

func (r *memberContractRepository) MemberContractBids(ctx context.Context, memberID, contractID uint, operators ...*athena.Operator) ([]*athena.MemberContractBid, error) {

	query, args, err := sq.Select(
		"member_id", "contract_id", "bid_id", "bidder", "amount", "bid_date", "created_at", "updated_at",
	).From(r.bids).Where(sq.Eq{"member_id": memberID, "contract_id": contractID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	var bids = make([]*athena.MemberContractBid, 0)
	err = r.db.SelectContext(ctx, &bids, query, args...)

	return bids, err

}

func (r *memberContractRepository) CreateMemberContractBids(ctx context.Context, memberID, contractID uint, bids []*athena.MemberContractBid) ([]*athena.MemberContractBid, error) {

	i := sq.Insert(r.bids).Columns(
		"member_id", "contract_id", "bid_id", "bidder", "amount", "bid_date", "created_at", "updated_at",
	)
	for _, bid := range bids {
		i = i.Values(
			memberID, contractID,
			bid.BidID, bid.BidderID,
			bid.Amount, bid.BidDate,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert records: %w", err)
	}

	return r.MemberContractBids(ctx, memberID, contractID)

}
