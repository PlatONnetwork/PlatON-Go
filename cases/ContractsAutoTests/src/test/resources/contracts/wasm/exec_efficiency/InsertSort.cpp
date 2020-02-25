#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 执行效率-插入排序
 * @author qcxiao
 **/

CONTRACT InsertSort : public platon::Contract {

    private:
	    platon::StorageType<"array_int8"_n, std::array<int8_t,10>> array_int8;
    public:
        ACTION void init(){}

        void swap(int &a, int &b)
        {
        	int temp = a;
        	a = b;
        	b = temp;
        }

        void insertSort(std::array<int8_t,10>& a, int length)
        {
        	for (int i = 1; i < length; i++)
        	{
        		for (int j = i - 1; j >= 0 && a[j + 1] < a[j]; j--)
        		{
        			swap(a[j], a[j + 1]);
        		}
        	}

        }

        ACTION void sort(std::array<int8_t, 10>& arr, int8_t length) {
            array_int8.self() = insertSort(arr, length);
        }

        CONST std::array<int8_t, 10> get_array() {
            return array_int8.self();
        }
};
PLATON_DISPATCH(InsertSort,(init)(sort)(get_array))
