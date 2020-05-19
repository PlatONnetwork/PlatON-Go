#define TESTNET
// Author: zjsunzone
#include <platon/platon.hpp>
#include <string>
using namespace platon;


class Token {
	public:
		// total amount of tokens
		platon::StorageType<"totalsupply"_n, uint64_t> totalSupply;
	
	public: // event
		// define: _from, _to, _value
		PLATON_EVENT2(Transfer, Address, Address, uint64_t);
		// define: _owner, _spender, _value
		PLATON_EVENT2(Approval, Address, Address, uint64_t);

	public:
		// @param _owner The address from which the balance will be retrieved
		// @return The balance.
		virtual uint64_t balanceOf(Address _owner) = 0;

		// @notice send '_value' token to `_to` from `msg.sender`
		// @param _to THe address of the recipient.
 		// @param _value The amount of token to be transferred.
		// @return Whether the transfer was successful or not.
		virtual bool transfer(Address _to, uint64_t _value) = 0;

		// @notice send `_value` token to `_to` from `_from` on the condition it is approved by `_from`
		// @param _from The address of the sender.
		// @param _to The address of the recepient.
		// @param _value The amount of token to be transferred.
		// @return Whether the transfer was successful or not.
		virtual bool transferFrom(Address _from, Address _to, uint64_t _value) = 0;

		// @notice `msg.sender` approves `_spender` to spend `_value` tokens
		// @param _spender The address of the account able to transfer the tokens
		// @param _value The amount of tokens to be approved for transfer
		// @return Whether thee approval was successful or not.
		virtual bool approve(Address _spender, uint64_t _value) = 0;

		// @param _owner The address of the account owning tokens
		// @param _spender The address of the account able to transfer the tokens
		// @return Amount of remaining tokens allowed to spent.
		virtual uint64_t allowance(Address _owner, Address _spender) = 0;

};

// You should inherit from StandardTOken or, for a token like you would want
// to deploy in something like MIst, see HumanStandardToken.cpp.
// (This implements ONLY the standard functions and NOTHING else.
// If you deploy this, you won't have anthing useful.)
class StandardToken: public Token
{

	protected: 
		platon::StorageType<"balances"_n, std::map<Address, uint64_t>> balances;
		platon::StorageType<"allowed"_n, std::map<Address, std::map<Address, uint64_t>>> allowed;

	public:
		CONST uint64_t balanceOf(Address _owner) {
			return balances.self()[_owner];		
		}

		ACTION bool transfer(Address _to, uint64_t _value){
			// Default assumes totalSupply can't be over max(2^64 - 1)
			// If your token leaves out totalSupply and can issue more tokens as time goes on,
			// you need to check if it doesn't wrap.
			// Replace the if with this on instead.
			Address sender = platon_caller();
			if (balances.self()[sender] >= _value && _value > 0) {
				balances.self()[sender] -= _value;
				balances.self()[_to] += _value;
				PLATON_EMIT_EVENT2(Transfer, sender, _to, _value);
				return true;
			} else {
				return false;			
			}
		}

		ACTION bool transferFrom(Address _from, Address _to, uint64_t _value) {
			// same as above. Replace this line with the following if you want to protect against
			// wrapping uints.
			Address sender = platon_caller();
			if(balances.self()[_from] >= _value 
				&& allowed.self()[_from][sender] >= _value && _value > 0){
				balances.self()[_to] += _value;
				balances.self()[_from] -= _value;
				PLATON_EMIT_EVENT2(Transfer, _from, _to, _value);
				return true;
			} else {
				return false;			
			}
		}

		ACTION bool approve(Address _spender, uint64_t _value){
			Address sender = platon_caller();			
			allowed.self()[sender][_spender] = _value;
			PLATON_EMIT_EVENT2(Approval, sender, _spender, _value);
			return true;		
		}

		CONST uint64_t allowance(Address _owner, Address _spender){
			return allowed.self()[_owner][_spender];		
		}
		
};


CONTRACT LATToken: public platon::Contract, public StandardToken
{
	
	public:
		platon::StorageType<"name"_n, std::string> name;		// fancy name: eg PlatON Token
		platon::StorageType<"decimals"_n, uint8_t> decimals;	// HOw many decimals to show.
		platon::StorageType<"symbol"_n, std::string> symbol;	// An identifier: eg LTT
		platon::StorageType<"version"_n, std::string> version;	// 0.1 standard. Just an arbitrary versioning scheme.

	public:
		ACTION void init(uint64_t _initialAmount, const std::string& _tokenName,
			uint8_t _decimalUnits, const std::string& _tokenSymbol)
		{
			Address sender = platon_caller();
			balances.self()[sender] = _initialAmount;		// Give the creator all initial tokens.
			totalSupply.self() = _initialAmount;			// Update total supply.
			name.self() = _tokenName;						// Set the name for display purposes
			decimals.self() = _decimalUnits;				// Amount of decimals for display purposes
			symbol.self() = _tokenSymbol;					// Set the symbol for display purposes.
		}

		CONST std::string getName(){
			return name.self();		
		}

        CONST Address getSender(){
            return platon_caller();
        }

		CONST uint8_t getDecimals(){
			return decimals.self();		
		}

		CONST std::string getSymbol(){
			return symbol.self();		
		}

		CONST uint64_t getTotalSupply(){
			return totalSupply.self();		
		}
		
		// Approves and then calls the receiving contract.
		ACTION bool approveAndCall(Address _spender, uint64_t _value, const bytes& _extraData) {
			Address sender = platon_caller();
			allowed.self()[sender][_spender] = _value;
			PLATON_EMIT_EVENT2(Approval, sender, _spender, _value);
			// call the receiveApproval function on the contract you want to be notified. This 
			// crafts the function signature manually so one doesn't have to include a contract 
			// in here just for this.
			// define: receiveApproval(Address _from, uint64_t _value, Address _tokenContract, bytes& _extraDaa)
			// it is assumed that when does this that the call should succeed.
			return true;
		}
};

PLATON_DISPATCH(LATToken,(init)(balanceOf)(transfer)(transferFrom)(approve)(allowance)
(getName)(getDecimals)(getSymbol)(getTotalSupply)(approveAndCall)(getSender))



