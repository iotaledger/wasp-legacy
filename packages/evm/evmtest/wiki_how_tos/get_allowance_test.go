package wiki_how_tos_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
)

// compile the solidity contract
//go:generate sh -c "solc --abi --bin --overwrite @iscmagic=`realpath ../../../vm/core/evm/iscmagic` GetAllowance.sol -o ."

var (
	//go:embed GetAllowance.abi
	GetAllowanceContractABI string
	//go:embed GetAllowance.bin
	GetAllowanceContractBytecodeHex string
	GetAllowanceContractBytecode    = common.FromHex(strings.TrimSpace(GetAllowanceContractBytecodeHex))
)

func TestGetAllowance(t *testing.T) {
	env := evmtest.InitEVMWithSolo(t, solo.New(t), true)
	privateKey, deployer := env.Chain.NewEthereumAccountWithL2Funds()

	instance := env.DeployContract(privateKey, GetAllowanceContractABI, GetAllowanceContractBytecode)
	value, err := instance.CallFn(nil, "getAllowanceFrom", deployer)
	assert.Nil(t, err)

	t.Log("value:", value)
}
