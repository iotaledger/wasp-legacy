package sbtests

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil/testmisc"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
)

func infiniteLoopRequest(ch *solo.Chain, gasBudget ...uint64) (*solo.CallParams, *cryptolib.KeyPair) {
	budget := uint64(math.MaxUint64)
	if len(gasBudget) > 0 {
		budget = gasBudget[0]
	}
	wallet, address := ch.Env.NewKeyPairWithFunds()
	baseTokensToSend := ch.Env.L1BaseTokens(address)

	req := solo.NewCallParams(ScName, sbtestsc.FuncInfiniteLoop.Name).
		AddBaseTokens(baseTokensToSend).
		WithGasBudget(budget)
	return req, wallet
}

func TestTxWithGasOverLimit(t *testing.T) { run2(t, testTxWithGasOverLimit) }
func testTxWithGasOverLimit(t *testing.T, w bool) {
	// create a TX that would use more than max gas limit, assert that only the maximum will be used
	if w { // TODO the WASM version of this must be tested.
		t.SkipNow()
	}
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil, w)

	req, wallet := infiniteLoopRequest(ch)
	_, err := ch.PostRequestSync(req, wallet)
	require.Error(t, err) // tx expected to run out of gas
	testmisc.RequireErrorToBe(t, err, vm.ErrGasBudgetExceeded)
	receipt := ch.LastReceipt()
	// assert that the submitted gas budget was limited to the max per call
	require.Less(t, receipt.GasBurned, req.GasBudget())
	require.GreaterOrEqual(t, receipt.GasBurned, ch.GetGasLimits().MaxGasPerRequest) // should exceed MaxGasPerRequest() by 1 operation
}

func TestBlockGasOverflow(t *testing.T) { run2(t, testBlockGasOverflow) }
func testBlockGasOverflow(t *testing.T, w bool) {
	// queue many transactions with enough gas to fill a block, assert that they are split across blocks
	if forceGoNoWasm && w {
		t.SkipNow()
	}
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil, w)
	initialBlockInfo := ch.GetLatestBlockInfo()

	// produce 1 request over the block gas limit (each request uses the maximum amount of gas a call can use)
	limits := ch.GetGasLimits()
	nRequests := int(limits.MaxGasPerBlock/limits.MaxGasPerRequest) + 1
	reqs := make([]isc.Request, nRequests)

	for i := 0; i < nRequests; i++ {
		req, wallet := infiniteLoopRequest(ch)
		iscReq, err := solo.ISCRequestFromCallParams(ch, req, wallet)
		require.NoError(t, err)
		reqs[i] = iscReq
	}

	ch.Env.AddRequestsToMempool(ch, reqs)
	ch.WaitUntilMempoolIsEmpty()

	// we should have produced 2 blocks
	require.EqualValues(t, initialBlockInfo.BlockIndex()+2, ch.LatestBlockIndex())

	fullGasBlockInfo, err := ch.GetBlockInfo(initialBlockInfo.BlockIndex() + 1)
	require.NoError(t, err)
	// the request number #{nRequests} should overflow the block and be moved to the next one
	require.Equal(t, nRequests-1, int(fullGasBlockInfo.TotalRequests))
	// gas burned will be slightly below the limit
	require.LessOrEqual(t, fullGasBlockInfo.GasBurned, limits.MaxGasPerBlock)

	// 1 requests should be moved to the next block
	followingBlockInfo, err := ch.GetBlockInfo(initialBlockInfo.BlockIndex() + 2)
	require.NoError(t, err)
	require.Equal(t, uint16(1), followingBlockInfo.TotalRequests)
}

func TestGasBudget(t *testing.T) { run2(t, testGasBudget) }
func testGasBudget(t *testing.T, w bool) {
	// create a TX with not enough gas, assert the receipt is as expected
	if w { // TODO the WASM version of this must be tested.
		t.SkipNow()
	}
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil, w)

	limits := ch.GetGasLimits()
	gasBudget := limits.MaxGasPerRequest / 2
	req, wallet := infiniteLoopRequest(ch, gasBudget)
	_, err := ch.PostRequestSync(req, wallet)
	require.Error(t, err) // tx expected to run out of gas
	testmisc.RequireErrorToBe(t, err, vm.ErrGasBudgetExceeded)
	receipt := ch.LastReceipt()
	require.EqualValues(t, gasBudget, receipt.GasBudget)

	// repeat with gas budget 0
	gasBudget = 0
	req, wallet = infiniteLoopRequest(ch, gasBudget)
	_, err = ch.PostRequestSync(req, wallet)
	require.Error(t, err) // tx expected to run out of gas
	testmisc.RequireErrorToBe(t, err, vm.ErrGasBudgetExceeded)
	receipt = ch.LastReceipt()
	require.EqualValues(t, limits.MinGasPerRequest, receipt.GasBudget) // gas budget should be adjusted to the minimum
}

func TestViewGasLimit(t *testing.T) { run2(t, testViewGasLimit) }
func testViewGasLimit(t *testing.T, w bool) {
	if forceGoNoWasm && w {
		t.SkipNow()
	}
	_, ch := setupChain(t, nil)
	setupTestSandboxSC(t, ch, nil, w)
	_, err := ch.CallView(sbtestsc.Contract.Name, sbtestsc.FuncInfiniteLoopView.Name)
	require.Error(t, err)
	testmisc.RequireErrorToBe(t, err, vm.ErrGasBudgetExceeded)
}
