package athena

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type MemberContractRepository interface {
	memberContractRepository
	memberContractItemRepository
	memberContractBidRepository
}

type memberContractRepository interface {
	MemberContract(ctx context.Context, memberID, contractID uint) (*MemberContract, error)
	MemberContracts(ctx context.Context, memberID uint, operators ...*Operator) ([]*MemberContract, error)
	CreateContracts(ctx context.Context, memberID uint, contracts []*MemberContract) ([]*MemberContract, error)
	UpdateContract(ctx context.Context, memberID uint, contracts []*MemberContract) ([]*MemberContract, error)
}

type memberContractItemRepository interface {
	MemberContractItems(ctx context.Context, memberID, contractID uint, operators ...*Operator) ([]*MemberContractItem, error)
	CreateMemberContractItems(ctx context.Context, memberID, contractID uint, items []*MemberContractItem) ([]*MemberContractItem, error)
}

type memberContractBidRepository interface {
	MemberContractBids(ctx context.Context, memberID, contractID uint, operators ...*Operator) ([]*MemberContractBid, error)
	CreateMemberContractBids(ctx context.Context, memberID, contractID uint, bids []*MemberContractBid) ([]*MemberContractBid, error)
}

type MemberContract struct {
	MemberID            uint                 `db:"member_id" json:"member_id"`
	ContractID          uint                 `db:"contract_id" json:"contract_id"`
	AcceptorID          null.Uint            `db:"acceptor_id" json:"acceptor_id"`
	AssigneeID          null.Uint            `db:"assignee_id" json:"assignee_id"`
	Availability        ContractAvailability `db:"availability" json:"availability"`
	Buyout              null.Float64         `db:"buyout,omitempty" json:"buyout,omitempty"`
	Collateral          null.Float64         `db:"collateral,omitempty" json:"collateral,omitempty"`
	DateAccepted        null.Time            `db:"date_accepted,omitempty" json:"date_accepted,omitempty"`
	DateCompleted       null.Time            `db:"date_completed,omitempty" json:"date_completed,omitempty"`
	DateExpired         time.Time            `db:"date_expired" json:"date_expired"`
	DateIssued          time.Time            `db:"date_issued" json:"date_issued"`
	DaysToComplete      null.Uint            `db:"days_to_complete,omitempty" json:"days_to_complete,omitempty"`
	EndLocationID       null.Uint64          `db:"end_location_id,omitempty" json:"end_location_id,omitempty"`
	ForCorporation      bool                 `db:"for_corporation" json:"for_corporation"`
	IssuerCorporationID uint                 `db:"issuer_corporation_id" json:"issuer_corporation_id"`
	IssuerID            uint64               `db:"issuer_id" json:"issuer_id"`
	Price               null.Uint            `db:"price,omitempty" json:"price,omitempty"`
	Reward              null.Uint            `db:"reward,omitempty" json:"reward,omitempty"`
	StartLocationID     null.Uint64          `db:"start_location_id,omitempty" json:"start_location_id,omitempty"`
	Status              ContractStatus       `db:"status" json:"status"`
	Title               null.String          `db:"title,omitempty" json:"title,omitempty"`
	Type                ContractType         `db:"type" json:"type"`
	Volume              null.Float64         `db:"volume,omitempty" json:"volume,omitempty"`
	CreatedAt           time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time            `db:"updated_at" json:"updated_at"`
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

type MemberContractBid struct {
	MemberID   uint `db:"member_id" json:"member_id"`
	ContractID uint `db:"contract_id" json:"contract_id"`

	// Unique ID for the bid
	BidID uint `db:"bid_id" json:"bid_id"`

	// Character ID of the bidder
	BidderID uint64 `db:"bidder" json:"bidder"`

	// The amount bid, in ISK
	Amount float64 `db:"amount" json:"amount"`

	// Datetime when the bid was placed
	BidDate time.Time `db:"bid_date" json:"bid_date"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MemberContractItem struct {
	MemberID   uint `db:"member_id" json:"member_id"`
	ContractID uint `db:"contract_id" json:"contract_id"`

	// UniqueID for the Item
	RecordID uint `db:"record_id" json:"record_id"`

	// Type ID for item
	TypeID uint `db:"type_id" json:"type_id"`

	// Number of items in the stack
	Quantity uint `db:"quantity" json:"quantity"`

	// -1 indicates that the item is a singleton (non-stackable).
	// If the item happens to be a Blueprint, -1 is an Original and -2 is a Blueprint Copy
	RawQuantity int `db:"raw_quantity" json:"raw_quantity"`

	// true if the contract issuer has submitted this item with the contract, false if the isser is asking for this item in the contract
	IsIncluded  bool `db:"is_included" json:"is_included"`
	IsSingleton bool `db:"is_singleton" json:"is_singleton"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
