#define TESTNET
#include <platon/platon.hpp>
using namespace platon;

CONTRACT SimpleStorage: public platon::Contract
{
	public:
		ACTION void init(){}
	
		ACTION void set(uint64_t input)
		{
			storedData.self() = input;		
		}
		
		CONST uint64_t get()
		{
			return storedData.self();
		}

		ACTION void set_address(const Address &addr)
		{
			migrateAddress.self() = addr;		
		}
		
		CONST Address get_address()
		{
			return migrateAddress.self();
		}

	private:
		platon::StorageType<"sstored"_n, uint64_t> storedData;
		platon::StorageType<"migrate_address"_n, Address> migrateAddress;
};

PLATON_DISPATCH(SimpleStorage,(init)(set)(get)(set_address)(get_address))