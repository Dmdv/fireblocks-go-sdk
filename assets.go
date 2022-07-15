package fireblocksdk

import "encoding/json"

type AssetTypeResponse struct {
	id              string
	name            string
	assetType       string
	contractAddress string
	nativeAsset     string
	decimals        *int64
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

// VaultAsset defines model for VaultAsset.
type VaultAsset struct {
	AllocatedBalances *[]AllocatedBalance `json:"allocatedBalances,omitempty"`
	// Funds available for transfer. Equals the blockchain balance minus any locked amounts
	Available string `json:"available,omitempty"`
	// Deprecated - replaced by "total"
	Balance     string `json:"balance,omitempty"`
	BlockHash   string `json:"blockHash,omitempty"`
	BlockHeight string `json:"blockHeight,omitempty"`
	// The cumulative frozen balance
	Frozen string `json:"frozen,omitempty"`
	Id     string `json:"id,omitempty"`
	// Funds in outgoing transactions that are not yet published to the network
	LockedAmount string `json:"lockedAmount,omitempty"`
	// The cumulative balance of all transactions pending to be cleared
	Pending              string       `json:"pending,omitempty"`
	PendingRefundCPU     string       `json:"pendingRefundCPU,omitempty"`
	PendingRefundNetwork string       `json:"pendingRefundNetwork,omitempty"`
	RewardsInfo          *RewardsInfo `json:"rewardsInfo,omitempty"`
	SelfStakedCPU        string       `json:"selfStakedCPU,omitempty"`
	SelfStakedNetwork    string       `json:"selfStakedNetwork,omitempty"`
	// Staked balance
	Staked string `json:"staked,omitempty"`
	// The total wallet balance. In EOS this value includes the network balance, self staking and pending refund. For all other coins it is the balance as it appears on the blockchain.
	Total              string `json:"total,omitempty"`
	TotalStakedCPU     string `json:"totalStakedCPU,omitempty"`
	TotalStakedNetwork string `json:"totalStakedNetwork,omitempty"`
}

// UnmarshalJSON handles deserialization of a VaultAsset.
// This custom unmarshalling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (c *VaultAsset) UnmarshalJSON(data []byte) error {
	type va VaultAsset
	var v va
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*c = VaultAsset(v)

	return nil
}
