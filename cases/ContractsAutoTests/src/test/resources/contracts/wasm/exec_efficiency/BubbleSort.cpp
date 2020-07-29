#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-冒泡排序
 * @author liweic
 **/

CONTRACT BubbleSort : public platon::Contract {
    private:
        platon::StorageType<"vecbubble"_n, std::vector<int64_t>> vector_bubble;
    public:
        ACTION void init(){}

        std::vector<int64_t>& bubbleSort(std::vector<int64_t>& a, int n)
        {
	        for (int i = 0; i < n - 1; i++)
            {
		        for (int j = 0; j<n - 1 - i; j++)
		        {
			        if (a[j]>a[j + 1])
			        {
				        int temp;
				        temp = a[j];
				        a[j] = a[j + 1];
				        a[j + 1] = temp;
			        }
		        }
	        }
            return a;
        }

        ACTION void sort(std::vector<int64_t>& arr, int n) {
            vector_bubble.self() = std::move(bubbleSort(arr, n));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_bubble.self();
        }

};
PLATON_DISPATCH(BubbleSort,(init)(sort)(get_array))