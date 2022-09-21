package multicall

import (
	"context"
	"log"
	"math/big"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client           *ethclient.Client
	ctx              context.Context
	multicallAddress common.Address
)

type Call struct {
	Address common.Address
	Name    string
	Params  []abi.Type
	Values  []any
	Output  abi.Type
	ParseFN func([]byte) any
	Result  any
}

type packedCall struct {
	Target   common.Address
	CallData []byte
}

type callResult struct {
	Success bool
	Data    []byte
}

type multicallResult struct {
	BlockNumber *big.Int
	BlockHash   common.Hash
	ReturnData  []callResult
}

func packCalls(calls []Call) ([]packedCall, error) {
	var multiCalls []packedCall

	for _, call := range calls {
		data := crypto.Keccak256([]byte(call.Name))[:4]

		var params abi.Arguments

		for _, p := range call.Params {
			params = append(params, abi.Argument{Type: p})
		}

		if len(call.Values) > 0 {
			packedArgs, err := params.Pack(call.Values...)

			if err != nil {
				log.Println(err)
				return nil, err
			}

			data = append(data, packedArgs...)
		}

		multiCalls = append(multiCalls, packedCall{
			Target:   call.Address,
			CallData: data,
		})
	}

	return multiCalls, nil
}

func decodeResponse(res []byte) (*multicallResult, error) {
	args := abi.Arguments{
		{
			Name: "BlockNumber",
			Type: Uint256,
		},
		{
			Name: "BlockHash",
			Type: Bytes32,
		},
		{
			Name: "Returns",
			Type: OutputArr,
		},
	}

	data, err := args.Unpack(res)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	decoded := &multicallResult{}

	decoded.BlockNumber = data[0].(*big.Int)

	blockHash := data[1].([32]byte)
	copy(decoded.BlockHash[:], blockHash[:])

	returnData := reflect.ValueOf(data[2])
	for i := 0; i < returnData.Len(); i++ {
		elem := returnData.Index(i)
		res := callResult{
			Success: elem.FieldByName("Success").Bool(),
			Data:    elem.FieldByName("Data").Bytes(),
		}
		decoded.ReturnData = append(decoded.ReturnData, res)
	}

	return decoded, nil
}

func parseResponse(res multicallResult, calls *[]Call) {
	for retIndex, retData := range res.ReturnData {
		if !retData.Success {
			continue
		}

		callOutputType := (*calls)[retIndex].Output.String()
		callParseFN := (*calls)[retIndex].ParseFN
		callRes := &((*calls)[retIndex].Result)

		if !reflect.ValueOf(callParseFN).IsNil() {
			*callRes = callParseFN(retData.Data)
			continue
		}

		switch callOutputType {
		case Bool.String():
			val, _ := strconv.ParseBool(string(retData.Data))
			*callRes = val
		case Int8.String(), Int16.String(), Int32.String(),
			Int64.String(), Int128.String(), Int256.String():
			val, _ := strconv.ParseInt(common.Bytes2Hex(retData.Data), 16, 0)
			*callRes = val
		case Uint8.String(), Uint16.String(), Uint32.String(),
			Uint64.String(), Uint128.String(), Uint256.String():
			val, _ := strconv.ParseUint(common.Bytes2Hex(retData.Data), 16, 0)
			*callRes = val
		case Address.String():
			val := common.BytesToAddress(retData.Data)
			*callRes = val
		case String.String():
			val := string(retData.Data)
			*callRes = val
		case Bytes.String():
			*callRes = retData.Data
		}
	}

}

func Init(c context.Context, cl *ethclient.Client, addr string) {
	ctx = c
	client = cl
	multicallAddress = common.HexToAddress(addr)
}

func MultiCall(calls *[]Call) (*big.Int, error) {
	packedCalls, err := packCalls(*calls)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// blockAndAggregate function signature
	data := common.Hex2Bytes("c3077fa9")
	args := abi.Arguments{{Type: InputArr, Name: "calls"}}
	aggregateArgs, _ := args.Pack(packedCalls)
	data = append(data, aggregateArgs...)

	res, err := client.CallContract(
		ctx,
		ethereum.CallMsg{
			To:   &multicallAddress,
			Data: data,
		},
		nil,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	decodedResponse, err := decodeResponse(res)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	parseResponse(*decodedResponse, calls)

	return decodedResponse.BlockNumber, nil
}
