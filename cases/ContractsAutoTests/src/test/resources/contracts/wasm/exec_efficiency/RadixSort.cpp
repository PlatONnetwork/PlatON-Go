#define TESTNET
#include <platon/platon.hpp>
#include <string>
#include<list>
using namespace platon;


/**
 * 执行效率-基数排序 支持正数排序
 * @author liweic
 **/

CONTRACT RadixSort : public platon::Contract {
    private:
        platon::StorageType<"vecradix"_n, std::vector<int64_t>> vector_radix;
    public:
        ACTION void init(){}

		int maxdigiit(std::vector<int64_t>& arr,int n)
		{
			int d = 1;
			int p = 10;
			for (int i = 0; i < n; ++i)
			{
				while (arr[i] >= p)
				{
					p *= 10;
					++d;
				}
			}
			return d;
		}


        std::vector<int64_t>& radixSort(std::vector<int64_t>& arr, int n)
        {
	        int digits = maxdigiit(arr,n);
			std::list<int> lists[10];
			int d,j,k,factor;
			for ( d = 1,factor=1; d <= digits;factor*=10, d++)
			{
				for ( j = 0; j < n; j++)
				{
					lists[(arr[j] / factor) % 10].push_back(arr[j]);
				}
				for (j = k = 0; j < 10; j++)
				{
					while (!lists[j].empty())
					{
						arr[k++] = lists[j].front();
						lists[j].pop_front();
					}
				}
			}
            return arr;
        }

        ACTION void sort(std::vector<int64_t>& arr, int n) {
            vector_radix.self() = std::move(radixSort(arr, n));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_radix.self();
        }

};
PLATON_DISPATCH(RadixSort,(init)(sort)(get_array))