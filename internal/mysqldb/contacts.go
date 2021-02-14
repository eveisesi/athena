package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type memberContactRepository struct {
	db               *sqlx.DB
	contacts, labels string
}

func NewMemberContactRepository(db *sql.DB) athena.MemberContactRepository {
	return &memberContactRepository{
		db:       sqlx.NewDb(db, "mysql"),
		contacts: "member_contacts",
		labels:   "member_contact_labels",
	}
}

func (r *memberContactRepository) MemberContact(ctx context.Context, memberID, contactID uint) (*athena.MemberContact, error) {

	query, args, err := sq.Select(
		"member_id", "contact_id", "contact_type",
		"is_blocked", "is_watched",
		// "label_ids",
		"standing", "created_at", "updated_at",
	).From(r.contacts).Where(sq.Eq{"member_id": memberID, "contact_id": contactID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var contact = new(athena.MemberContact)
	err = r.db.GetContext(ctx, contact, query, args...)

	return contact, err

}

func (r *memberContactRepository) MemberContacts(ctx context.Context, memberID uint) ([]*athena.MemberContact, error) {

	query, args, err := sq.Select(
		"member_id", "contact_id", "contact_type",
		"is_blocked", "is_watched",
		// "label_ids",
		"standing", "created_at", "updated_at",
	).From(r.contacts).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var contacts = make([]*athena.MemberContact, 0)
	err = r.db.SelectContext(ctx, &contacts, query, args...)

	return contacts, err

}

func (r *memberContactRepository) CreateMemberContacts(ctx context.Context, memberID uint, contacts []*athena.MemberContact) ([]*athena.MemberContact, error) {

	i := sq.Insert(r.contacts).
		Columns(
			"member_id", "contact_id", "contact_type",
			"is_blocked", "is_watched",
			//  "label_ids",
			"standing", "created_at", "updated_at",
		)
	for _, contact := range contacts {
		i = i.Values(
			memberID,
			contact.ContactID, contact.ContactType,
			contact.IsBlocked, contact.IsWatched,
			// contact.LabelIDs,
			contact.Standing,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
	}

	return r.MemberContacts(ctx, memberID)

}

func (r *memberContactRepository) UpdateMemberContact(ctx context.Context, memberID uint, contact *athena.MemberContact) (*athena.MemberContact, error) {

	query, args, err := sq.Update(r.contacts).
		Set("contact_type", contact.ContactType).
		Set("is_blocked", contact.IsBlocked).
		Set("is_watched", contact.IsWatched).
		// Set("label_ids", contact.LabelIDs).
		Set("standing", contact.Standing).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID, "contact_id": contact.ContactID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
	}

	return r.MemberContact(ctx, memberID, contact.ContactID)
}

func (r *memberContactRepository) DeleteMemberContacts(ctx context.Context, memberID uint, contacts []*athena.MemberContact) (bool, error) {

	for _, contact := range contacts {
		query, args, err := sq.Delete(r.contacts).Where(sq.Eq{"member_id": memberID, "contact_id": contact.ContactID}).ToSql()
		if err != nil {
			return false, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return false, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
		}

	}
	return true, nil

}

func (r *memberContactRepository) MemberContactLabels(ctx context.Context, memberID uint) ([]*athena.MemberContactLabel, error) {

	query, args, err := sq.Select(
		"member_id", "label_id", "label_name",
		"created_at", "updated_at",
	).From(r.labels).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	var labels = make([]*athena.MemberContactLabel, 0)
	err = r.db.SelectContext(ctx, &labels, query, args...)

	return labels, err

}

func (r *memberContactRepository) CreateMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, error) {

	i := sq.Insert(r.labels).
		Columns(
			"member_id", "label_id", "label_name",
			"created_at", "updated_at",
		)
	for _, label := range labels {
		i = i.Values(
			memberID,
			label.LabelID, label.LabelName,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
	}

	return r.MemberContactLabels(ctx, memberID)

}

func (r *memberContactRepository) UpdateMemberContactLabel(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, error) {

	for _, label := range labels {
		query, args, err := sq.Update(r.labels).
			Set("label_name", label.LabelName).
			Set("updated_at", sq.Expr(`NOW()`)).
			Where(sq.Eq{"member_id": memberID, "labnel_id": label.LabelID}).ToSql()
		if err != nil {
			return nil, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
		}
	}

	return r.MemberContactLabels(ctx, memberID)

}

func (r *memberContactRepository) DeleteMemberContactLabels(ctx context.Context, memberID uint, labels []*athena.MemberContactLabel) (bool, error) {

	for _, label := range labels {
		query, args, err := sq.Delete(r.labels).Where(sq.Eq{"member_id": memberID, "label_id": label.LabelID}).ToSql()
		if err != nil {
			return false, fmt.Errorf("[Contact Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return false, fmt.Errorf("[Contact Repository] Failed to insert records: %w", err)
		}
	}

	return true, nil

}
