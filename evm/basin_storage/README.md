# Basin Storage Contracts

This project is built using foundry [framework](https://book.getfoundry.sh/)

## Usage

If the commands are run from project root pass the `--root` flag to `forge` to tell it where the contracts and libs are located.

### Compile contracts

```shell
forge compile --root evm/basin_storage
```

### Test

```shell
forge test --root evm/basin_storage -vvv --gas-report
```

### Gas Snapshots

```shell
forge snapshot --root evm/basin-storage
```

### Deploying Basin storage

```shell
PRIVATE_KEY=<your private key> forge script evm/basin_storage/script/BasinStorage.s.sol:BasinStorageScript --root evm/basin_storage --broadcast --rpc-url <rpc url> --skip-simulation --gas-estimate-multiplier 5000 --retries 10
```

When deploying to Filecoin Calibration use the following RPC url:
`https://api.calibration.node.glif.io/rpc/v1`

### Using cast

Create a Pub

```sh
cast send --private-key <your private key> --rpc-url <rpc url> <contract addr> "createPub(address,string)" "<owner adder>" "<nsname>.<relname>"
```

Fetch pubs of an owner

```sh
 cast call <contract addr> --rpc-url <rpc url> "pubsOfOwner(address)(string[])" <owner addr>
```
