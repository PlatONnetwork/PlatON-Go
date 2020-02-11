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
		
		/// 获取消息携带的value(fix) 
		/// define: u256 platon_call_value();
		CONST std::string value() { 
			u256 val = platon::platon_call_value();
			return to_string(val);		
		}

		/// sha3操作
		/// define: h256 platon_sha3(const bytes& data);
		CONST std::string sha3(const std::string& str) {
			bytes data;
			data.insert(data.begin(), str.begin(), str.end());
			h256 hash = platon::platon_sha3(data);
			return hash.toString();
		} 

		/// 设置函数返回值
		/// define: template <typename T> void platon_return(const T& t);
		CONST void rreturn() {
			std::string str = "hello";
			platon::platon_return(str);
		}

		/// 终止交易 panic, 交易完成，合约执行异常
		/// define: void platon_panic();
		ACTION void panic() {
			platon::platon_panic();		
		}

		/// 终止交易 revert
		/// define: void platon_revert();
		ACTION void revert(int64_t flag) {
			if(flag == 1){
				platon::platon_revert();
			}		
		} 

		/// 合约销毁 destroy, 销毁后检测余额
		/// define: bool platon_destroy(const Address& addr);
		ACTION void destroy(const std::string& addr) {
			platon::platon_destroy(Address(addr));		
		}
		
		/// 消息的原始发送者origin
		/// define: Address platon_origin();
		CONST std::string origin() {
			Address ori = platon::platon_origin();
			return ori.toString();		
		}

		/// compile test
		/// summary: compile success.
		std::string compile(){
			return "compile";		
		}
	
};

PLATON_DISPATCH(InnerFunction, (init)(compile)(gas_price)(block_number)(gas_limit)(timestamp)(gas)(nonce)(block_hash)(coinbase)(transfer)(value)(sha3)(rreturn)(panic)(revert)(destroy)(origin))


