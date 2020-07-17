#define TESTNET
#include <platon/platon.hpp>
using namespace platon;

class message {
   public:
      std::string head;
      PLATON_SERIALIZE(message, (head))
};


CONTRACT ContractMigrateTypes: public platon::Contract {
	public:
		PLATON_EVENT1(transfer,std::string,std::string)
		ACTION void init() {
		  
    }

		ACTION std::string migrate_contract(const bytes &init_arg, uint64_t transfer_value, uint64_t gas_value) {
        DEBUG("init_arg is :", toHex(init_arg));
        Address return_address;
        platon_migrate_contract(return_address, init_arg, transfer_value, gas_value);
        PLATON_EMIT_EVENT1(transfer,return_address.toString(),return_address.toString());
        DEBUG("new contract address:", return_address.toString());
        return return_address.toString();
    }

		ACTION void setMessage(const message &msg) {
      sMessage.self() = msg;
		}
		CONST message getMessage() {
      return sMessage.self();
		}

		ACTION void pushVector(uint16_t element) {
		  sVector.self().push_back(element);
		}
		CONST uint16_t getVectorElement(uint64_t index) {
		  return sVector.self()[index];
		}

    ACTION void setMap(std::string key, std::string value) {
      sMap.self()[key] = value;
		}
		CONST std::string getMapElement(std::string key) {
		  return sMap.self()[key];
		}

	private:
	  platon::StorageType<"message"_n, message> sMessage;
		platon::StorageType<"vectorvar"_n, std::vector<uint16_t>> sVector;
		platon::StorageType<"mapvar"_n,std::map<std::string,std::string>> sMap;
		
};

PLATON_DISPATCH(ContractMigrateTypes,(init)(migrate_contract)(setMessage)(getMessage)
	(pushVector)(getVectorElement)(setMap)(getMapElement))

