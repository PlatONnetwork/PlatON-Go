pragma solidity ^0.5.13;
/**
 * blockhash功能测试
 */
contract Blockhash {

    //取当前块高之前300个块hash(理论上取不到)
    function getBlockhashbefore257() public view returns(bytes32){
        bytes32 blockhs = blockhash(block.number - 257);

        return blockhs;
    }

    //取当前块高之前256个块hash(理论上取不到)
    function getBlockhashbefore256() public view returns(bytes32){
        bytes32 blockhs = blockhash(block.number - 256);

        return blockhs;
    }

    //取当前块高之前255个块hash(理论上取不到)
    function getBlockhashbefore255() public view returns(bytes32){
        bytes32 blockhs = blockhash(block.number - 255);

        return blockhs;
    }


    //取当前块高之前30个块hash(理论上可以取到)
    function getBlockhashbefore30() public view returns(bytes32){
        bytes32 blockhs = blockhash(block.number - 30);

        return blockhs;
    }

    //取当前块高之前30个块hash(理论上可以取到)
    function getBlockhashbefore0() public view returns(bytes32){
        bytes32 blockhs = blockhash(block.number);

        return blockhs;
    }

}