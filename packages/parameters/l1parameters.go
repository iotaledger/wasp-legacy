package parameters

import (
	"os"
	"strings"
	"sync"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/tpkg"
)

// L1Params describes parameters coming from the L1Params node
type L1Params struct {
	MaxPayloadSize int                        `json:"maxPayloadSize" swagger:"required"`
	Protocol       *iotago.ProtocolParameters `json:"protocol" swagger:"required"`
	BaseToken      *BaseToken                 `json:"baseToken" swagger:"required"`
}

type BaseToken struct {
	Name            string `json:"name" swagger:"desc(The base token name),required"`
	TickerSymbol    string `json:"tickerSymbol" swagger:"desc(The ticker symbol),required"`
	Unit            string `json:"unit" swagger:"desc(The token unit),required"`
	Subunit         string `json:"subunit" swagger:"desc(The token subunit),required"`
	Decimals        uint32 `json:"decimals" swagger:"desc(The token decimals),required"`
	UseMetricPrefix bool   `json:"useMetricPrefix" swagger:"desc(Whether or not the token uses a metric prefix),required"`
}

const MaxPayloadSize = iotago.BlockBinSerializedMaxSize - // BlockSizeMax
	serializer.OneByte - // ProtocolVersion
	serializer.OneByte - // ParentCount
	(iotago.BlockMaxParents * iotago.BlockIDLength) - // Parents
	serializer.UInt32ByteSize - // PayloadLength
	serializer.UInt64ByteSize // Nonce

var (
	l1ParamsMutex = &sync.RWMutex{}
	l1Params      *L1Params

	L1ForTesting = &L1Params{
		// There are no limits on how big from a size perspective an essence can be, so it is just derived from 32KB - Message fields without payload = max size of the payload
		MaxPayloadSize: MaxPayloadSize,
		Protocol: &iotago.ProtocolParameters{
			Version:     tpkg.TestProtoParas.Version,
			NetworkName: tpkg.TestProtoParas.NetworkName,
			Bech32HRP:   tpkg.TestProtoParas.Bech32HRP,
			MinPoWScore: tpkg.TestProtoParas.MinPoWScore,
			RentStructure: iotago.RentStructure{
				VByteCost:    10,
				VBFactorData: 1,
				VBFactorKey:  1,
			},
			TokenSupply: tpkg.TestProtoParas.TokenSupply,
		},
		BaseToken: &BaseToken{
			Name:            "Iota",
			TickerSymbol:    "MIOTA",
			Unit:            "MIOTA",
			Subunit:         "IOTA",
			Decimals:        6,
			UseMetricPrefix: false,
		},
	}

	L1ForHack = &L1Params{
		// There are no limits on how big from a size perspective an essence can be, so it is just derived from 32KB - Message fields without payload = max size of the payload
		MaxPayloadSize: MaxPayloadSize,
		Protocol: &iotago.ProtocolParameters{
			Version:     2,
			NetworkName: "iota",
			Bech32HRP:   "iota",
			MinPoWScore: 0,
			RentStructure: iotago.RentStructure{
				VByteCost:    500,
				VBFactorData: 1,
				VBFactorKey:  10,
			},
			TokenSupply: 2779530283277761,
		},
		BaseToken: &BaseToken{
			Name:            "Iota",
			TickerSymbol:    "MIOTA",
			Unit:            "MIOTA",
			Subunit:         "IOTA",
			Decimals:        6,
			UseMetricPrefix: false,
		},
	}

	l1ParamsLazyInit func()
)

func isTestContext() bool {
	return strings.HasSuffix(os.Args[0], ".test") ||
		strings.HasSuffix(os.Args[0], ".test.exe") ||
		strings.Contains(os.Args[0], "__debug_bin")
}

func L1() *L1Params {
	l1ParamsMutex.Lock()
	defer l1ParamsMutex.Unlock()
	return L1NoLock()
}

func L1NoLock() *L1Params {
	if l1Params == nil {
		if isTestContext() {
			l1Params = L1ForTesting
		} else if l1ParamsLazyInit != nil {
			l1ParamsLazyInit()
		}
	}
	if l1Params == nil {
		panic("call InitL1() first")
	}
	return l1Params
}

// InitL1Lazy sets a function to be called the first time L1() is called.
// The function must call InitL1().
func InitL1Lazy(f func()) {
	l1ParamsLazyInit = f
}

func InitL1(l1 *L1Params) {
	l1Params = l1
}
