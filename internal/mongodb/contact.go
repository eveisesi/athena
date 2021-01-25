package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/athena"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type memberContactRepository struct {
	contacts *mongo.Collection
	labels   *mongo.Collection
}

func NewMemberContactRepository(d *mongo.Database) (athena.MemberContactRepository, error) {

	var ctx = context.Background()

	contacts := d.Collection("member_contacts")
	contactsIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"contact_id": 1,
			"member_id":  1,
		},
		Options: &options.IndexOptions{
			Name:   newString("member_contacts_member_id_contact_id_unique"),
			Unique: newBool(true),
		},
	}

	_, err := contacts.Indexes().CreateOne(ctx, contactsIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to create index on contact repository: %w", err)
	}

	labels := d.Collection("member_contact_labels")
	labelsIdxMod := mongo.IndexModel{
		Keys: bson.M{
			"member_id": 1,
		},
		Options: &options.IndexOptions{
			Name: newString("member_contact_labels_member_id_idx"),
		},
	}
	_, err = labels.Indexes().CreateOne(ctx, labelsIdxMod)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to create index on contact labels repository: %w", err)
	}

	return &memberContactRepository{
		contacts: contacts,
		labels:   labels,
	}, nil

}

func (r *memberContactRepository) MemberContact(ctx context.Context, memberID string, contactID int) (*athena.MemberContact, error) {

	contact := new(athena.MemberContact)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	err = r.contacts.FindOne(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "contact_id", Value: contactID}}).Decode(contact)

	return contact, err

}

func (r *memberContactRepository) MemberContacts(ctx context.Context, memberID string) ([]*athena.MemberContact, error) {

	contacts := make([]*athena.MemberContact, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	results, err := r.contacts.Find(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}})
	if err != nil {
		return nil, err
	}

	return contacts, results.All(ctx, &contacts)

}

func (r *memberContactRepository) CreateMemberContacts(ctx context.Context, memberID string, contacts []*athena.MemberContact) ([]*athena.MemberContact, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	documents := make([]interface{}, len(contacts))
	now := time.Now()
	for i, contact := range contacts {
		contact.MemberID = pid
		contact.CreatedAt = now
		contact.UpdatedAt = now

		documents[i] = contact

	}

	_, err = r.contacts.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to insert record into the member contacts collection: %w", err)
	}

	return contacts, nil

}

func (r *memberContactRepository) UpdateMemberContact(ctx context.Context, memberID string, contact *athena.MemberContact) (*athena.MemberContact, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "contact_id", Value: contact.ContactID}}
	update := primitive.D{primitive.E{Key: "$set", Value: contact}}

	_, err = r.contacts.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to update record in the member contacts collection: %w", err)
	}

	return contact, nil

}

func (r *memberContactRepository) DeleteMemberContact(ctx context.Context, memberID string, contactID int) (bool, error) {
	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "contact_id", Value: pid}}

	results, err := r.contacts.DeleteOne(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contact Repository] Failed to delete record from the member contacts collection: %w", err)
	}

	return results.DeletedCount > 0, err
}

func (r *memberContactRepository) DeleteMemberContacts(ctx context.Context, memberID string, contacts []*athena.MemberContact) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	contactIDs := make([]int, len(contacts))
	for i, contact := range contacts {
		contactIDs[i] = contact.ContactID
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
	}

	if len(contactIDs) > 0 {
		filter = append(filter, primitive.E{
			Key: "contact_id",
			Value: primitive.D{
				primitive.E{
					Key:   "$in",
					Value: contactIDs,
				},
			},
		})
	}

	results, err := r.contacts.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contact Repository] Failed to delete records from the member contacts collection: %w", err)
	}

	return results.DeletedCount > 0, err

}

func (r *memberContactRepository) MemberContactLabels(ctx context.Context, memberID string) ([]*athena.MemberContactLabel, error) {

	labels := make([]*athena.MemberContactLabel, 0)

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to cast id to objectID: %w", err)
	}

	results, err := r.labels.Find(ctx, primitive.D{primitive.E{Key: "member_id", Value: pid}})
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to fetch records from the member_contact_labels collection: %w", err)
	}

	return labels, results.All(ctx, &labels)

}

func (r *memberContactRepository) CreateMemberContactLabels(ctx context.Context, memberID string, labels []*athena.MemberContactLabel) ([]*athena.MemberContactLabel, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	documents := make([]interface{}, len(labels))
	now := time.Now()
	for i, label := range labels {
		label.MemberID = pid
		label.CreatedAt = now
		label.UpdatedAt = now

		documents[i] = label

	}

	_, err = r.labels.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to insert record into the member_contact_labels collection: %w", err)
	}

	return labels, nil

}

func (r *memberContactRepository) UpdateMemberContactLabel(ctx context.Context, memberID string, label *athena.MemberContactLabel) (*athena.MemberContactLabel, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	filter := primitive.D{primitive.E{Key: "member_id", Value: pid}, primitive.E{Key: "label_id", Value: label.LabelID}}
	update := primitive.D{primitive.E{Key: "$set", Value: label}}

	_, err = r.labels.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("[Contact Repository] Failed to update record in the member_contact_labels collection: %w", err)
	}

	return label, nil

}

func (r *memberContactRepository) DeleteMemberContactLabels(ctx context.Context, memberID string, labels []*athena.MemberContactLabel) (bool, error) {

	pid, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return false, fmt.Errorf("[Clone Repository] Failed to cast id to objectID: %w", err)
	}

	labelIDs := make([]int64, len(labels))
	for i, label := range labels {
		labelIDs[i] = label.LabelID
	}

	filter := primitive.D{
		primitive.E{
			Key:   "member_id",
			Value: pid,
		},
	}

	if len(labelIDs) > 0 {
		filter = append(filter, primitive.E{
			Key: "label_id",
			Value: primitive.D{
				primitive.E{
					Key:   "$in",
					Value: labelIDs,
				},
			},
		})
	}

	results, err := r.contacts.DeleteMany(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("[Contact Repository] Failed to delete record from the member_contact_labels collection: %w", err)
	}

	return results.DeletedCount > 0, err

}
