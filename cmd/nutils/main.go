package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	app "github.com/cosmos/sdk-application-tutorial"
	amino "github.com/tendermint/go-amino"
)

type (
	// DecodeReq defines a tx decoding request.
	DecodeReq struct {
		Tx string `json:"tx"`
	}

	// DecodeResp defines a tx decoding response.
	DecodeResp authtypes.StdTx
)

// txDecodeRespStr implements a simple Stringer wrapper for a decoded tx.
type txDecodeRespTx authtypes.StdTx

func (tx txDecodeRespTx) String() string {
	return tx.String()
}

// GetDecodeCommand returns the decode command to take Amino-serialized bytes and turn it into
// a JSONified transaction
// Use:   "decode [amino-byte-string]",
// Short: "Decode an amino-encoded transaction string",

func myMain(codec *amino.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(codec)
	cliCtx.OutputFormat = "json"

	txBytesBase64 := os.Args[1]
	fmt.Println(txBytesBase64)

	txBytes, err := base64.StdEncoding.DecodeString(txBytesBase64)
	if err != nil {
		return err
	}

	var stdTx authtypes.StdTx
	err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(txBytes, &stdTx)
	if err != nil {
		return err
	}

	response := txDecodeRespTx(stdTx)
	// fmt.Printf("resp: %v\n", response)
	_ = cliCtx.PrintOutput(response)

	return nil
}

func main() {
	cdc := app.MakeCodec()
	err := myMain(cdc)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println("All good")

}
