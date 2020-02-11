#include <platon/platon.hpp>
#include <string>
using namespace platon;

/**
 * 针对系统链上函数的调用
 */
CONTRACT InnerFunction:public platon::Contract{
	public:
		ACTION void init() {}
		
		/// 获取GasPrice
		CONST uint64_t gas_price(){
			return platon::platon_gas_price();
		}

		/// 获取区块高度
		CONST uint64_t block_number(){
			return platon::platon_block_number();	
		}

		/// 获取gasLimit
		CONST uint64_t gas_limit() {
			return platon::platon_gas_limit();		
		}
		
		/// 获取当前交易发送的Gas
		CONST uint64_t gas() {
			return platon::platon_gas();		
		}
	
		/// 获取当前块的时间戳
		CONST int64_t timestamp() {
			return platon::platon_timestamp();		
		}

		/// 获取消息发送者的nonce
		CONST uint64_t nonce() {
			return platon::platon_caller_nonce();		
		}

		/// 获取指定区块高度的哈希
		CONST std::string block_hash(int64_t bn) {
			h256 bhash = platon::platon_block_hash(bn);
			return bhash.toString();	
		}
			
		/// 获取当前旷工地址
		CONST std::string coinbase() {
			return platon::platon_coinbase().toString();		
		}
		
		/// 获取指定地址的余额(bug)
		CONST std::string balanceOf(const std::string& addr) {
			Energon e = platon::platon_balance(Address(addr));
			return to_string(e.Get());		
		}

		/// 主币转账
		/// define: int32_t platon_transfer(const Address& addr, const Energon& amount);
		ACTION void transfer(const std::string& addr, uint64_t amount) {
			if(amount == 0){
				DEBUG("Transfer failed", "address", addr, "amount", amount);
				return;
			}		
			platon::platon_transfer(Address(addr), Energon(amount));
		}
		
		
	
};

PLATON_DISPATCH(InnerFunction, (init)(gas_price)(block_number)(gas_limit)(timestamp)(gas)(nonce)(block_hash)(coinbase)(transfer)(value)(sha3)(rreturn)(panic)(revert)(destroy)(origin))


