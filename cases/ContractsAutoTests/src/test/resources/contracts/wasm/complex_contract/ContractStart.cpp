#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

/**
 * https://github.com/sun-asterisk-research/EOS-2048/blob/master/contracts/eos/mycontract/main.cpp
 *
 */

struct dummy_action_hello {
  std::string vaccount;
  uint64_t b;
  uint64_t c;
  dummy_action_hello(){}

  PLATON_SERIALIZE( dummy_action_hello, (vaccount)(b)(c) )
};

struct shardbucket {
   std::vector<byte> shard_uri;
   uint64_t shard;
   uint64_t primary_key() const { return shard; }

   PLATON_SERIALIZE( shardbucket, (shard_uri)(shard))
};

struct TKEY{
    std::string pubkey;
    std::string vaccount;
    uint64_t nonce;
    uint64_t primary_key;

    PLATON_SERIALIZE( TKEY, (pubkey)(vaccount)(nonce)(primary_key))
};


CONTRACT ContractStart : public platon::Contract{

    private:
        platon::StorageType<"stat"_n, uint64_t> stat;
        platon::StorageType<"vaccounts"_n, shardbucket> cold_accounts_t_abi;

    public:
        ACTION void init(){
        }

        ACTION bool Timer_callback(std::vector<byte> payload,uint64_t seconds){
            if(stat.self() == 0 ){
                stat.self() = seconds;
            }

            auto reschedule = false;
            if(seconds++ < 3){
              reschedule = true;
            }
            stat.self() = reschedule;
            return reschedule;

        }

        ACTION void testschedule() {
            std::vector<byte> payload;

        }

        ACTION void hello(dummy_action_hello payload){
            platon_assert(payload.vaccount.length() > 0, "invalid vaccount");
            DEBUG("hello from ",payload.vaccount," ",payload.b,payload.c);

        }

        ACTION void hello2(dummy_action_hello payload){
            DEBUG("hello2(default action) from ",payload.vaccount," ",payload.b,payload.c);

        }

        ACTION void transfer(Address from,Address to,uint64_t quantity,std::string memo){
            platon_assert(from.toString().length() < 5, "invalid from");
            platon_assert(from.toString() == to.toString(), "from can not same as to");
            if(memo.size() > 0){
                Address to_act = Address(memo.c_str());
                from = to_act;
            }
            platon_assert(quantity > 0, "quantity can not less zero");
            add_cold_balance(from,quantity);

        }

        ACTION void add_cold_balance(Address owner,uint64_t value){
            platon_transfer(owner, Energon(value));
        }

};

PLATON_DISPATCH(ContractStart, (init)(Timer_callback)(testschedule)(hello)(hello2)(transfer)(add_cold_balance))

