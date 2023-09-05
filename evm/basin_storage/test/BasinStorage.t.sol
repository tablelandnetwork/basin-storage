// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.21;

import {Test, console2} from "forge-std/Test.sol";
import {BasinStorage} from "../src/BasinStorage.sol";

contract BasinStorageAddDealsTest is Test {
    BasinStorage public basinStorage;

    event DealAdded(
        uint256 indexed dealId,
        string indexed pub,
        address indexed owner
    );

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a deal
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function testAddDealUnauthorized() public {
        // Call the AddDeal function with an unauthorized account
        vm.prank(address(0));

        // Define the input parameters
        string memory pub = "123456";
        // BasinStorage.DealInfo memory dealInfo[] =
        BasinStorage.DealInfo[] memory deals = new BasinStorage.DealInfo[](1);
        deals[0] = BasinStorage.DealInfo({
            id: 1,
            selectorPath: "path/to/selector1"
        });

        string memory reason = string.concat(
            "AccessControl: account ",
            "0x0000000000000000000000000000000000000000"
            " is missing role ",
            "0xafda658ee731b8f86292e3b52a311534cd93642b12a698012439316e0c3a0995"
        );
        vm.expectRevert(bytes(reason));
        basinStorage.addDeals(pub, deals);
    }

    // Test the CreateDealInfo function
    function testAddDealWithoutAddingPub() public {
        // Define the input parameters
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = new BasinStorage.DealInfo[](1);
        deals[0] = BasinStorage.DealInfo({
            id: 1,
            selectorPath: "path/to/selector1"
        });

        vm.expectRevert(
            abi.encodeWithSelector(BasinStorage.PubDoesNotExist.selector, pub)
        );
        basinStorage.addDeals(pub, deals);
    }

    function testAddDealEmptyInput() public {
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = new BasinStorage.DealInfo[](0);

        basinStorage.createPub(address(this), pub);
        basinStorage.addDeals(pub, deals);
        deals = basinStorage.latestNDeals(pub, 2);
        assertEq(deals.length, 0, "Number of deals should be 2");
    }

    function testAddDealSuccess() public {
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = new BasinStorage.DealInfo[](2);
        deals[0] = BasinStorage.DealInfo({
            id: 1,
            selectorPath: "path/to/selector1"
        });
        deals[1] = BasinStorage.DealInfo({
            id: 2,
            selectorPath: "path/to/selector2"
        });

        basinStorage.createPub(address(this), pub);

        // check that the 1st event is emitted
        vm.expectEmit(address(basinStorage));
        emit BasinStorage.DealAdded(1, pub, address(this));

        // check that the 2nd event is emitted
        vm.expectEmit(address(basinStorage));
        emit BasinStorage.DealAdded(2, pub, address(this));

        basinStorage.addDeals(pub, deals);
        deals = basinStorage.latestNDeals(pub, 2);
        assertEq(deals.length, 2, "Number of deals should be 2");
        assertEq(deals[0].id, 1, "Deal ID should be correct");
        assertEq(
            deals[0].selectorPath,
            "path/to/selector1",
            "Deal selector should be correct"
        );
        assertEq(deals[1].id, 2, "Deal ID should be correct");
        assertEq(
            deals[1].selectorPath,
            "path/to/selector2",
            "Deal selector should be correct"
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
        vm.roll(0);
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = new BasinStorage.DealInfo[](3);
        deals[0] = BasinStorage.DealInfo({
            id: 1,
            selectorPath: "path/to/selector1"
        });
        deals[1] = BasinStorage.DealInfo({
            id: 2,
            selectorPath: "path/to/selector2"
        });
        deals[2] = BasinStorage.DealInfo({
            id: 3,
            selectorPath: "path/to/selector3"
        });

        basinStorage.createPub(address(this), pub);
        basinStorage.addDeals(pub, deals);

        // Create deals for pub 2, owner 1 (current contract)
        vm.roll(100);
        pub = "654321";
        deals[0] = BasinStorage.DealInfo({
            id: 4,
            selectorPath: "path/to/selector4"
        });
        deals[1] = BasinStorage.DealInfo({
            id: 5,
            selectorPath: "path/to/selector5"
        });
        deals[2] = BasinStorage.DealInfo({
            id: 6,
            selectorPath: "path/to/selector6"
        });
        basinStorage.createPub(address(0x123), pub);
        basinStorage.addDeals(pub, deals);

        // Create deals for pub 3, owner 2 (address 0x123)
        pub = "111111";
        vm.roll(150);
        deals[0] = BasinStorage.DealInfo({
            id: 7,
            selectorPath: "path/to/selector7"
        });
        deals[1] = BasinStorage.DealInfo({
            id: 8,
            selectorPath: "path/to/selector8"
        });
        deals[2] = BasinStorage.DealInfo({
            id: 9,
            selectorPath: "path/to/selector9"
        });
        basinStorage.createPub(address(this), pub);
        basinStorage.addDeals(pub, deals);

        // same pub as 123456 but on a different block
        pub = "123456";
        vm.roll(200);
        deals[0] = BasinStorage.DealInfo({
            id: 10,
            selectorPath: "path/to/selector10"
        });
        deals[1] = BasinStorage.DealInfo({
            id: 11,
            selectorPath: "path/to/selector11"
        });
        deals[2] = BasinStorage.DealInfo({
            id: 12,
            selectorPath: "path/to/selector12"
        });
        // no creating pub if it already exists
        basinStorage.addDeals(pub, deals);
    }
}

contract BasinStoragePaginatedDealsTest is Test, HelperContract {
    BasinStorage public basinStorage;

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a deal
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function setUp() public {
        HelperContract.setUpHelper(basinStorage);
    }

    function testDealsPaginated(uint16 offset, uint16 limit) public {
        vm.assume(offset >= 0 && offset <= 200);

        // get latest N deals where N > total deals for an epoch range
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = basinStorage.dealsPaginated(
            pub,
            offset,
            limit
        );

        if (limit >= 6 && offset == 200) {
            assertEq(deals.length, 6, "Deals count should be 6");
            // from 2nd (newer) batch
            assertEq(deals[0].id, 10, "Deal id should be 10");
            assertEq(deals[1].id, 11, "Deal id should be 11");
            assertEq(deals[2].id, 12, "Deal id should be 12");
            // from 1st (older) batch
            assertEq(deals[3].id, 1, "Deal id should be 1");
            assertEq(deals[4].id, 2, "Deal id should be 2");
            assertEq(deals[5].id, 3, "Deal id should be 3");
        } else if (limit >= 6 && offset < 200) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 1, "Deal id should be 1");
            assertEq(deals[1].id, 2, "Deal id should be 2");
            assertEq(deals[2].id, 3, "Deal id should be 3");
        } else if (limit > 3 && offset < 200) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 1, "Deal id should be 1");
            assertEq(deals[1].id, 2, "Deal id should be 2");
            assertEq(deals[2].id, 3, "Deal id should be 3");
        } else if (limit < 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }

        pub = "654321"; // same owner different pub
        deals = basinStorage.dealsPaginated(pub, offset, limit);
        // all 3 deals were added @ block number 100
        if (offset < 100) {
            assertEq(deals.length, 0, "Deals count should be 0");
        } else if (offset >= 100 && limit > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 4, "Deal id should be 4");
            assertEq(deals[1].id, 5, "Deal id should be 5");
            assertEq(deals[2].id, 6, "Deal id should be 6");
        } else if (offset >= 100 && limit <= 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }

        pub = "111111"; // pub of 0x123
        deals = basinStorage.dealsPaginated(pub, offset, limit);
        // all 3 deals were added @ block number 150
        if (offset < 150) {
            assertEq(deals.length, 0, "Deals count should be 0");
        } else if (offset >= 150 && limit > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 7, "Deal id should be 7");
            assertEq(deals[1].id, 8, "Deal id should be 8");
            assertEq(deals[2].id, 9, "Deal id should be 9");
        } else if (offset >= 150 && limit <= 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }
    }
}

