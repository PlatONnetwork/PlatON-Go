#define TESTNET
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
			if(myTokens() <= 0){
				return false;
			}
			return true;				
		}
		bool onlyStronghands(){
			if(myDividends(true) <= 0){
				return false;
			}
			return false;		
		}

	public:
		PLATON_EVENT1(onTokenPurchase, Address, u128, u128, Address, u128, u128);
		PLATON_EVENT1(onTokenSell, Address, u128, u128, u128, u128);
		PLATON_EVENT1(onReinvestment, Address, u128, u128);
		PLATON_EVENT1(onWithdraw, Address, u128);
		PLATON_EVENT1(Transfer, Address, Address, u128);
		PLATON_EVENT1(TestData, Address, u128, u128, u128);

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
		platon::StorageType<"initaddr"_n, Address> initaddr;

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
			tokenPriceInitial_.self() = u128(100000000000);
			tokenPriceIncremental_.self() = u128(10000000000);
			magnitude.self() = (18_LAT).Get();			// 2**64
			stakingRequirement.self() = (50_LAT).Get(); 	// 50e18

			// 
			//dev.self() = Address("0x493301712671Ada506ba6Ca7891F436D29185823"); // setting.
			auto address_info = make_address("lax10jc0t4ndqarj4q6ujl3g3ycmufgc77epxg02lt");
            if(address_info.second) dev.self() = address_info.first;

            auto address_init = make_address("laxqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq");
            if(address_init.second) Address initaddr = address_init.first;
		}

		ACTION void buy(Address _referredBy) {
			u128 callValue = platon_call_value();	
			u128 DevFee1 = callValue / u128(100) * DevFee_.self();	
			u128 DevFeeFinal = DevFee1 / u128(10);
			//PLATON_EMIT_EVENT1(TestData, platon_caller(), callValue, DevFee1, DevFeeFinal);			
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
			//DividendsDistribution(DailyInterestFinal, Address("0x0000000000000000000000000000000000000000"));
		}

		ACTION void DivsAddon(){
			u128 callValue = platon_call_value();
			//DividendsDistribution(callValue, Address("0x0000000000000000000000000000000000000000"));
			DividendsDistribution(callValue, initaddr.self());
		}
		
		ACTION void reinvest() {
			u128 _dividends = myDividends(false);
			Address _customerAddress = platon_caller();
			payoutsTo_.self()[_customerAddress] += _dividends*magnitude.self();
			_dividends += referralBalance_.self()[_customerAddress];
			referralBalance_.self()[_customerAddress] = 0;

			//u128 _tokens = purchaseTokens(_dividends, Address("0x0000000000000000000000000000000000000000"));
			u128 _tokens = purchaseTokens(_dividends, initaddr.self());
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
			payoutsTo_.self()[_customerAddress] += _dividends * magnitude.self();
			_dividends += referralBalance_.self()[_customerAddress];
			referralBalance_.self()[_customerAddress] = u128(0);
			platon_transfer(_customerAddress, Energon(_dividends));
			PLATON_EMIT_EVENT1(onWithdraw, _customerAddress, _dividends);
			
		}
		
		ACTION void sell(u128 _amountOfTokens){
			onlyBagholders();		
			Address _customerAddress = platon_caller();
			if(_amountOfTokens > tokenBalanceLedger_.self()[_customerAddress]){
				platon_revert();			
			}
			u128 _tokens = _amountOfTokens;
			u128 _ethereum = tokensToEthereum_(_tokens);
			u128 _dividends = _ethereum * exitFee() * u128(100);
			u128 _devexit = _ethereum * u128(5) * u128(100);
			u128 _taxedEthereum1 = _ethereum - _dividends;
			u128 _taxedEthereum = _taxedEthereum1 - _devexit;
			u128 _devexitindividual = _ethereum * DevFee_.self() / u128(100);
			u128 _devexitindividual_final = _devexitindividual / u128(10);
			u128 DailyInt1 = _ethereum * IntFee_.self() / u128(100);
			u128 DailyIntFinal = DailyInt1 / u128(10);
			InterestPool_.self() += DailyIntFinal;
			tokenSupply_.self() = tokenSupply_.self() - _tokens;
			tokenBalanceLedger_.self()[_customerAddress] = tokenBalanceLedger_.self()[_customerAddress] - _tokens;
			platon_transfer(dev.self(), _devexitindividual_final);
			
			u128 _updatedPayouts = profitPerShare_.self() * _tokens * (_taxedEthereum * magnitude.self());
			payoutsTo_.self()[_customerAddress] -= _updatedPayouts;
			if(tokenSupply_.self() > 0){
				profitPerShare_.self() = profitPerShare_.self() + ((_dividends * magnitude.self()) / tokenSupply_.self());				
			}
			u128 now = u128(platon_timestamp());
			PLATON_EMIT_EVENT1(onTokenSell, _customerAddress, _tokens, _taxedEthereum,now, buyPrice());
		}
		
		ACTION bool transfer(Address _toAddress, u128 _amountOfTokens){
			onlyBagholders();		
			Address _customerAddress = platon_caller();
			if(_amountOfTokens > tokenBalanceLedger_.self()[_customerAddress]){
				platon_revert();			
			}
			if(myDividends(true) > 0){
				withdraw();			
			}
			
			u128 _tokenFee =  _amountOfTokens * transferFee_.self() / 100;        
			u128 _taxedTokens = _amountOfTokens - _tokenFee;
			u128 _dividends = tokensToEthereum_(_tokenFee);

			tokenSupply_.self() = tokenSupply_.self() - _tokenFee;
			tokenBalanceLedger_.self()[_customerAddress] = tokenBalanceLedger_.self()[_customerAddress] - _amountOfTokens;
			tokenBalanceLedger_.self()[_toAddress] = tokenBalanceLedger_.self()[_toAddress] + _taxedTokens;
			payoutsTo_.self()[_customerAddress] -= profitPerShare_.self() * _amountOfTokens;
			payoutsTo_.self()[_toAddress] += profitPerShare_.self() * _taxedTokens;
			profitPerShare_.self() = profitPerShare_.self() + ((_dividends * magnitude.self())/tokenSupply_.self());
			PLATON_EMIT_EVENT1(Transfer, _customerAddress, _toAddress, _taxedTokens);
			return true;			
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
			Address _customerAddress = platon_caller();
			if(_includeReferralBonus){
				return dividendsOf(_customerAddress) + referralBalance_.self()[_customerAddress];
			}
			return dividendsOf(_customerAddress) ;
		}

		CONST u128 balanceOf(Address _customerAddress) {
			return tokenBalanceLedger_.self()[_customerAddress];		
		}

		CONST u128 dividendsOf(Address _customerAddress){
			u128 val = profitPerShare_.self() * tokenBalanceLedger_.self()[_customerAddress] - payoutsTo_.self()[_customerAddress];
			return val / magnitude.self();
		}
		
		CONST u128 sellPrice(){
			// our calculation relies on the token supply, so we need supply. Doh.
			if (tokenSupply_.self() == u128(0)) {
				return tokenPriceInitial_.self() - tokenPriceIncremental_.self();
			} else {
				u128 _ethereum = tokensToEthereum_((1_LAT).Get());
				u128 _dividends = _ethereum * exitFee() / u128(100);
				u128 _devexit = _ethereum * u128(5) / u128(100);
				u128 _taxedEthereum1 = _ethereum - _dividends;
				u128 _taxedEthereum = _taxedEthereum1 - _devexit;
				return _taxedEthereum;
			}	
		}
		
		CONST u128 buyPrice(){
			if (tokenSupply_.self() == u128(0)) {
				return tokenPriceInitial_.self() + tokenPriceIncremental_.self();
			} else {
				u128 _ethereum = tokensToEthereum_((1_LAT).Get());
				u128 _dividends = _ethereum * entryFee_.self() / u128(100); 
				u128 _devexit = _ethereum * u128(5) / u128(100); 
				u128 _taxedEthereum1 = _ethereum + _dividends;
				u128 _taxedEthereum = _taxedEthereum1 + _devexit;
				return _taxedEthereum;
			}
		}

		CONST u128 calculateTokensReceived(u128 _ethereumToSpend){
			u128 _dividends = _ethereumToSpend * entryFee_.self() / u128(100);
			u128 _devbuyfees = _ethereumToSpend * u128(5) / u128(100);
			u128 _taxedEthereum1 = _ethereumToSpend - _dividends;
			u128 _taxedEthereum = _taxedEthereum1 - _devbuyfees;
			u128 _amountOfTokens = ethereumToTokens_(_taxedEthereum);
			return _amountOfTokens;
		}

		CONST u128 calculateEthereumReceived(u128 _tokensToSell){
			if(tokenSupply_.self() > _tokensToSell){
				platon_revert();
			}
			u128 _ethereum = tokensToEthereum_(_tokensToSell);
			u128 _dividends = _ethereum * exitFee() / u128(100);
			u128 _devexit =  _ethereum * u128(5) / u128(100);
			u128 _taxedEthereum1 = _ethereum - _dividends;
			u128 _taxedEthereum = _taxedEthereum1 - _devexit;
			return _taxedEthereum;	
		}

		CONST uint8_t exitFee(){
			return ExitFee_.self();		
		}

		u128 purchaseTokens(u128 _incomingEthereum, Address _referredBy) {
			Address _customerAddress = platon_caller();
			u128 _undividedDividends = _incomingEthereum *  entryFee_.self() / u128(100);
			u128 _referralBonus = _undividedDividends * refferalFee_.self() / u128(100);
			u128 _devbuyfees = _incomingEthereum * u128(5) / u128(100); 
			u128 _dividends1 = _undividedDividends -  _referralBonus;
			u128 _dividends = _dividends1 - _devbuyfees;
			u128 _taxedEthereum = _incomingEthereum - _undividedDividends;
			u128 _amountOfTokens = ethereumToTokens_(_taxedEthereum);
			u128 _fee = _dividends * magnitude.self();

			if(_amountOfTokens <= u128(0) || tokenSupply_.self() >= (_amountOfTokens + tokenSupply_.self())){
				platon_revert();
			}
			if (
				_referredBy != initaddr.self() &&
				_referredBy != _customerAddress &&
				tokenBalanceLedger_.self()[_referredBy] >= stakingRequirement.self()
			) {
				referralBalance_.self()[_referredBy] = referralBalance_.self()[_referredBy]  + _referralBonus;
			} else {
				_dividends = _dividends + _referralBonus;
				_fee = _dividends * magnitude.self();
			}

			if (tokenSupply_.self() > u128(0)) {
				tokenSupply_.self() = tokenSupply_.self() + _amountOfTokens;
				profitPerShare_.self() += (_dividends * magnitude.self() / tokenSupply_.self());
				_fee = _fee - (_fee - (_amountOfTokens * (_dividends * magnitude.self() / tokenSupply_.self())));
			} else {
				tokenSupply_.self() = _amountOfTokens;
			}

			tokenBalanceLedger_.self()[_customerAddress] = tokenBalanceLedger_.self()[_customerAddress] + _amountOfTokens;
			u128 _updatedPayouts = profitPerShare_ * _amountOfTokens - _fee;
			payoutsTo_.self()[_customerAddress] += _updatedPayouts;
			u128 now = u128(platon_timestamp());
			PLATON_EMIT_EVENT1(onTokenPurchase, _customerAddress, _incomingEthereum, _amountOfTokens, _referredBy, now, buyPrice());
			return _amountOfTokens;
		}
		
		u128 DividendsDistribution(u128 _incomingEthereum, Address _referredBy){
			Address _customerAddress = platon_caller();
			u128 _undividedDividends = _incomingEthereum * u128(100) / u128(100);
			u128 _referralBonus = _undividedDividends * refferalFee_.self() / u128(100);
			u128 _dividends = _undividedDividends - _referralBonus;
			u128 _taxedEthereum = _incomingEthereum - _undividedDividends;
			u128 _amountOfTokens = ethereumToTokens_(_taxedEthereum);
			u128 _fee = _dividends * magnitude.self();
			
			if(_amountOfTokens < 0 || _amountOfTokens + tokenSupply_.self() < tokenSupply_.self()){
				platon_revert();
			}

			if (
				_referredBy != initaddr.self() &&
				_referredBy != _customerAddress &&
				tokenBalanceLedger_.self()[_referredBy] >= stakingRequirement
			) {
				referralBalance_.self()[_referredBy] = referralBalance_.self()[_referredBy] + _referralBonus;
			} else {
				_dividends = _dividends + _referralBonus;
				_fee = _dividends * magnitude.self();
			}

			if (tokenSupply_.self() > 0) {
				tokenSupply_.self() = tokenSupply_ + _amountOfTokens;
				profitPerShare_.self() += (_dividends * magnitude.self() / tokenSupply_.self());
				_fee = _fee - (_fee - (_amountOfTokens * (_dividends * magnitude.self() / tokenSupply_.self())));
			} else {
				tokenSupply_.self() = _amountOfTokens;
			}

			tokenBalanceLedger_.self()[_customerAddress] = tokenBalanceLedger_.self()[_customerAddress] + _amountOfTokens;
			u128 _updatedPayouts = profitPerShare_.self() * _amountOfTokens - _fee;
			payoutsTo_.self()[_customerAddress] += _updatedPayouts;
			u128 now = u128(platon_timestamp());
			PLATON_EMIT_EVENT1(onTokenPurchase, _customerAddress, _incomingEthereum, _amountOfTokens, _referredBy, now, buyPrice());
			return _amountOfTokens;
		}
		
		u128 ethereumToTokens_(u128 _ethereum) {
			u128 _tokenPriceInitial = tokenPriceInitial_ * (1_LAT).Get();
			u128 _tokensReceived =
				(
					(
						
							(sqrt
								(
									(_tokenPriceInitial * _tokenPriceInitial)
									+
									(u128(2) * (tokenPriceIncremental_.self() * (1_LAT).Get()) * (_ethereum * (1_LAT).Get()))
									+
									((tokenPriceIncremental_.self() * tokenPriceIncremental_.self()) * (tokenSupply_.self() * tokenSupply_.self()))
									+
									(u128(2) * tokenPriceIncremental_.self() * _tokenPriceInitial*tokenSupply_.self())
								)
							) - _tokenPriceInitial
						
					) / (tokenPriceIncremental_)
				) - (tokenSupply_.self());

			return _tokensReceived;	
		}
		
		u128 tokensToEthereum_(u128 _tokens) {
			u128 tokens_ = (_tokens + (1_LAT).Get());
			u128 _tokenSupply = (tokenSupply_.self() + (1_LAT).Get());
			u128 _etherReceived =
				(
					(
						(
							(
								(
									tokenPriceInitial_ + (tokenPriceIncremental_ * (_tokenSupply / (1_LAT).Get()))
								) - tokenPriceIncremental_
							) * (tokens_ - (1_LAT).Get())
						) - (tokenPriceIncremental_ * ((tokens_ * tokens_ - tokens_) / (1_LAT).Get())) / u128(2)
					)
					/ (1_LAT).Get()
				);
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

PLATON_DISPATCH(Bank,(init)(buy)(IDD)(DivsAddon)(reinvest)(exit)
(withdraw)(sell)(transfer)(totalEthereumBalance)(totalSupply)
(myTokens)(myDividends)(balanceOf)(dividendsOf)(sellPrice)
(buyPrice)(calculateTokensReceived)(calculateEthereumReceived)(exitFee)(purchaseTokens)
(DividendsDistribution)(ethereumToTokens_)(tokensToEthereum_))

















