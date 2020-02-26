#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-空间复杂度
 * @author qcxiao
 **/
CONTRACT QuickSort : public platon::Contract {
    private:
        platon::StorageType<"vector1"_n, std::vector<clothes>> vector_clothes;
    public:
        ACTION void init(){}
        std::vector<clothes> quickSort(std::vector<clothes>& array, int start, int last)
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

        ACTION void sort(std::vector<clothes>& arr, int8_t start, int8_t last) {
            vector_clothes.self() = quickSort(arr, start, last);
        }

        CONST std::vector<clothes> get_array() {
            return vector_clothes.self();
        }

};
PLATON_DISPATCH(SpaceComplexity,(init)(sort)(get_array))