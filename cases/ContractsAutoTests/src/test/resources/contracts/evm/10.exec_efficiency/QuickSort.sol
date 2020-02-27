pragma solidity ^0.4.12;

/**
 * 快速排序
 ×
 × @author zjsunzone
 */
contract QuickSort {

    uint256[] public arr;
    uint256 public p;
    
    function QuickSort(uint256[] _arr){
        for(uint256 i = 0; i < _arr.length; i++){
            arr.push(_arr[i]);
        }
    }
    
    function sort(uint256 low, uint256 high) public {
        quick_sort(low, high);
    }

    function quick_sort(uint256 low, uint256 high) internal {
        if(low < high){
            uint256 i = partition(low, high);
            if(i != 0){
                quick_sort(low, i-1);
            }
            quick_sort(i+1, high);
            p = i;
        }
    }
    
    function partition(uint256 low, uint256 high) internal returns (uint256) {
        uint256 temp = arr[low];
        uint i = low;
        uint j = high;
        while(i != j) {
            while(i < j && arr[j] > temp){
                j = j - 1;
            }
            if(i < j){
                arr[i] = arr[j];
                i = i + 1;
            }
            while(i < j && arr[i] < temp){
                i = i+1;
            }
            if(i < j){
                arr[j] = arr[i];
                j = j - 1;
            }
            arr[i] = temp;
        }
        arr[i] = temp;
        return i;
    }
}
