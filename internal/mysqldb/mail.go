package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/jmoiron/sqlx"
)

type mailRepository struct {
	db *sqlx.DB
	mail,
	recipients,
	member,
	labels,
	memberLists,
	lists string
}

func NewMailRepository(db *sql.DB) athena.MailRepository {
	return &mailRepository{
		db:          sqlx.NewDb(db, "mysql"),
		mail:        "mail_headers",
		recipients:  "mail_header_recipients",
		member:      "member_mail_headers",
		labels:      "member_mail_labels",
		memberLists: "member_mailing_lists",
		lists:       "mailing_lists",
	}
}

func (r *mailRepository) MailHeader(ctx context.Context, mailID uint) (*athena.MailHeader, error) {

	mails, err := r.MailHeaders(ctx, athena.NewEqualOperator("mail_id", mailID))
	if err != nil {
		return nil, err
	}

	if len(mails) == 1 {
		return mails[0], nil
	}

	return nil, nil

}

func (r *mailRepository) MailHeaders(ctx context.Context, operators ...*athena.Operator) ([]*athena.MailHeader, error) {

	query, args, err := BuildFilters(sq.Select(
		"id", "sender_id", "sender_type", "subject", "body", "sent", "created_at",
	).From(r.mail), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var headers = make([]*athena.MailHeader, 0)
	err = r.db.SelectContext(ctx, &headers, query, args...)

	return headers, err

}

func (r *mailRepository) CreateMailHeaders(ctx context.Context, headers []*athena.MailHeader) ([]*athena.MailHeader, error) {

	i := sq.Insert(r.mail).Options("IGNORE").Columns(
		"id", "sender_id", "sender_type", "subject", "body", "sent", "created_at",
	)
	mailIDs := make([]interface{}, 0)
	for _, header := range headers {
		mailIDs = append(mailIDs, header.MailID)
		i = i.Values(
			header.MailID, header.Sender, header.SenderType,
			header.Subject, header.Body, header.Timestamp.Format(time.RFC3339),
			sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
	}

	return r.MailHeaders(ctx, athena.NewInOperator("id", mailIDs...))

}

func (r *mailRepository) MemberMailHeaders(ctx context.Context, operators ...*athena.Operator) ([]*athena.MemberMailHeader, error) {

	query, args, err := BuildFilters(sq.Select(
		"member_id", "mail_id", "labels", "is_read", "created_at", "updated_at",
	).From(r.member), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var headers = make([]*athena.MemberMailHeader, 0)
	err = r.db.SelectContext(ctx, &headers, query, args...)

	return headers, err

}

func (r *mailRepository) CreateMemberMailHeaders(ctx context.Context, memberID uint, headers []*athena.MemberMailHeader) ([]*athena.MemberMailHeader, error) {

	i := sq.Insert(r.member).Columns(
		"member_id", "mail_id", "labels", "is_read", "created_at", "updated_at",
	)

	mailIDs := make([]interface{}, 0)
	for _, header := range headers {
		mailIDs = append(mailIDs, header.MailID)
		i = i.Values(
			memberID, header.MailID,
			header.Labels, header.IsRead,
			sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
	}

	return r.MemberMailHeaders(ctx, athena.NewEqualOperator("member_id", memberID), athena.NewInOperator("mail_id", mailIDs...))

}

func (r *mailRepository) UpdateMemberMailHeaders(ctx context.Context, memberID uint, headers []*athena.MemberMailHeader) ([]*athena.MemberMailHeader, error) {

	mailIDs := make([]interface{}, 0)
	for _, header := range headers {
		query, args, err := sq.Update(r.member).
			Set("member_id", header.MemberID).
			Set("mail_id", header.MailID).
			Set("labels", header.Labels).
			Set("is_read", header.IsRead).
			Set("created_at", header.CreatedAt).
			Set("updated_at", header.UpdatedAt).
			Where(sq.Eq{"member_id": memberID, "mail_id": header.MailID}).ToSql()
		if err != nil {
			return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
		}

		_, err = r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
		}
		mailIDs = append(mailIDs, header.MailID)

	}

	return r.MemberMailHeaders(ctx, athena.NewEqualOperator("member_id", memberID), athena.NewInOperator("mail_id", mailIDs...))

}

func (r *mailRepository) MailRecipients(ctx context.Context, operators ...*athena.Operator) ([]*athena.MailRecipient, error) {

	query, args, err := BuildFilters(sq.Select(
		"mail_id", "recipient_id", "recipient_type", "created_at",
	).From(r.recipients), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var recipients = make([]*athena.MailRecipient, 0)
	err = r.db.SelectContext(ctx, &recipients, query, args...)

	return recipients, err

}

func (r *mailRepository) CreateMailRecipients(ctx context.Context, recipients []*athena.MailRecipient) ([]*athena.MailRecipient, error) {

	i := sq.Insert(r.recipients).Options("IGNORE").Columns(
		"mail_id", "recipient_id", "recipient_type", "created_at",
	)

	mailIDs := make([]interface{}, 0)
	for _, recipient := range recipients {
		mailIDs = append(mailIDs, recipient.MailID)
		i = i.Values(
			recipient.MailID,
			recipient.RecipientID,
			recipient.RecipientType,
			sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
	}

	return r.MailRecipients(ctx, athena.NewInOperator("mail_id", mailIDs...))

}

func (r *mailRepository) MemberMailLabels(ctx context.Context, memberID uint) (*athena.MemberMailLabels, error) {

	query, args, err := sq.Select(
		"member_id", "labels", "total_unread_count", "created_at", "updated_at",
	).From(r.labels).Where("member_id", memberID).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var labels = new(athena.MemberMailLabels)
	err = r.db.SelectContext(ctx, labels, query, args...)

	return labels, err
}

func (r *mailRepository) CreateMemberMailLabel(ctx context.Context, memberID uint, labels *athena.MemberMailLabels) (*athena.MemberMailLabels, error) {

	query, args, err := sq.Insert(r.labels).Columns(
		"member_id", "labels", "total_unread_count", "created_at", "updated_at",
	).Values(
		memberID, labels.Labels, labels.TotalUnreadCount,
		sq.Expr(`NOW()`), sq.Expr(`NOW()`),
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
	}

	return r.MemberMailLabels(ctx, memberID)

}

func (r *mailRepository) UpdateMemberMailLabel(ctx context.Context, memberID uint, labels *athena.MemberMailLabels) (*athena.MemberMailLabels, error) {

	query, args, err := sq.Update(r.labels).
		Set("labels", labels.Labels).
		Set("total_unread_count", labels.TotalUnreadCount).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to insert records: %w", err)
	}

	return r.MemberMailLabels(ctx, memberID)

}

func (r *mailRepository) MemberMailingLists(ctx context.Context, memberID uint) ([]*athena.MemberMailingList, error) {

	query, args, err := sq.Select(
		"member_id", "mailing_list_id", "created_at",
	).From(r.memberLists).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var lists = make([]*athena.MemberMailingList, 0)
	err = r.db.SelectContext(ctx, &lists, query, args...)

	return lists, err

}

func (r *mailRepository) CreateMemberMailingLists(ctx context.Context, memberID uint, lists []*athena.MemberMailingList) ([]*athena.MemberMailingList, error) {

	i := sq.Insert(r.memberLists).Options("IGNORE").Columns(
		"member_id", "mailing_list_id", "created_at",
	)
	for _, list := range lists {
		i = i.Values(
			memberID, list.MailingListID, sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Respository] Failed to insert records: %w", err)
	}

	return r.MemberMailingLists(ctx, memberID)

}

func (r *mailRepository) DeleteMemberMailingListsAll(ctx context.Context, memberID uint) (bool, error) {

	query, args, err := sq.Delete(r.memberLists).Where(sq.Eq{"member_id": memberID}).ToSql()
	if err != nil {
		return false, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("[Mail Respository] Failed to delete records: %w", err)
	}

	return true, nil

}

func (r *mailRepository) MailingList(ctx context.Context, mailingListID uint) (*athena.MailingList, error) {

	query, args, err := sq.Select(
		"mailing_list_id", "name", "created_at", "updated_at",
	).From(r.lists).Where(sq.Eq{"mailing_list_id": mailingListID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var list = new(athena.MailingList)
	err = r.db.GetContext(ctx, list, query, args...)

	return list, err

}

func (r *mailRepository) MailingLists(ctx context.Context, operators ...*athena.Operator) ([]*athena.MailingList, error) {

	query, args, err := BuildFilters(sq.Select(
		"mailing_list_id", "name", "created_at", "updated_at",
	).From(r.lists), operators...).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	var lists = make([]*athena.MailingList, 0)
	err = r.db.SelectContext(ctx, &lists, query, args...)

	return lists, err

}

func (r *mailRepository) CreateMailingLists(ctx context.Context, lists []*athena.MailingList) ([]*athena.MailingList, error) {

	i := sq.Insert(r.lists).Options("IGNORE").Columns(
		"mailing_list_id", "name", "created_at", "updated_at",
	)
	ids := make([]interface{}, 0, len(lists))
	for _, list := range lists {
		ids = append(ids, list.MailingListID)
		i = i.Values(
			list.MailingListID, list.Name, sq.Expr(`NOW()`), sq.Expr(`NOW()`),
		)
	}

	query, args, err := i.ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Respository] Failed to insert records: %w", err)
	}

	return r.MailingLists(ctx, athena.NewInOperator("mailing_list_id", ids...))

}

func (r *mailRepository) UpdateMailingList(ctx context.Context, mailingListID uint, list *athena.MailingList) (*athena.MailingList, error) {

	query, args, err := sq.Update(r.lists).
		Set("name", list.Name).
		Set("updated_at", sq.Expr(`NOW()`)).
		Where(sq.Eq{"mailing_list_id": list.MailingListID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[Mail Repository] Failed to generate query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("[Mail Respository] Failed to insert records: %w", err)
	}

	return r.MailingList(ctx, list.MailingListID)

}
