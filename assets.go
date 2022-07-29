package fireblocksdk

type AssetTypeResponse struct {
	ID              string
	Name            string
	AssetType       string
	ContractAddress string
	NativeAsset     string
	Cecimals        *int64
}

// CreateVaultAssetResponse defines model for CreateVaultAssetResponse.
type CreateVaultAssetResponse struct {
	ActivationTxID    string `json:"activationTxId,omitempty"`
	Address           string `json:"address,omitempty"`
	EnterpriseAddress string `json:"enterpriseAddress,omitempty"`
	EosAccountName    string `json:"eosAccountName,omitempty"`
	ID                string `json:"id,omitempty"`
	LegacyAddress     string `json:"legacyAddress,omitempty"`
	Status            string `json:"status,omitempty"`
	Tag               string `json:"tag,omitempty"`
}

// AllocatedBalance defines model for AllocatedBalance.
type AllocatedBalance struct {
	Affiliation         string `json:"affiliation,omitempty"`
	AllocationID        string `json:"allocationId,omitempty"`
	Available           string `json:"available,omitempty"`
	Frozen              string `json:"frozen,omitempty"`
	Locked              string `json:"locked,omitempty"`
	Pending             string `json:"pending,omitempty"`
	Staked              string `json:"staked,omitempty"`
	ThirdPartyAccountID string `json:"thirdPartyAccountId,omitempty"`
	Total               string `json:"total,omitempty"`
	VirtualType         string `json:"virtualType,omitempty"`
}

// RewardsInfo defines model for RewardsInfo.
type RewardsInfo struct {
	// Amount that is pending for rewards
	PendingRewards string `json:"pendingRewards,omitempty"`
}

/*
export interface AssetResponse {
    id: string;
    total: string;
	// Deprecated, replaced by total
    balance?: string;
    lockedAmount?: string;
    available: string;
    pending: string;
    selfStakedCPU?: string;
    selfStakedNetwork?: string;
    pendingRefundCPU?: string;
    pendingRefundNetwork?: string;
    totalStakedCPU?: string;
    totalStakedNetwork?: string;
    rewardInfo?: BalanceRewardInfo;
    blockHeight?: string;
    blockHash?: string;
    allocatedBalances?: {
        allocationId: string;
        thirdPartyAccountId?: string;
        affiliation?: VirtualAffiliation;
        virtualType?: VirtualType;
        total: string;
        available: string;
        pending?: string;
        frozen?: string;
        locked?: string;
    }[];
}
*/

// AssetResponse defines model for VaultAsset.
type AssetResponse struct {
	ID                   string              `json:"id,omitempty"`
	Total                string              `json:"total,omitempty"`   // The total wallet balance. In EOS this value includes the network balance, self staking and pending refund. For all other coins it is the balance as it appears on the blockchain.
	Balance              *string             `json:"balance,omitempty"` // Deprecated - replaced by "total"
	AllocatedBalances    *[]AllocatedBalance `json:"allocatedBalances,omitempty"`
	Available            string              `json:"available,omitempty"` // Funds available for transfer. Equals the blockchain balance minus any locked amounts
	BlockHash            string              `json:"blockHash,omitempty"`
	BlockHeight          string              `json:"blockHeight,omitempty"`
	Frozen               string              `json:"frozen,omitempty"`               // The cumulative frozen balance by AML policy
	LockedAmount         *string             `json:"lockedAmount,omitempty"`         // Funds in outgoing transactions that are not yet published to the network
	Pending              string              `json:"pending,omitempty"`              // The cumulative balance of all transactions pending to be cleared
	PendingRefundCPU     string              `json:"pendingRefundCPU,omitempty"`     // Deprecated
	PendingRefundNetwork string              `json:"pendingRefundNetwork,omitempty"` // Deprecated
	RewardsInfo          *RewardsInfo        `json:"rewardsInfo,omitempty"`
	SelfStakedCPU        *string             `json:"selfStakedCPU,omitempty"`      // Deprecated
	SelfStakedNetwork    *string             `json:"selfStakedNetwork,omitempty"`  // Deprecated
	Staked               string              `json:"staked,omitempty"`             // Staked funds, returned only for DOT
	TotalStakedCPU       *string             `json:"totalStakedCPU,omitempty"`     // Deprecated
	TotalStakedNetwork   *string             `json:"totalStakedNetwork,omitempty"` // Deprecated
}
