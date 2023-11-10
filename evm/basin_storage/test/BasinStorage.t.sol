// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.21;

import {Test, console2} from "forge-std/Test.sol";
import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageAddDealsTest is Test {
    BasinStorage public basinStorage;

    event CIDAdded(
        string indexed cid,
        string indexed pub,
        address indexed owner
    );

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a deal
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function testAddCIDUnauthorized() public {
        // Call the AddDeal function with an unauthorized account
        vm.prank(address(0));

        // Define the input parameters
        string memory pub = "123456";
        uint256 epoch = block.timestamp;
        string memory cid = "bafyfoobar1";

        string memory reason = string.concat(
            "AccessControl: account ",
            "0x0000000000000000000000000000000000000000"
            " is missing role ",
            "0xafda658ee731b8f86292e3b52a311534cd93642b12a698012439316e0c3a0995"
        );
        vm.expectRevert(bytes(reason));
        basinStorage.addCID(pub, cid, epoch);
    }

    // Test the CreateDealInfo function
    function testAddCIDWithoutAddingPub() public {
        // Define the input parameters
        string memory pub = "123456";
        string memory cid = "bafyfoobar1";
        uint256 epoch = block.timestamp;
        vm.expectRevert(
            abi.encodeWithSelector(BasinStorage.PubDoesNotExist.selector, pub)
        );
        basinStorage.addCID(pub, cid, epoch);
    }

    function testAddCIDSuccess() public {
        string memory pub = "123456";
        uint256 epoch = block.timestamp;
        string memory cid1 = "bafyfoobar1";
        string memory cid2 = "bafyfoobar2";
        basinStorage.createPub(address(this), pub);

        // check that the 1st event is emitted
        vm.expectEmit(address(basinStorage));
        emit BasinStorage.CIDAdded("bafyfoobar1", pub, address(this));
        basinStorage.addCID(pub, cid1, epoch);

        // check that the 2nd event is emitted
        vm.expectEmit(address(basinStorage));
        emit BasinStorage.CIDAdded("bafyfoobar2", pub, address(this));
        basinStorage.addCID(pub, cid2, epoch);

        string[] memory cids = basinStorage.cidsAtTimestamp(pub, 1);

        assertEq(cids.length, 2, "Number of deals should be 2");
        assertEq(
            cids[0],
            "bafyfoobar1",
            "content identifier should be correct"
        );
        assertEq(
            cids[1],
            "bafyfoobar2",
            "content identifier should be correct"
        );
    }
}

contract BasinStoragePubsTest is Test {
    BasinStorage public basinStorage;

    event PubCreated(string indexed pub, address indexed owner);

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a deal
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function testCreatePubUnauthorized() public {
        vm.prank(address(0)); // this address is not permitted to create a pub
        string memory pub = "123456";
        string memory reason = string.concat(
            "AccessControl: account ",
            "0x0000000000000000000000000000000000000000"
            " is missing role ",
            "0xafda658ee731b8f86292e3b52a311534cd93642b12a698012439316e0c3a0995"
        );
        vm.expectRevert(bytes(reason));
        basinStorage.createPub(address(0), pub);
    }

    function testCreatePubSuccess() public {
        string memory pub = "123456";
        vm.expectEmit(address(basinStorage));
        emit BasinStorage.PubCreated(pub, address(this));
        basinStorage.createPub(address(this), pub);
        assertEq(basinStorage.pubsOfOwner(address(this)).length, 1);
        assertEq(basinStorage.pubsOfOwner(address(this))[0], pub);
    }

    function testPubsOfOwner() public {
        string memory pub1 = "123456";
        string memory pub2 = "654321";
        string memory pub3 = "111111";

        basinStorage.createPub(address(this), pub1);
        basinStorage.createPub(address(0x123), pub2);
        basinStorage.createPub(address(this), pub3);

        string[] memory pubs = basinStorage.pubsOfOwner(address(this));
        assertEq(pubs.length, 2, "Number of pubs should be 3");
        assertEq(pubs[0], pub1, "Pub should be correct");
        assertEq(pubs[1], pub3, "Pub should be correct");

        pubs = basinStorage.pubsOfOwner(address(0x123));
        assertEq(pubs.length, 1, "Number of pubs should be 3");
        assertEq(pubs[0], pub2, "Pub should be correct");
    }
}

