package fireblocksdk

type AssetTypeResponse struct {
	id              string
	name            string
	assetType       string
	contractAddress string
	nativeAsset     string
	decimals        *int64
}
