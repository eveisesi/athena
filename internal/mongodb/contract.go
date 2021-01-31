package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type memberContractRepository struct {
	contracts *mongo.Collection
	items     *mongo.Collection
	bids      *mongo.Collection
}

func NewMemberContractRepository(d *mongo.Database) (athena.MemberContractRepository, error) {

	contracts := d.Collection("member_contracts")
	contractIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "contract_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_contracts_member_id_contract_id_unique"),
				Unique: newBool(true),
			},
		},
	}

	_, err := contracts.Indexes().CreateMany(context.TODO(), contractIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository]: Failed to create index on member_contracts collection: %w", err)
	}

	items := d.Collection("member_contract_items")
	itemIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "contract_id",
					Value: 1,
				},
				primitive.E{
					Key:   "record_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_contract_items_member_id_contract_id_record_id_unique"),
				Unique: newBool(true),
			},
		},
	}

	_, err = items.Indexes().CreateMany(context.TODO(), itemIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository]: Failed to create index on member_contract_items collection: %w", err)
	}

	bids := d.Collection("member_contract_bids")
	bidIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "contract_id",
					Value: 1,
				},
				primitive.E{
					Key:   "bid_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_contract_items_member_id_contract_id_bid_id_unique"),
				Unique: newBool(true),
			},
		},
	}

	_, err = bids.Indexes().CreateMany(context.TODO(), bidIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository]: Failed to create index on member_contract_bids collection: %w", err)
	}

	return &memberContractRepository{
		contracts: contracts,
		items:     items,
		bids:      bids,
	}, err

}

func (r *memberContractRepository) MemberContract(ctx context.Context, memberID string, contractID int) (*athena.MemberContract, error) {

	contract := new(athena.MemberContract)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filters := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "contract_id", Value: contractID}}

	err = r.contracts.FindOne(ctx, filters).Decode(contract)

	return contract, err

}

func (r *memberContractRepository) Contracts(ctx context.Context, memberID string, operators ...*athena.Operator) ([]*athena.MemberContract, error) {

	contracts := make([]*athena.MemberContract, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	operators = append(operators, athena.NewEqualOperator("member_id", pid))

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	results, err := r.contracts.Find(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	return contracts, results.All(ctx, &contracts)

}

func (r *memberContractRepository) CreateContracts(ctx context.Context, memberID string, contracts []*athena.MemberContract) ([]*athena.MemberContract, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	documents := make([]interface{}, len(contracts))
	now := time.Now()
	for i, contract := range contracts {
		contract.MemberID = pid
		contract.CreatedAt = now
		contract.UpdatedAt = now

		documents[i] = contract

	}

	_, err = r.contracts.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert record into the member_contracts collection: %w", err)
	}

	return contracts, nil

}

func (r *memberContractRepository) UpdateContract(ctx context.Context, memberID string, contractID int, contract *athena.MemberContract) (*athena.MemberContract, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "contract_id", Value: contract.ContractID}}
	update := primitive.D{primitive.E{Key: "$set", Value: contract}}

	_, err = r.contracts.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to update record in the member_contracts collection: %w", err)
	}

	return contract, nil

}

func (r *memberContractRepository) DeleteContracts(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
	}

	results, err := r.contracts.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to delete records from the member_contracts collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberContractRepository) MemberContractItems(ctx context.Context, memberID string, contractID int, operators ...*athena.Operator) ([]*athena.MemberContractItem, error) {

	items := make([]*athena.MemberContractItem, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	operators = append(operators, athena.NewEqualOperator("member_id", pid), athena.NewEqualOperator("contract_id", contractID))

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	results, err := r.items.Find(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	return items, results.All(ctx, &items)

}

func (r *memberContractRepository) CreateMemberContractItems(ctx context.Context, memberID string, contractID int, items []*athena.MemberContractItem) ([]*athena.MemberContractItem, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	documents := make([]interface{}, len(items))
	now := time.Now()
	for i, item := range items {
		item.MemberID = pid
		item.ContractID = contractID
		item.CreatedAt = now
		item.UpdatedAt = now

		documents[i] = item

	}

	_, err = r.items.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert record into the member_contract_items collection: %w", err)
	}

	return items, nil

}

func (r *memberContractRepository) DeleteMemberContractItems(ctx context.Context, memberID string, contractID int) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
		primitive.E{
			Key:   "contract_id",
			Value: contractID,
		},
	}

	results, err := r.items.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to delete records from the member_contract_items collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberContractRepository) DeleteMemberContractItemsAll(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
	}

	results, err := r.items.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to delete records from the member_contract_items collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberContractRepository) MemberContractBids(ctx context.Context, memberID string, contractID int, operators ...*athena.Operator) ([]*athena.MemberContractBid, error) {

	bids := make([]*athena.MemberContractBid, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	operators = append(operators, athena.NewEqualOperator("member_id", pid), athena.NewEqualOperator("contract_id", contractID))

	filters := BuildFilters(operators...)
	options := BuildFindOptions(operators...)

	results, err := r.bids.Find(ctx, filters, options)
	if err != nil {
		return nil, err
	}

	return bids, results.All(ctx, &bids)

}

func (r *memberContractRepository) CreateMemberContractBids(ctx context.Context, memberID string, contractID int, bids []*athena.MemberContractBid) ([]*athena.MemberContractBid, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	documents := make([]interface{}, len(bids))
	now := time.Now()
	for i, item := range bids {
		item.MemberID = pid
		item.ContractID = contractID
		item.CreatedAt = now
		item.UpdatedAt = now

		documents[i] = item

	}

	_, err = r.bids.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Contract Repository] Failed to insert record into the member_contract_bids collection: %w", err)
	}

	return bids, nil

}

func (r *memberContractRepository) DeleteMemberContractBids(ctx context.Context, memberID string, contractID int) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
		primitive.E{
			Key:   "contract_id",
			Value: contractID,
		},
	}

	results, err := r.bids.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to delete records from the member_contract_bids collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberContractRepository) DeleteMemberContractBidsAll(ctx context.Context, memberID string) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
	}

	results, err := r.bids.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contract Repository] Failed to delete records from the member_contract_bids collection: %w", err)
	}

	return results.DeletedCount > 0, err

}
