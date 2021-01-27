package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberContractRepository interface {
	MemberContract(ctx context.Context, memberID string, contractID int) (*MemberContract, error)
	Contracts(ctx context.Context, memberID string, operators ...*Operator) ([]*MemberContract, error)
	CreateContracts(ctx context.Context, memberID string, contracts []*MemberContract) ([]*MemberContract, error)
	UpdateContract(ctx context.Context, memberID string, contractID int, contract *MemberContract) (*MemberContract, error)
	DeleteContract(ctx context.Context, memberID string) (bool, error)
}

type MemberContract struct {
	MemberID            primitive.ObjectID   `bson:"member_id" json:"member_id"`
	ContractID          int                  `bson:"contract_id" json:"contract_id"`
	AcceptorID          null.Int             `bson:"acceptor_id" json:"acceptor_id"`
	AssigneeID          null.Int             `bson:"assignee_id" json:"assignee_id"`
	Availability        ContractAvailability `bson:"availability" json:"availability"`
	Buyout              null.Float64         `bson:"buyout,omitempty" json:"buyout,omitempty"`
	Collateral          null.Float64         `bson:"collateral,omitempty" json:"collateral,omitempty"`
	DateAccepted        null.Time            `bson:"date_accepted,omitempty" json:"date_accepted,omitempty"`
	DateCompleted       null.Time            `bson:"date_completed,omitempty" json:"date_completed,omitempty"`
	DateExpired         time.Time            `bson:"date_expired" json:"date_expired"`
	DateIssued          time.Time            `bson:"date_issued" json:"date_issued"`
	DaysToComplete      null.Int             `bson:"days_to_complete,omitempty" json:"days_to_complete,omitempty"`
	EndLocationID       null.Int64           `bson:"end_location_id,omitempty" json:"end_location_id,omitempty"`
	ForCorporation      bool                 `bson:"for_corporation" json:"for_corporation"`
	IssuerCorporationID int                  `bson:"issuer_corporation_id" json:"issuer_corporation_id"`
	IssuerID            int64                `bson:"issuer_id" json:"issuer_id"`
	Price               null.Int             `bson:"price,omitempty" json:"price,omitempty"`
	Reward              null.Int             `bson:"reward,omitempty" json:"reward,omitempty"`
	StartLocationID     null.Int64           `bson:"start_location_id,omitempty" json:"start_location_id,omitempty"`
	Status              ContractStatus       `bson:"status" json:"status"`
	Title               null.String          `bson:"title,omitempty" json:"title,omitempty"`
	Type                ContractType         `bson:"type" json:"type"`
	Volume              null.Float64         `bson:"volume,omitempty" json:"volume,omitempty"`
	CreatedAt           time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time            `bson:"updated_at" json:"updated_at"`
}

type ContractAvailability string

const (
	ContractAvailabilityPublic      ContractAvailability = "public"
	ContractAvailabilityPersonal    ContractAvailability = "personal"
	ContractAvailabilityCorporation ContractAvailability = "corporation"
	ContractAvailabilityAlliance    ContractAvailability = "alliance"
)

var AllContractAvailabilities = []ContractAvailability{
	ContractAvailabilityPublic,
	ContractAvailabilityPersonal,
	ContractAvailabilityCorporation,
	ContractAvailabilityAlliance,
}

func (c ContractAvailability) Valid() bool {
	for _, v := range AllContractAvailabilities {
		if v == c {
			return true
		}
	}

	return true
}

type ContractStatus string

const (
	ContractStatusOutstanding        ContractStatus = "outstanding"
	ContractStatusInProgress         ContractStatus = "in_progress"
	ContractStatusFinishedIssuer     ContractStatus = "finished_issuer"
	ContractStatusFinishedContractor ContractStatus = "finished_contractor"
	ContractStatusFinished           ContractStatus = "finished"
	ContractStatusConcelled          ContractStatus = "cancelled"
	ContractStatusRejected           ContractStatus = "rejected"
	ContractStatusFailed             ContractStatus = "failed"
	ContractStatusDeleted            ContractStatus = "deleted"
	ContractStatusRevered            ContractStatus = "reversed"
)

var AllContractStatuses = []ContractStatus{
	ContractStatusOutstanding,
	ContractStatusInProgress,
	ContractStatusFinishedIssuer,
	ContractStatusFinishedContractor,
	ContractStatusFinished,
	ContractStatusConcelled,
	ContractStatusRejected,
	ContractStatusFailed,
	ContractStatusDeleted,
	ContractStatusRevered,
}

func (c ContractStatus) Valid() bool {
	for _, v := range AllContractStatuses {
		if v == c {
			return true
		}
	}

	return true
}

type ContractType string

const (
	ContractTypeUnknown      ContractType = "unknown"
	ContractTypeItemExchange ContractType = "item_exchange"
	ContractTypeAuction      ContractType = "auction"
	ContractTypeCourier      ContractType = "courier"
	ContractTypeLoan         ContractType = "loan"
)

var AllContractTypes = []ContractType{
	ContractTypeUnknown,
	ContractTypeItemExchange,
	ContractTypeAuction,
	ContractTypeCourier,
	ContractTypeLoan,
}

func (c ContractType) Valid() bool {
	for _, v := range AllContractTypes {
		if v == c {
			return true
		}
	}

	return true
}
