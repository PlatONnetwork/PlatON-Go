#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;
/**
 * 执行效率-快速排序
 * @author qcxiao
 **/
CONTRACT QuickSort : public platon::Contract {
    private:
       platon::StorageType<"vector1"_n, std::vector<int64_t>> vector_clothes;
    public:
       ACTION void init(){}
       std::vector<int64_t>& quickSort(std::vector<int64_t>& array, long start, long last)
       {
            long i = start;
            long j = last;
            long temp = array[i];
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

       ACTION void sort(std::vector<int64_t>& arr, int64_t start, int64_t last) {
           vector_clothes.self() = std::move(quickSort(arr, start, last));
       }

       CONST std::vector<int64_t> get_array() {
           return vector_clothes.self();
       }

};
PLATON_DISPATCH(QuickSort,(init)(sort)(get_array))