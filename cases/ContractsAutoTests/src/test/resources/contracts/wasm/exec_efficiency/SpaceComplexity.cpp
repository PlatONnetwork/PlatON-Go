#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-空间复杂度
 * @author qcxiao
 **/

CONTRACT SpaceComplexity : public platon::Contract {

    private:
	    platon::StorageType<"test"_n, uint64_t> sum;
    public:
        ACTION void init(){}

        int partition(int a[], int p, int r) {
          int key = a[r];//取最后一个
          int i = p - 1;
          for (int j = p; j < r; j++)
          {
            if (a[j] <= key)
            {
              i++;
              //i一直代表小于key元素的最后一个索引，当发现有比key小的a[j]时候，i+1 后交换
              exchange(&a[i], &a[j]);
            }
          }
          exchange(&a[i + 1], &a[r]);//将key切换到中间来，左边是小于key的，右边是大于key的值。
          return i + 1;
        }

        void quickSort(int a[], int p, int r) {
          int position = 0;
          if (p<r)
          {
            position = partition(a,p,r);//返回划分元素的最终位置
            quickSort(a,p,position-1);//划分左边递归
            quickSort(a, position + 1,r);//划分右边递归
          }
        }

        int *generateRandomArray(int n, int rangeL, int rangeR) {
        	assert(rangeL <= rangeR);

        	int *arr = new int[n]; // 创建一个 n个元素的数组
        	
        	srand(time(NULL)); // 随机种子
        	for (int i = 0; i < n; i++)
        	    arr[i] = rand() % (rangeR - rangeL + 1) + rangeL;
        	return arr;
        }

        ACTION void sort(uint64_t n) {
            int a[] = generateRandomArray(n, 0, n*n);
            quickSort(a, 0, n);
        }

};
PLATON_DISPATCH(SpaceComplexity,(init)(sort))