contract BasinStorageLatestDealsTest is Test, HelperContract {
    BasinStorage public basinStorage;

    constructor() {
        basinStorage = new BasinStorage();
        // give the contract the PUB_ADMIN_ROLE before adding a deal
        basinStorage.grantRole(basinStorage.PUB_ADMIN_ROLE(), address(this));
    }

    function setUp() public {
        HelperContract.setUpHelper(basinStorage);
    }

    function testLatestNDeals(uint16 n) public {
        // get latest N deals where N > total deals
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = basinStorage.latestNDeals(
            pub,
            n
        );

        if (n > 6) {
            assertEq(deals.length, 6, "Deals count should be 6");
            // from 2nd (newer) batch
            assertEq(deals[0].id, 10, "Deal id should be 10");
            assertEq(deals[1].id, 11, "Deal id should be 11");
            assertEq(deals[2].id, 12, "Deal id should be 12");
            // from 1st (older) batch
            assertEq(deals[3].id, 1, "Deal id should be 1");
            assertEq(deals[4].id, 2, "Deal id should be 2");
            assertEq(deals[5].id, 3, "Deal id should be 3");
        } else {
            assertEq(deals.length, n, "Deals count should be: n");
        }

        pub = "111111"; // pub of 0x123
        deals = basinStorage.latestNDeals(pub, n);
        if (n > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 7, "Deal id should be 10");
            assertEq(deals[1].id, 8, "Deal id should be 11");
            assertEq(deals[2].id, 9, "Deal id should be 12");
        } else {
            assertEq(deals.length, n, "Deals count should be: n");
        }

        pub = "654321";
        deals = basinStorage.latestNDeals(pub, n);
        if (n > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 4, "Deal id should be 4");
            assertEq(deals[1].id, 5, "Deal id should be 5");
            assertEq(deals[2].id, 6, "Deal id should be 6");
        } else {
            assertEq(deals.length, n, "Deals count should be: n");
        }
    }

    function testDealsPaginated(uint16 offset, uint16 limit) public {
        vm.assume(offset >= 0 && offset <= 200);

        // get latest N deals where N > total deals for an epoch range
        string memory pub = "123456";
        BasinStorage.DealInfo[] memory deals = basinStorage.dealsPaginated(
            pub,
            offset,
            limit
        );

        if (limit >= 6 && offset == 200) {
            assertEq(deals.length, 6, "Deals count should be 6");
            // from 2nd (newer) batch
            assertEq(deals[0].id, 10, "Deal id should be 10");
            assertEq(deals[1].id, 11, "Deal id should be 11");
            assertEq(deals[2].id, 12, "Deal id should be 12");
            // from 1st (older) batch
            assertEq(deals[3].id, 1, "Deal id should be 1");
            assertEq(deals[4].id, 2, "Deal id should be 2");
            assertEq(deals[5].id, 3, "Deal id should be 3");
        } else if (limit >= 6 && offset < 200) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 1, "Deal id should be 1");
            assertEq(deals[1].id, 2, "Deal id should be 2");
            assertEq(deals[2].id, 3, "Deal id should be 3");
        } else if (limit > 3 && offset < 200) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 1, "Deal id should be 1");
            assertEq(deals[1].id, 2, "Deal id should be 2");
            assertEq(deals[2].id, 3, "Deal id should be 3");
        } else if (limit < 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }

        pub = "654321"; // same owner different pub
        deals = basinStorage.dealsPaginated(pub, offset, limit);
        // all 3 deals were added @ block number 100
        if (offset < 100) {
            assertEq(deals.length, 0, "Deals count should be 0");
        } else if (offset >= 100 && limit > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 4, "Deal id should be 4");
            assertEq(deals[1].id, 5, "Deal id should be 5");
            assertEq(deals[2].id, 6, "Deal id should be 6");
        } else if (offset >= 100 && limit <= 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }

        pub = "111111"; // pub of 0x123
        deals = basinStorage.dealsPaginated(pub, offset, limit);
        // all 3 deals were added @ block number 150
        if (offset < 150) {
            assertEq(deals.length, 0, "Deals count should be 0");
        } else if (offset >= 150 && limit > 3) {
            assertEq(deals.length, 3, "Deals count should be 3");
            assertEq(deals[0].id, 7, "Deal id should be 7");
            assertEq(deals[1].id, 8, "Deal id should be 8");
            assertEq(deals[2].id, 9, "Deal id should be 9");
        } else if (offset >= 150 && limit <= 3) {
            assertEq(deals.length, limit, "Deals count should be == limit");
        }
    }
}
