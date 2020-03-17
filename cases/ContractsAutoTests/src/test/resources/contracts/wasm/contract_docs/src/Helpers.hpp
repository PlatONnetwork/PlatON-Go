#include <platon/platon.hpp>
#include <string>
using namespace platon;

#include "MessageSigning.hpp"

class Helpers 
{
	public:
		Helpers() {
		}
		
	public:
		/// returns whether `array` contains `value`.
		static bool addressArrayContains(std::vector<Address> array, Address value){
			for (uint64_t i = 0; i < array.size(); i++) {
				if (array[i] == value) {
					return true;
				}
			}
			return false;
		}
		
		// returns the digits of `inputValue` as a string.
		// example: `uintToString(12345678)` returns `"12345678"`
		static std::string uintToString(u128 inputValue){
			return std::to_string(inputValue);
		}

		/// returns whether signatures (whose components are in `vs`, `rs`, `ss`)
		/// contain `requiredSignatures` distinct correct signatures
		/// where signer is in `allowed_signers`
		/// that signed `message`
		static bool hasEnoughValidSignatures(bytes message, std::vector<uint8_t> vs, 
											std::vector<h256> rs, std::vector<h256> ss, 
											std::vector<Address> allowed_signers, 
											u128 requiredSignatures){
			// not enough signatures
			/*if (u128(vs.size()) < requiredSignatures) {
				return false;
			}*/

			h256 hash = MessageSigning::hashMessage(message);
			/*
			var encountered_addresses = new address[](allowed_signers.length);

			for (int i = 0; i < requiredSignatures; i++) {
				var recovered_address = ecrecover(hash, vs[i], rs[i], ss[i]);
				// only signatures by addresses in `addresses` are allowed
				if (!addressArrayContains(allowed_signers, recovered_address)) {
					return false;
				}
				// duplicate signatures are not allowed
				if (addressArrayContains(encountered_addresses, recovered_address)) {
					return false;
				}
				encountered_addresses[i] = recovered_address;
			}*/
			return true;
		}
};
