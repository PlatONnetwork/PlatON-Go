#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-希尔排序
 * @author liweic
 **/

CONTRACT ShellSort : public platon::Contract {
    private:
        platon::StorageType<"vecshell"_n, std::vector<int64_t>> vector_shell;
    public:
        ACTION void init(){}
        std::vector<int64_t>& shellSort(std::vector<int64_t>& a, int n)
        {
            int i, j, gap;
            // gap为步长，每次减为原来的一半。
            for (gap = n / 2; gap > 0; gap /= 2)
            {
                // 共gap个组，对每一组都执行直接插入排序
                for (i = 0; i < gap; i++)
                    {
                        for (j = i + gap; j < n; j += gap)
                        {
                            // 如果a[j] < a[j-gap]，则寻找a[j]位置，并将后面数据的位置都后移。
                            if (a[j] < a[j - gap])
                                {
                                    int tmp = a[j];
                                    int k = j - gap;
                                    while (k >= 0 && a[k] > tmp)
                                    {
                                        a[k + gap] = a[k];
                                        k -= gap;
                                    }
                                    a[k + gap] = tmp;
                                }
                        }
                    }
            }
            return a;
        }

        ACTION void sort(std::vector<int64_t>& arr, int n) {
            vector_shell.self() = std::move(shellSort(arr, n));
        }

        CONST std::vector<int64_t> get_array() {
            return vector_shell.self();
        }

};
PLATON_DISPATCH(ShellSort,(init)(sort)(get_array))