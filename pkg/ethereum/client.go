package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/textileio/go-tableland/pkg/wallet"
)

// BasinStorage is an interface that defines the methods to interact with the BasinStorage smart contract.
type BasinStorage interface {
	EstimateGas(ctx context.Context,
		pub string,
		deals []BasinStorageDealInfo) (*bind.TransactOpts, error)
	GetRecentDeals(ctx context.Context,
		pub string) (map[uint64]BasinStorageDealInfo, error)
	AddDeals(ctx context.Context,
		pub string,
		deals []BasinStorageDealInfo, txOpts *bind.TransactOpts) error
}

// Client is the Ethereum implementation of the registry client.
type Client struct {
	contract     *Contract
	contractAddr common.Address
	backend      bind.ContractBackend
	rpcBackend   bind.DeployBackend
	wallet       *wallet.Wallet
	chainID      uint64
}

// NewClient creates a new Client.
func NewClient(
	contractBackend bind.ContractBackend,
	rpcBackend bind.DeployBackend,
	chainID uint64,
	contractAddr common.Address,
	wallet *wallet.Wallet,
) (*Client, error) {
	contract, err := NewContract(contractAddr, contractBackend)
	if err != nil {
		return nil, fmt.Errorf("creating contract: %v", err)
	}
	return &Client{
		contract:     contract,
		contractAddr: contractAddr,
		backend:      contractBackend,
		rpcBackend:   rpcBackend,
		wallet:       wallet,
		chainID:      chainID,
	}, nil
}

// EstimateGas estimates the gas required to execute the AddDeals function of the BasinStorage smart contract.
func (c *Client) EstimateGas(
	ctx context.Context,
	pub string,
	deals []BasinStorageDealInfo,
) (*bind.TransactOpts, error) {
	txOpts, err := bind.NewKeyedTransactorWithChainID(
		c.wallet.PrivateKey(),
		big.NewInt(int64(c.chainID)),
	)
	if err != nil {
		return &bind.TransactOpts{}, fmt.Errorf("failed to initialize tx opts: %v", err)
	}

	gasTipCap, err := c.backend.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed while suggesting gas tip cap: %v", err)
	}

	BasinStorageABI, err := abi.JSON(strings.NewReader(ContractMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	data, err := BasinStorageABI.Pack("addDeals", []interface{}{pub, deals}...)
	if err != nil {
		return nil, fmt.Errorf("failed to abi pack: %v", err)
	}

	gasLimit, err := c.backend.EstimateGas(ctx, eth.CallMsg{
		From: txOpts.From,
		To:   &c.contractAddr,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("error while calling EstimateGas rpc: %v", err)
	}

	return &bind.TransactOpts{
		Context:   ctx,
		Signer:    txOpts.Signer,
		From:      txOpts.From,
		GasTipCap: gasTipCap.Mul(gasTipCap, big.NewInt(500)),
		GasLimit:  gasLimit * 4,
	}, nil
}

// GetRecentDeals returns the latest 10 deals added to the BasinStorage smart contract for the given publisher.
func (c *Client) GetRecentDeals(ctx context.Context, pub string) (map[uint64]BasinStorageDealInfo, error) {
	callOpts := &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	latestDeals, err := c.contract.LatestNDeals(callOpts, pub, big.NewInt(10))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent deals: %v", err)
	}

	// index recent deals in a map
	recentDeals := make(map[uint64]BasinStorageDealInfo)
	for _, d := range latestDeals {
		recentDeals[d.Id] = d
	}

	return recentDeals, nil
}

// AddDeals adds the given deals to the BasinStorage smart contract for the given pub.
func (c *Client) AddDeals(ctx context.Context,
	pub string,
	deals []BasinStorageDealInfo,
	txOpts *bind.TransactOpts,
) error {
	// TODO: implement retry logic
	tx, err := c.contract.AddDeals(txOpts, pub, deals)
	if err != nil {
		return fmt.Errorf("failed to add deals: %v", err)
	}
	fmt.Printf("tx sent: %v \n", tx.Hash())
	time.Sleep(150 * time.Second)
	receipt, err := c.rpcBackend.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get tx receipt: %v", err)
	}
	fmt.Printf("got tx receipt: %v \n", receipt)

	return nil
}
