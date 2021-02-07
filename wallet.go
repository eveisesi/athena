package athena

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/volatiletech/null"
)

type MemberWalletRepository interface {
	memberWalletBalanceRepository
	memberWalletTransactionRepository
	memberWalletJournalRepository
}

type memberWalletBalanceRepository interface {
	MemberWalletBalance(ctx context.Context, memberID uint) (*MemberWalletBalance, error)
	CreateMemberWalletBalance(ctx context.Context, memberID uint, balance float64) (*MemberWalletBalance, error)
	UpdateMemberWalletBalance(ctx context.Context, memberID uint, balance float64) (*MemberWalletBalance, error)
}

type memberWalletTransactionRepository interface {
	MemberWalletTransactions(ctx context.Context, operators ...*Operator) ([]*MemberWalletTransaction, error)
	CreateMemberWalletTransactions(ctx context.Context, memberID uint, transactions []*MemberWalletTransaction) ([]*MemberWalletTransaction, error)
}

type memberWalletJournalRepository interface {
	MemberWalletJournals(ctx context.Context, operators ...*Operator) ([]*MemberWalletJournal, error)
	CreateMemberWalletJournals(ctx context.Context, memberID uint, entries []*MemberWalletJournal) ([]*MemberWalletJournal, error)
}

type MemberWalletBalance struct {
	MemberID  uint      `db:"member_id" json:"member_id"`
	Balance   float64   `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type MemberWalletTransaction struct {
	MemberID           uint         `db:"member_id" json:"member_id"`
	TransactionID      uint64       `db:"transaction_id" json:"transaction_id"`
	JournalReferenceID uint64       `db:"journal_ref_id" json:"journal_ref_id"`
	ClientID           uint32       `db:"client_id" json:"client_id"`
	ClientType         ClientType   `db:"client_type" json:"client_type"`
	LocationID         uint64       `db:"location_id" json:"location_id"`
	LocationType       LocationType `db:"location_type" json:"location_type"`
	TypeID             uint         `db:"type_id" json:"type_id"`
	Quantity           uint         `db:"quantity" json:"quantity"`
	UnitPrice          float64      `db:"unit_price" json:"unit_price"`
	IsBuy              bool         `db:"is_buy" json:"is_buy"`
	IsPersonal         bool         `db:"is_personal" json:"is_personal"`
	Date               time.Time    `db:"date" json:"date"`
	CreatedAt          time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time    `db:"updated_at" json:"updated_at"`
}

type ClientType string

const (
	ClientTypeCharacter   ClientType = "Character"
	ClientTypeCorporation ClientType = "Corporation"
)

var AllClientTypes = []ClientType{
	ClientTypeCharacter,
	ClientTypeCorporation,
}

func (i ClientType) Valid() bool {
	for _, v := range AllClientTypes {
		if i == v {
			return true
		}
	}

	return false
}

func (i ClientType) String() string {
	return string(i)
}

type LocationType string

const (
	LocationTypeStation   LocationType = "Station"
	LocationTypeStructure LocationType = "Structure"
)

var AllLocationTypes = []LocationType{
	LocationTypeStation,
	LocationTypeStructure,
}

func (i LocationType) Valid() bool {
	for _, v := range AllLocationTypes {
		if i == v {
			return true
		}
	}

	return false
}

func (i LocationType) String() string {
	return string(i)
}

type MemberWalletJournal struct {
	MemberID        uint                  `db:"member_id" json:"member_id"`
	JournalID       uint64                `db:"journal_id" json:"id"`
	RefType         RefType               `db:"ref_type" json:"ref_type"`
	ContextID       null.Int64            `db:"context_id,omitempty" json:"context_id,omitempty"`
	ContextType     NullableContextIDType `db:"context_id_type,omitempty" json:"context_id_type,omitempty"`
	Description     string                `db:"description" json:"description"`
	Reason          null.String           `db:"reason,omitempty" json:"reason,omitempty"`
	FirstPartyID    null.Int              `db:"first_party_id,omitempty" json:"first_party_id,omitempty"`
	FirstPartyType  null.String           `db:"first_party_type,omitempty" json:"first_party_type,omitempty"`
	SecondPartyID   null.Int              `db:"second_party_id,omitempty" json:"second_party_id,omitempty"`
	SecondPartyType null.String           `db:"second_party_type,omitempty" json:"second_party_type,omitempty"`
	Amount          null.Float64          `db:"amount,omitempty" json:"amount,omitempty"`
	Balance         null.Float64          `db:"balance,omitempty" json:"balance,omitempty"`
	Tax             null.Float64          `db:"tax,omitempty" json:"tax,omitempty"`
	TaxReceiverID   null.Int              `db:"tax_receiver_id,omitempty" json:"tax_receiver_id,omitempty"`
	Date            time.Time             `db:"date" json:"date"`
	CreatedAt       time.Time             `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time             `db:"updated_at" json:"updated_at"`
}

