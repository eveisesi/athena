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

	return &memberContractRepository{
		contracts: contracts,
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

func (r *memberContractRepository) DeleteContract(ctx context.Context, memberID string) (bool, error) {

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
