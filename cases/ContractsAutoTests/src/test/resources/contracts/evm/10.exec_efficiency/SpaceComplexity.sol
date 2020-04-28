pragma solidity ^0.5.13;
/**
 * 空间复杂度验证
 * （1）申请storage内存，测试申请内存的极限
 * （2）申请memory内存，测试申请内存的极限
 * （3）申请calldata内存，测试申请内存的极限
 * @author Albedo
 * @dev 2019/12/19
 **/
contract SpaceComplexity {

    string _name = "qcxiao";

    // 测试storage并多次修改
    function testStorage(uint n) external {
        for (uint i = 0; i < n; i++) {
            if (n % 2 == 0) { // 所有字符修改为大写
                modifyOfStorage1(_name);
            } else { // 所有字符修改为小写
                modifyOfStorage2(_name);
            }
        }
    }

    // 声明为storage，只能是internal
    function modifyOfStorage1(string storage name) internal {
        bytes(name)[0] = "Q";
        bytes(name)[1] = "C";
        bytes(name)[2] = "X";
        bytes(name)[3] = "I";
        bytes(name)[4] = "A";
        bytes(name)[5] = "O";
    }

    function modifyOfStorage2(string storage name) internal {
        bytes(name)[0] = "q";
        bytes(name)[1] = "c";
        bytes(name)[2] = "x";
        bytes(name)[3] = "i";
        bytes(name)[4] = "a";
        bytes(name)[5] = "o";
    }

    function name() public view returns (string memory) {
        return _name;
    }

    function testBigObjectOfStorage(uint n) public {

    }

    //memory

}