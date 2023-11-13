// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.21;

import {Script, console2} from "forge-std/Script.sol";
import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageScript is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddress = vm.addr(deployerPrivateKey);
        console2.log("deployerAddress: %s", deployerAddress);

        vm.startBroadcast(deployerPrivateKey);

        BasinStorage basinStorage = new BasinStorage();

        console2.log("Contract Address: %s", address(basinStorage));
        // grant PUB_ADMIN_ROLE to required accounts (basin staging wallet)
        address pubAdimin2 = 0x1D3888b19E973E3341960a1938e51e40875a6A15;
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), deployerAddress);
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), pubAdimin2);

        vm.stopBroadcast();
    }
}

