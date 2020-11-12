#define TESTNET
#include <platon/platon.hpp>
using namespace platon;

CONTRACT ContractMigrateOld: public platon::Contract {
	public:
		PLATON_EVENT1(transfer,std::string,std::string)
		ACTION void init(uint8_t input) {
		  tUint8.self() = input;
    }

		ACTION std::string migrate_contract(const bytes &init_arg, uint64_t transfer_value, uint64_t gas_value) {
        DEBUG("init_arg is :", toHex(init_arg));
        Address return_address;
        platon_migrate_contract(return_address, init_arg, transfer_value, gas_value);
        PLATON_EMIT_EVENT1(transfer,return_address.toString(),return_address.toString());
        DEBUG("new contract address:", return_address.toString());
        return return_address.toString();
    }

		ACTION void setUint8(uint8_t input)
		{
			tUint8.self() = input;
		}
		CONST uint8_t getUint8()
		{
			return tUint8.self();
		}

	private:
	  platon::StorageType<"suint"_n, uint8_t> tUint8;
};

PLATON_DISPATCH(ContractMigrateOld,(init)(migrate_contract)(setUint8)(getUint8))

