package main

import (
	"context"
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	bitcointypes "github.com/goatnetwork/goat/x/bitcoin/types"
	relayertypes "github.com/goatnetwork/goat/x/relayer/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	grpcConn, err := grpc.NewClient(
		"localhost:9090", // Replace with the actual gRPC server address
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close()

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	clientCtx := client.Context{}.
		WithGRPCClient(grpcConn).
		WithCodec(marshaler).
		WithChainID("goat-testnet-1") // Replace with the actual chain ID

	// Create MsgBlockHeader message
	msg := &bitcointypes.MsgNewBlockHashes{
		Proposer:         "goat1...",
		StartBlockNumber: 100,
		BlockHash: [][]byte{
			[]byte("hash1"),
			[]byte("hash2"),
		},
		Vote: &relayertypes.Votes{
			Signature: []byte("signature"),
		},
	}

	// Create transaction
	txFactory := tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithGas(200000).
		WithFees("1000ugoat").
		WithKeybase(clientCtx.Keyring).
		WithAccountRetriever(clientCtx.AccountRetriever)

	// Sign transaction
	txBuilder, err := txFactory.BuildUnsignedTx(msg)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v", err)
	}

	err = tx.Sign(context.Background(), txFactory, "your_key_name", txBuilder, true) // Replace with your key name
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Broadcast transaction
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Fatalf("Failed to encode transaction: %v", err)
	}

	// Use TxClient to broadcast transaction
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}

	log.Printf("Transaction successfully broadcasted: %s", res.TxHash)
}
