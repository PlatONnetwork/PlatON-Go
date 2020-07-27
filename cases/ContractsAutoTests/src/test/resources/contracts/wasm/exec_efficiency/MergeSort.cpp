#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-归并排序
 * @author liweic
 **/

CONTRACT MergeSort : public platon::Contract {
    private:
        platon::StorageType<"vecmerge"_n, std::vector<int64_t>> vector_merge;
    public:
        ACTION void init(){}

        void merge(std::vector<int64_t>& a, int start, int mid, int end)
        {
            int *tmp = new int[end - start + 1];    // tmp是汇总2个有序区的临时区域
            int i = start;            // 第1个有序区的索引
            int j = mid + 1;        // 第2个有序区的索引
            int k = 0;                // 临时区域的索引

            while (i <= mid && j <= end)
            {
                if (a[i] <= a[j])
                    tmp[k++] = a[i++];
                else
                    tmp[k++] = a[j++];
            }
            while (i <= mid)
                tmp[k++] = a[i++];

            while (j <= end)
                tmp[k++] = a[j++];

            // 将排序后的元素，全部都整合到数组a中。
            for (i = 0; i < k; i++)
                a[start + i] = tmp[i];

            delete[] tmp;
        }

        std::vector<int64_t>& mergeSort(std::vector<int64_t>& a, int start, int end)
        {
            if (start >= end)
                return a;

            int mid = (end + start) / 2;
	        mergeSort(a, start, mid);		// 递归排序a[start...mid]
	        mergeSort(a, mid + 1, end);		// 递归排序a[mid+1...end]

            // a[start...mid] 和 a[mid...end]是两个有序空间，
            // 将它们排序成一个有序空间a[start...end]
            merge(a, start, mid, end);
            return a;
        }

        ACTION void sort(std::vector<int64_t>& arr, int start, int end) {
            vector_merge.self() = std::move(mergeSort(arr, start, end));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_merge.self();
        }

};
PLATON_DISPATCH(MergeSort,(init)(sort)(get_array))