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
	Limit              int64  `json:"limit,omitempty"`   // Returns the maximum number of vault accounts in a single response. The default value is 300 and maximum value is 500.
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

type VaultAccountRequest struct {
	Name          string `json:"name"`
	CustomerRefID string `json:"customerRefId,omitempty"`
	HiddenOnUI    *bool  `json:"hiddenOnUi,omitempty"`
	AutoFuel      *bool  `json:"autoFuel,omitempty"`
}

// VaultAccountResponse defines model for VaultAccount.
type VaultAccountResponse struct {
	ID            string           `json:"id,omitempty"`
	Name          string           `json:"name,omitempty"`
	Assets        []*AssetResponse `json:"assets,omitempty"`
	CustomerRefID *string          `json:"customerRefId,omitempty"`
	AutoFuel      *bool            `json:"autoFuel,omitempty"`
	HiddenOnUI    *bool            `json:"hiddenOnUI,omitempty"`
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
	PreviousURL string `json:"previousUrl,omitempty"`
	NextURL     string `json:"nextUrl,omitempty"`
}

/*
export interface DepositAddressResponse {
    assetId: string;
    address: string;
    tag?: string;
    description?: string;
    type: string;
    customerRefId?: string;
    addressFormat: string;
    legacyAddress?: string;
    enterpriseAddress?: string;
}
*/

type DepositAddressResponse struct {
	AssetID           string `json:"assetId,omitempty"`           // The ID of the asset
	Address           string `json:"address,omitempty"`           // Address of the asset in a Vault Account, for BTC/LTC the address is in Segwit (Bech32) format, for BCH cash format
	Tag               string `json:"tag,omitempty"`               // Destination tag for XRP, used as memo for EOS/XLM, for the fiat providers (Signet by Signature, SEN by Silvergate, BLINC by BCB Group), it is the Bank Transfer
	Description       string `json:"description,omitempty"`       // Description of the address
	TypeAddress       string `json:"type,omitempty"`              // Address type
	CustomerRefID     string `json:"customerRefId,omitempty"`     // [optional] The ID for AML providers to associate the owner of funds with transactions
	LegacyAddress     string `json:"legacyAddress,omitempty"`     // For BTC/LTC/BCH the legacy format address
	AddressFormat     string `json:"addressFormat,omitempty"`     //
	EnterpriseAddress string `json:"enterpriseAddress,omitempty"` //
	Bip44AddressIndex int    `json:"bip44AddressIndex,omitempty"` // [optional] The address_index, addressFormat, and enterpriseAddress in the derivation path of this address based on BIP44
}

/*
export interface PublicKeyInfoForVaultAccountArgs {
    assetId: string;
    vaultAccountId: number;
    change: number;
    addressIndex: number;
    compressed?: boolean;
}
*/

/*type PublicKeyInfoForVaultAccountArgs struct {
	AssetId        string // The ID of the asset
	VaultAccountId int64  // The ID of the vault account which address should be retrieved, or 'default' for the default vault account
	Change         int64  // Whether the address should be derived internal (change) or not
	AddressIndex   int64  // The index of the address for the derivation path
	Compressed     bool   // Boolean, whether the returned key should be in compressed format or not, false by default
}*/

type PublicKeyInfoResponse struct {
	Status         int     `json:"status,omitempty"`
	PublicKey      string  `json:"publicKey,omitempty"`
	Algorithm      string  `json:"algorithm,omitempty"`
	DerivationPath []int64 `json:"derivationPath,omitempty"`
}

/*
export interface GenerateAddressResponse {
    address: string;
    tag?: string;
    legacyAddress?: string;
    enterpriseAddress?: string;
}
*/

type GenerateAddressResponse struct {
	Address           string  `json:"address,omitempty"`
	Tag               *string `json:"tag,omitempty"`
	LegacyAddress     *string `json:"legacy_address,omitempty"`
	EnterpriseAddress *string `json:"enterpriseAddress,omitempty"`
	Bip44AddressIndex *int    `json:"bip44AddressIndex,omitempty"` // The address_index in the derivation path of this address based on BIP44
}

type PostOptions struct {
	Description   string `json:"description,omitempty"`
	CustomerRefID string `json:"customerRefId,omitempty"`
}

