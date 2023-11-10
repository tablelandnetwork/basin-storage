# Basin Storage Contract

**Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.**

Foundry consists of:

-   **Forge**: Ethereum testing framework (like Truffle, Hardhat and DappTools).
-   **Cast**: Swiss army knife for interacting with EVM smart contracts, sending transactions and getting chain data.
-   **Anvil**: Local Ethereum node, akin to Ganache, Hardhat Network.
-   **Chisel**: Fast, utilitarian, and verbose solidity REPL.

## Documentation

https://book.getfoundry.sh/

## Usage

### Build

```shell
$ forge build
```

### Test

```shell
$ forge test
```

### Format

```shell
$ forge fmt
```

### Gas Snapshots

```shell
$ forge snapshot
```

### Anvil

```shell
$ anvil
```

### Deploy

```shell
$ forge script script/Counter.s.sol:CounterScript --rpc-url <your_rpc_url> --private-key <your_private_key>
```

### Cast commands

#### Add CID

Add a CID to the contract

```shell
cast send \
--private-key <your private key> \
--rpc-url https://api.calibration.node.glif.io/rpc/v1 \
<contract adddress> \
"addCID(string,string,uint256)" \
"pub.name" \
"bafyexamplecid" \
1699615822
```

#### Get pubs of owner

Get pubs of an owner

```shell
cast call <contract address> \
--rpc-url "https://api.calibration.node.glif.io/rpc/v1" \
"pubsOfOwner(address)(string[])" \
<data owner's address>
```

#### Get cid at timestamp

Get cid at timestamp (change the timestamp accordingly)

```shell
cast call <contract address> \
--rpc-url "https://api.calibration.node.glif.io/rpc/v1" \
"cidsAtTimestamp(string,uint256)(string[])" \
"pub.name" \
1699615822
```

#### Get cid in time range

Get cid in time range (change timestamps accordingly)

```shell
cast call <contract address> \
--rpc-url "https://api.calibration.node.glif.io/rpc/v1" \
"cidsInRange(string,uint256,uint256)(string[])" \
"testavichalp.data" 1699615821 1699615823
```