abstract contract HelperContract is Test {
    function setUpHelper(BasinStorage basinStorage) public {
        // Create deals for pub 1, owner 1 (current contract)
        string memory pub = "123456";
        string memory cid1 = "bafyfoobar1";
        string memory cid2 = "bafyfoobar2";
        string memory cid3 = "bafyfoobar3";
        uint256 epoch1 = block.timestamp;

        basinStorage.createPub(address(this), pub);
        basinStorage.addCID(pub, cid1, epoch1);
        basinStorage.addCID(pub, cid2, epoch1);
        basinStorage.addCID(pub, cid3, epoch1);

        // Create deals for pub 2, owner 1 (current contract)
        pub = "654321";
        string memory cid4 = "bafyfoobar4";
        string memory cid5 = "bafyfoobar5";
        string memory cid6 = "bafyfoobar6";
        uint256 epoch2 = block.timestamp + 1;

        basinStorage.createPub(address(0x123), pub);
        basinStorage.addCID(pub, cid4, epoch2);
        basinStorage.addCID(pub, cid5, epoch2);
        basinStorage.addCID(pub, cid6, epoch2);

        // Create deals for pub 3, owner 2 (address 0x123)
        pub = "111111";
        string memory cid7 = "bafyfoobar7";
        string memory cid8 = "bafyfoobar8";
        string memory cid9 = "bafyfoobar9";
        uint256 epoch3 = block.timestamp + 2;

        basinStorage.createPub(address(this), pub);
        basinStorage.addCID(pub, cid7, epoch3);
        basinStorage.addCID(pub, cid8, epoch3);
        basinStorage.addCID(pub, cid9, epoch3);

        // same pub as 123456 but on a different block
        pub = "123456";
        uint256 epoch4 = block.timestamp + 3;
        string memory cid10 = "bafyfoobar10";
        string memory cid11 = "bafyfoobar11";
        string memory cid12 = "bafyfoobar12";

        // no creating pub if it already exists
        basinStorage.addCID(pub, cid10, epoch4);
        basinStorage.addCID(pub, cid11, epoch4);
        basinStorage.addCID(pub, cid12, epoch4);
    }
}

contract BasinStorageCidsInRangeTest is Test, HelperContract {
    BasinStorage public basinStorage;

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a cid
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function setUp() public {
        HelperContract.setUpHelper(basinStorage);
    }

    function testcidsInRange() public {
        string memory pub = "123456";
        // after 0, before 4, excluding both 0 and 5
        string[] memory deals = basinStorage.cidsInRange(pub, 0, 4);

        assertEq(deals.length, 3, "Deals count should be 3");
        assertEq(deals[0], "bafyfoobar1", "cid should be bafyfoobar1");
        assertEq(deals[1], "bafyfoobar2", "cid should be bafyfoobar2");
        assertEq(deals[2], "bafyfoobar3", "cid should be bafyfoobar3");

        // after 0, before 5, excluding both 0 and 5
        deals = basinStorage.cidsInRange(pub, 0, 5);
        assertEq(deals.length, 6, "Deals count should be 6");
        assertEq(deals[0], "bafyfoobar1", "cid should be bafyfoobar1");
        assertEq(deals[1], "bafyfoobar2", "cid should be bafyfoobar2");
        assertEq(deals[2], "bafyfoobar3", "cid should be bafyfoobar3");
        assertEq(deals[3], "bafyfoobar10", "cid should be bafyfoobar10");
        assertEq(deals[4], "bafyfoobar11", "cid should be bafyfoobar11");
        assertEq(deals[5], "bafyfoobar12", "cid should be bafyfoobar12");

        deals = basinStorage.cidsInRange(pub, 4, 5);
        assertEq(deals.length, 0, "Deals count should be 0");

        // after == before raises error
        vm.expectRevert(
            abi.encodeWithSelector(BasinStorage.IncorrectRange.selector, 1, 1)
        );
        deals = basinStorage.cidsInRange(pub, 1, 1);
        assertEq(deals.length, 0, "Deals count should be 0");

        // after > before raises error
        vm.expectRevert(
            abi.encodeWithSelector(BasinStorage.IncorrectRange.selector, 5, 4)
        );
        deals = basinStorage.cidsInRange(pub, 5, 4);

        pub = "654321"; // same owner different pub
        deals = basinStorage.cidsInRange(pub, 0, 5);
        assertEq(deals.length, 3, "Deals count should be 3");
        assertEq(deals[0], "bafyfoobar4", "cid should be bafyfoobar4");
        assertEq(deals[1], "bafyfoobar5", "cid should be bafyfoobar5");
        assertEq(deals[2], "bafyfoobar6", "cid should be bafyfoobar6");

        deals = basinStorage.cidsInRange(pub, 1, 3);
        assertEq(deals.length, 3, "Deals count should be 3");

        deals = basinStorage.cidsInRange(pub, 1, 2);
        assertEq(deals.length, 0, "Deals count should be 0");

        pub = "111111"; // pub of diff owner: 0x123
        deals = basinStorage.cidsInRange(pub, 2, 5);
        assertEq(deals.length, 3, "Deals count should be 3");
        assertEq(deals[0], "bafyfoobar7", "cid should be bafyfoobar7");
        assertEq(deals[1], "bafyfoobar8", "cid should be bafyfoobar8");
        assertEq(deals[2], "bafyfoobar9", "cid should be bafyfoobar9");
    }
}
