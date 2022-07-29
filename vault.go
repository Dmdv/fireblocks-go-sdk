package fireblocksdk

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Filters

// VaultAccountsFilter defines parameters for GetVaultAccountsPaged.
type VaultAccountsFilter struct {
	NamePrefix         string `json:"namePrefix,omitempty"`
	NameSuffix         string `json:"nameSuffix,omitempty"`
	MinAmountThreshold string `json:"minAmountThreshold,omitempty"`
	AssetID            string `json:"assetId,omitempty"`
}

// PagedVaultAccountsRequestFilters defines parameters for GetVaultAccountsPaged.
type PagedVaultAccountsRequestFilters struct {
	NamePrefix         string `json:"namePrefix,omitempty"`
	NameSuffix         string `json:"nameSuffix,omitempty"`
	MinAmountThreshold string `json:"minAmountThreshold,omitempty"`
	AssetID            string `json:"assetId,omitempty"`
	OrderBy            string `json:"orderBy,omitempty"`
	Before             string `json:"before,omitempty"`
	After              string `json:"after,omitempty"`
	Limit              int64  `json:"limit,omitempty"`
}

// VaultAccount endpoint

/*
export interface VaultAccountResponse {
    id: string;
    name: string;
    hiddenOnUI: boolean;
    assets: AssetResponse[];
    customerRefId?: string;
    autoFuel: boolean;
}
*/

// VaultAccountResponse defines model for VaultAccount.
type VaultAccountResponse struct {
	ID            string           `json:"id,omitempty"`
	Name          string           `json:"name,omitempty"`
	Assets        []*AssetResponse `json:"assets,omitempty"`
	CustomerRefID *string          `json:"customerRefId,omitempty"`
	AutoFuel      bool             `json:"autoFuel,omitempty"`
	HiddenOnUI    bool             `json:"hiddenOnUI,omitempty"`
}

// GetVaultAccounts Gets all assets that are currently supported by Fireblocks
func (sdk *FireblocksSDK) GetVaultAccounts(q *VaultAccountsFilter) (resp []*VaultAccountResponse, err error) {
	query := BuildQuery(q).UrlValues()
	body, status, err := sdk.client.DoGetRequest("/vault/accounts", query)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}
