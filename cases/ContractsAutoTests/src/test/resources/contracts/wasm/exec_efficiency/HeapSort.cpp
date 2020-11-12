#define TESTNET
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

        void MaxSort(std::vector<int64_t>& a, int root, int n)
        {
            int parent = root;
            int child = parent*2+1; // 左孩子
            while( child < n )
            {
                if( (child+1) < n && a[child+1] > a[child] )
                {
                    ++child;
                }
                if( a[child] > a[parent] )
                {
                    std::swap(a[child],a[parent]);
                    parent = child;
                    child = parent*2+1;
                }
                else
                    break;
            }
        }

        std::vector<int64_t>& heapSort(std::vector<int64_t>& a, int n)
        {
            assert(a);
            for( int i = (n-2)/2; i >=0 ; i-- )
            {
                MaxSort(a,i,n);
            }

            int end = n-1;
            while( end > 0 ){
                std::swap(a[0],a[end]);
                MaxSort(a,0,end); // end其实就是不算后面的一个元素，原因是最后一个节点已经是最大的
                end--;
            }
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