package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	DocTransactions string = "bookkeeping/transactions"
	DocBlocks       string = "bookkeeping/blocks"
)

type BlockBook struct {
	Height string `json:"height"`
}

type TxBook struct {
	TxHash string `json:"txhash"`
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

type Txs struct {
	TotalCount string `json:"total_count"`
	Count      string `json:"count"`
	PageNumber string `json:"page_number"`
	PageTotal  string `json:"page_total"`
	Limit      string `json:"limit"`
	Txs        []struct {
		Height string `json:"height"`
		Txhash string `json:"txhash"`
		RawLog string `json:"raw_log"`
		Logs   []struct {
			MsgIndex int    `json:"msg_index"`
			Success  bool   `json:"success"`
			Log      string `json:"log"`
		} `json:"logs"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
		Tx struct {
			Type  string `json:"type"`
			Value struct {
				Msg []struct {
					Type  string `json:"type"`
					Value struct {
						FromAddress string `json:"from_address"`
						ToAddress   string `json:"to_address"`
						Amount      []struct {
							Denom  string `json:"denom"`
							Amount string `json:"amount"`
						} `json:"amount"`
					} `json:"value"`
				} `json:"msg"`
				Fee struct {
					Amount []interface{} `json:"amount"`
					Gas    string        `json:"gas"`
				} `json:"fee"`
				Signatures []struct {
					PubKey struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					Signature string `json:"signature"`
				} `json:"signatures"`
				Memo string `json:"memo"`
			} `json:"value"`
		} `json:"tx"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"txs"`
}

func processTxs(app *firebase.App, height int) {
	res, err := http.Get(fmt.Sprintf("http://localhost:1317/txs?tx.height=%d", height))
	if err != nil {
		log.Fatal(err)
	}
	jsonBlob, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var txs Txs
	err = json.Unmarshal(jsonBlob, &txs)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("num transactions: %d\n", len(txs.Txs))
	for _, tx := range txs.Txs {
		fmt.Printf("transaction hash: %s\n", tx.Txhash)
		updateTransactions(app, tx.Txhash)
	}
}

func processSequentially(app *firebase.App) error {
	processed, err := processedSofar(app)
	if err != nil {
		return err
	}
	var i int
	for {
		if latestHeight() < processed {
			time.Sleep(2 * time.Second)
			log.Printf("sleeping")
		}
		i++
		processTxs(app, processed)
		if i%100 == 0 {
			if err = updateBookkeeping(strconv.Itoa(processed), app); err != nil {
				log.Printf("Some error: %v", err)
			}
		}
		processed++
	}
	return nil
}

func updateTransactions(app *firebase.App, txhash string) error {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	bd := client.Doc(DocTransactions)

	// Check if document does not exist and create with new value
	if _, err = bd.Get(ctx); err != nil && grpc.Code(err) == codes.NotFound {
		log.Printf("creating doc from scratch")
		_, err := bd.Create(ctx, &TxBook{txhash})
		return err
	}

	_, err = bd.Set(ctx, &TxBook{txhash})
	return err
}

func updateBookkeeping(height string, app *firebase.App) error {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	bd := client.Doc(DocBlocks)

	// Check if document does not exist and create with new value
	if _, err = bd.Get(ctx); err != nil && grpc.Code(err) == codes.NotFound {
		log.Printf("creating doc from scratch")
		_, err := bd.Create(ctx, &BlockBook{"0"})
		return err
	}

	_, err = bd.Set(ctx, &BlockBook{height})
	return err
}

func processedSofar(app *firebase.App) (int, error) {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	bd := client.Doc(DocBlocks)
	docsnap, err := bd.Get(ctx)
	if err != nil {
		return 0, err
	}
	var myData BlockBook
	if err := docsnap.DataTo(&myData); err != nil {
		return 0, err
	}
	return strconv.Atoi(myData.Height)
}

func latestHeight() int {
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
		log.Fatal(err)
	}
	height := block.BlockMeta.Header.Height
	fmt.Printf("latest-block-height: %s\n", height)
	intHeight, err := strconv.Atoi(height)
	if err != nil {
		log.Fatal(err)
	}
	return intHeight
}

func main() {

	opt := option.WithCredentialsFile("worldcoin-dev-firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleting document %s to restart from scratch", DocBlocks)
	if _, err = client.Doc(DocBlocks).Delete(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Recreating document %s with height set to 0", DocBlocks)
	if err = updateBookkeeping("0", app); err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleting document %s to restart from scratch", DocTransactions)
	if _, err = client.Doc(DocTransactions).Delete(ctx); err != nil {
		log.Fatal(err)
	}
	if err = processSequentially(app); err != nil {
		log.Fatal(err)
	}
}
