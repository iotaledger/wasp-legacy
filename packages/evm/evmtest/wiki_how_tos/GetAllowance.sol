// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iscmagic/ISC.sol";

contract GetAllowance {
    event AllowanceFrom(ISCAssets assets);

    function getAllowanceFrom(address _address) public {
        ISCAssets memory assets = ISC.sandbox.getAllowanceFrom(_address);
        emit AllowanceFrom(assets);
    }
}
