package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberWalletRepository struct {
	db *sqlx.DB
	balance,
	transactions,
	journals string
}

func NewMemberWalletRepository(db *sql.DB) athena.MemberWalletRepository {
	return &memberWalletRepository{
		db:           sqlx.NewDb(db, "mysql"),
		balance:      "member_wallet_balance",
		transactions: "member_wallet_transactions",
		journals:     "member_wallet_journals",
	}
}

func (r *memberWalletRepository) MemberWalletBalance(ctx context.Context, memberID uint) (*athena.MemberWalletBalance, error) {

	query, args, err := sq.Select(
		"member_id", "balance", "created_at", "updated_at",
	).From(r.balance).Where(sq.Eq{"member_id": memberID}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var balance = new(athena.MemberWalletBalance)
	err = r.db.GetContext(ctx, balance, query, args...)

	return balance, err

}

func (r *memberWalletRepository) CreateMemberWalletBalance(ctx context.Context, memberID uint, balance float64) (*athena.MemberWalletBalance, error) {

	i := sq.Insert(r.balance).
		Columns(
			"member_id", "balance", "created_at", "updated_at",
		).Values(
		memberID, balance,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	)
	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberWalletBalance(ctx, memberID)

}

func (r *memberWalletRepository) UpdateMemberWalletBalance(ctx context.Context, memberID uint, balance float64) (*athena.MemberWalletBalance, error) {

	query, args, err := sq.Update(r.balance).
		Set("balance", balance).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberWalletBalance(ctx, memberID)

}

func (r *memberWalletRepository) MemberWalletTransactions(ctx context.Context, operators ...*athena.Operator) ([]*athena.MemberWalletTransaction, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "transaction_id", "journal_ref_id",
		"client_id", "client_type", "location_id",
		"location_type", "type_id", "quantity",
		"unit_price", "is_buy", "is_personal",
		"date", "created_at", "updated_at",
	).From(r.transactions), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var transactions = make([]*athena.MemberWalletTransaction, 0)
	err = r.db.SelectContext(ctx, transactions, query, args...)

	return transactions, err

}

func (r *memberWalletRepository) CreateMemberWalletTransactions(ctx context.Context, memberID uint, transactions []*athena.MemberWalletTransaction) ([]*athena.MemberWalletTransaction, error) {

	i := sq.Insert(r.transactions).Columns(
		"member_id", "transaction_id", "journal_ref_id",
		"client_id", "client_type", "location_id",
		"location_type", "type_id", "quantity",
		"unit_price", "is_buy", "is_personal",
		"date", "created_at", "updated_at",
	)
	transactionIDs := make([]uint64, 0)
	for _, transaction := range transactions {
		i = i.Values(
			memberID,
			transaction.TransactionID, transaction.JournalReferenceID, transaction.ClientID,
			transaction.ClientType, transaction.LocationID, transaction.LocationType,
			transaction.TypeID, transaction.Quantity, transaction.UnitPrice,
			transaction.IsBuy, transaction.IsPersonal, transaction.Date,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)

		transactionIDs = append(transactionIDs, transaction.TransactionID)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberWalletTransactions(ctx, athena.NewInOperator("transaction_id", transactionIDs))

}

func (r *memberWalletRepository) MemberWalletJournals(ctx context.Context, operators ...*athena.Operator) ([]*athena.MemberWalletJournal, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "journal_id", "ref_type",
		"context_id", "context_id_type", "description",
		"reason", "first_party_id", "first_party_type",
		"second_party_id", "second_party_type", "amount",
		"balance", "tax", "tax_receiver_id",
		"date", "created_at", "updated_at",
	).From(r.journals), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	var journals = make([]*athena.MemberWalletJournal, 0)
	err = r.db.SelectContext(ctx, journals, query, args...)

	return journals, err

}

func (r *memberWalletRepository) CreateMemberWalletJournals(ctx context.Context, memberID uint, journals []*athena.MemberWalletJournal) ([]*athena.MemberWalletJournal, error) {

	i := sq.Insert(r.transactions).Columns(
		"member_id", "journal_id", "ref_type",
		"context_id", "context_id_type", "description",
		"reason", "first_party_id", "first_party_type",
		"second_party_id", "second_party_type", "amount",
		"balance", "tax", "tax_receiver_id",
		"date", "created_at", "updated_at",
	)
	journalIDs := make([]uint64, 0)
	for _, journal := range journals {
		i = i.Values(
			memberID,
			journal.JournalID, journal.RefType, journal.ContextID,
			journal.ContextType, journal.Description, journal.Reason,
			journal.FirstPartyID, journal.FirstPartyType, journal.SecondPartyID,
			journal.SecondPartyType, journal.Amount, journal.Balance,
			journal.Tax, journal.TaxReceiverID, journal.Date,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)

		journalIDs = append(journalIDs, journal.JournalID)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Clones Repository] Failed to insert records: %w", err)
	}

	return r.MemberWalletJournals(ctx, athena.NewInOperator("journal_id", journalIDs))

}
