#define TESTNET
#include <platon/platon.hpp>
#include <string>
using namespace platon;

#include "src/BridgeDeploymentAddressStorage.hpp"
#include "src/ForeignBridgeGasConsumptionLimitsStorage.hpp"
#include "src/MessageSigning.hpp"


class ForeignBridge: public platon::Contract, public BridgeDeploymentAddressStorage, public ForeignBridgeGasConsumptionLimitsStorage
{
	public:
		 /// Number of authorities signatures required to withdraw the money.
	    ///
	    /// Must be less than number of authorities.
	    platon::StorageType<"initaddr"_n, Address> initaddr;

	    platon::StorageType<"requiredSignatures"_n, u128> requiredSignatures;

	    platon::StorageType<"estimatedGasCostOfWithdraw"_n, u128> estimatedGasCostOfWithdraw;

	    // Original parity-bridge assumes that anyone could forward final
	    // withdraw confirmation to the HomeBridge contract. That's why
	    // they need to make sure that no one is trying to steal funds by
	    // setting a big gas price of withdraw transaction. So,
	    // funds sender is responsible to limit this by setting gasprice
	    // as part of withdraw request.
	    // Since it is not the case for POA CCT bridge, gasprice is set
	    // to 1 Gwei which is minimal gasprice for POA network.
	 	// 1000000000 wei
	    platon::StorageType<"homeGasPrice"_n, u128> homeGasPrice;

	    /// Contract authorities.
	    platon::StorageType<"authorities"_n, std::map<Address, bool>> authorities;

	    /// Pending mesages
	    platon::StorageType<"messages"_n, std::map<h256, bytes>> messages;
	    /// ???
	    platon::StorageType<"signatures"_n, std::map<h256, bytes>> signatures;
	    
	    /// Pending deposits and authorities who confirmed them
	    platon::StorageType<"messages_signed"_n, std::map<h256, bool>> messages_signed;
	    platon::StorageType<"num_messages_signed"_n, std::map<h256, u128>> num_messages_signed;

	    /// Pending deposits and authorities who confirmed them
	    platon::StorageType<"deposits_signed"_n, std::map<h256, bool>> deposits_signed;
	    platon::StorageType<"num_deposits_signed"_n, std::map<h256, u128>> num_deposits_signed;

	    /// Token to work with
	    platon::StorageType<"erc20token"_n, Address> erc20token;
	    //ERC20 public erc20token;

	    /// List of authorities confirmed to set up ERC-20 token address
	    platon::StorageType<"tokenAddressAprroval_signs"_n, std::map<h256, bool>> tokenAddressAprroval_signs;
	    platon::StorageType<"num_tokenAddressAprroval_signs"_n, std::map<Address, u128>> num_tokenAddressAprroval_signs;

	public:
	    /// triggered when relay of deposit from HomeBridge is complete
	    PLATON_EVENT1(Deposit, Address, u128);

	    /// Event created on money withdraw.
	    PLATON_EVENT1(Withdraw, Address, u128, u128);

	    /// Collected signatures which should be relayed to home chain.
	    // params: address authorityResponsibleForRelay, bytes32 messageHash, uint256 NumberOfCollectedSignatures
	    PLATON_EVENT1(CollectedSignatures, Address, h256, u128);

	    /// Event created when new token address is set up.
	    // params: address token
	    PLATON_EVENT0(TokenAddress, Address);

	public:
		/// Constructor.
	    ACTION void init(u128 _requiredSignatures, std::vector<Address> _authorities, u128 _estimatedGasCostOfWithdraw)
	    {
	    	if(_requiredSignatures == u128(0)){
	    		platon_revert();
	    	}
	    	if(u128(_authorities.size()) < _requiredSignatures){
	    		platon_revert();
	    	}
	        requiredSignatures.self() = _requiredSignatures;

	        for (int i = 0; i < _authorities.size(); i++) {
	            authorities.self()[_authorities[i]] = true;
	        }
	        estimatedGasCostOfWithdraw.self() = _estimatedGasCostOfWithdraw;

	        homeGasPrice.self() = u128(1000000000);

	        auto address_init = make_address("laxqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq");
            if(address_init.second) Address initaddr = address_init.first;
	    }

		/// require that sender is an authority
	    void onlyAuthority(){
	    	Address sender = platon_caller();
	    	if(!authorities.self()[sender]){
	    		platon_revert();
	    	}
	    } 

