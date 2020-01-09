pragma solidity 0.5.13;
/**
 * 验证delete关键字，验证各种类型的delete,包括bool,uint,address,bytes,string,enum,变长数组
 * 验证删除struct, struct中的映射不会被删除其他的会被删除
 * @author liweic
 * @dev 2020/01/09 16:10
 */

contract DeleteDemo{

    bool public b  = true;
    uint public i = 1;
    address public addr = msg.sender;
    bytes public varByte = "123";
    string  public str = "abc";
    enum Color{RED,GREEN,YELLOW}
    Color public color = Color.GREEN;

    struct S{
        uint a;
        string r;
    }

    S s;

    struct MapStruct{
        mapping(address => uint) m;
        uint n;
    }

    MapStruct ms;

    function delMapping() payable public{
        ms = MapStruct(200);
        ms.m[msg.sender] = 2000;

        delete ms;
    }

    function getdelMapping() view public returns(uint,uint){
        return (ms.m[msg.sender],ms.n);
    }

    function delStruct() payable public returns(uint, string memory){
        DeleteDemo.S(10, "abc");
        delete s;

        return (s.a,s.r);
    }

    function deleteAttr() public {
        delete b; // false
        delete i; // 0
        delete addr; // 0x0
        delete varByte; // 0x
        delete str; // ""
        delete color;//Color.RED
    }

    function getbool() view public returns(bool){
        return b;
    }

    function getunit() view public returns(uint){
        return i;
    }

    function getaddress() view public returns(address){
        return addr;
    }

    function getbytes() view public returns(bytes memory){
        return varByte;
    }

    function getstr() view public returns(string memory){
        return str;
    }

    function getenum() view public returns(Color){
        return color;
    }

    function getstruct() view public returns(uint, string memory){
        return (s.a,s.r);
    }

    function delDynamicArray() view public returns(uint){
        uint[] memory a = new uint[](7);
        a[0] = 100;
        a[1] = 200;
        delete a;
        return (a.length);
    }
}