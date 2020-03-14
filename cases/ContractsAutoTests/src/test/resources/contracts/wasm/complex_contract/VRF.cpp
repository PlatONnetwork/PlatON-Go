#include <platon/platon.hpp>
#include <vector>
#include <string>

using namespace platon;

/**
 * VRF.sol合约翻译
 * */

CONTRACT VRF : public platon::Contract{

    private:
        platon::StorageType<"groupOrder"_n, u128> groupOrder;
        platon::StorageType<"fieldSize"_n, uint64_t> fieldSize;
        platon::StorageType<"wordlength"_n, uint64_t> wordLengthBytes;
        platon::StorageType<"bigmod"_n,std::array<uint64_t,6>> bigModExpContractInputs;
        platon::StorageType<"output"_n, uint64_t> output;
        platon::StorageType<"sqrtpower"_n, uint64_t> sqrtPower;
        platon::StorageType<"newcandi"_n,std::array<uint64_t,2>> newCandiArr;
        platon::StorageType<"rv"_n,std::array<uint64_t,2>> rv;
        platon::StorageType<"scalar"_n, uint64_t> SCALAR_FROM_CURVE_POINTS_HASH_PREFIX;
        platon::StorageType<"scalaprr"_n, uint64_t> VRF_RANDOM_OUTPUT_HASH_PREFIX;
        platon::StorageType<"prooflen"_n, uint64_t> PROOF_LENGTH;
        platon::StorageType<"addstr"_n, std::string> address_str;

    public:
        ACTION void init(){
            groupOrder.self() = "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"_u128; 
            fieldSize.self() =  "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"_u128;
            wordLengthBytes.self() = 32; //0x20;
            SCALAR_FROM_CURVE_POINTS_HASH_PREFIX.self() = 2;
            VRF_RANDOM_OUTPUT_HASH_PREFIX.self() = 3;
            address_str.self() = "0x000000000000000000000000000000000000002";
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

        uint64_t bigModExp(uint64_t base,uint64_t exponent){
            uint64_t callResult;
            bigModExpContractInputs.self()[0] = wordLengthBytes.self();
            bigModExpContractInputs.self()[1] = wordLengthBytes.self();
            bigModExpContractInputs.self()[2] = wordLengthBytes.self();
            bigModExpContractInputs.self()[3] = base;
            bigModExpContractInputs.self()[4] = exponent;
            bigModExpContractInputs.self()[5] = fieldSize.self();

            uint64_t value = 12;
            uint64_t gas = 4712388;
            std::string addr = "0x000000000000000000000000000000000000005";
            platon::bytes params = platon::cross_call_args("data", bigModExpContractInputs.self());

            if (platon_call(Address(addr), params, value, gas)) {
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);
                base = std::stoul(str);
            }
            return base;

        }

        uint64_t getSqrtPower(){
            return (fieldSize.self()+1) >>2;
        }

        uint64_t squareRoot(uint64_t x){
            return bigModExp(x,getSqrtPower());
        }

        uint64_t ySquared(uint64_t x){
            uint64_t value1 = x*x%fieldSize.self();
            uint64_t xCubed = x*value1%fieldSize.self();
            return (xCubed+7)%fieldSize.self();
        }

        bool isOnCurve(uint64_t p1,uint64_t p2){
            return ySquared(p1) == p2*p2%fieldSize.self();
        }

        //c++实现keccak256算法实现未开发
        // uint64_t fieldHash(bytes &b){
        uint64_t fieldHash(bytes &bt){
            uint64_t base;
            platon::bytes res = keccak256(bt);
            std::string str = toHex(res); //byte转string
            base = std::stoul(str);//string转uint64
            while(base >= fieldSize.self()){
                res = keccak256(res);
                std::string str = toHex(res);
                base = std::stoul(str);
            }
//            std::string str = toHex(res);
            return base;
        }

        bytes keccak256(bytes &bt){
            uint64_t value = 12;
            uint64_t gas = 4712388;
            std::string target_addr = "0x000000000000000000000000000000000000002";
            platon::bytes ret;

            if (platon_call(Address(target_addr), bt, value, gas)) {
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
            }else{
                platon_panic();
            }
            return ret;
        }

        // std::array<uint64_t,2> newCandidateSecp256k1Point(bytes &bt){
        std::array<uint64_t,2> newCandidateSecp256k1Point(bytes bt){
            newCandiArr.self()[0] = fieldHash(bt);
            newCandiArr.self()[1] = squareRoot(ySquared(newCandiArr.self()[0]));
            if (newCandiArr.self()[1] % 2 ==1){
                newCandiArr.self()[1] = fieldSize.self() - newCandiArr.self()[1];
            }
            return newCandiArr.self();
        }

        
        std::array<uint64_t,2>  hashToCurve(uint64_t &pk1,uint64_t pk2,uint64_t input){
            // uint64_t pk1 = 0;//test
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

        bool ecmulVerify(std::array<uint64_t,2> multiplicand,uint64_t scalar,std::array<uint64_t,2> product){
            platon_assert(scalar != 0, "bad value");
            uint64_t x = multiplicand[0];
            uint64_t v = multiplicand[1] % 2 == 0 ? 27 : 28;
            // platon::bytes scalarTimesX = fromHex(scalar*x%groupOrder.self());
            // Address actual = ecrecover(fromHex(0),v,fromHex(x),groupOrder.self());//待处理
            // Address exponent = Address(uint256(keccak256(abi.encodePacked(product))));//待处理
            Address actual = Address(address_str.self());
            Address exponent = Address(address_str.self());
            return (actual == exponent);
        }

        std::array<uint64_t,2> projectiveSub(uint64_t x1,uint64_t z1,uint64_t x2,uint64_t z2){
            uint64_t num1 = z2*x1%fieldSize.self();
            uint64_t num2 = (fieldSize.self()-x2)*z1%fieldSize.self();
            std::array<uint64_t,2> arrays;
            arrays[0] = (num1+num2)%fieldSize;
            arrays[1] = (z1*z2)%fieldSize.self();
            
            return arrays;
        }

        // Returns x1/z1*x2/z2=(x1x2)/(z1z2), in projective coordinates on P¹(𝔽ₙ)
        std::array<uint64_t,2> projectiveMul(uint64_t x1,uint64_t z1,uint64_t x2,uint64_t z2){
            std::array<uint64_t,2> arr;
            arr[0] = x1*x2%fieldSize.self();
            arr[1] = z1*z2%fieldSize.self();  

            return arr;
        }

       std::array<uint64_t,3> projectiveECAdd(uint64_t px,uint64_t py,uint64_t qx,uint64_t qy){
           std::array<uint64_t,3> resArr;

        //    std::tuple<uint64_t, uint64_t> (z1,z2) = make_tuple(1,1);
           uint64_t z1 =1;
           uint64_t z2 =1;
           uint64_t lx = (qx+fieldSize.self()-py)%fieldSize.self();
           uint64_t lz = (qx+fieldSize.self()-px)%fieldSize.self();
           uint64_t dx;

           std::array<uint64_t,2> sxdx1 = projectiveMul(lx, lz, lx, lz); // ((qy-py)/(qx-px))^2
           std::array<uint64_t,2> sxdx2 = projectiveSub(resArr[0], dx, px, z1); // ((qy-py)/(qx-px))^2-px
           std::array<uint64_t,2> sxdx3 = projectiveSub(resArr[0], dx, qx, z2); // ((qy-py)/(qx-px))^2-px-qx 

           uint64_t dy;
           std::array<uint64_t,2> sydy1 = projectiveSub(px, z1, sxdx3[0], dx); // px-sx
           std::array<uint64_t,2> sydy2 = projectiveMul(resArr[1], dy, lx, lz); // ((qy-py)/(qx-px))(px-sx)
           std::array<uint64_t,2> sydy3 = projectiveSub(resArr[1], dy, py, z1); // ((qy-py)/(qx-px))(px-sx)-py
           

           
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

       std::array<uint64_t,2> affineECAdd(std::array<uint64_t,2> p1,std::array<uint64_t,2> p2,uint64_t invZ){
         uint64_t x;
         uint64_t y;
         uint64_t z;
         std::array<uint64_t,3> arr3 = projectiveECAdd(p1[0], p1[1], p2[0], p2[1]);

        platon_assert( z*invZ%fieldSize.self() == 1, "bad value");

        std::array<uint64_t,2> resArr;
        resArr[0] = x*invZ%fieldSize.self();
        resArr[1] = y*invZ%fieldSize.self();
        return resArr;
       }

       bool verifyLinearCombinationWithGenerator(uint64_t c,std::array<uint64_t,2> p,uint64_t s,Address lcWitness){
           uint64_t value = 12;
           uint64_t gas = 4712388;

           platon_assert( lcWitness != Address(0), "bad value");
           uint64_t v = (p[1] % 2 == 0) ? 27 : 28;

            std::string datastr = std::to_string((groupOrder.self() - p[0]*s%groupOrder.self()));
            bytes data;
            data.insert(data.begin(), datastr.begin(), datastr.end());

            platon::bytes pseudoHash = data; // -s*p[0]

            std::string data1str = std::to_string(c*p[0]%groupOrder.self());
            bytes data1;
            data1.insert(data1.begin(), data1str.begin(), data1str.end());

            platon::bytes pseudoSignature = data1; // c*p[0]
//            Address computed = ecrecover(pseudoHash, v, fromHex(p[0]), pseudoSignature);
            Address computed = Address(0);
            if (platon_call(Address(address_str.self()), pseudoSignature, value, gas)) {
                platon::bytes ret;
                size_t len = platon_get_call_output_length();
                ret.resize(len);
                platon_get_call_output(ret.data());
                std::string str = toHex(ret);
                computed = Address(str);
            }

           return computed == lcWitness;
       }

       std::array<uint64_t,2> linearCombination(uint64_t c,std::array<uint64_t,2> p1,std::array<uint64_t,2> cp1Witness,
       uint64_t s, std::array<uint64_t,2> p2,std::array<uint64_t,2> sp2Witness,uint64_t zInv){
           std::array<uint64_t,2> arr1;
        //    arr1[0] = p1;
        //    arr1[1] = c;
        //    arr1[2] = cp1Witness;
           platon_assert( ecmulVerify(p1,c,cp1Witness), "bad value");
           platon_assert( ecmulVerify(p2,s,sp2Witness), "bad value");
           return affineECAdd(cp1Witness, sp2Witness, zInv);
       }

       uint64_t scalarFromCurvePoints(std::array<uint64_t,2> hash,std::array<uint64_t,2> pk,std::array<uint64_t,2> gamma,
       Address uWitness,std::array<uint64_t,2> v){
            uint64_t res = 0;
            //uint256(keccak256(abi.encodePacked(SCALAR_FROM_CURVE_POINTS_HASH_PREFIX,hash, pk, gamma, v, uWitness)));
            return res;
       }

       void verifyVRFProof(std::array<uint64_t,2> pk,std::array<uint64_t,2> gamma,uint64_t c,uint64_t s,
       uint64_t seed,Address uWitness,std::array<uint64_t,2> cGammaWitness,
       std::array<uint64_t,2> sHashWitness,uint64_t zInv){
           platon_assert(isOnCurve(pk[0],pk[1]), "public key is not on curve");
           platon_assert(isOnCurve(gamma[0],gamma[1]), "gamma is not on curve");
           platon_assert(isOnCurve(cGammaWitness[0],cGammaWitness[1]), "cGammaWitness is not on curve");
           platon_assert(isOnCurve(sHashWitness[0],sHashWitness[1]), "sHashWitness is not on curve");
           platon_assert(verifyLinearCombinationWithGenerator(c, pk, s, uWitness),"addr(c*pk+s*g)≠_uWitness");

           std::array<uint64_t,2> hash = hashToCurve(pk[0],pk[1],seed);
           std::array<uint64_t,2> v = linearCombination(c, gamma, cGammaWitness, s, hash, sHashWitness, zInv);
           uint64_t derivedC = scalarFromCurvePoints(hash, pk, gamma, uWitness, v);
           platon_assert(c == derivedC, "invalid proof");
       }

       uint64_t randomValueFromVRFProof(bytes proof){
            platon_assert(sizeof(proof) == PROOF_LENGTH.self(), "wrong proof length");

            std::array<uint64_t,2> pk; // parse proof contents into these variables
            std::array<uint64_t,2> gamma;
            std::array<uint64_t,3> cSSeed;
            Address uWitness;
            std::array<uint64_t,2> cGammaWitness;
            std::array<uint64_t,2> sHashWitness;
            uint64_t zInv;
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
//             output = uint256(keccak256(abi.encode(VRF_RANDOM_OUTPUT_HASH_PREFIX, gamma)));
             platon::bytes gammaBytes;
             uint64_t output = std::stoul(toHex(keccak256(gammaBytes)));
             return zInv;
       }
};

PLATON_DISPATCH(VRF, (init)(newCandidateSecp256k1Point))

