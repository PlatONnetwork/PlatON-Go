// Author: zjsunzone
#include <platon/platon.hpp>
#include <string>

#include "src/Ownable.hpp"
using namespace platon;


CONTRACT Bank: public platon::Contract, public Ownable
{

	// define event
	public:
		bool onlyBagholders(){	
			return false;				
		}
		bool onlyStronghands(){
			return false;		
		}

	public:
		PLATON_EVENT1(onTokenPurchase, Address, u128, u128, Address, u128, u128);
		PLATON_EVENT1(onTokenSell, Address, u128, u128, u128, u128);
		PLATON_EVENT1(onReinvestment, Address, u128, u128);
		PLATON_EVENT1(onWithdraw, Address, u128);
		PLATON_EVENT1(Transfer, Address, Address, u128);

	public:
		platon::StorageType<"name"_n, std::string> name;
		platon::StorageType<"symbol"_n, std::string> symbol;
		platon::StorageType<"decimals"_n, uint8_t> decimals;
		platon::StorageType<"entryFee"_n, uint8_t> entryFee_;
		platon::StorageType<"transferFee_"_n, uint8_t> transferFee_;
		platon::StorageType<"exitFee"_n, uint8_t> ExitFee_;
		platon::StorageType<"refferalFee_"_n, uint8_t> refferalFee_;
		platon::StorageType<"DevFee_"_n, uint8_t> DevFee_;
		platon::StorageType<"DailyINterest_"_n, uint8_t> DailyInterest_;
		platon::StorageType<"IntFee_"_n, uint8_t> IntFee_;
		platon::StorageType<"InterestPool_"_n, u128> InterestPool_;
		platon::StorageType<"TokenPriceInitial_"_n, u128> tokenPriceInitial_;
		platon::StorageType<"tokenPriceIncremental_"_n, u128> tokenPriceIncremental_;
		platon::StorageType<"magnitude"_n, u128> magnitude;
		platon::StorageType<"stakingRequirement"_n, u128> stakingRequirement;
		
		// map 
		platon::StorageType<"tokenBalanceLedger_"_n, std::map<Address, u128>> tokenBalanceLedger_;
		platon::StorageType<"referralBalance_"_n, std::map<Address, u128>> referralBalance_;
		platon::StorageType<"payoutsTo_"_n, std::map<Address, u128>> payoutsTo_;

		platon::StorageType<"tokenSupply_"_n, u128> tokenSupply_;
		platon::StorageType<"profitPerShare_"_n, u128> profitPerShare_;
		platon::StorageType<"dev"_n, Address> dev;

	public:
		ACTION void init()
		{
			name.self() = "Cypher Bank";
			symbol.self() = "CBT";
			decimals.self() = 18;
			entryFee_.self() = 15;
			transferFee_.self() = 1;
			ExitFee_.self() = 20;
			refferalFee_.self() = 8;
			DevFee_.self() = 25;
			DailyInterest_.self() = 1;
			IntFee_.self() = 25;
			InterestPool_.self() = 0;
			tokenPriceInitial_.self() = u128("100000000000");
			tokenPriceIncremental_.self() = u128("10000000000");
			magnitude.self() = u128(2);			// 2**64
			stakingRequirement.self() = u128("50"); 	// 50e18

			// 
			dev.self() = Address(""); // setting.
		}

		ACTION void buy(Address _referredBy) {
			u128 callValue = platon_call_value();
			u128 DevFee1 = callValue / u128(100) * DevFee_.self();	
			u128 DevFeeFinal = DevFee1 / u128(10);
			platon_transfer(dev.self(), Energon(DevFeeFinal));
			
			//
			u128 DailyInt1 = callValue/ u128(100) * IntFee_.self();
			u128 DailyIntFinal = DailyInt1 / u128(10);
			InterestPool_.self() += DailyIntFinal;
			purchaseTokens(callValue, _referredBy);
		}

		ACTION void IDD(){
			Address sender = platon_caller();
			if(sender != owner.self()){
				platon_revert();			
			}		
			Energon cbalance = platon_balance(platon_address());
			u128 Contract_Bal = cbalance.Get() - InterestPool_.self();
			u128 DailyInterest1 = Contract_Bal * InterestPool_.self() / u128(100);
			u128 DailyInterestFinal = DailyInterest1 / u128(10);
			InterestPool_.self() -= DailyInterestFinal;
			//
			DividendsDistribution(DailyInterestFinal, Address("0x0"));
		}

		ACTION void DivsAddon(){
			u128 callValue = platon_call_value();
			DividendsDistribution(callValue, Address("0x0"));
		}
		
		ACTION void reinvest() {
			u128 _dividends = myDividends(false);
			Address _customerAddress = platon_caller();
			payoutsTo_.self()[_customerAddress] += _dividends*magnitude.self();
			_dividends += referralBalance_.self()[_customerAddress];
			referralBalance_.self()[_customerAddress] = 0;
			u128 _tokens = purchaseTokens(_dividends, Address("0x0"));
			PLATON_EMIT_EVENT1(onReinvestment, _customerAddress, _dividends, _tokens);
		}

		ACTION void exit(){
			Address _customerAddress = platon_caller();
			u128 _tokens = tokenBalanceLedger_.self()[_customerAddress];
			if (_tokens > 0){
				sell(_tokens);			
			}
			withdraw();
		}

		ACTION void withdraw(){
			onlyStronghands();		
			Address _customerAddress = platon_caller();
			u128 _dividends = myDividends(false);
			//payoutsTo_.self()[_customerAddress] += 
		}
		
		ACTION void sell(u128 _amountOfTokens){
			onlyBagholders();		
		}
		
		ACTION void transfer(Address _toAddress, u128 _amountOfTOkens){
			onlyBagholders();		
		}
	
		CONST u128 totalEthereumBalance(){
			return platon_balance(platon_address()).Get();		
		}
		
		CONST u128 totalSupply() {
			return tokenSupply_.self();		
		}
		
		CONST u128 myTokens(){
			Address _customerAddress = platon_caller();
			return balanceOf(_customerAddress);		
		}

		CONST u128 myDividends(bool _includeReferralBonus){
			return u128(0);
		}

		CONST u128 balanceOf(Address _customerAddress) {
			return tokenBalanceLedger_.self()[_customerAddress];		
		}

		CONST u128 dividendsOf(Address _customerAddress){
			return u128(0);		
		}
		
		CONST u128 sellPrice(){
			return u128(0);		
		}
		
		CONST u128 buyPrice(){
			return u128(0);		
		}

		CONST u128 calculateTokensReceived(u128 _ethereumToSpend){
			return u128(0);	
		}

		CONST u128 calculateEthereumReceived(u128 _tokensToSell){
			return u128(0);			
		}

		CONST uint8_t exitFee(){
			return ExitFee_.self();		
		}

		u128 purchaseTokens(u128 _incomingEthereum, Address _referredBy) {
			return u128(0);		
		}
		
		u128 DividendsDistribution(u128 _incomingEthereum, Address _referredBy){
			Address _customerAddress = platon_caller();
			u128 _undividedDividends = (_incomingEthereum * u128(100)) / u128(100);
			u128 _referralBonus = (_undividedDividends * refferalFee_.self()) / u128(100);
			u128 _dividends = _undividedDividends - _referralBonus;
			u128 _taxedEthereum = _incomingEthereum - _undividedDividends;
			u128 _amountOfTokens = ethereumToTokens_(_taxedEthereum);
			u128 _fee = _dividends * magnitude.self();
			
			if(_amountOfTokens >= 0 && (_amountOfTokens + tokenSupply_.self()) >= tokenSupply_.self()){
				platon_revert();			
			}

			//
			if(_referredBy != Address("0x0000000000000000000000000000000000000000") &&
				_referredBy != _customerAddress &&
				tokenBalanceLedger_.self()[_referredBy] >= stakingRequirement.self()			
			){
				referralBalance_.self()[_referredBy] = referralBalance_.self()[_referredBy] + _referralBonus;
			}
			return u128(0);
		}
		
		u128 ethereumToTokens_(u128 _ethereum) {
			return u128(0);		
		}
		
		u128 tokensToEthereum(u128 _tokens) {
			u128 tokens_ = _tokens + u128("1000000000000000000");
			u128 _tokenSupply = tokenSupply_.self() + u128("1000000000000000000");
			u128 _etherReceived = ((tokenPriceInitial_.self() + (tokenPriceIncremental_.self() * (_tokenSupply / u128("1000000000000000000")))) * (tokens_ - u128("1000000000000000000")) - (tokenPriceIncremental_.self() * ((tokens_*u128(4) - tokens_) / u128("1000000000000000000")) / u128(2))) / u128("1000000000000000000");
			return _etherReceived;		
		}
		
		u128 sqrt(u128 x) {
			u128 z = (x + u128(1)) / u128(2);
			u128 y = x;
			while(z < y){
				y = z;
				z = (x / z + z) / u128(2);			
			}
			return y;			
		}











};

PLATON_DISPATCH(Bank,(init))

















