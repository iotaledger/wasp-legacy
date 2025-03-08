package apiextensions

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/packages/isc"
)

func AssetsFromAPIResponse(assetsResponse *apiclient.AssetsResponse) (*isc.Assets, error) {
	assets := isc.NewEmptyAssets()

	baseTokens, err := iotago.DecodeUint64(assetsResponse.BaseTokens)
	if err != nil {
		return nil, err
	}

	assets.BaseTokens = baseTokens

	for _, nativeToken := range assetsResponse.NativeTokens {
		nativeTokenIDHex, err2 := iotago.DecodeHex(nativeToken.Id)
		if err2 != nil {
			return nil, err2
		}

		nativeTokenID, err2 := isc.NativeTokenIDFromBytes(nativeTokenIDHex)
		if err2 != nil {
			return nil, err2
		}

		amount, err2 := iotago.DecodeUint256(nativeToken.Amount)
		if err2 != nil {
			return nil, err2
		}

		assets.NativeTokens = append(assets.NativeTokens, &iotago.NativeToken{
			ID:     nativeTokenID,
			Amount: amount,
		})
	}

	return assets, nil
}
