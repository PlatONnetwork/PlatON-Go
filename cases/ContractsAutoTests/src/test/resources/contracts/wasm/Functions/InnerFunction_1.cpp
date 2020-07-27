#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 针对系统链上函数的调用
 */
CONTRACT InnerFunction_1:public platon::Contract{
	public:
		ACTION void init() {}
		
		/// 获取当前交易发送的Gas
		CONST uint64_t gas() {
			return platon_gas();		
		}

		/// 获取消息发送者的nonce
		CONST uint64_t nonce() {
			return platon_caller_nonce();		
		}

		/// 获取指定区块高度的哈希
		CONST std::string block_hash(uint64_t bn) {
			h256 bhash = platon_block_hash(bn);
			return bhash.toString();	
		}
			
		/// 获取当前旷工地址
		CONST Address coinbase() {
			return platon_coinbase();		
		}

		/// 获取指定地址的余额(bug)
		CONST std::string balanceOf(const std::string& addr) {
		    auto address_info = make_address(addr);
		    if(address_info.second){
		        Energon e = platon_balance(address_info.first);
                return std::to_string(e.Get());
		    }else{
		        return "";
		    }

		}
		
};


PLATON_DISPATCH(InnerFunction_1, (init)(gas)(nonce)(block_hash)(coinbase)(balanceOf))


