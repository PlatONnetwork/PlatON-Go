pragma solidity ^0.5.13;

/**
 * EVM 插入排序算法复杂度验证
 **/
contract InsertSort{

    function OuputArrays(int[] memory arr, uint n) public payable returns(int[] memory){
        uint i;
        uint k;
        uint j;
        for(i=1;i<n;i++)
        {
            uint k;
            int temp=arr[i];
            j=i;
            while(j>=1 && temp<arr[j-1])
            {
                arr[j]=arr[j-1];
                j--;
            }
            arr[j]=temp;
        }
        return arr;
    }
}