package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	mc "github.com/radchukd/multicall-go/src"
)

func main() {
	calls := []mc.Call{
		{
			Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			Name:    "ownerOf(uint256)",
			Params:  []abi.Type{mc.Uint256},
			Values:  []any{big.NewInt(6834)},
			Output:  mc.Address,
		},
		{
			Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			Name:    "name()",
			Output:  mc.String,
		},
		{
			Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			Name:    "totalSupply()",
			Output:  mc.Uint256,
		},
		{
			Address: common.HexToAddress("0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"),
			Name:    "balanceOf(address)",
			Params:  []abi.Type{mc.Address},
			Values:  []any{common.HexToAddress("0xf4893542E4ec7C33356579F91bF22E8FA7CD06dc")},
			Output:  mc.Uint256,
		},
	}

	ctx := context.Background()
	rpcURL := os.Getenv("RPC_URL")
	client, _ := ethclient.DialContext(ctx, rpcURL)
	mcAddressMainnet := "0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696"

	mc.Init(ctx, client, mcAddressMainnet)
	blockNumber, err := mc.MultiCall(&calls)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Results for block #%s:\n", blockNumber.String())

	for _, call := range calls {
		fmt.Println(call)
	}
}
