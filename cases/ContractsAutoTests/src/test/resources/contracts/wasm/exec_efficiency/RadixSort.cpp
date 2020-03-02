#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-基数排序
 * @author liweic
 **/

CONTRACT RadixSort : public platon::Contract {
    private:
        platon::StorageType<"vecradix"_n, std::vector<int64_t>> vector_radix;
    public:
        ACTION void init(){}

        //基数排序：计算数组里最大位数，即为进行排序次数
        int maxbit(std::vector<int64_t>& a, int n)
        {
	        int d = 1; //保存最大的位数
	        int p = 10;
	        for (int i = 0; i < n; i++)
	        {
		        while (a[i] >= p)
		        {
			        p *= 10;
			        d++;
		        }
	        }
	        return d;
        }

        std::vector<int64_t>& radixSort(std::vector<int64_t>& a, int n)
        {
	        int d = maxbit(a,n);
	        int *tmp = new int[n];
	        int count[10]; //计数器
	        int i, j, k;
	        int radix = 1;
	        for (i = 1; i <= d; i++) //进行d次排序
	        {
		        for (j = 0; j < 10; j++)
			        count[j] = 0; //每次分配前清空计数器
		        for (j = 0; j < n; j++)
		        {
			        k = (a[j] / radix) % 10; //统计每个桶中的记录数
			        count[k]++;
		        }
		        for (j = 1; j < 10; j++)
			        count[j] = count[j - 1] + count[j]; //这时count里的值表示在tmp中的位置（减一为tmp里的存储下标）
		        for (j = n - 1; j >= 0; j--) //将所有桶中记录依次收集到tmp中
		        {
			        k = (a[j] / radix) % 10;
			        tmp[count[k] - 1] = a[j];
			        count[k]--;			//用于一个桶中有多个数，减一为桶中前一个数在tmp里的位置
		        }
		        for (j = 0; j < n; j++) //将临时数组的内容复制到data中
			        a[j] = tmp[j];
		        radix = radix * 10;
	        }
	        delete[] tmp;
            return a;
        }

        ACTION void sort(std::vector<int64_t>& arr, int n) {
            vector_radix.self() = std::move(radixSort(arr, n));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_radix.self();
        }

};
PLATON_DISPATCH(RadixSort,(init)(sort)(get_array))