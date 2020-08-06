#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-插入排序
 * @author qcxiao
 **/

CONTRACT InsertSort : public platon::Contract {

    private:
            platon::StorageType<"vector1"_n, std::vector<int64_t>> vector_clothes;
    public:
        ACTION void init(){}

        std::vector<int64_t> insertSort(std::vector<int64_t>& arr, long n)
        {
        	long i,j,key;
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

        ACTION void sort(std::vector<int64_t>& arr, int64_t length) {
            vector_clothes.self() = insertSort(arr, length);
        }

        CONST std::vector<int64_t> get_array() {
            return vector_clothes.self();
        }
};
PLATON_DISPATCH(InsertSort,(init)(sort)(get_array))
