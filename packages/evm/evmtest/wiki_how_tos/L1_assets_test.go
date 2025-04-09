package wiki_how_tos_test

import (
	_ "embed"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	//"github.com/iotaledger/wasp/packages/isc"
	//"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"

	//"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
	"github.com/stretchr/testify/require"
)

//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` L1Assets.sol -o ."
var (
	//go:embed L1Assets.abi
	L1AssetsContractABI string
	//go:embed L1Assets.bin
	L1AssetsContractBytecodeHex string
	L1AssetsContractBytecode    = common.FromHex(strings.TrimSpace(L1AssetsContractBytecodeHex))
)

func TestWithdraw(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	_, receiver := env.Chain.Env.NewKeyPair()

	// Deploy L1Assets contract
	instance := env.DeployContract(privateKey, L1AssetsContractABI, L1AssetsContractBytecode)

	require.Zero(t, env.Chain.Env.L1BaseTokens(receiver))
	// const storageDeposit uint64 = 10_000

	contractAgentID := isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)
	env.Chain.GetL2FundsFromFaucet(contractAgentID, 2000)

	{
		const baseTokensDepositFee = 500
		k, _ := env.Chain.Env.NewKeyPairWithFunds(env.Chain.Env.NewSeedFromIndex(1))
		//k, _ := env.solo.NewKeyPairWithFunds(env.solo.NewSeedFromIndex(1))
		err := env.Chain.SendFromL1ToL2AccountBaseTokens(baseTokensDepositFee, 1*isc.Million, contractAgentID, k)
		require.NoError(t, err)
		require.EqualValues(t, 1*isc.Million, env.Chain.L2BaseTokens(contractAgentID))
	}

	// create a new native token on L1
	foundry, tokenID, err := env.Chain.NewNativeTokenParams(100000000000000).CreateFoundry()
	require.NoError(t, err)
	// the token id in bytes, used to call the contract
	nativeTokenIDBytes := isc.NativeTokenIDToBytes(tokenID)

	nativeTokenID := iscmagic.NativeTokenID{
		Data: nativeTokenIDBytes,
	}

	// mint some native tokens to the chain originator
	err = env.Chain.MintTokens(foundry, 10000000, env.Chain.OriginatorPrivateKey)
	require.NoError(t, err)

	// Create ISCAssets with native tokens
	amount := big.NewInt(500)
	assets := iscmagic.ISCAssets{
		NativeTokens: []iscmagic.NativeToken{{ID: nativeTokenID, Amount: amount}},
	}

	// Allow the L1Assets contract to withdraw the funds
	_, err = instance.CallFn(nil, "allow", deployer, assets)
	require.NoError(t, err)

	// Withdraw funds to receiver using the withdraw function of L1Assets contract
	_, err = instance.CallFn(nil, "withdraw", iscmagic.WrapL1Address(receiver))
	require.NoError(t, err)
	//require.GreaterOrEqual(t, env.Chain.Env.L1BaseTokens(receiver))

	// Verify balances
	//require.LessOrEqual(t, env.Chain.L2BaseTokens(isc.NewEthereumAddressAgentID(env.Chain.ChainID, deployer)))
}
