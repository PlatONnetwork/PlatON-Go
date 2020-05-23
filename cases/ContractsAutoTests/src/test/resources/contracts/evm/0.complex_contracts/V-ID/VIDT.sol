pragma solidity ^0.4.24;
/* 
	This solidity version supports auto-generated getter functions for all public state variables. 
	The default ERC20 getter functions for name(), symbol(), decimals() and totalSupply() are created by the compiler.
*/

contract ERC20 {
/* 
	Full ERC20 EIP20 compliancy
	Credits https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md 
*/

	uint256 public totalSupply;

	function balanceOf(address who) public view returns (uint256 balance);

	function allowance(address owner, address spender) public view returns (uint256 remaining);

	function transfer(address to, uint256 value) public returns (bool success);

	function approve(address spender, uint256 value) public returns (bool success);

	function transferFrom(address from, address to, uint256 value) public returns (bool success);

	event Transfer(address indexed _from, address indexed _to, uint256 _value);

	event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}

library SafeMath {
/* 
	SafeMath library with small extension to the assert statements, adds or subtracts two numbers, reverts on overflow
*/

	function sub(uint256 a, uint256 b) internal pure returns (uint256 c) {
		c = a - b;
		assert(b <= a && c <= a);
		return c;
	}

	function add(uint256 a, uint256 b) internal pure returns (uint256 c) {
		c = a + b;
		assert(c >= a && c>=b);
		return c;
	}
}

library SafeERC20 {
/* 
	SafeERC20 library just here to transfer any non V-ID tokens accidentally sent to the contract, can be used to reimburse senders 
	Credits https://github.com/OpenZeppelin/openzeppelin-solidity/blob/master/contracts/token/ERC20/SafeERC20.sol (not identical, tailored to V-ID)
*/

	function safeTransfer(ERC20 _token, address _to, uint256 _value) internal {
		require(_token.transfer(_to, _value));
	}
}

contract Owned {
/*
	The contract is Owned so ownership can be transferred and some functions can only be invoked by the owner
	Credits https://github.com/OpenZeppelin/openzeppelin-solidity/blob/master/contracts/ownership/Ownable.sol (not identical, tailored to V-ID)
*/

	address public owner;

	constructor() public {
		owner = msg.sender;
	}

	modifier onlyOwner {
		require(msg.sender == owner,"O1- Owner only function");
		_;
	}

	function setOwner(address newOwner) onlyOwner public {
		owner = newOwner;
	}
}

contract Pausable is Owned {
/*
	The contract is Pauseable so token transfers can be paused in case of a security risk or Ethereum protocol problem
	Credits https://github.com/OpenZeppelin/openzeppelin-solidity/blob/master/contracts/token/ERC20/PausableToken.sol (not identical, tailored to V-ID)
*/

	event Pause();
	event Unpause();

	bool public paused = false;

	modifier whenNotPaused() {
		require(!paused);
		_;
	}

	modifier whenPaused() {
		require(paused);
		_;
	}

	function pause() public onlyOwner whenNotPaused {
		paused = true;
		emit Pause();
	}

	function unpause() public onlyOwner whenPaused {
		paused = false;
		emit Unpause();
	}
}

