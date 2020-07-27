#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

CONTRACT call_precompile : public platon::Contract {
    public:
        ACTION void init(){}


        CONST const std::string  cross_call_ecrecover (platon::bytes msgh, uint8_t v, platon::bytes r, platon::bytes s, uint64_t value, uint64_t gas) {

                    // uint8_t to bytes32
                    std::vector<byte> vbytes;
                    vbytes.resize(32);
                    memset(vbytes.data(), 0, 32);
                    vbytes[31] = v;

                    // v append to msgh
                    std::copy(vbytes.begin(), vbytes.end(), std::back_inserter(msgh));

                    // append r
                    std::copy(r.begin(), r.end(), std::back_inserter(msgh));
                    // append s
                    std::copy(s.begin(), s.end(), std::back_inserter(msgh));

                    platon::bytes input = msgh;

                    std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxxvxfq";
                    auto address_info = make_address(addr);
                    if(address_info.second){
                        if (platon_call(address_info.first, input, value, gas)) {
                            DEBUG("cross call contract ecrecover success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                          std::string str = toHex(ret);

                          DEBUG("cross call contract ecrecover success", "acc", str);

                          return str;
                        }
                    }


                    DEBUG("cross call contract ecrecover fail", "address", addr);
                    return "";
                }



                CONST const std::string  cross_call_sha256hash (std::string &in, uint64_t value, uint64_t gas) {


                    platon::bytes  input = fromHex(in);

                    std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzg4es8l";

                     auto address_info = make_address(addr);
                     if(address_info.second){
                        if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract sha256hash success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract sha256hash success", "hash", str);
                           return str;
                         }
                     }


                    DEBUG("cross call contract sha256hash fail", "address", addr);
                    return "";
                }


                 CONST const std::string  cross_call_ripemd160hash (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqr4rd96d";

                      auto address_info = make_address(addr);
                      if(address_info.second){
                         if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract ripemd160hash success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract ripemd160hash success", "hash", str);
                           return str;
                         }
                      }


                     DEBUG("cross call contract ripemd160hash fail", "address", addr);
                     return "";
                 }


                 CONST const std::string  cross_call_dataCopy (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqy5664mg";

                     auto address_info = make_address(addr);
                     if(address_info.second){
                         if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract dataCopy success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract dataCopy success", "hash", str);
                           return str;
                         }
                     }

                     DEBUG("cross call contract dataCopy fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bigModExp (platon::bytes base, platon::bytes exponent, platon::bytes modulus, uint64_t value, uint64_t gas) {


                      // uint8_t to bytes32
                      std::vector<byte> len;
                      len.resize(32);
                      memset(len.data(), 0, 32);

                      uint8_t l = 32;
                      len[31] = l;

                      platon::bytes input;

                      // [32]byte(baseLen) + [32]byte(expLen) + [32]byte(modLen)
                      std::copy(len.begin(), len.end(), std::back_inserter(input));
                      std::copy(len.begin(), len.end(), std::back_inserter(input));
                      std::copy(len.begin(), len.end(), std::back_inserter(input));

                      // append base
                      std::copy(base.begin(), base.end(), std::back_inserter(input));
                      // append exponent
                      std::copy(exponent.begin(), exponent.end(), std::back_inserter(input));
                      // append modulus
                      std::copy(modulus.begin(), modulus.end(), std::back_inserter(input));

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq9fvwqx6";

                     auto address_info = make_address(addr);
                     if(address_info.second){
                          if (platon_call(address_info.first, input, value, gas)) {
                              DEBUG("cross call contract bigModExp success", "address", addr);

                              platon::bytes ret;
                              size_t len = platon_get_call_output_length();

                              ret.resize(len);
                              platon_get_call_output(ret.data());

                              std::string str = toHex(ret);
                              DEBUG("cross call contract bigModExp success", "hash", str);
                              return str;
                          }
                     }

                     DEBUG("cross call contract bigModExp fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bn256Add (platon::bytes ax, platon::bytes ay, platon::bytes bx, platon::bytes by, uint64_t value, uint64_t gas) {


                     platon::bytes input;

                     // ax + ay + bx + by
                     std::copy(ax.begin(), ax.end(), std::back_inserter(input));
                     std::copy(ay.begin(), ay.end(), std::back_inserter(input));
                     std::copy(bx.begin(), bx.end(), std::back_inserter(input));
                     std::copy(by.begin(), by.end(), std::back_inserter(input));

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqx8lmkg9";
                     auto address_info = make_address(addr);
                     if(address_info.second){
                         if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract bn256Add success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract bn256Add success", "hash", str);
                           return str;
                         }
                     }

                     DEBUG("cross call contract bn256Add fail", "address", addr);
                     return "";
                 }



                 CONST const std::string  cross_call_bn256ScalarMul (platon::bytes x, platon::bytes y, platon::bytes scalar, uint64_t value, uint64_t gas) {


                    platon::bytes input;

                     // x + y + scalar
                     std::copy(x.begin(), x.end(), std::back_inserter(input));
                     std::copy(y.begin(), y.end(), std::back_inserter(input));
                     std::copy(scalar.begin(), scalar.end(), std::back_inserter(input));

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq86f0r4h";

                     auto address_info = make_address(addr);
                     if(address_info.second){
                         if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract bn256ScalarMul success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract bn256ScalarMul success", "hash", str);
                           return str;
                         }
                     }


                     DEBUG("cross call contract bn256ScalarMul fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bn256Pairing (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqg9yul20";

                     auto address_info = make_address(addr);
                     if(address_info.second){
                         if (platon_call(address_info.first, input, value, gas)) {
                             DEBUG("cross call contract bn256Pairing success", "address", addr);

                           platon::bytes ret;
                           size_t len = platon_get_call_output_length();

                           ret.resize(len);
                           platon_get_call_output(ret.data());

                           std::string str = toHex(ret);
                           DEBUG("cross call contract bn256Pairing success", "hash", str);
                           return str;
                         }
                     }

                     DEBUG("cross call contract bn256Pairing fail", "address", addr);
                     return "";
                 }



};

PLATON_DISPATCH(call_precompile, (init)(cross_call_ecrecover)(cross_call_sha256hash)(cross_call_ripemd160hash)(cross_call_dataCopy)(cross_call_bigModExp)(cross_call_bn256Add)(cross_call_bn256ScalarMul)(cross_call_bn256Pairing))