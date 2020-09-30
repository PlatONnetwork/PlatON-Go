#define TESTNET
#include <platon/platon.hpp>
using namespace platon;

CONTRACT ContractMigrateNew: public platon::Contract {
		public:
		PLATON_EVENT1(transfer,std::string,std::string)
		ACTION void init(uint8_t input, uint16_t input2) {
		  tUint8.self() = input;
			tUint16.self() = input2;
    }

		ACTION void setUint8New(uint8_t input)
		{
			tUint8.self() = input;
		}
		CONST uint8_t getUint8New()
		{
			return tUint8.self();
		}

		ACTION void setUint16(uint16_t input)
		{
			tUint16.self() = input;
		}
		CONST uint16_t getUint16()
		{
			return tUint16.self();
		}

	private:
	  platon::StorageType<"suintone"_n, uint8_t> tUint8;
		platon::StorageType<"suinttwo"_n, uint16_t> tUint16;
};

PLATON_DISPATCH(ContractMigrateNew,(init)(setUint8New)(getUint8New)(setUint16)(getUint16))


