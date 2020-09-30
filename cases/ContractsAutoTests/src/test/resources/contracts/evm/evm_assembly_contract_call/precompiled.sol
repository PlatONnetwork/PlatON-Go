pragma solidity ^0.5.13;

/*
network.platon.test.evm 跨合约调用 系统合约
*/
contract precompiled {

    bytes callDatacopyValue;
    bytes32 callBigModExpValue;
    uint256[2] callBn256AddValues;
    bytes32[2] callBn256ScalarMulValues;
    bytes32 callBn256PairingValue;
    address feePayerAddr;
    bool validateSenderFlg;

    //Address 0x01
    function callEcrecover(bytes32 hash, uint8 v, bytes32 r, bytes32 s) public pure returns (address) {
        return ecrecover(hash, v, r, s);
    }


    //Address 0x02: sha256(data)
    function callSha256(bytes memory data) public pure returns(bytes32 result){
        return sha256(data);
    }


    //Address 0x03: ripemd160(data)
    function callRipemd160(bytes memory data) public pure returns(bytes32 result){
        return ripemd160(data);
    }


    //pass 0x04
    function callDatacopy(bytes memory data) public  returns (bytes memory) {
        bytes memory ret = new bytes(data.length);
        assembly {
            let len := mload(data)
            if iszero(call(gas, 0x04, 0, add(data, 0x20), len, add(ret,0x20), len)) {
                invalid()
            }
        }
        callDatacopyValue = ret;
        return ret;
    }

    function getCallDatacopyValue() public view returns(bytes memory){
        return callDatacopyValue;
    }

    //pass 0x05
    function callBigModExp(bytes32 base, bytes32 exponent, bytes32 modulus) public returns (bytes32 result) {
        assembly {
        // free memory pointer
            let memPtr := mload(0x40)

        // length of base, exponent, modulus
            mstore(memPtr, 0x20)
            mstore(add(memPtr, 0x20), 0x20)
            mstore(add(memPtr, 0x40), 0x20)

        // assign base, exponent, modulus
            mstore(add(memPtr, 0x60), base)
            mstore(add(memPtr, 0x80), exponent)
            mstore(add(memPtr, 0xa0), modulus)

        // call the precompiled contract BigModExp (0x05)
            let success := call(gas, 0x05, 0x0, memPtr, 0xc0, memPtr, 0x20)
            switch success
            case 0 {
                revert(0x0, 0x0)
            } default {
                result := mload(memPtr)
            }
        }
        callBigModExpValue = result;
    }


    function getCallBigModExpValue() public view returns(bytes32 ){
        return callBigModExpValue;
    }

    //pass 0x06
    function callBn256Add(uint256 ax, uint256 ay, uint256 bx, uint256 by) public returns (uint256[2] memory result) {
        uint256[4] memory input;
        input[0] = ax;
        input[1] = ay;
        input[2] = bx;
        input[3] = by;
        assembly {
            let success := call(gas, 0x06, 0, input, 0x80, result, 0x40)
            switch success
            case 0 {
                revert(0,0)
            }
        }
        callBn256AddValues = result;
    }


    function getCallBn256AddValues() public view returns(uint256[2] memory result ){
        return callBn256AddValues;
    }


    //0x07
    function callBn256ScalarMul(bytes32 x, bytes32 y, bytes32 scalar) public returns (bytes32[2] memory result) {
        bytes32[3] memory input;
        input[0] = x;
        input[1] = y;
        input[2] = scalar;
        assembly {
            let success := call(gas, 0x07, 0, input, 0x60, result, 0x40)
            switch success
            case 0 {
                revert(0,0)
            }
        }
        callBn256ScalarMulValues = result;
    }

    function getCallBn256ScalarMulValues() public view returns(bytes32[2] memory result ){
        return callBn256ScalarMulValues;
    }

    //Address 0x08: bn256Pairing(a1, b1, a2, b2, a3, b3, ..., ak, bk)
    function callBn256Pairing(bytes memory input) public returns (bytes32 result) {
        // input is a serialized bytes stream of (a1, b1, a2, b2, ..., ak, bk) from (G_1 x G_2)^k
        uint256 len = input.length;
        require(len % 192 == 0);
        assembly {
            let memPtr := mload(0x40)
            let success := call(gas, 0x08, 0, add(input, 0x20), len, memPtr, 0x20)
            switch success
            case 0 {
                revert(0,0)
            } default {
                result := mload(memPtr)
            }
        }
        callBn256PairingValue =result;
    }

    function getCallBn256PairingValue() public view returns(bytes32 result ){
        return callBn256PairingValue;
    }

}