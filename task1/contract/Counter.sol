// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0; 

contract Counter {

    event Increment(uint256 count);
    uint public count;

    function get() public view returns (uint) {
        return count;
    }

    function inc() public {
        count += 1;
        emit Increment(count);
    }
}