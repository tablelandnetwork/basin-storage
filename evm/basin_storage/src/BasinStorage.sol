// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.21;

import {Ownable} from "openzeppelin/access/Ownable.sol";
import {AccessControl} from "openzeppelin/access/AccessControl.sol";

contract BasinStorage is AccessControl {
    bytes32 public constant PUB_ADMIN_ROLE = keccak256("PUB_ADMIN_ROLE");

    // Pub address by pub name
    mapping(string => address) private _pubs;

    // Pubs by owner, a reverse mapping of _pubs
    mapping(address => string[]) private _ownerPubs;

    // CID count by pub
    mapping(string => uint256) pubCIDCount;

    // CID storage indexes by pub, indexed by epoch.
    mapping(string pub => mapping(uint256 epoch => string[])) private _cids;

    // Event to log when a CID is indexed
    event CIDAdded(
        string indexed cid,
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

    // IncorrectRange is returned when the timestamp range is incorrect
    error IncorrectRange(uint256 aftr, uint256 before);

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

    /// @dev Adds the CID for the pub and timestamp.
    ///      Can only be called by the Pub Admin.
    /// @param pub The publication to add the CID for.
    /// @param cid The content id for object in the Filecoin deal.
    /// @param timestamp The timestamp provided by data owner.
    function addCID(
        string calldata pub,
        string calldata cid,
        uint256 timestamp
    ) external onlyRole(PUB_ADMIN_ROLE) {
        address owner = _pubs[pub];
        // Pub must already exist
        if (owner == address(0)) {
            revert PubDoesNotExist(pub);
        }
        _cids[pub][timestamp].push(cid);
        pubCIDCount[pub]++;
        emit CIDAdded(cid, pub, owner);
    }

    /// @dev Returns the pubs of a given data owner.
    /// @param owner The owner address to get the pubs for.
    /// @return The pubs for the given data owner.
    function pubsOfOwner(address owner) public view returns (string[] memory) {
        return _ownerPubs[owner];
    }

    /// @dev Returns the CIDs for a given pub in the given epoch range.
    /// @param pub The pub to get the CIDs for.
    /// @param aftr CIDs to fetch after _this_ ts.
    /// @param before CIDs to fetch before _this_ ts.
    function cidsInRange(
        string calldata pub,
        uint256 aftr,
        uint256 before
    ) external view returns (string[] memory) {
        if (aftr >= before) {
            revert IncorrectRange(aftr, before);
        }
        string[] memory cids = new string[](_pubCIDCount[pub]);
        uint256 lastCIDFetchedIdx = 0;

        uint256 epoch = aftr + 1;
        while (epoch < before) {
            string[] memory epochCIDs = _cids[pub][epoch];
            for (uint256 i = 0; i < epochCIDs.length; i++) {
                cids[lastCIDFetchedIdx] = epochCIDs[i];
                lastCIDFetchedIdx++;
            }
            epoch++;
        }

        // remove the trailing elements if they are empty.
        // there can be empty array elements towards the end.
        assembly {
            mstore(cids, lastCIDFetchedIdx)
        }
        return cids;
    }

    /// @dev Returns the CIDs for a given pub at a given epoch.
    /// @param pub The pub to get the cids for.
    /// @param epoch The epoch to get the cids for.
    /// @return cids The CIDs for the given pub at the given epoch.
    function cidsAtTimestamp(
        string calldata pub,
        uint256 epoch
    ) external view returns (string[] memory) {
        return _cids[pub][epoch];
    }
}
