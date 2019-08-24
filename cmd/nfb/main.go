package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type BlockBook struct {
	Height string `json:"height"`
}

type Block struct {
	BlockMeta struct {
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total string `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
		Header struct {
			Version struct {
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"version"`
			ChainID     string    `json:"chain_id"`
			Height      string    `json:"height"`
			Time        time.Time `json:"time"`
			NumTxs      string    `json:"num_txs"`
			TotalTxs    string    `json:"total_txs"`
			LastBlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total string `json:"total"`
					Hash  string `json:"hash"`
				} `json:"parts"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			NextValidatorsHash string `json:"next_validators_hash"`
			ConsensusHash      string `json:"consensus_hash"`
			AppHash            string `json:"app_hash"`
			LastResultsHash    string `json:"last_results_hash"`
			EvidenceHash       string `json:"evidence_hash"`
			ProposerAddress    string `json:"proposer_address"`
		} `json:"header"`
	} `json:"block_meta"`
	Block struct {
		Header struct {
			Version struct {
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"version"`
			ChainID     string    `json:"chain_id"`
			Height      string    `json:"height"`
			Time        time.Time `json:"time"`
			NumTxs      string    `json:"num_txs"`
			TotalTxs    string    `json:"total_txs"`
			LastBlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total string `json:"total"`
					Hash  string `json:"hash"`
				} `json:"parts"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			NextValidatorsHash string `json:"next_validators_hash"`
			ConsensusHash      string `json:"consensus_hash"`
			AppHash            string `json:"app_hash"`
			LastResultsHash    string `json:"last_results_hash"`
			EvidenceHash       string `json:"evidence_hash"`
			ProposerAddress    string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs interface{} `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence interface{} `json:"evidence"`
		} `json:"evidence"`
		LastCommit struct {
			BlockID struct {
				Hash  string `json:"hash"`
				Parts struct {
					Total string `json:"total"`
					Hash  string `json:"hash"`
				} `json:"parts"`
			} `json:"block_id"`
			Precommits []struct {
				Type    int    `json:"type"`
				Height  string `json:"height"`
				Round   string `json:"round"`
				BlockID struct {
					Hash  string `json:"hash"`
					Parts struct {
						Total string `json:"total"`
						Hash  string `json:"hash"`
					} `json:"parts"`
				} `json:"block_id"`
				Timestamp        time.Time `json:"timestamp"`
				ValidatorAddress string    `json:"validator_address"`
				ValidatorIndex   string    `json:"validator_index"`
				Signature        string    `json:"signature"`
			} `json:"precommits"`
		} `json:"last_commit"`
	} `json:"block"`
}

func main() {

	res, err := http.Get("http://localhost:1317/blocks/latest")
	if err != nil {
		log.Fatal(err)
	}
	jsonBlob, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var block Block
	err = json.Unmarshal(jsonBlob, &block)
	if err != nil {
		fmt.Println("error:", err)
	}
	height := block.BlockMeta.Header.Height
	fmt.Printf("latest-block-height: %s\n", height)

	opt := option.WithCredentialsFile("worldcoin-dev-firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
		return
	}
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	bd := client.Doc("bookkeeping/blocks")
	docsnap, err := bd.Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {

			log.Printf("Doc does not exist!\n")
			_, err := bd.Create(ctx, map[string]string{"height": height})
			if err != nil {
				// TODO: Handle error.
				log.Fatal(err)
			}
			docsnap, err = bd.Get(ctx)
			if err != nil {
				log.Fatal(err)
			}

		}
		log.Printf("err: %v\n", err)
	}
	var myData BlockBook
	if err := docsnap.DataTo(&myData); err != nil {
		// TODO: Handle error.
		log.Fatal(err)
	}
	log.Printf("docsnap: %v\n", myData)
	// _, _, err = client.Collection("Bookkeeping").Add(ctx, map[string]interface{"blocks": "10"})
	// if err != nil {
	// 	log.Fatalf("Failed adding block: %v", err)
	// }
	// _, _, err = client.Collection("blocks").Add(ctx, block)
	// if err != nil {
	// 	log.Fatalf("Failed adding block: %v", err)
	// }
	fmt.Printf("All good\n")

}