type ContextIDType string

const (
	ContextIDTypeStructureID         ContextIDType = "structure_id"
	ContextIDTypeStationID           ContextIDType = "station_id"
	ContextIDTypeMarketTransactionID ContextIDType = "market_transaction_id"
	ContextIDTypeCharacterID         ContextIDType = "character_id"
	ContextIDTypeCorporationID       ContextIDType = "corporation_id"
	ContextIDTypeAllianceID          ContextIDType = "alliance_id"
	ContextIDTypeEveSystem           ContextIDType = "eve_system"
	ContextIDTypeIndustryJobID       ContextIDType = "industry_job_id"
	ContextIDTypeContractID          ContextIDType = "contract_id"
	ContextIDTypePlanetID            ContextIDType = "planet_id"
	ContextIDTypeSystemID            ContextIDType = "system_id"
	ContextIDTypeTypeID              ContextIDType = "type_id"
)

var AllContextIDType = []ContextIDType{
	ContextIDTypeStructureID,
	ContextIDTypeStationID,
	ContextIDTypeMarketTransactionID,
	ContextIDTypeCharacterID,
	ContextIDTypeCorporationID,
	ContextIDTypeAllianceID,
	ContextIDTypeEveSystem,
	ContextIDTypeIndustryJobID,
	ContextIDTypeContractID,
	ContextIDTypePlanetID,
	ContextIDTypeSystemID,
	ContextIDTypeTypeID,
}

func (i ContextIDType) Valid() bool {
	for _, v := range AllContextIDType {
		if i == v {
			return true
		}
	}

	return false
}

func (i ContextIDType) String() string {
	return string(i)
}

type RefType string

