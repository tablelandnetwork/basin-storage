// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console2} from "forge-std/Script.sol";

import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageScript is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        BasinStorage basinStorage = new BasinStorage();
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));

        // todo grant INDEXER_ROLE to other accounts
        console2.log("deployer addr:", address(this));

        vm.stopBroadcast();
    }
}
