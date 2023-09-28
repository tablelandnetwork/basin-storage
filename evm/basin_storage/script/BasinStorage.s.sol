// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console2} from "forge-std/Script.sol";

import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageScript is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddress = vm.addr(deployerPrivateKey);
        vm.startBroadcast(deployerPrivateKey);

        BasinStorage basinStorage = new BasinStorage();

        // grant PUB_ADMIN_ROLE to required accounts
        address pubAdimin2 = 0x0eed5C7ac9D867239A5F550cF94E740f515659Ab;
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), deployerAddress);
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), pubAdimin2);

        vm.stopBroadcast();
    }
}

contract BasinStorageCreatePub is Script {
    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        BasinStorage basinStorage = BasinStorage(
            vm.envAddress("BASIN_STORAGE")
        );

        string memory pubName = vm.envString("PUB");
        address owner = vm.envAddress("OWNER");
        basinStorage.createPub(owner, pubName);

        vm.stopBroadcast();
    }
}
