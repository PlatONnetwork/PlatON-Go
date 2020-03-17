#include <platon/platon.hpp>
#include <string>
using namespace platon;

class Message 
{	
	public:
		 // layout of message :: bytes:
	    // offset  0: 32 bytes :: uint256 - message length
	    // offset 32: 20 bytes :: address - recipient address
	    // offset 52: 32 bytes :: uint256 - value
	    // offset 84: 32 bytes :: bytes32 - transaction hash
	    // offset 116: 32 bytes :: uint256 - home gas price

	    // bytes 1 to 32 are 0 because message length is stored as little endian.
	    // mload always reads 32 bytes.
	    // so we can and have to start reading recipient at offset 20 instead of 32.
	    // if we were to read at 32 the address would contain part of value and be corrupted.
	    // when reading from offset 20 mload will read 12 zero bytes followed
	    // by the 20 recipient address bytes and correctly convert it into an address.
	    // this saves some storage/gas over the alternative solution
	    // which is padding address to 32 bytes and reading recipient at offset 32.
	    // for more details see discussion in:
	    // https://github.com/paritytech/parity-bridge/issues/61
	    static Address getRecipient(bytes message) {
	        /*address recipient;
	        // solium-disable-next-line security/no-inline-assembly
	        assembly {
	            recipient := mload(add(message, 20))
	        }
	        return recipient;*/
	        Address sender = platon_caller();
	        return sender;
	    }

	    static u128 getValue(bytes message) {
	        return platon_call_value();
	    }

	    static h256 getTransactionHash(bytes message) {
	        /*bytes32 hash;
	        // solium-disable-next-line security/no-inline-assembly
	        assembly {
	            hash := mload(add(message, 84))
	        }
	        return hash;*/
	        h256 hash = platon_block_hash(platon_block_number()-1);
	        return hash;
	    }

	    static u128 getHomeGasPrice(bytes message) {
	        /*uint256 gasPrice;
	        // solium-disable-next-line security/no-inline-assembly
	        assembly {
	            gasPrice := mload(add(message, 116))
	        }
	        return gasPrice;*/
	        u128 gasPrice = platon_gas_price();
	        return gasPrice;
	    }
};
