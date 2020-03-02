#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-堆排序
 * @author liweic
 **/

CONTRACT HeapSort : public platon::Contract {
    private:
        platon::StorageType<"vecheap"_n, std::vector<int64_t>> vector_heap;
    public:
        ACTION void init(){}

        void MaxSort(std::vector<int64_t>& a, int i, int n)
        {
            int j = 2*i+1;
            int temp = a[i];
            while(j < n)
            {
                if(j+1 <n && a[j] < a[j+1])
                    ++j;
                if(temp > a[j])
                    break;
                else
                {
                    a[i] = a[j];
                    i = j;
                    j = 2*i+1;
                }
            }
            a[i] = temp;
        }

        std::vector<int64_t>& heapSort(std::vector<int64_t>& a, int n)
        {
            for(int i= n/2-1;i>=0;i--)//从最后一个结点的父结点开始“向前遍历”
                MaxSort(a,i,n);
            for(int i=n-1;i>=1;i--)
            {
                MaxSort(a,0,i);
            }//逆序
            return a;
        }

        ACTION void sort(std::vector<int64_t>& arr, int n) {
            vector_heap.self() = std::move(heapSort(arr, n));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_heap.self();
        }

};
PLATON_DISPATCH(HeapSort,(init)(sort)(get_array))