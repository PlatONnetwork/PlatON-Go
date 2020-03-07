pragma solidity 0.5.13;

contract ReferenceDataTypeArrayComplexContract {


   uint256[] arrayUint = [1,2]; 

    /**
     * 验证含数组运算逻辑合约
     * 此函数是加和运算限制：
     * 1)、数组长度不可以超过10；
     * 2)、数组值将大于1000过滤掉不进行计算；
     * 3)、总和大于500将不再进行加和运算
     */
    function sumComplexArray(uint256[] calldata array) external pure returns (uint256) {

        uint i = 0;
        uint sum = 0;
        while(i < array.length){
            uint idx = i;
            //1、数组长度大于10，停止执行计算
            if(idx > 10)
                break;
            //获取数组中值
            uint x = array[idx];
            //2、数组值将大于1000过滤掉不进行计算
            if(x >= 1000){
                 i += 1;
                 continue; 
            }else{
                sum += x;
            } 
            //3、总和大于500将不再进行加和运算
            if(sum >= 500)
                return sum;
            
            i++;
        }
        return sum;
    }
}
