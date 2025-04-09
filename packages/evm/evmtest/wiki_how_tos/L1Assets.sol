// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@iscmagic/ISC.sol";
import "@iscmagic/ISCTypes.sol";

contract L1Assets {

  function allow(address _address, ISCAssets memory _assets ) public {
    // NativeTokenID[] memory nativeTokenIds = new NativeTokenID[](1);
    // nativeTokenIds[0] = NativeTokenID.wrap(_nativeTokenId);
    // ISCAssets memory assets;
    // assets.nativeTokens = _nativeTokenIds;
    ISC.sandbox.allow(_address, _assets);
  }

  function withdraw(L1Address memory to) public {
    ISCAssets memory allowance = ISC.sandbox.getAllowanceFrom(msg.sender);
    ISC.sandbox.takeAllowedFunds(msg.sender, allowance);

    ISCSendMetadata memory metadata;
    ISCSendOptions memory options;
    ISC.sandbox.send(to, allowance, false, metadata, options);
  }
}
