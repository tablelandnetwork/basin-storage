// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console2} from "forge-std/Script.sol";

import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageScript is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddress = vm.addr(deployerPrivateKey);
        console2.log("deployer: %s", deployerAddress);
        vm.startBroadcast(deployerPrivateKey);

        BasinStorage basinStorage = new BasinStorage();

        // TODO: grant INDEXER_ROLE to other accounts
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), deployerAddress);

        // Create Pub
        // basinStorage.createPub(address(0x123), "test_pub");

        // console2.log("deployer addr:", address(this));

        vm.stopBroadcast();
    }
}