contract VIDToken is Owned, Pausable, ERC20 {
/*
	The V-ID contract itself, using all the before mentioned defaults and libraries
*/
	
// Applying the SafeMath library to all integers used
	using SafeMath for uint256; 

// Applying the SafeERC20 library to all ERC20 instances 
 	using SafeERC20 for ERC20;

// Default ERC20 array for storing the balances
	mapping (address => uint256) public balances; 

// Default ERC20 array for storing the allowances between accounts
 	mapping (address => mapping (address => uint256)) public allowed;

// A custom array to store accounts that are (temporarily) frozen, to be able to prevent abuse
	mapping (address => bool) public frozenAccount; 

// A custom array to store all verified publisher accounts
	mapping (address => bool) public verifyPublisher; 

// A custom array to store all verified validation wallets
	mapping (address => bool) public verifyWallet; 

// Simple structure for the storage of the hashes
	struct fStruct { uint256 index; } 

// The hashes are stored in this array
	mapping(string => fStruct) private fileHashes; 

// The array of hashes is indexed in this array for queries
	string[] private fileIndex; 

// The name of the token
	string public constant name = "V-ID Token"; 

// The number of decimals, 18 for most broadly used implementation
	uint8 public constant decimals = 18; 

// The preferred ticker / symbol of the token
	string public constant symbol = "VIDT"; 

// The initial supply (lacking the 18 decimals) 100.000.000 VIDT
	uint256 public constant initialSupply = 100000000; 

// The default validation price per hash
	uint256 public validationPrice = 7 * 10 ** uint(decimals); 

// The default validation wallet
	address public validationWallet = address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j");

// The constructor is invoked on contract creation, it sets the contract owner and default wallet

	constructor() public {
		validationWallet = msg.sender;
		verifyWallet[msg.sender] = true;
		totalSupply = initialSupply * 10 ** uint(decimals);
		balances[msg.sender] = totalSupply;
		emit Transfer(address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j"),owner,initialSupply);
	}

// The fallback function which will be invoked when the function signature is not found, the revert() prevents the contract from receiving Ethereum

	function () public payable {
		revert();
	}

// The default ERC20 transfer() function is extended to be pausable and to prevent transfers from or to frozen accounts

	function transfer(address _to, uint256 _value) public whenNotPaused returns (bool success) {
		require(_to != msg.sender,"T1- Recipient can not be the same as sender");
		require(_to != address("lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a") && _to != address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j"),"T2- Please check the recipient address");
		require(balances[msg.sender] >= _value,"T3- The balance of sender is too low");
		require(!frozenAccount[msg.sender],"T4- The wallet of sender is frozen");
		require(!frozenAccount[_to],"T5- The wallet of recipient is frozen");

		balances[msg.sender] = balances[msg.sender].sub(_value);
		balances[_to] = balances[_to].add(_value);

		emit Transfer(msg.sender, _to, _value);

		return true;
	}

// The default ERC20 transferFrom() function is extended to be pausable and to prevent transfers from or to frozen accounts

	function transferFrom(address _from, address _to, uint256 _value) public whenNotPaused returns (bool success) {
		require(_to != address("lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a") && _to != address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j"),"TF1- Please check the recipient address");
		require(balances[_from] >= _value,"TF2- The balance of sender is too low");
		require(allowed[_from][msg.sender] >= _value,"TF3- The allowance of sender is too low");
		require(!frozenAccount[_from],"TF4- The wallet of sender is frozen");
		require(!frozenAccount[_to],"TF5- The wallet of recipient is frozen");

		balances[_from] = balances[_from].sub(_value);
		balances[_to] = balances[_to].add(_value);

		allowed[_from][msg.sender] = allowed[_from][msg.sender].sub(_value);

		emit Transfer(_from, _to, _value);

		return true;
	}

// Default ERC20 balanceOf() returns the current balance of _owner

	function balanceOf(address _owner) public view returns (uint256 balance) {
		return balances[_owner];
	}

// Default ERC20 approve() extended to prevent abuse under race conditions by requiring a 0 allowance before approving

	function approve(address _spender, uint256 _value) public whenNotPaused returns (bool success) {
		require((_value == 0) || (allowed[msg.sender][_spender] == 0),"A1- Reset allowance to 0 first");

		allowed[msg.sender][_spender] = _value;

		emit Approval(msg.sender, _spender, _value);

		return true;
	}

// The increaseApproval() function is here to efficiently increase any allowance in 1 transaction, protected against overflow bug using SafeMath

	function increaseApproval(address _spender, uint256 _addedValue) public whenNotPaused returns (bool) {
		allowed[msg.sender][_spender] = allowed[msg.sender][_spender].add(_addedValue);

		emit Approval(msg.sender, _spender, allowed[msg.sender][_spender]);

		return true;
	}

// The decreaseApproval() function is here to efficiently decrease any allowance in 1 transaction, protected against overflow bug using SafeMath

	function decreaseApproval(address _spender, uint256 _subtractedValue) public whenNotPaused returns (bool) {
		allowed[msg.sender][_spender] = allowed[msg.sender][_spender].sub(_subtractedValue);

		emit Approval(msg.sender, _spender, allowed[msg.sender][_spender]);

		return true;
	}

// Default ERC20 allowance() returns the current allowance of _spender using tokens of _owner 

	function allowance(address _owner, address _spender) public view returns (uint256 remaining) {
		return allowed[_owner][_spender];
	}

// The struct is a named array to be used in the tokenFallback() function

	struct TKN { address sender; uint256 value; bytes data; bytes4 sig; }

// Though the tokenFallback() is a pure bogey, it's existence confirms the contract's ability to transfer any token (ERC223 compatibility)

	function tokenFallback(address _from, uint256 _value, bytes _data) public pure returns (bool) {
		TKN memory tkn;
		tkn.sender = _from;
		tkn.value = _value;
		tkn.data = _data;
		uint32 u = uint32(_data[3]) + (uint32(_data[2]) << 8) + (uint32(_data[1]) << 16) + (uint32(_data[0]) << 24);
		tkn.sig = bytes4(u);
		return true;
	}

// The transferToken() is the function that is actually used by the contract owner to transfer tokens to the owner wallet so these can be reimbursed to any sender

	function transferToken(address tokenAddress, uint256 tokens) public onlyOwner {
		ERC20(tokenAddress).safeTransfer(owner,tokens);
	}

// The V-ID tokenomics model is set to burn() a minimum of 10% of all tokens received in payment for validations and implementations, tokens are burned by sending them to the (0) address

	function burn(uint256 _value) public onlyOwner returns (bool) {
		require(_value <= balances[msg.sender],"B1- The balance of burner is too low");

		balances[msg.sender] = balances[msg.sender].sub(_value);
		totalSupply = totalSupply.sub(_value);

		emit Burn(msg.sender, _value);

		emit Transfer(msg.sender, address("lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqmscn5j"), _value);

		return true;
	}

// The freeze() function is used to prevent abuse or protecting funds by freezing accounts

	function freeze(address _address, bool _state) public onlyOwner returns (bool) {
		frozenAccount[_address] = _state;

		emit Freeze(_address, _state);

		return true;
	}

// All publishers are identified and validated by wallet / account and are named in the validation transaction

	function validatePublisher(address Address, bool State, string Publisher) public onlyOwner returns (bool) {
		verifyPublisher[Address] = State;

		emit ValidatePublisher(Address,State,Publisher);

		return true;
	}

// Multiple wallets can be used to receive the V-ID Token payments, the contract owner / admin can (in)validate wallets using this function

	function validateWallet(address Address, bool State, string Wallet) public onlyOwner returns (bool) {
		verifyWallet[Address] = State;

		emit ValidateWallet(Address,State,Wallet);

		return true;
	}

/*
	The essence of the validation process is the validation of a files fingerprint / hash using the validateFile() function. 
	The hash is alway stored in the transaction log (input) and optionally (though default) in the contract itself and in the event logs.
	In this way the validateFile() function provides the use of all 3 tiers in the data storage logic of the Ethereum blockchain (transaction, state and receipt). 
	
	Herein also lies the logic of using the token as method of most efficient payment for each micro transaction.
	As we will be using automated validation tools (also API's) it is easy to vary the price per validation, therefore the validationPrice variable only represents the minimum amount.
*/

	function validateFile(address To, uint256 Payment, bytes Data, bool cStore, bool eLog) public whenNotPaused returns (bool) {
		require(Payment>=validationPrice,"V1- Insufficient payment provided");
		require(verifyPublisher[msg.sender],"V2- Unverified publisher address");
		require(!frozenAccount[msg.sender],"V3- The wallet of publisher is frozen");
		require(Data.length == 64,"V4- Invalid hash provided");

		if (!verifyWallet[To] || frozenAccount[To]) {
			To = validationWallet;
		}

		uint256 index = 0;
		string memory fileHash = string(Data);

		if (cStore) {
			if (fileIndex.length > 0) {
				require(fileHashes[fileHash].index == 0,"V5- This hash was previously validated");
			}

			fileHashes[fileHash].index = fileIndex.push(fileHash)-1;
			index = fileHashes[fileHash].index;
		}

		if (allowed[To][msg.sender] >= Payment) {
			allowed[To][msg.sender] = allowed[To][msg.sender].sub(Payment);
		} else {
			balances[msg.sender] = balances[msg.sender].sub(Payment);
			balances[To] = balances[To].add(Payment);
		}

		emit Transfer(msg.sender, To, Payment);

		if (eLog) {
			emit ValidateFile(index,fileHash);
		}

		return true;
	}

/* 
	Verifying a file is available at zero cost as it does not require a transaction, just a query. 
	The verifyFile() function matches twice, first by retrieving the index of a file by its hash (string) and again by comparing the strings 1 character at a time.
*/

	function verifyFile(string fileHash) public view returns (bool) {
		if (fileIndex.length == 0) {
			return false;
		}

		bytes memory a = bytes(fileIndex[fileHashes[fileHash].index]);
		bytes memory b = bytes(fileHash);

		if (a.length != b.length) {
			return false;
		}

		for (uint256 i = 0; i < a.length; i ++) {
			if (a[i] != b[i]) {
				return false;
			}
		}

		return true;
	}

// The contract owner / V-ID can set the price for each validation using the setPrice() function

	function setPrice(uint256 newPrice) public onlyOwner {
		validationPrice = newPrice;
	}

// The contract owner / V-ID can set the default wallet for payments using the setWallet() function

	function setWallet(address newWallet) public onlyOwner {
		validationWallet = newWallet;
	}

// The listFiles() function is to provide direct access to the data for audits and consistency checks

	function listFiles(uint256 startAt, uint256 stopAt) onlyOwner public returns (bool) {
		if (fileIndex.length == 0) {
			return false;
		}

		require(startAt <= fileIndex.length-1,"L1- Please select a valid start");

		if (stopAt > 0) {
			require(stopAt > startAt && stopAt <= fileIndex.length-1,"L2- Please select a valid stop");
		} else {
			stopAt = fileIndex.length-1;
		}

		for (uint256 i = startAt; i <= stopAt; i++) {
			emit LogEvent(i,fileIndex[i]);
		}

		return true;
	}

// All possible custom events are listed here for their signatures

	event Burn(address indexed burner, uint256 value);
	event Freeze(address target, bool frozen);

	event ValidateFile(uint256 index, string data);
	event ValidatePublisher(address indexed publisherAddress, bool state, string indexed publisherName);
	event ValidateWallet(address indexed walletAddress, bool state, string indexed walletName);

	event LogEvent(uint256 index, string data) anonymous;
}
