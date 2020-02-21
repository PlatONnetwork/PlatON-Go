#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

CONTRACT call_precompile : public platon::Contract {
    public:
        ACTION void init(){}

/*
        ACTION uint64_t cross_call_ecrecover (std::string &msg_hash, std::string &v, std::string &r, std::string &s, uint64_t value, uint64_t gas) {

            // uint256[4] memory input;
            // input[0] = uint256(msgh);
            // input[1] = v;
            // input[2] = uint256(r);
            // input[3] = uint256(s);
            //
            //
            //
            // dataHash: "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d", this hash is not txHash
            //V = 27
            //R = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe"
            //S = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6"
            //address: "0x8a9B36694F1eeeb500c84A19bB34137B05162EC5"
            h256[4] in;
            in[0] = h256(msg_hash);
            in[1] = h256(v);
            in[2] = h256(r);
            in[3] = h256(s);

            // platon::bytes  input = fromHex(in);

            std::string addr = "0x01"

            if (platon_call(Address(addr), in.toBytes(), value, gas)) {
                DEBUG("cross call contract success", "address", addr);
            } else {
                DEBUG("cross call contract fail", "address", addr);
            }
        }
*/

        CONST const std::string  cross_call_ecrecover (std::string &in, uint64_t value, uint64_t gas) {

                    // uint256[4] memory input;
                    // input[0] = uint256(msgh);
                    // input[1] = v;
                    // input[2] = uint256(r);
                    // input[3] = uint256(s);
                    //
                    //
                    //
                    // dataHash: "0xe281eaa11e6e37e6f53aade5d6c5b7201ef1c66162ec42ccc3215a0c4349350d", this hash is not txHash
                    //V = 27
                    //R = "0x55b60cadd4b4a3ea4fc368ef338f97e12e7328dd6e9e969a3fd8e5c10be855fe"
                    //S = "0x2b42cee2585a16ea537efcb88009c1aeac693c28b59aa6bbff0baf22730338f6"
                    //address: "0x8a9B36694F1eeeb500c84A19bB34137B05162EC5"

                    platon::bytes  input = fromHex(in);

                    std::string addr = "0x0000000000000000000000000000000000000001";

                    if (platon_call(Address(addr), input, value, gas)) {
                        DEBUG("cross call contract ecrecover success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                      std::string str = toHex(ret);

                      //  Address as = Address("0000000000000000000000008a9b36694f1eeeb500c84a19bb34137b05162ec5");
                       Address bs = Address(ret);

                       // std::string str =  acc.toString(); // toHex(ret);
                      DEBUG("cross call contract ecrecover success", "acc", str);
                       // DEBUG("cross call contract success", "as", as.toString());
                        DEBUG("cross call contract ecrecover success", "bs", bs.toString());
                       // return  str;
                       return str;
                    }
                    DEBUG("cross call contract ecrecover fail", "address", addr);
                    return "";
                }



                CONST const std::string  cross_call_sha256hash (std::string &in, uint64_t value, uint64_t gas) {


                    platon::bytes  input = fromHex(in);

                    std::string addr = "0x0000000000000000000000000000000000000002";

                    if (platon_call(Address(addr), input, value, gas)) {
                        DEBUG("cross call contract sha256hash success", "address", addr);

                      platon::bytes ret;
                      size_t len = platon_get_call_output_length();

                      ret.resize(len);
                      platon_get_call_output(ret.data());

                      std::string str = toHex(ret);
                      DEBUG("cross call contract sha256hash success", "hash", str);
                      return str;
                    }
                    DEBUG("cross call contract sha256hash fail", "address", addr);
                    return "";
                }


                 CONST const std::string  cross_call_ripemd160hash (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000003";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract ripemd160hash success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract ripemd160hash success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract ripemd160hash fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_dataCopy (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000004";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract dataCopy success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract dataCopy success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract dataCopy fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bigModExp (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000005";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract bigModExp success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract bigModExp success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract bigModExp fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bn256Add (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000006";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract bn256Add success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract bn256Add success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract bn256Add fail", "address", addr);
                     return "";
                 }



                 CONST const std::string  cross_call_bn256ScalarMul (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000007";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract bn256ScalarMul success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract bn256ScalarMul success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract bn256ScalarMul fail", "address", addr);
                     return "";
                 }

                 CONST const std::string  cross_call_bn256Pairing (std::string &in, uint64_t value, uint64_t gas) {


                     platon::bytes  input = fromHex(in);

                     std::string addr = "0x0000000000000000000000000000000000000008";

                     if (platon_call(Address(addr), input, value, gas)) {
                         DEBUG("cross call contract bn256Pairing success", "address", addr);

                       platon::bytes ret;
                       size_t len = platon_get_call_output_length();

                       ret.resize(len);
                       platon_get_call_output(ret.data());

                       std::string str = toHex(ret);
                       DEBUG("cross call contract bn256Pairing success", "hash", str);
                       return str;
                     }
                     DEBUG("cross call contract bn256Pairing fail", "address", addr);
                     return "";
                 }



};

PLATON_DISPATCH(call_precompile, (init)(cross_call_ecrecover)(cross_call_sha256hash)(cross_call_ripemd160hash)(cross_call_dataCopy)(cross_call_bigModExp)(cross_call_bn256Add)(cross_call_bn256ScalarMul)(cross_call_bn256Pairing))