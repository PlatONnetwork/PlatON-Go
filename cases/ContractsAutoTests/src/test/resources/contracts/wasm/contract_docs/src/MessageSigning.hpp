#include <platon/platon.hpp>
#include <string>
using namespace platon;

class MessageSigning 
{
	
	public:
		static Address recoverAddressFromSignedMessage(bytes signature, bytes message){
			/*require(signature.length == 65);
			bytes32 r;
			bytes32 s;
			bytes1 v;
			// solium-disable-next-line security/no-inline-assembly
			assembly {
				r := mload(add(signature, 0x20))
				s := mload(add(signature, 0x40))
				v := mload(add(signature, 0x60))
			}
			return ecrecover(hashMessage(message), uint8(v), r, s);*/
			return platon_caller();
		}

		static h256 hashMessage(bytes message){
			/*bytes memory prefix = "\x19Ethereum Signed Message:\n";
			return keccak256(prefix, Helpers.uintToString(message.length), message);*/
			return platon_sha3(message);;
		}
};