// GetVaultAccounts Deprecated, Gets all assets that are currently supported by Fireblocks,
func (sdk *FireblocksSDK) GetVaultAccounts(q *VaultAccountsFilter) (resp []*VaultAccountResponse, err error) {
	query := BuildQuery(q).URLValues()
	body, status, err := sdk.client.DoGetRequest("/vault/accounts", query)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// GetVaultAccountsWithPageInfo Gets all assets that are currently supported by Fireblocks
// Retrieves all vault accounts in your workspace. This endpoint returns a limited amount of results and quick response time.
func (sdk *FireblocksSDK) GetVaultAccountsWithPageInfo(q *PagedVaultAccountsRequestFilters) (resp *PagedVaultAccountsResponse, err error) {
	query := BuildQuery(q).URLValues()
	body, status, err := sdk.client.DoGetRequest("/vault/accounts_paged", query)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

func (sdk *FireblocksSDK) GetVaultAccountsByID(vaultAccountID string) (resp *VaultAccountResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf("/vault/accounts/%s", vaultAccountID), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// GetVaultAccountAsset Retrieves a wallet of a specific asset under a Fireblocks Vault Account.
func (sdk *FireblocksSDK) GetVaultAccountAsset(vaultAccountID, assetID string) (resp *AssetResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf("/vault/accounts/%s/%s", vaultAccountID, assetID), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// GetDepositAddresses Retrieves a wallet of a specific asset under a Fireblocks Vault Account.
func (sdk *FireblocksSDK) GetDepositAddresses(vaultAccountID, assetID string) (resp []*DepositAddressResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf("/vault/accounts/%s/%s/addresses", vaultAccountID, assetID), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// Response type supposed to be this - test it

/*
{
    "input": {
        "txHash": "string",
        "index": "number"
    },
    "address": "string",
    "amount": "string",
    "confirmations": "number",
    "status": "string"
}
*/

// GetUnspentInputs Returns unspent inputs of the requested asset in the Vault Account.
// Gets utxo list for an asset in a vault account
// vaultAccountId - The vault account ID
// assetId - The ID of the asset for which to get the utxo list
func (sdk *FireblocksSDK) GetUnspentInputs(vaultAccountID, assetID string) (resp []DepositAddressResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf("/vault/accounts/%s/%s/unspent_inputs", vaultAccountID, assetID), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

// GetPublicKeyInfoForVaultAccount Get the public key information for a vault account
func (sdk *FireblocksSDK) GetPublicKeyInfoForVaultAccount(
	vaultAccountID, assetID string,
	change int,
	addressIndex int,
) (resp *PublicKeyInfoResponse, err error) {
	body, status, err := sdk.client.DoGetRequest(fmt.Sprintf(
		"/vault/accounts/%s/%s/%v/%v/public_key_info",
		vaultAccountID,
		assetID,
		change,
		addressIndex,
	), nil)
	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

func (sdk *FireblocksSDK) GenerateNewAddress(vaultAccountID, assetID, description, customerRefID string, opts ...func(*PostRequestOption)) (resp *GenerateAddressResponse, err error) {
	post := &PostOptions{
		Description:   description,
		CustomerRefID: customerRefID,
	}

	body, status, err := sdk.client.DoPostRequest(
		fmt.Sprintf("/vault/accounts/%s/%s/addresses", vaultAccountID, assetID),
		post,
		opts...,
	)

	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}

func (sdk *FireblocksSDK) CreateVaultAccount(name string, customerRefID string, hiddenOnUI *bool, autoFuel *bool, opts ...func(*PostRequestOption)) (resp VaultAccountResponse, err error) {
	post := &VaultAccountRequest{
		Name:          name,
		HiddenOnUI:    hiddenOnUI,
		CustomerRefID: customerRefID,
		AutoFuel:      autoFuel,
	}

	body, status, err := sdk.client.DoPostRequest(
		"/vault/accounts",
		post,
		opts...,
	)

	if err == nil && status == http.StatusOK {
		err = json.Unmarshal(body, &resp)
		return
	}

	return resp, errors.Wrap(err, "failed to make request")
}
