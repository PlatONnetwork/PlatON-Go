pragma solidity ^0.5.13;

/**
 * EVM 冒泡排序算法复杂度验证
 **/

contract BubbleSort{

    int[] result_arr;

    function BubbleArrays(int[] memory arr, uint n) public payable{
        for (uint i = 0; i < n - 1; i++)
        {
            for (uint j = 0; j<n - 1 - i; j++)
            {
                if (arr[j]>arr[j + 1])
                {
                    int temp;
                    temp = arr[j];
                    arr[j] = arr[j + 1];
                    arr[j + 1] = temp;
                }
            }
        }

        result_arr = arr;
    }

    function get_arr() public view returns(int[] memory){
        return result_arr;
    }

}