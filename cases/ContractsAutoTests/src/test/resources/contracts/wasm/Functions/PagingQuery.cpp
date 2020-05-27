#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace std;
using namespace platon;

/**
* vector分页查询的实现测试
* @author liweic
*/

CONTRACT PagingQuery : public platon::Contract{

    private:
      platon::StorageType<"vecstorage"_n, std::vector<std::string>> storage_vector_string;

    public:
        ACTION void init(){}

        ACTION void insertVectorValue(const std::string &my_value)
     {
        storage_vector_string.self().push_back(my_value);
     }

     //vector大小
        CONST uint64_t getVectorSize()
     {
       return storage_vector_string.self().size();
     }

        CONST std::string getPagingQuery(uint64_t CurrentPage, uint64_t PageMaxSize)
     {
        int vecSize = 0;
        int pages = 0;
        vecSize = storage_vector_string.self().size();
        pages = vecSize / PageMaxSize;
        if(0 != vecSize % PageMaxSize)
        {
            pages += 1;
        }

        if(CurrentPage <= 0)
        {
        	return "";
        }
        else if(CurrentPage > pages)
        {
            return "";
        }

        // 计算当前页对应的开始与结束index
        int nStartIndex = (CurrentPage - 1) * PageMaxSize;
        int nEndIndex = nStartIndex + PageMaxSize - 1;

        // 超出范围,index为最后一个元素的下标
        if(nEndIndex >= vecSize)
        {
        	nEndIndex = vecSize - 1;
        }

        std::string strVecInfo = "{";
        strVecInfo += "\"";
        strVecInfo += "PageTotal";
        strVecInfo += "\":";
        strVecInfo += to_string(pages);
        strVecInfo += ",";
        strVecInfo += "\"";
        strVecInfo += "Data";
        strVecInfo += "\":[";

        for(int i = nStartIndex; i <= nEndIndex; i++)
        {
            std::string strTmp = "";
            strTmp = storage_vector_string.self().at(i);
            strVecInfo += strTmp + ",";
        }

        strVecInfo.pop_back();
        strVecInfo += "]}";

        char* buf = (char*)malloc(strVecInfo.size() + 1);
        memset(buf, 0, strVecInfo.size()+1);
        strcpy(buf, strVecInfo.c_str());
        return buf;
     }

};

PLATON_DISPATCH(PagingQuery, (init)(insertVectorValue)(getVectorSize)(getPagingQuery))
