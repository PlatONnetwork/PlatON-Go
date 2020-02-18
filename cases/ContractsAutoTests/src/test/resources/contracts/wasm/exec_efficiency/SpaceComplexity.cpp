

#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-空间复杂度
 * @author qcxiao
 **/


CONTRACT SpaceComplexity : public platon::Contract {

    private:
        platon::StorageType<"storage_array_int64"_n, std::array<int64_t,10>> storage_array_int64;
    public:
        ACTION void init(){}

        void quickSort(std::array<int64_t, 10>& array, int start, int last)
        {
            int i = start;
            int j = last;
            int temp = array[i];
            if (i < j)
            {
                while (i < j)
                {
                    //
                    while (i < j &&  array[j]>=temp )
                        j--;
                    if (i < j)
                    {
                        array[i] = array[j];
                        i++;
                    }

                    while (i < j && temp > array[i])
                        i++;
                    if (i < j)
                    {
                        array[j] = array[i];
                        j--;
                    }

                }
                //把基准数放到i位置
                array[i] = temp;
                //递归方法
                quickSort(array, start, i - 1);
                quickSort(array, i + 1, last);
            }
        }

        ACTION void sort(std::array<int64_t, 10> arr, int64_t start, int64_t end) {
            storage_array_int64.self() = arr;
            quickSort(storage_array_int64.self(), start, end);
        }

        CONST std::array<int64_t, 10> get_array() {
            return storage_array_int64.self();
        }

};
PLATON_DISPATCH(SpaceComplexity,(init)(sort)(get_array))