const (
	RefTypeAccelerationGateFee                       RefType = "acceleration_gate_fee"
	RefTypeAdvertisementListingFee                   RefType = "advertisement_listing_fee"
	RefTypeAgentDonation                             RefType = "agent_donation"
	RefTypeAgentLocationServices                     RefType = "agent_location_services"
	RefTypeAgentMiscellaneous                        RefType = "agent_miscellaneous"
	RefTypeAgentMissionCollateralPaid                RefType = "agent_mission_collateral_paid"
	RefTypeAgentMissionCollateralRefunded            RefType = "agent_mission_collateral_refunded"
	RefTypeAgentMissionReward                        RefType = "agent_mission_reward"
	RefTypeAgentMissionRewardCorporationTax          RefType = "agent_mission_reward_corporation_tax"
	RefTypeAgentMissionTimeBonusReward               RefType = "agent_mission_time_bonus_reward"
	RefTypeAgentMissionTimeBonusRewardCorporationTax RefType = "agent_mission_time_bonus_reward_corporation_tax"
	RefTypeAgentSecurityServices                     RefType = "agent_security_services"
	RefTypeAgentServicesRendered                     RefType = "agent_services_rendered"
	RefTypeAgentsPreward                             RefType = "agents_preward"
	RefTypeAllianceMaintainanceFee                   RefType = "alliance_maintainance_fee"
	RefTypeAllianceRegistrationFee                   RefType = "alliance_registration_fee"
	RefTypeAssetSafetyRecoveryTax                    RefType = "asset_safety_recovery_tax"
	RefTypeBounty                                    RefType = "bounty"
	RefTypeBountyPrize                               RefType = "bounty_prize"
	RefTypeBountyPrizeCorporationTax                 RefType = "bounty_prize_corporation_tax"
	RefTypeBountyPrizes                              RefType = "bounty_prizes"
	RefTypeBountyReimbursement                       RefType = "bounty_reimbursement"
	RefTypeBountySurcharge                           RefType = "bounty_surcharge"
	RefTypeBrokersFee                                RefType = "brokers_fee"
	RefTypeCloneActivation                           RefType = "clone_activation"
	RefTypeCloneTransfer                             RefType = "clone_transfer"
	RefTypeContrabandFine                            RefType = "contraband_fine"
	RefTypeContractAuctionBid                        RefType = "contract_auction_bid"
	RefTypeContractAuctionBidCorp                    RefType = "contract_auction_bid_corp"
	RefTypeContractAuctionBidRefund                  RefType = "contract_auction_bid_refund"
	RefTypeContractAuctionSold                       RefType = "contract_auction_sold"
	RefTypeContractBrokersFee                        RefType = "contract_brokers_fee"
	RefTypeContractBrokersFeeCorp                    RefType = "contract_brokers_fee_corp"
	RefTypeContractCollateral                        RefType = "contract_collateral"
	RefTypeContractCollateralDepositedCorp           RefType = "contract_collateral_deposited_corp"
	RefTypeContractCollateralPayout                  RefType = "contract_collateral_payout"
	RefTypeContractCollateralRefund                  RefType = "contract_collateral_refund"
	RefTypeContractDeposit                           RefType = "contract_deposit"
	RefTypeContractDepositCorp                       RefType = "contract_deposit_corp"
	RefTypeContractDepositRefund                     RefType = "contract_deposit_refund"
	RefTypeContractDepositSalesTax                   RefType = "contract_deposit_sales_tax"
	RefTypeContractPrice                             RefType = "contract_price"
	RefTypeContractPricePaymentCorp                  RefType = "contract_price_payment_corp"
	RefTypeContractReversal                          RefType = "contract_reversal"
	RefTypeContractReward                            RefType = "contract_reward"
	RefTypeContractRewardDeposited                   RefType = "contract_reward_deposited"
	RefTypeContractRewardDepositedCorp               RefType = "contract_reward_deposited_corp"
	RefTypeContractRewardRefund                      RefType = "contract_reward_refund"
	RefTypeContractSalesTax                          RefType = "contract_sales_tax"
	RefTypeCopying                                   RefType = "copying"
	RefTypeCorporateRewardPayout                     RefType = "corporate_reward_payout"
	RefTypeCorporateRewardTax                        RefType = "corporate_reward_tax"
	RefTypeCorporationAccountWithdrawal              RefType = "corporation_account_withdrawal"
	RefTypeCorporationBulkPayment                    RefType = "corporation_bulk_payment"
	RefTypeCorporationDividendPayment                RefType = "corporation_dividend_payment"
	RefTypeCorporationLiquidation                    RefType = "corporation_liquidation"
	RefTypeCorporationLogoChangeCost                 RefType = "corporation_logo_change_cost"
	RefTypeCorporationPayment                        RefType = "corporation_payment"
	RefTypeCorporationRegistrationFee                RefType = "corporation_registration_fee"
	RefTypeCourierMissionEscrow                      RefType = "courier_mission_escrow"
	RefTypeCSPA                                      RefType = "cspa"
	RefTypeCSPAOfflineRefund                         RefType = "cspaofflinerefund"
	RefTypeDatacoreFee                               RefType = "datacore_fee"
	RefTypeDnaModificationFee                        RefType = "dna_modification_fee"
	RefTypeDockingFee                                RefType = "docking_fee"
	RefTypeDuelWagerEscrow                           RefType = "duel_wager_escrow"
	RefTypeDuelWagerPayment                          RefType = "duel_wager_payment"
	RefTypeDuelWagerRefund                           RefType = "duel_wager_refund"
	RefTypeEssEscrowTransfer                         RefType = "ess_escrow_transfer"
	RefTypeFactorySlotRentalFee                      RefType = "factory_slot_rental_fee"
	RefTypeGmCashTransfer                            RefType = "gm_cash_transfer"
	RefTypeIndustryJobTax                            RefType = "industry_job_tax"
	RefTypeInfrastructureHubMaintenance              RefType = "infrastructure_hub_maintenance"
	RefTypeInheritance                               RefType = "inheritance"
	RefTypeInsurance                                 RefType = "insurance"
	RefTypeItemTraderPayment                         RefType = "item_trader_payment"
	RefTypeJumpCloneActivationFee                    RefType = "jump_clone_activation_fee"
	RefTypeJumpCloneInstallationFee                  RefType = "jump_clone_installation_fee"
	RefTypeKillRightFee                              RefType = "kill_right_fee"
	RefTypeLpStore                                   RefType = "lp_store"
	RefTypeManufacturing                             RefType = "manufacturing"
	RefTypeMarketEscrow                              RefType = "market_escrow"
	RefTypeMarketFinePaid                            RefType = "market_fine_paid"
	RefTypeMarketTransaction                         RefType = "market_transaction"
	RefTypeMedalCreation                             RefType = "medal_creation"
	RefTypeMedalIssued                               RefType = "medal_issued"
	RefTypeMissionCompletion                         RefType = "mission_completion"
	RefTypeMissionCost                               RefType = "mission_cost"
	RefTypeMissionExpiration                         RefType = "mission_expiration"
	RefTypeMissionReward                             RefType = "mission_reward"
	RefTypeOfficeRentalFee                           RefType = "office_rental_fee"
	RefTypeOperationBonus                            RefType = "operation_bonus"
	RefTypeOpportunityReward                         RefType = "opportunity_reward"
	RefTypePlanetaryConstruction                     RefType = "planetary_construction"
	RefTypePlanetaryExportTax                        RefType = "planetary_export_tax"
	RefTypePlanetaryImportTax                        RefType = "planetary_import_tax"
	RefTypePlayerDonation                            RefType = "player_donation"
	RefTypePlayerTrading                             RefType = "player_trading"
	RefTypeProjectDiscoveryReward                    RefType = "project_discovery_reward"
	RefTypeProjectDiscoveryTax                       RefType = "project_discovery_tax"
	RefTypeReaction                                  RefType = "reaction"
	RefTypeReleaseOfImpoundedProperty                RefType = "release_of_impounded_property"
	RefTypeRepairBill                                RefType = "repair_bill"
	RefTypeReprocessingTax                           RefType = "reprocessing_tax"
	RefTypeResearchingMaterialProductivity           RefType = "researching_material_productivity"
	RefTypeResearchingTechnology                     RefType = "researching_technology"
	RefTypeResearchingTimeProductivity               RefType = "researching_time_productivity"
	RefTypeResourceWarsReward                        RefType = "resource_wars_reward"
	RefTypeReverseEngineering                        RefType = "reverse_engineering"
	RefTypeSecurityProcessingFee                     RefType = "security_processing_fee"
	RefTypeShares                                    RefType = "shares"
	RefTypeSkillPurchase                             RefType = "skill_purchase"
	RefTypeSovereignityBill                          RefType = "sovereignity_bill"
	RefTypeStorePurchase                             RefType = "store_purchase"
	RefTypeStorePurchaseRefund                       RefType = "store_purchase_refund"
	RefTypeStructureGateJump                         RefType = "structure_gate_jump"
	RefTypeTransactionTax                            RefType = "transaction_tax"
	RefTypeUpkeepAdjustmentFee                       RefType = "upkeep_adjustment_fee"
	RefTypeWarAllyContract                           RefType = "war_ally_contract"
	RefTypeWarFee                                    RefType = "war_fee"
	RefTypeWarFeeSurrender                           RefType = "war_fee_surrender"
)