	    /// Set up the token address. It allows to set up or change
	    /// the ERC20 token address only if authorities confirmed this.
	    ///
	    /// Usage maps instead of arrey allows to reduce gas consumption
	    ///
	    /// token address (address)
	    ACTION void setTokenAddress(Address token) {

	    	onlyAuthority();

	        // Duplicated deposits
	        std::string str = token.toString();
	        bytes data;
	        data.insert(data.begin(), str.begin(), str.end());
	        h256 token_sender = platon_sha3(data);

		    if(tokenAddressAprroval_signs.self()[token_sender]){
		    	platon_revert();
		    }
		    tokenAddressAprroval_signs.self()[token_sender]= true;

	        u128 _signed = num_tokenAddressAprroval_signs.self()[token] + u128(1);
	        num_tokenAddressAprroval_signs.self()[token] = _signed;

	        // TODO: this may cause troubles if requriedSignatures len is changed
	        if (_signed == requiredSignatures.self()) {
	        	//todo: 还有严重问题
	            erc20token.self() = token;
	            PLATON_EMIT_EVENT0(TokenAddress, token);
	        }
	    }

	    /// Used to transfer tokens to the `recipient`.
	    /// The bridge contract must own enough tokens to release them for 
	    /// recipients. Tokens must be transfered to the bridge contract BEFORE
	    /// the first deposit will be performed.
	    ///
	    /// Usage maps instead of array allows to reduce gas consumption
	    /// from 91169 to 89348 (solc 0.4.19). 
	    ///
	    /// deposit recipient (bytes20)
	    /// deposit value (uint256)
	    /// mainnet transaction hash (bytes32) // to avoid transaction duplication
	    ACTION void deposit(Address recipient, u128 value, h256 transactionHash) {
	    	onlyAuthority();
	    	if(erc20token.self() == initaddr.self()){
	    		platon_revert();
	    	}
	        // Protection from misbehaing authority
	        // todo: 此处需要调整
	        //h256 hash_msg = keccak256(recipient, value, transactionHash);
	    	std::string recipientStr = recipient.toString();
	    	std::string valueStr = std::to_string(value);
	    	std::string transactionHashStr = transactionHash.toString();
	    	std::string hashMsgStr = "";
	    	hashMsgStr.append(recipientStr);
	    	hashMsgStr.append(valueStr);
	    	hashMsgStr.append(transactionHashStr);
	    	bytes hashMsgData;
	    	hashMsgData.insert(hashMsgData.begin(), hashMsgStr.begin(), hashMsgStr.end());
	    	h256 hash_msg = platon_sha3(hashMsgData);

	        //h256 hash_sender = keccak256(msg.sender, hash_msg);
	        Address sender = platon_caller();
	        std::string senderStr = sender.toString();
	        std::string hashMessageStr = hash_msg.toString();
	        std::string hashSenderStr = "";
	        hashSenderStr.append(senderStr);
	        hashSenderStr.append(hashMessageStr);
	        bytes hashSenderData;
	        hashSenderData.insert(hashSenderData.begin(), hashSenderStr.begin(), hashSenderStr.end());
	        h256 hash_sender = platon_sha3(hashSenderData);

	        // Duplicated deposits
	        if(deposits_signed.self()[hash_sender]){
	        	platon_revert();
	        }
	        deposits_signed.self()[hash_sender] = true;

	        u128 _signed = num_deposits_signed.self()[hash_msg] + u128(1);
	        num_deposits_signed.self()[hash_msg] = _signed;

	        // TODO: this may cause troubles if requriedSignatures len is changed
	        if (_signed == requiredSignatures.self()) {
	            // If the bridge contract does not own enough tokens to transfer
	            // it will couse funds lock on the home side of the bridge
	            // todo: 此处要调用ERC20的代币接口
	            //erc20token.transfer(recipient, value);
	            PLATON_EMIT_EVENT1(Deposit, recipient, value);
	        }
	    }

