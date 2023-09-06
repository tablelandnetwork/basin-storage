// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.21;

import {Ownable} from "openzeppelin/access/Ownable.sol";
import {AccessControl} from "openzeppelin/access/AccessControl.sol";
import {MarketAPI} from "filecoin-solidity/contracts/v0.8/MarketAPI.sol";
import {MarketTypes} from "filecoin-solidity/contracts/v0.8/types/MarketTypes.sol";
import {CommonTypes} from "filecoin-solidity/contracts/v0.8/types/CommonTypes.sol";

contract BasinStorage is AccessControl {
    bytes32 public constant PUB_ADMIN_ROLE = keccak256("PUB_ADMIN_ROLE");

    // DealInfo contains metadata about a Filecoin deal
    struct DealInfo {
        uint64 id;
        string selectorPath;
    }

    // Pub address by pub name
    mapping(string => address) private _pubs;

    // Pubs by owner, a reverse mapping of _pubs
    mapping(address => string[]) private _ownerPubs;

    // deal count by pub
    mapping(string => uint256) private _pubDealCount;

    // Deal storage indexes by pub, indexed by epoch.
    mapping(string pub => mapping(uint256 epoch => DealInfo[])) private _deals;    

    // Event to log when a deal is added or updated
    event DealAdded(
        uint256 indexed dealId,
        string indexed pub,
        address indexed owner
    );

    // Event to log when a pub is created
    event PubCreated(string indexed pub, address indexed owner);

    // Error messages

    // PubAlreadyExists is returned when a pub already exists
    error PubAlreadyExists(string reason);

    // PubDoesNotExist is returned when a pub doesn't exist
    error PubDoesNotExist(string reason);

    // DealEpochAlreadyExists is returned when a deal already exists for an epoch
    error DealEpochAlreadyExists(uint256 epoch);

    constructor() {
        // Set the deployer as the default admin role
        // the default admin shall grant INDEXER roles to other accounts
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    // @dev Creates a new pub for a given owner.
    function createPub(
        address owner,
        string calldata pub
    ) external onlyRole(PUB_ADMIN_ROLE) {
        // Check if pub doesn't already exist
        if (_pubs[pub] != address(0)) {
            revert PubAlreadyExists(pub);
        }

        _pubs[pub] = owner;

        // Add the pub to the owner's list of pubs
        _ownerPubs[owner].push(pub);

        emit PubCreated(pub, owner);
    }

    /// @dev Adds the deal info, given deals and a pub.
    ///      Can only be called by the Pub Admin.
    /// @param deals The Filecoin deal object.
    /// @param pub The publication to add the deal info for.
    function addDeals(
        string calldata pub,
        DealInfo[] calldata deals
    ) external onlyRole(PUB_ADMIN_ROLE) {
        uint256 epoch = block.number;
        address owner = _pubs[pub];

        // Pub must already exist
        if (owner == address(0)) {
            revert PubDoesNotExist(pub);
        }

        // Loop through the input deals and add them against an epoch
        for (uint256 i = 0; i < deals.length; i++) {
            _deals[pub][epoch].push(deals[i]);
            _pubDealCount[pub]++;
            emit DealAdded(deals[i].id, pub, owner);
        }
    }

    /// @dev Returns the pubs of a given data owner.
    /// @param owner The owner address to get the deals for.
    /// @return deals The deals for the given data owner.
    function pubsOfOwner(address owner) public view returns (string[] memory) {
        return _ownerPubs[owner];
    }

    /// @dev Returns the given number of deals for a
    ///      given pub starting from a given epoch.
    /// @param pub The pub to get the deals for.
    /// @param startEpoch The epoch to start from.
    /// @param numDealsToFetch The number of deals to fetch.
    /// @return deals The deals for the given pub.
    function _getDeals(
        string calldata pub,
        uint256 startEpoch,
        uint256 numDealsToFetch
    ) internal view returns (DealInfo[] memory) {
        uint256 lastDealFetchedIdx = 0;
        DealInfo[] memory deals = new DealInfo[](numDealsToFetch);

        uint256 epoch = startEpoch;
        // walk backwards until block 0
        while (epoch >= 0 && lastDealFetchedIdx < numDealsToFetch) {
            DealInfo[] memory epochDeals = _deals[pub][epoch];
            for (uint256 i = 0; i < epochDeals.length; i++) {
                // return if we have fetched the required number of deals
                if (lastDealFetchedIdx >= numDealsToFetch) {
                    break;
                }
                deals[lastDealFetchedIdx] = epochDeals[i];
                lastDealFetchedIdx++;
            }

            if (epoch == 0) {
                break;
            }

            epoch--;
        }

        // remove the trailing elements if they are empty.
        // there can be empty array elements towards the end
        // if NumDealsToFetch is greater than the total deals
        // available within [epoch-0]
        assembly {
            mstore(deals, lastDealFetchedIdx)
        }
        return deals;
    }

    /// @dev Returns the latest N deals for a given pub.
    /// @param pub The pub to get the deals for.
    /// @param n The number of deals to fetch.
    /// @return deals The deals for the given pub.
    function latestNDeals(
        string calldata pub,
        uint256 n
    ) external view returns (DealInfo[] memory) {
        uint256 totalDeals = _pubDealCount[pub];
        // if n is greater than total deals indexed for the pub
        // set n to total deals
        n = n > totalDeals ? totalDeals : n;

        return _getDeals(pub, block.number, n);
    }

    /// @dev Returns the `limit` number of deals
    ///      for a given pub, starting from `offset` epoch.
    /// @param pub The pub to get the deals for.
    /// @param offset The epoch to start from.
    /// @param limit The number of deals to fetch.
    function dealsPaginated(
        string calldata pub,
        uint256 offset,
        uint256 limit
    ) external view returns (DealInfo[] memory) {
        uint256 totalDeals = _pubDealCount[pub];
        // if limit is greater than total deals indexed for the pub
        // set limit to total deals
        limit = limit > totalDeals ? totalDeals : limit;

        return _getDeals(pub, offset, limit);
    }

    // MARKET API Wrappers

    /// @dev returns the client id for a given deal
    /// @param dealID the deal id
    /// @return the client id
    function dealClient(uint64 dealID) public view returns (uint64) {
        return MarketAPI.getDealClient(dealID);
    }

    /// @dev returns the provider id for a given deal
    /// @param dealID the deal id
    /// @return the provider id
    function dealProvider(uint64 dealID) public view returns (uint64) {
        return MarketAPI.getDealProvider(dealID);
    }

    /// @dev returns the label for a deal
    /// @param dealID the deal id
    /// @return the label and if label isString for a deal
    function dealLabel(uint64 dealID) public view returns (bytes memory, bool) {
        CommonTypes.DealLabel memory label = MarketAPI.getDealLabel(dealID);
        return (label.data, label.isString);
    }

    /// @dev returns the start and end epoch for a deal
    /// @param dealID the deal id
    /// @return the start and end epoch for a deal
    function dealTerm(uint64 dealID) public view returns (int64, int64) {
        MarketTypes.GetDealTermReturn memory term = MarketAPI.getDealTerm(
            dealID
        );
        int64 start = CommonTypes.ChainEpoch.unwrap(term.start);
        int64 end = CommonTypes.ChainEpoch.unwrap(term.end);
        return (start, end);
    }

    /// @dev returns the total price paid for a deal
    /// @param dealID the deal id
    /// @return the total price paid for a deal
    function dealTotalPrice(
        uint64 dealID
    ) public view returns (CommonTypes.BigInt memory) {
        return MarketAPI.getDealTotalPrice(dealID);
    }

    /// @dev gives the client's collateral amount for a deal
    /// @param dealID the deal id
    /// @return the client collateral for a deal
    function dealClientCollateral(
        uint64 dealID
    ) public view returns (CommonTypes.BigInt memory) {
        return MarketAPI.getDealClientCollateral(dealID);
    }

    /// @dev gives the provider's collateral amount for a deal
    /// @param dealID the deal id
    /// @return the provider collateral for a deal
    function dealProviderCollateral(
        uint64 dealID
    ) public view returns (CommonTypes.BigInt memory) {
        return MarketAPI.getDealProviderCollateral(dealID);
    }

    /// @dev returns true if a deal is verified
    /// @param dealID the deal id
    /// @return verified
    function dealVerified(uint64 dealID) public view returns (bool) {
        return MarketAPI.getDealVerified(dealID);
    }

    /// @dev gives the activation period for a deal
    /// @param dealID the deal id
    /// @return start and end epoch for a deal
    function dealActivation(uint64 dealID) public view returns (int64, int64) {
        MarketTypes.GetDealActivationReturn memory activation = MarketAPI
            .getDealActivation(dealID);

        int64 activated = CommonTypes.ChainEpoch.unwrap(activation.activated);
        int64 terminated = CommonTypes.ChainEpoch.unwrap(activation.terminated);
        return (activated, terminated);
    }
}
