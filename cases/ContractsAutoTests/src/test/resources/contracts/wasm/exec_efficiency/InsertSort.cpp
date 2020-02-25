#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-插入排序
 * @author qcxiao
 **/

CONTRACT InsertSort : public platon::Contract {

    private:
	    platon::StorageType<"arrayint"_n, std::array<int8_t,10>> arrayint;
    public:
        ACTION void init(){}

        std::array<int8_t, 10> insertSort(std::array<int8_t,10>& arr, int n)
        {
        	int i,j,key;
            for(i = 1; i < n; i++)
            {
                key = arr[i];
                j = i - 1;
                while(j >= 0 && arr[j] > key)
                {
                    arr[j+1] = arr[j];
                    j--;
                }
                arr[j+1] = key;
            }
            return arr;
        }

        ACTION void sort(std::array<int8_t, 10>& arr, int8_t length) {
            arrayint.self() = insertSort(arr, length);
        }

        CONST std::array<int8_t, 10> get_array() {
            return arrayint.self();
        }
};
PLATON_DISPATCH(InsertSort,(init)(sort)(get_array))
