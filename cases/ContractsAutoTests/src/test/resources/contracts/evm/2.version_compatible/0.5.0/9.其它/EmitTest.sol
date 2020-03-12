pragma solidity ^0.5.0;
/**
 * 09-其它
 * 1-引入emit做为触发事件的关键字而不是做为一个修饰符
 *
 * @author hudenian
 * @dev 2019/12/24 09:57
 */

contract EmitTest {
  
   event EventName(address bidder, uint amount); 

   function testEvent() payable public {
    // emit将触发一个事件（0.5.0必须使用emit关键字）
     emit EventName(msg.sender, msg.value);
   }

//0.4.0版本触发事件可以不需要使用emit关键字
//    function testEvent() payable public {
//     // emit将触发一个事件
//      EventName(msg.sender, msg.value); 
//    }
}     
