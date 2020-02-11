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

		
	
};

PLATON_DISPATCH(InnerFunction, (init)(gas_price)(block_number)(gas_limit)(timestamp)(gas)(nonce)(block_hash)(coinbase)(transfer)(value)(sha3)(rreturn)(panic)(revert)(destroy)(origin))


