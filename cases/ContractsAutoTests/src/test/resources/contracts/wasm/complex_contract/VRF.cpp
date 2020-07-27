#define TESTNET
#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

/**
 * VRF.sol合约
 * */

CONTRACT VRF : public platon::Contract{

    private:
        platon::StorageType<"groupOrder"_n, u128> groupOrder;
        platon::StorageType<"fieldSize"_n, u128> fieldSize;
        platon::StorageType<"wordlength"_n, u128> wordLengthBytes;
        platon::StorageType<"bigmod"_n,std::array<u128,6>> bigModExpContractInputs;
        platon::StorageType<"output"_n, u128> output;
        platon::StorageType<"sqrtpower"_n, u128> sqrtPower;
        platon::StorageType<"newcandi"_n,std::array<u128,2>> newCandiArr;
        platon::StorageType<"rv"_n,std::array<u128,2>> rv;
        platon::StorageType<"scalar"_n, u128> SCALAR_FROM_CURVE_POINTS_HASH_PREFIX;
        platon::StorageType<"scalaprr"_n, u128> VRF_RANDOM_OUTPUT_HASH_PREFIX;
        platon::StorageType<"prooflen"_n, u128> PROOF_LENGTH;
        platon::StorageType<"addstr"_n, std::string> address_str;
        platon::StorageType<"findhash"_n, u128> fieldHash128;

    public:
        ACTION void init(){
            groupOrder.self() = "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"_u128; 
            fieldSize.self() =  "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"_u128;
            wordLengthBytes.self() = 32; //0x20;
            SCALAR_FROM_CURVE_POINTS_HASH_PREFIX.self() = 2;
            VRF_RANDOM_OUTPUT_HASH_PREFIX.self() = 3;
            address_str.self() = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzg4es8l";
            PROOF_LENGTH.self() =  64 + // PublicKey (uncompressed format.)
                            64 + // Gamma
                            32 + // C
                            32 + // S
                            32 + // Seed
                            0 + // Dummy entry: The following elements are included for gas efficiency:
                            32 + // uWitness (gets padded to 256 bits, even though it's only 160)
                            64 + // cGammaWitness
                            64 + // sHashWitness
                            32; // zInv  (Leave Output out, because that can be efficiently calculated)
            
        }

        u128 bigModExp(u128 base,u128 exponent){
            u128 callResult;
            bigModExpContractInputs.self()[0] = wordLengthBytes.self();
            bigModExpContractInputs.self()[1] = wordLengthBytes.self();
            bigModExpContractInputs.self()[2] = wordLengthBytes.self();
            bigModExpContractInputs.self()[3] = base;
            bigModExpContractInputs.self()[4] = exponent;
            bigModExpContractInputs.self()[5] = fieldSize.self();

            u128 value = 0; //转账金额
            u128 gas = 4712388; //预估的手续费
            std::string addr = "0x05";//目标合约地址
            platon::bytes params = platon::cross_call_args("data", bigModExpContractInputs.self());

            u128 resulth128;

            Address addr3;
            auto address_info3 = make_address(addr);
            if(address_info3.second){
              addr3 = address_info3.first;
            }

            if (platon_call(addr3, params, value, gas)) {
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);
                resulth128 = std::stoull(str);
            }
            return resulth128;

        }

        u128 getSqrtPower(){
            return (fieldSize.self()+1) >>2;
        }

        u128 squareRoot(u128 x){
            return bigModExp(x,getSqrtPower());
        }

        u128 ySquared(u128 x){
            u128 value1 = x*x%fieldSize.self();
            u128 xCubed = x*value1%fieldSize.self();
            return (xCubed+7)%fieldSize.self();
        }

        bool isOnCurve(u128 p1,u128 p2){
            return ySquared(p1) == p2*p2%fieldSize.self();
        }

        //c++实现keccak256算法实现未开发
        // u128 fieldHash(bytes &b){
        u128 fieldHash(bytes &bt){
//            h256 res = platon_sha256(bt);//h256转u128
//            fieldHash128.self() = toHex(res.toString())
//            std::string res1 = std::to_string(res);
//            while(res >= fieldSize.self()){
//                res = platon_sha256(res);
//                std::string str = toHex(res);
//                base = std::stoul(str);
//            }
            return fieldHash128.self();
        }

        bytes keccak256(bytes &bt){
            u128 value = 12;
            u128 gas = 4712388;
            std::string target_addr = "lax1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzg4es8l";
            platon::bytes ret;

            auto address_info = make_address(target_addr);
            if(address_info.second){
              Address addr = address_info.first;
              if (platon_call(addr, bt, value, gas)) {
                  size_t len = platon_get_call_output_length();
                  ret.resize(len);
                  platon_get_call_output(ret.data());
              }else{
                  platon_panic();
              }
            }


            return ret;
        }

        // std::array<u128,2> newCandidateSecp256k1Point(bytes &bt){
        ACTION std::array<u128,2> newCandidateSecp256k1Point(bytes bt){
            newCandiArr.self()[0] = fieldHash(bt);
            newCandiArr.self()[1] = squareRoot(ySquared(newCandiArr.self()[0]));
            if (newCandiArr.self()[1] % 2 ==1){
                newCandiArr.self()[1] = fieldSize.self() - newCandiArr.self()[1];
            }
            return newCandiArr.self();
        }

        
        std::array<u128,2>  hashToCurve(u128 &pk1,u128 pk2,u128 input){
            // u128 pk1 = 0;//test
            std::string datapk1str = std::to_string(pk1);
            bytes datapk1;
            datapk1.insert(datapk1.begin(), datapk1str.begin(), datapk1str.end());

            rv.self() = newCandidateSecp256k1Point(datapk1);

            while(!isOnCurve(rv.self()[0],rv.self()[1])){
                // rv.self() = newCandidateSecp256k1Point(fromHex(pk1));
                rv.self() = newCandidateSecp256k1Point(datapk1);
            }
            return rv.self();
        }

        bool ecmulVerify(std::array<u128,2> multiplicand,u128 scalar,std::array<u128,2> product){
            platon_assert(scalar != 0, "bad value");
            u128 x = multiplicand[0];
            u128 v = multiplicand[1] % 2 == 0 ? 27 : 28;
            // platon::bytes scalarTimesX = fromHex(scalar*x%groupOrder.self());
            // Address actual = ecrecover(fromHex(0),v,fromHex(x),groupOrder.self());//待处理
            // Address exponent = Address(uint256(keccak256(abi.encodePacked(product))));//待处理

            return (address_str.self() == address_str.self());
        }

        std::array<u128,2> projectiveSub(u128 x1,u128 z1,u128 x2,u128 z2){
            u128 num1 = z2*x1%fieldSize.self();
            u128 num2 = (fieldSize.self()-x2)*z1%fieldSize.self();
            std::array<u128,2> arrays;
            arrays[0] = (num1+num2)%fieldSize;
            arrays[1] = (z1*z2)%fieldSize.self();
            
            return arrays;
        }

        // Returns x1/z1*x2/z2=(x1x2)/(z1z2), in projective coordinates on P¹(𝔽ₙ)
        std::array<u128,2> projectiveMul(u128 x1,u128 z1,u128 x2,u128 z2){
            std::array<u128,2> arr;
            arr[0] = x1*x2%fieldSize.self();
            arr[1] = z1*z2%fieldSize.self();  

            return arr;
        }

       std::array<u128,3> projectiveECAdd(u128 px,u128 py,u128 qx,u128 qy){
           std::array<u128,3> resArr;

        //    std::tuple<u128, u128> (z1,z2) = make_tuple(1,1);
           u128 z1 =1;
           u128 z2 =1;
           u128 lx = (qx+fieldSize.self()-py)%fieldSize.self();
           u128 lz = (qx+fieldSize.self()-px)%fieldSize.self();
           u128 dx;

           std::array<u128,2> sxdx1 = projectiveMul(lx, lz, lx, lz); // ((qy-py)/(qx-px))^2
           std::array<u128,2> sxdx2 = projectiveSub(resArr[0], dx, px, z1); // ((qy-py)/(qx-px))^2-px
           std::array<u128,2> sxdx3 = projectiveSub(resArr[0], dx, qx, z2); // ((qy-py)/(qx-px))^2-px-qx

           u128 dy;
           std::array<u128,2> sydy1 = projectiveSub(px, z1, sxdx3[0], dx); // px-sx
           std::array<u128,2> sydy2 = projectiveMul(resArr[1], dy, lx, lz); // ((qy-py)/(qx-px))(px-sx)
           std::array<u128,2> sydy3 = projectiveSub(resArr[1], dy, py, z1); // ((qy-py)/(qx-px))(px-sx)-py
           

           
           if(sxdx1[0] !=sydy1[0] || sxdx1[1] !=sydy1[1] || sxdx2[0] !=sydy2[0] || sxdx2[1] !=sydy2[1] || sxdx3[0] !=sydy3[0] || sxdx3[1] !=sydy3[1]){
              resArr[0] = (sxdx3[1]*dy%fieldSize.self());
              resArr[1] = (sydy3[1]*dx%fieldSize.self());
              resArr[2] = (dx*dy%fieldSize.self());
           }else{
               resArr[0] = sxdx3[0];
               resArr[1] = sydy3[0];
               resArr[2] = sxdx3[0];
           }
           return resArr;
       }

       std::array<u128,2> affineECAdd(std::array<u128,2> p1,std::array<u128,2> p2,u128 invZ){
         u128 x;
         u128 y;
         u128 z;
         std::array<u128,3> arr3 = projectiveECAdd(p1[0], p1[1], p2[0], p2[1]);

        platon_assert( z*invZ%fieldSize.self() == 1, "bad value");

        std::array<u128,2> resArr;
        resArr[0] = x*invZ%fieldSize.self();
        resArr[1] = y*invZ%fieldSize.self();
        return resArr;
       }

       bool verifyLinearCombinationWithGenerator(u128 c,std::array<u128,2> p,u128 s,Address lcWitness){
           u128 value = 12;
           u128 gas = 4712388;

           Address addr;
           auto address_info = make_address("lax1zqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvj4t2u");
           if(address_info.second){
             addr = address_info.first;
           }

           platon_assert( lcWitness != addr, "bad value");
           u128 v = (p[1] % 2 == 0) ? 27 : 28;

            std::string datastr = std::to_string((groupOrder.self() - p[0]*s%groupOrder.self()));
            bytes data;
            data.insert(data.begin(), datastr.begin(), datastr.end());

            platon::bytes pseudoHash = data; // -s*p[0]

            std::string data1str = std::to_string(c*p[0]%groupOrder.self());
            bytes data1;
            data1.insert(data1.begin(), data1str.begin(), data1str.end());

            platon::bytes pseudoSignature = data1; // c*p[0]
//            Address computed = ecrecover(pseudoHash, v, fromHex(p[0]), pseudoSignature);
            Address computed = addr;

            Address addr2;
            auto address_info2 = make_address(address_str.self());
            if(address_info2.second){
              addr2 = address_info2.first;
            }

            if (platon_call(addr2, pseudoSignature, value, gas)) {
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);

                Address addr3;
                auto address_info3 = make_address(address_str.self());
                if(address_info3.second){
                  addr3 = address_info3.first;
                }

                computed = addr3;
            }

           return computed == lcWitness;
       }

       std::array<u128,2> linearCombination(u128 c,std::array<u128,2> p1,std::array<u128,2> cp1Witness,
       u128 s, std::array<u128,2> p2,std::array<u128,2> sp2Witness,u128 zInv){
           std::array<u128,2> arr1;
        //    arr1[0] = p1;
        //    arr1[1] = c;
        //    arr1[2] = cp1Witness;
           platon_assert( ecmulVerify(p1,c,cp1Witness), "bad value");
           platon_assert( ecmulVerify(p2,s,sp2Witness), "bad value");
           return affineECAdd(cp1Witness, sp2Witness, zInv);
       }

       u128 scalarFromCurvePoints(std::array<u128,2> hash,std::array<u128,2> pk,std::array<u128,2> gamma,
       Address uWitness,std::array<u128,2> v){
            u128 res = 0; //输入参数进行转换
            platon::bytes inputBytes;

            std::vector<byte> result;
            result.resize(32);
            platon_sha256(inputBytes, result.data());

            //uint256(keccak256(abi.encodePacked(SCALAR_FROM_CURVE_POINTS_HASH_PREFIX,hash, pk, gamma, v, uWitness)));
            return res;
       }

       void verifyVRFProof(std::array<u128,2> pk,std::array<u128,2> gamma,u128 c,u128 s,
       u128 seed,Address uWitness,std::array<u128,2> cGammaWitness,
       std::array<u128,2> sHashWitness,u128 zInv){
           platon_assert(isOnCurve(pk[0],pk[1]), "public key is not on curve");
           platon_assert(isOnCurve(gamma[0],gamma[1]), "gamma is not on curve");
           platon_assert(isOnCurve(cGammaWitness[0],cGammaWitness[1]), "cGammaWitness is not on curve");
           platon_assert(isOnCurve(sHashWitness[0],sHashWitness[1]), "sHashWitness is not on curve");
           platon_assert(verifyLinearCombinationWithGenerator(c, pk, s, uWitness),"addr(c*pk+s*g)≠_uWitness");

           std::array<u128,2> hash = hashToCurve(pk[0],pk[1],seed);
           std::array<u128,2> v = linearCombination(c, gamma, cGammaWitness, s, hash, sHashWitness, zInv);
           u128 derivedC = scalarFromCurvePoints(hash, pk, gamma, uWitness, v);
           platon_assert(c == derivedC, "invalid proof");
       }

       ACTION h256 randomValueFromVRFProof(bytes proof){
            platon_assert(sizeof(proof) == PROOF_LENGTH.self(), "wrong proof length");

            std::array<u128,2> pk; // parse proof contents into these variables
            std::array<u128,2> gamma;
            std::array<u128,3> cSSeed;
            Address uWitness;
            std::array<u128,2> cGammaWitness;
            std::array<u128,2> sHashWitness;
            u128 zInv;
            //FIXME
            // (pk, gamma, cSSeed, uWitness, cGammaWitness, sHashWitness, zInv) = abi.decode( //不好翻译
            //     proof, (uint256[2], uint256[2], uint256[3], address, uint256[2],
            //             uint256[2], uint256));
             verifyVRFProof(
                 pk,
                 gamma,
                 cSSeed[0], // c
                 cSSeed[1], // s
                 cSSeed[2], // seed
                 uWitness,
                 cGammaWitness,
                 sHashWitness,
                 zInv
             );
             platon::bytes gammaBytes;
             h256 output;
             std::vector<byte> result;
             result.resize(32);
             platon_sha256(gammaBytes, result.data());
             return output;
       }
};

PLATON_DISPATCH(VRF, (init)(newCandidateSecp256k1Point)(randomValueFromVRFProof))