	    /// Used to transfer `value` of tokens from `_from`s balance on local 
	    /// (`foreign`) chain to the same address (`_from`) on `home` chain.
	    /// Transfer of tokens within local (`foreign`) chain performed by usual
	    /// way through transfer method of the token contract.
	    /// In order to swap tokens to coins the owner (`_from`) must allow this
	    /// explicitly in the token contract by calling approveAndCall with address
	    /// of the bridge account.
	    /// The method locks tokens and emits a `Withdraw` event which will be
	    /// picked up by the bridge authorities.
	    /// Bridge authorities will then sign off (by calling `submitSignature`) on
	    /// a message containing `value`, the recipient (`_from`) and the `hash` of
	    /// the transaction on `foreign` containing the `Withdraw` event.
	    /// Once `requiredSignatures` are collected a `CollectedSignatures` event
	    /// will be emitted.
	    /// An authority will pick up `CollectedSignatures` an call
	    /// `HomeBridge.withdraw` which transfers `value - relayCost` to the
	    /// recipient completing the transfer.
	    ACTION bool receiveApproval(Address _from, u128 _value, Address _tokenContract, bytes _msg) {
	        /*require(erc20token != address(0x0));
	        require(msg.sender == address(erc20token));
	        require(erc20token.allowance(_from, this) >= _value);
	        erc20token.transferFrom(_from, this, _value);*/
	        PLATON_EMIT_EVENT1(Withdraw, _from, _value, homeGasPrice.self());

	        return true;
	    }

	    /// Should be used as sync tool
	    ///
	    /// Message is a message that should be relayed to main chain once authorities sign it.
	    ///
	    /// Usage several maps instead of structure allows to reduce gas consumption
	    /// from 265102 to 242334 (solc 0.4.19). 
	    ///
	    /// for withdraw message contains:
	    /// withdrawal recipient (bytes20)
	    /// withdrawal value (uint256)
	    /// foreign transaction hash (bytes32) // to avoid transaction duplication
	    ACTION void submitSignature(bytes signature, bytes message) {
	    	onlyAuthority();

	        // ensure that `signature` is really `message` signed by `msg.sender`
	        Address sender = platon_caller(); 
	        if(sender != MessageSigning::recoverAddressFromSignedMessage(signature, message)){
	        	platon_revert();
	        }
	        if(message.size() != 116){
	        	platon_revert();
	        }
	        h256 hash = platon_sha3(message);
	        std::string sendrStr = sender.toString();
	        std::string hashStr = hash.toString();
	        std::string hashByteStr = "";
	        hashByteStr.append(sendrStr);
	        hashByteStr.append(hashStr);
	        bytes hashSenderData;
	        hashSenderData.insert(hashSenderData.begin(), hashByteStr.begin(), hashByteStr.end());
			h256 hash_sender = platon_sha3(hashSenderData);

	        u128 _signed = num_messages_signed.self()[hash_sender] + u128(1);

	        if (_signed > u128(1)) {
	            // Duplicated signatures
	            if(messages_signed.self()[hash_sender]) {
	            	platon_revert();
	            }
	        }
	        else {
	            // check if it will really reduce gas usage in case of the second transaction
	            // with the same hash
	            messages.self()[hash] = message;
	        }
	        messages_signed.self()[hash_sender] = true;

	        //h256 sign_idx = keccak256(hash, (signed-u128(1));
	        std::string signIdStr = "";
	        signIdStr.append(hashStr);
	        signIdStr.append(std::to_string(_signed - u128(1)));
	        bytes signIdxData;
	        signIdxData.insert(signIdxData.begin(), signIdStr.begin(), signIdStr.end());
	        h256 sign_idx = platon_sha3(signIdxData);
	        signatures.self()[sign_idx]= signature;

	        num_messages_signed.self()[hash_sender] = _signed;

	        // TODO: this may cause troubles if requiredSignatures len is changed
	        if (_signed == requiredSignatures.self()) {
	            PLATON_EMIT_EVENT1(CollectedSignatures, platon_caller(), hash, _signed);
	        }
	    }

	    /// Get signature
	    CONST bytes signature(h256 hash, u128 index) {
	    	std::string signatureStr = "";
	    	signatureStr.append(hash.toString());
	    	signatureStr.append(std::to_string(index));
	    	bytes data;
	    	data.insert(data.begin(), signatureStr.begin(), signatureStr.end());
	    	h256 sign_idx = platon_sha3(data);
	        return signatures.self()[sign_idx];
	    }

	    /// Get message
	    CONST bytes message(h256 hash){
	        return messages.self()[hash];
	    }

};

PLATON_DISPATCH(ForeignBridge,(init)(setTokenAddress)(deposit)(receiveApproval)
	(submitSignature)(signature)(message));
