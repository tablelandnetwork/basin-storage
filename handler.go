package handler

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	"github.com/textileio/go-tableland/pkg/wallet"

	w3s "github.com/web3-storage/go-w3s-client"

	ethereum "github.com/tablelandnetwork/basin-storage/pkg/ethereum"
	bstorage "github.com/tablelandnetwork/basin-storage/pkg/storage"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("Uploader", Uploader)
	functions.HTTP("StatusChecker", StatusChecker)
}

// Uploader is the CloudEvent function that is called by the Functions Framework.
// It is triggered by a CloudEvent that is published by the GCS bucket.
// The CloudEvent contains the name of the bucket and the name of the file.
// The file is downloaded from GCS and uploaded to web3.storage.
func Uploader(ctx context.Context, e event.Event) error {
	// Set a timeout of 60 minutes, thats the max time a function can run on GCP (gen2)
	// we want to ensure larger files can be uploaded
	cctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	// Read w3s token and db conn string from environment variables
	web3StorageToken := os.Getenv("WEB3STORAGE_TOKEN")
	crdbConnStr := os.Getenv("CRDB_CONN_STRING")

	// Initialize GCS client to download file
	// bucket name and file name are passed in the CloudEvent
	storageClient, err := bstorage.NewGCSClient(ctx, e.Data())
	if err != nil {
		return fmt.Errorf("failed to initialize storage client: %v", err)
	}

	// Initialize web3.storage client to upload file
	w3sOpts := []w3s.Option{
		w3s.WithToken(web3StorageToken),
		w3s.WithHTTPClient(
			&http.Client{
				Timeout: 0, // no timeout
			},
		),
	}
	w3sClient, err := w3s.NewClient(w3sOpts...)
	if err != nil {
		return fmt.Errorf("failed to initialize web3.storage client: %v", err)
	}

	// Initialize cockroachdb client to store metadata
	dbClient, err := bstorage.NewDB(crdbConnStr)
	if err != nil {
		return fmt.Errorf("failed to initialize cockroachdb client: %v", err)
	}
	defer dbClient.DB.Close()

	u := &bstorage.FileUploader{
		StorageClient: storageClient,
		DealClient:    w3sClient,
		DBClient:      dbClient,
	}
	err = u.Upload(cctx)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

func AddDeals(ctx context.Context,
	backend *ethclient.Client,
	pub string,
	deals []ethereum.BasinStorageDealInfo,
) error {
	contractAddrStr := "0x4b1f4d8100e51afe644b189d77784dec225e0596"
	addr, err := common.NewMixedcaseAddressFromString(contractAddrStr)
	if err != nil {
		panic(err)
	}
	commonAddr := addr.Address()
	contract, err := ethereum.NewContract(commonAddr, backend)
	if err != nil {
		panic(err)
	}

	gasTipCap, err := backend.SuggestGasTipCap(ctx)
	if err != nil {
		panic(err)
	}

	wallet, err := wallet.NewWallet("f4212951e2639746613d34dd912713c25df1be3d08e87680d4c9f147c81ed8b6")
	if err != nil {
		panic(err)
	}

	chainId := 314159
	auth, err := bind.NewKeyedTransactorWithChainID(
		wallet.PrivateKey(),
		big.NewInt(int64(chainId)),
	)
	if err != nil {
		panic(err)
	}

	/* 	BasinStorageABI, err := abi.JSON(strings.NewReader(ethereum.ContractMetaData.ABI))
	   	if err != nil {
	   		panic(err)
	   	}
	*/
	/* deal1 := ethereum.BasinStorageDealInfo{
		Id:           3,
		SelectorPath: "/ipf3",
	}
	deal2 := ethereum.BasinStorageDealInfo{
		Id:           3,
		SelectorPath: "/ipf4",
	}
	deals := []ethereum.BasinStorageDealInfo{deal1, deal2}

	*/ /* data, err := BasinStorageABI.Pack("addDeals", []interface{}{"foobar", deals}...)
	fmt.Println(data)
	if err != nil {
		panic(err)
	}
	*/
	/* gasLimit, err := backend.EstimateGas(ctx, eth.CallMsg{
		From: auth.From,
		To:   &commonAddr,
		Data: data,
	})
	if err != nil {
		panic(err)
	}
	*/
	/* 	nonce, err := backend.PendingNonceAt(ctx, wallet.Address())
	   	if err != nil {
	   		return fmt.Errorf("getting nonce: %s", err)
	   	} */

	fmt.Println("from address", auth.From)
	opts := &bind.TransactOpts{
		Context: ctx,
		Signer:  auth.Signer,
		From:    auth.From,
		//Nonce:     big.NewInt(0).SetInt64(nonce),
		GasTipCap: gasTipCap.Mul(gasTipCap, big.NewInt(500)),
		GasLimit:  27_719_768, //gasLimit * 4,
	}

	pubsOfOwner, err := contract.PubsOfOwner(
		&bind.CallOpts{
			Pending: true,
			Context: ctx,
		}, commonAddr)
	if err != nil {
		panic(err)
	}

	ownerByPub, err := contract.GetOwnerByPub(&bind.CallOpts{
		Pending: true,
		Context: ctx,
	}, pub)
	if err != nil {
		panic(err)
	}

	fmt.Println("owner by pub", ownerByPub)

	latestDeals, err := contract.LatestNDeals(&bind.CallOpts{
		Pending: true,
		Context: ctx,
	}, "foobar", big.NewInt(2))
	if err != nil {
		panic(err)
	}
	for _, d := range latestDeals {
		fmt.Println("latest N deals", d.Id, d.SelectorPath)
	}

	fmt.Println("pubs of owner", pubsOfOwner)

	tx, err := contract.AddDeals(opts, pub, deals) // CreateTable(opts, owner, statement)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tx sent: %v \n", tx)

	time.Sleep(120 * time.Second)

	// get receipt
	receipt, err := backend.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v \n", receipt)

	for _, l := range receipt.Logs {
		fmt.Println(l.Data)
	}

	// fmt.Println(tx.)
	if err != nil {
		panic(err)
	}

	return nil
}

// StatusChecker is the HTTP function that is called by the Functions Framework.
func StatusChecker(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	backend, err := ethclient.DialContext(ctx, "https://api.calibration.node.glif.io/rpc/v1")
	if err != nil {
		panic(err)
	}
	fmt.Println("prepared provider", backend)
	pub := "foobar" //fmt.Sprintf("0x%x", []byte("foobar"))

	// Read w3s token and db conn string from environment variables
	web3StorageToken := os.Getenv("WEB3STORAGE_TOKEN")
	crdbConnStr := os.Getenv("CRDB_CONN_STRING")

	// Initialize web3.storage client to upload file
	w3sOpts := []w3s.Option{
		w3s.WithToken(web3StorageToken),
		w3s.WithHTTPClient(
			&http.Client{
				Timeout: 0, // no timeout
			},
		),
	}
	w3sClient, err := w3s.NewClient(w3sOpts...)
	if err != nil {
		panic(err)
	}

	// Initialize cockroachdb client to store metadata
	dbClient, err := bstorage.NewDB(crdbConnStr)
	if err != nil {
		panic(err)
	}
	defer dbClient.DB.Close()

	unfinihedJobs, err := dbClient.UnfinishedJobs(r.Context())
	if err != nil {
		panic(err)
	}

	for _, job := range unfinihedJobs {
		fmt.Println("====> job", job)
		jobCid, err := cid.Parse(job)
		if err != nil {
			panic(err)
		}
		status, err := w3sClient.Status(r.Context(), jobCid)
		if err != nil {
			panic(err)
		}
		if status != nil && len(status.Deals) > 0 {
			// 1. send transaction to DataStorage contract
			// -- both should be done in a transaction, if FVM update fails, job status should not be updated
			deals := []ethereum.BasinStorageDealInfo{}
			for _, d := range status.Deals {
				fmt.Println("Activation ===", d.Activation)
				deals = append(deals, ethereum.BasinStorageDealInfo{
					Id:           d.DealID,
					SelectorPath: d.DataModelSelector,
				})
			}

			// TODO before adding deals, check if deals already exist
			// by calling contract.LatestNDeals() and comparing with deal ids
				

			fmt.Println("deals ready to be prepared", deals)

			err := AddDeals(ctx, backend, pub, deals)
			if err != nil {
				panic(err)
			}

			// contract.AddDeals()

			// TODO
			// -- check revert error from contract is "deal already exists" then update job status too
			// -- if not, then raise errors

			// 2. update job status in db
			/* dbClient.UpdateJobStatus(
				r.Context(), jobCid.String(), status.Deals[0].Activation,
			) */
		}

	}

}
