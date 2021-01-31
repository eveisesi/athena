package mongodb

import (
	"context"
	"fmt"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type memberWalletRepository struct {
	balance      *mongo.Collection
	transactions *mongo.Collection
	journals     *mongo.Collection
}

func NewMemberWalletRepository(ctx context.Context, d *mongo.Database) (athena.MemberWalletRepository, error) {

	balance := d.Collection("member_wallet_balance")
	balanceIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_wallet_balance_member_id_unique"),
				Unique: newBool(true),
			},
		},
	}

	_, err := balance.Indexes().CreateMany(ctx, balanceIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Wallet Repository]: Failed to create index on member_wallet_balance collection: %w", err)
	}

	transactions := d.Collection("member_wallet_transactions")
	transactionIdxMods := []mongo.IndexModel{
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "member_id",
					Value: 1,
				},
				primitive.E{
					Key:   "transaction_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name:   newString("member_wallet_transaction_member_id_transaction_id_unique"),
				Unique: newBool(true),
			},
		},
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "transaction_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_wallet_transaction_transaction_id_idx"),
			},
		},
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "journal_id",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_wallet_transaction_journal_id_idx"),
			},
		},
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "client_id",
					Value: 1,
				},
				primitive.E{
					Key:   "client_type",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_wallet_transaction_client_id_client_type_idx"),
			},
		},
		{
			Keys: primitive.D{
				primitive.E{
					Key:   "location_id",
					Value: 1,
				},
				primitive.E{
					Key:   "location_type",
					Value: 1,
				},
			},
			Options: &options.IndexOptions{
				Name: newString("member_wallet_transaction_location_id_location_type_idx"),
			},
		},
	}

	_, err = transactions.Indexes().CreateMany(ctx, transactionIdxMods)
	if err != nil {
		return nil, fmt.Errorf("[Wallet Repository]: Failed to create index on member_wallet_transactions collection: %w", err)
	}

	// journals := d.Collection("member_wallet_journal")

	return nil, nil

}
