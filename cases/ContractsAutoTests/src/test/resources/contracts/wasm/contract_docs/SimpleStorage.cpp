// Author: zjsunzone
// 简单的存储合约
#include <platon/platon.hpp>
#include <string>
using namespace platon;


CONTRACT SimpleStorage: public platon::Contract
{
	public:
		ACTION void init()
		{
			
		}
	
		ACTION void set(uint64_t input)
		{
			storedData.self() = input;		
		}
		
		CONST uint64_t get()
		{
			return storedData.self();
		}

	private:
		platon::StorageType<"suint64"_n, uint64_t> storedData;
};

PLATON_DISPATCH(SimpleStorage,(init)(set)(get))



