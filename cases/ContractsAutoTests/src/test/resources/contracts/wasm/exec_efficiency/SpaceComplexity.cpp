#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-空间复杂度
 * @author qcxiao
 **/
CONTRACT SpaceComplexity : public platon::Contract {
    private:
        platon::StorageType<"arrayint"_n, std::array<int8_t,10>> arrayint;
    public:
        ACTION void init(){}
        std::array<int8_t, 10> quickSort(std::array<int8_t, 10>& array, int start, int last)
        {
            int i = start;
            int j = last;
            int temp = array[i];
            if (i < j)
            {
                while (i < j)
                {
                    while (i < j && array[j]>=temp )
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
                array[i] = temp;
                quickSort(array, start, i - 1);
                quickSort(array, i + 1, last);
            }
            return array;
        }

        ACTION void sort(std::array<int8_t, 10>& arr, int8_t start, int8_t last) {
            arrayint.self() = quickSort(arr, start, last);
        }

        CONST std::array<int8_t, 10> get_array() {
            return arrayint.self();
        }

};
PLATON_DISPATCH(SpaceComplexity,(init)(sort)(get_array))