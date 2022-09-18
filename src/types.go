package multicall

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	Bool, _     = abi.NewType("bool", "", nil)
	Int8, _     = abi.NewType("int8", "", nil)
	Int16, _    = abi.NewType("int16", "", nil)
	Int32, _    = abi.NewType("int32", "", nil)
	Int64, _    = abi.NewType("int64", "", nil)
	Int128, _   = abi.NewType("int128", "", nil)
	Int256, _   = abi.NewType("int256", "", nil)
	Uint8, _    = abi.NewType("uint8", "", nil)
	Uint16, _   = abi.NewType("uint16", "", nil)
	Uint32, _   = abi.NewType("uint32", "", nil)
	Uint64, _   = abi.NewType("uint64", "", nil)
	Uint128, _  = abi.NewType("uint128", "", nil)
	Uint256, _  = abi.NewType("uint256", "", nil)
	Address, _  = abi.NewType("address", "", nil)
	String, _   = abi.NewType("string", "", nil)
	Bytes, _    = abi.NewType("bytes", "", nil)
	Bytes32, _  = abi.NewType("bytes32", "", nil)
	InputArr, _ = abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Type: "address", Name: "Target"},
		{Type: "bytes", Name: "CallData"},
	})
	OutputArr, _ = abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "Success", Type: "bool"},
		{Name: "Data", Type: "bytes"},
	})
)