var AllRefTypes = []RefType{
	RefTypeAccelerationGateFee, RefTypeAdvertisementListingFee, RefTypeAgentDonation, RefTypeAgentLocationServices,
	RefTypeAgentMiscellaneous, RefTypeAgentMissionCollateralPaid, RefTypeAgentMissionCollateralRefunded, RefTypeAgentMissionReward,
	RefTypeAgentMissionRewardCorporationTax, RefTypeAgentMissionTimeBonusReward, RefTypeAgentMissionTimeBonusRewardCorporationTax, RefTypeAgentSecurityServices,
	RefTypeAgentServicesRendered, RefTypeAgentsPreward, RefTypeAllianceMaintainanceFee, RefTypeAllianceRegistrationFee,
	RefTypeAssetSafetyRecoveryTax, RefTypeBounty, RefTypeBountyPrize, RefTypeBountyPrizeCorporationTax,
	RefTypeBountyPrizes, RefTypeBountyReimbursement, RefTypeBountySurcharge, RefTypeBrokersFee,
	RefTypeCloneActivation, RefTypeCloneTransfer, RefTypeContrabandFine, RefTypeContractAuctionBid,
	RefTypeContractAuctionBidCorp, RefTypeContractAuctionBidRefund, RefTypeContractAuctionSold, RefTypeContractBrokersFee,
	RefTypeContractBrokersFeeCorp, RefTypeContractCollateral, RefTypeContractCollateralDepositedCorp, RefTypeContractCollateralPayout,
	RefTypeContractCollateralRefund, RefTypeContractDeposit, RefTypeContractDepositCorp, RefTypeContractDepositRefund,
	RefTypeContractDepositSalesTax, RefTypeContractPrice, RefTypeContractPricePaymentCorp, RefTypeContractReversal,
	RefTypeContractReward, RefTypeContractRewardDeposited, RefTypeContractRewardDepositedCorp, RefTypeContractRewardRefund,
	RefTypeContractSalesTax, RefTypeCopying, RefTypeCorporateRewardPayout, RefTypeCorporateRewardTax,
	RefTypeCorporationAccountWithdrawal, RefTypeCorporationBulkPayment, RefTypeCorporationDividendPayment, RefTypeCorporationLiquidation,
	RefTypeCorporationLogoChangeCost, RefTypeCorporationPayment, RefTypeCorporationRegistrationFee, RefTypeCourierMissionEscrow,
	RefTypeCSPA, RefTypeCSPAOfflineRefund, RefTypeDatacoreFee, RefTypeDnaModificationFee,
	RefTypeDockingFee, RefTypeDuelWagerEscrow, RefTypeDuelWagerPayment, RefTypeDuelWagerRefund,
	RefTypeEssEscrowTransfer, RefTypeFactorySlotRentalFee, RefTypeGmCashTransfer, RefTypeIndustryJobTax,
	RefTypeInfrastructureHubMaintenance, RefTypeInheritance, RefTypeInsurance, RefTypeItemTraderPayment,
	RefTypeJumpCloneActivationFee, RefTypeJumpCloneInstallationFee, RefTypeKillRightFee, RefTypeLpStore,
	RefTypeManufacturing, RefTypeMarketEscrow, RefTypeMarketFinePaid, RefTypeMarketTransaction,
	RefTypeMedalCreation, RefTypeMedalIssued, RefTypeMissionCompletion, RefTypeMissionCost,
	RefTypeMissionExpiration, RefTypeMissionReward, RefTypeOfficeRentalFee, RefTypeOperationBonus,
	RefTypeOpportunityReward, RefTypePlanetaryConstruction, RefTypePlanetaryExportTax, RefTypePlanetaryImportTax,
	RefTypePlayerDonation, RefTypePlayerTrading, RefTypeProjectDiscoveryReward, RefTypeProjectDiscoveryTax,
	RefTypeReaction, RefTypeReleaseOfImpoundedProperty, RefTypeRepairBill, RefTypeReprocessingTax,
	RefTypeResearchingMaterialProductivity, RefTypeResearchingTechnology, RefTypeResearchingTimeProductivity, RefTypeResourceWarsReward,
	RefTypeReverseEngineering, RefTypeSecurityProcessingFee, RefTypeShares, RefTypeSkillPurchase,
	RefTypeSovereignityBill, RefTypeStorePurchase, RefTypeStorePurchaseRefund, RefTypeStructureGateJump,
	RefTypeTransactionTax, RefTypeUpkeepAdjustmentFee, RefTypeWarAllyContract, RefTypeWarFee,
	RefTypeWarFeeSurrender,
}

func (i RefType) Valid() bool {
	for _, v := range AllRefTypes {
		if i == v {
			return true
		}
	}

	return false
}

func (i RefType) String() string {
	return string(i)
}

type NullableContextIDType struct {
	Valid         bool
	ContextIDType ContextIDType
}

// MarshalJSON implements json.Marshaler.
func (i NullableContextIDType) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return null.NullBytes, nil
	}
	return []byte(i.ContextIDType.String()), nil
}

func (i *NullableContextIDType) UnmarshalJSON(data []byte) error {

	if bytes.Equal(data, null.NullBytes) {
		i.Valid = false
		i.ContextIDType = ""
		return nil
	}

	if err := json.Unmarshal(data, &i.ContextIDType); err != nil {
		return err
	}

	i.Valid = true
	return nil

}

// IsZero returns true for invalid NullableContextIDType's, for future omitempty support (Go 1.4?)
func (i NullableContextIDType) IsZero() bool {
	return !i.Valid
}
