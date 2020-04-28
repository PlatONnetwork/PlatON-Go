pragma solidity ^0.5.13;
/**
 *  控制结构
 *  1. if...else
 *  2. do...while
 *  3. for循环
 *  4. for循环包含break
 *  5. for循环包含continue
 *  6. for循环包含return
 *  7. 三目运算符
 *
 *
 * @author hudenian
 * @dev 2019/12/25 11:09
 */


contract Control {

    //1.if控制结构执行结果值
    string public ifControlResult;

    //2.doWhile控制结构执行结果值
    uint public doWhileControlResult;

    //3.forControl控制结构执行结果值
    uint public forControlResult;

    //4.forBreakControl控制结构执行结果值
    uint public forBreakControlResult;

    //5.forContinueControl控制结构执行结果值
    uint public forContinueControlResult;

    //6.forReturnControl控制结构执行结果值
    uint public forReturnControlResult;

    //forControl控制结构执行结果值
    string public forThreeControlControlResult;


    /**
     *1. if...else
     */
    function ifControl(uint age) public {
        if(age < 20){
            ifControlResult = "you are a young man";
        }else if (age < 60){
            ifControlResult = "you are a middle man";
        }else {
            ifControlResult = "you are a old man";
        }
    }

    function getIfControlResult() view public returns(string memory){
        return ifControlResult;
    }

    /**
     *2. do...while
     */
    function doWhileControl() public returns (uint) {
        doWhileControlResult = 0;
        uint i = 0;
        do{
            doWhileControlResult +=i;
            ++i;
        }while(i <10);
        return doWhileControlResult;
    }

    function getdoWhileResult() view public returns(uint){
        return doWhileControlResult;
    }



    /**
     *3. for循环
     */
    function forControl() public returns (uint) {
        forControlResult = 0;
        for(uint i=0;i<10;i++){
            forControlResult +=i;
        }
        return forControlResult;
    }

    function getForControlResult() view public returns(uint){
        return forControlResult;
    }

    /**
     *4. for循环包含break
     * 满足条件就退出for循环
     */
    function forBreakControl() public returns (uint) {
        forBreakControlResult = 0;
        for(uint i=1;i<10;i++){
            if(i % 2 == 0){
                break;
            }
            forBreakControlResult +=i;
        }
        return forBreakControlResult;
    }

    function getForBreakControlResult() view public returns(uint){
        return forBreakControlResult;
    }

    /**
     * 5. for循环包含continue
     * 满足条件的就跳过
     */
    function forContinueControl() public returns (uint) {
        forContinueControlResult = 0;
        for(uint i=0;i<10;i++){
            if(i % 2 == 0){
                continue;
            }
            forContinueControlResult +=i;
        }
        return forContinueControlResult;
    }

    function getForContinueControlResult() view public returns(uint){
        return forContinueControlResult;
    }

    /**
     * 6. for循环包含return
     *    满足条件就返回
     */
    function forReturnControl() public returns (uint) {
        forReturnControlResult = 0;
        for(uint i=1;i<10;i++){
            if(i % 5 == 0){
                return forReturnControlResult;
            }
            forReturnControlResult +=i;
        }
        return forReturnControlResult;
    }

    function getForReturnControlResult() view public returns(uint){
        return forReturnControlResult;
    }


    /**
     *  7. 三目运算符
     */
    function forThreeControlControl(int age) public {
        forThreeControlControlResult = age> 20?"less than 20":"more than 20";
    }

    function getForThreeControlControlResult() view public returns(string memory){
        return forThreeControlControlResult;
    }

}