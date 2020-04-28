pragma solidity ^0.5.0;
//pragma solidity ^0.4.25;
/**
 * 09-其它
 * 1-do...while循环里的continue不再跳转到循环体内,而是跳转到while处判断循环条件,
 *   若条件为假,就退出循环。这一修改更符合一般编程语言的设计风格
 *
 * @author hudenian
 * @dev 2019/12/19 09:57
 */

contract DoWhileLogicAnd99Style {
    uint forSum;
    uint doWhileSum;


    function dowhile(uint x) public{
        uint y = x+10;
        uint z = x+9;
        do {
            x += 1;
            /**
             * 0.4.25版本'Continue' 将会转到循环体内，从而导致死循环
             * 0.5.0版本'Continue'跳转到while处判断循环条件
             */
            if (x > z) continue;
        } while (x < y);
        doWhileSum = x;
    }

    /**
     * c-99语法风格
     * 块内声明的变量，在块外不能使用
     */
    function forsum(uint x) public{

        for(uint i=0;i<x;i++){
            forSum=forSum + i;
        }
        //i = i++; 块内声明的变量，在块外不能使用
    }


    //获取for循环的结果
    function getForSum() view public returns(uint){
        return forSum;
    }

    //获取do...while的结果
    function getDoWhileSum() view public returns(uint){
        return doWhileSum;
    }




    /**
    *0.4.25变量可以在声明块外使用
    *
    */
    //   function getsum() public returns(uint sum){

    //         for(uint i=0;i<10;i++){
    //             sum=sum + i;
    //         }
    //         sum = sum+100;  //块内声明的变量，在块外可以使用
    //         return sum;
    //     }

}
