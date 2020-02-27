 pragma solidity ^0.4.12;

/**
 * 快速排序
 ×
 × @author zjsunzone
 */
contract QuickSort {

    int256[] public arr;
    
    function QuickSort(){
        
    }
    
    function sort(int256[] _arr, uint256 low, uint256 high) public {
        quick_sort(_arr, low, high);
        for(uint256 i = 0; i < _arr.length; i++){
            arr.push(_arr[i]);
        }
    }

    function quick_sort(int256[] _arr, uint256 low, uint256 high) internal {
        if(low < high){
            uint256 i = partition(_arr, low, high);
            if(i != 0){
                quick_sort(_arr, low, i-1);
            }
            quick_sort(_arr, i+1, high);
        }
    }
    
    function partition(int256[] _arr, uint256 low, uint256 high) internal returns (uint256) {
        int256 temp = _arr[low];
        uint i = low;
        uint j = high;
        while(i != j) {
            while(i < j && _arr[j] > temp){
                j = j - 1;
            }
            if(i < j){
                _arr[i] = _arr[j];
                i = i + 1;
            }
            while(i < j && _arr[i] < temp){
                i = i+1;
            }
            if(i < j){
                _arr[j] = _arr[i];
                j = j - 1;
            }
            _arr[i] = temp;
        }
        _arr[i] = temp;
        return i;
    }
}
