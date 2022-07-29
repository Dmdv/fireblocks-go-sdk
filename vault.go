package fireblocksdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// VaultAccount endpoint

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
	OrderBy            string `json:"orderBy,omitempty"` // ASC | DESC, The results are ordered by the creation time of the vault account.
	Before             string `json:"before,omitempty"`  // [optional] cursor string, if specified then we give the next results after this cursor
	After              string `json:"after,omitempty"`   // [optional] cursor string, if specified then we give the next results before this cursor
	Limit              int64  `json:"limit,omitempty"`   // Returns the maximum number of valut accounts in a single response. The default value is 300 and maximum value is 500.
}

// Responses

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

/*
export interface PagedVaultAccountsResponse {
    accounts: VaultAccountResponse[];
    paging: {
        before: string;
        after: string;
    };
    previousUrl: string;
    nextUrl: string;
}
*/

type PagedVaultAccountsResponse struct {
	Accounts []VaultAccountResponse `json:"accounts,omitempty"`
	Paging   struct {
		Before string `json:"before,omitempty"`
		After  string `json:"after,omitempty"`
	} `json:"paging"`
	PreviousUrl string `json:"previousUrl,omitempty"`
	NextUrl     string `json:"nextUrl,omitempty"`
}

// GetVaultAccounts Deprecated, Gets all assets that are currently supported by Fireblocks,
func (sdk *FireblocksSDK) GetVaultAccounts(q *VaultAccountsFilter) (resp []*VaultAccountResponse, err error) {
	query := BuildQuery(q).UrlValues()
	body, status, err := sdk.client.DoGetRequest("/vault/accounts", query)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// GetVaultAccountsWithPageInfo Gets all assets that are currently supported by Fireblocks
func (sdk *FireblocksSDK) GetVaultAccountsWithPageInfo(q *PagedVaultAccountsRequestFilters) (resp *PagedVaultAccountsResponse, err error) {
	query := BuildQuery(q).UrlValues()
	body, status, err := sdk.client.DoGetRequest("/vault/accounts_paged", query)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

func (sdk *FireblocksSDK) GetVaultAccountsByID(accountID string) (resp *VaultAccountResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf("/vault/accounts/%s", accountID), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}
