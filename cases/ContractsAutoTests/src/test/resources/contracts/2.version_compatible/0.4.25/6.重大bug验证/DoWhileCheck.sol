pragma solidity ^0.4.14;

/**
 * 1. 0.4.25版本do...while循环里的continue跳转到循环体内，可能会产生死循环验证
 * 2. 0.4.25版本局部变量上级作用域生效验证
 * @author Albedo
 * @dev 2019/12/24
 **/
contract DoWhileCheck {
    function doWhileCheck() public view returns (uint256,uint256){
        uint count=0;
        uint a=0;
        //0.4.25版本do...while循环里的continue跳转到循环体内，可能会产生死循环验证
        do {
            count++;
            a++;
            if(a > 20){
                break;
            }
            //因为continue会直接跳回到do处，所以count和a都会自增，直到a到21时通过break跳出整个循环
            //因此，此时的count的值也是21
            if(count>10) continue;
        }while(count<30);

        //0.4.25版本局部变量上级作用域生效验证
        if(count==21){
            uint c=12;
        }
        c=14;
        return (count,c); //输出(21,14)
    }
}