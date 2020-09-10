#define TESTNET
#include <platon/platon.hpp>
#include "Gas.h"

using namespace platon;

CONTRACT GasPrice : public platon::Contract{
    public:
    ACTION void init() {}

    ACTION void platonGasPrice(){
        uint8_t gas_price[16] = {};
        platon_gas_price(gas_price);
    }

    ACTION void platonBlockHash(int64_t num){
        uint8_t hash[32] = {};
        platon_block_hash(num, hash);
    }

    ACTION void platonBlockNumber(){
        platon_block_number();
    }

    ACTION void platonGasLimit(){
        platon_gas_limit();
    }

    ACTION void platonGas(){
        platon_gas();
    }

    ACTION void platonTimestamp(){
        platon_timestamp();
    }

    ACTION void platonCoinbase(){
        uint8_t addr[20] = {};
        platon_coinbase(addr);
    }

    ACTION void platonBalance(const std::array<uint8_t, 20> &addr){
        uint8_t balance[16] = {};
        platon_balance(addr.data(), balance);
    }

    ACTION void platonOrigin(){
        uint8_t addr[20] = {};
        platon_coinbase(addr);
    }

    ACTION void platonCaller(){
        uint8_t addr[20] = {};
        platon_caller(addr);
    }

    ACTION void platonCallValue(){
        uint8_t addr[20] = {};
        platon_call_value(addr);
    }

    ACTION void platonAddress(){
        uint8_t addr[20] = {};
        platon_address(addr);
    }

   ACTION void platonSha3(const std::vector<uint8_t> & src){
       uint8_t dest[32] = {};
       platon_sha3(src.data(), src.size(), dest, sizeof(dest));
   }

    ACTION void platonCallerNonce(){
        platon_caller_nonce();
    }

    ACTION void platonTransfer(const std::array<uint8_t, 20> &to){
        uint8_t amount = 100;
        platon_transfer(to.data(), &amount, 1);
    }

   ACTION void platonSetState(const std::vector<uint8_t> &key, const std::vector<uint8_t> &value){
       platon_set_state(key.data(), key.size(), value.data(), value.size());
   }

   ACTION void platonGetStateLength(const std::vector<uint8_t> &key){
       platon_get_state_length(key.data(), key.size());
   }

   ACTION void platonGetState(const std::vector<uint8_t> &key, size_t length){
       uint8_t *value = (uint8_t *) malloc(length);
       platon_get_state(key.data(), key.size(), value, length);
   }

    ACTION void platonGetInputLength(){
        platon_get_input_length();
    }

    ACTION void platonGetInput() {
        uint8_t input[100] = {};
        platon_get_input(input);
    }

    ACTION void platonGetCallOutputLength() {
        platon_get_call_output_length();
    }

    ACTION void platonGetCallOutput() {
        uint8_t value[100] = {};
        platon_get_call_output(value);
    }

    ACTION void platonReturn(size_t length) {
        uint8_t *value = (uint8_t *) malloc(length);
        platon_return(value, length);
    }

    ACTION void platonRevert() {
        platon_revert();
    }

    ACTION void platonPanic() {
        platon_panic();
    }

    ACTION void platonDebug(size_t length) {
        uint8_t *src = (uint8_t *) malloc(length);
        platon_debug(src, length);
    }

    ACTION void platonCall(const Address & contract_address, const std::string &method){
        bytes paras = cross_call_args(method);
        bytes value_bytes = value_to_bytes(0U);
        bytes gas_bytes = value_to_bytes(0U);
        platon_call(contract_address.data(), paras.data(), paras.size(), value_bytes.data(),
                    value_bytes.size(), gas_bytes.data(), gas_bytes.size());
    }

    ACTION void platonDelegateCall(const Address & contract_address, const std::string &method){
        bytes paras = cross_call_args(method);
        bytes gas_bytes = value_to_bytes(0U);
        platon_delegate_call(contract_address.data(), paras.data(), paras.size(),
                             gas_bytes.data(), gas_bytes.size());
    }

    ACTION void platonDestory(const Address & to) {
        platon_destroy(to.data());
    }

    ACTION void platonMigrate(const Address & contract_address){
        size_t code_length = platon_contract_code_length(contract_address.data());
        if(0 == code_length) return;
        bytes contract_code(code_length);
        int32_t get_result = platon_contract_code(contract_address.data(), contract_code.data(), code_length);
        if(get_result <= 0 ) return;

        // value and gas
        bytes value_bytes = value_to_bytes(0U);
        bytes gas_bytes = value_to_bytes(0U);

        // deploy agrs
        std::vector<byte> magic_number = {0x00, 0x61, 0x73, 0x6d};
        bytes init_rlp = cross_call_args("init");
        RLPSize rlps;
        rlps << contract_code << init_rlp;
        RLPStream stream;
        stream.appendPrefix(magic_number);
        stream.reserve(rlps.size() + 4);
        stream.appendList(2);
        stream << contract_code << init_rlp;
        bytesRef result = stream.out();

        // migrate
        Address return_address;
        platon_migrate(return_address.data(), result.data(), result.size(), value_bytes.data(), value_bytes.size(), gas_bytes.data(), gas_bytes.size());
    }

    ACTION void platonMigrateClone(const Address & contract_address){
        // value and gas
        bytes value_bytes = value_to_bytes(0U);
        bytes gas_bytes = value_to_bytes(0U);

        // init args
        bytes init_rlp = cross_call_args("init");

        // migrate
        Address return_address;
        platon_clone_migrate(contract_address.data(), return_address.data(), init_rlp.data(), init_rlp.size(), value_bytes.data(), value_bytes.size(), gas_bytes.data(), gas_bytes.size());
    }

    ACTION void platonEvent(const std::vector<uint8_t> &topic, const std::vector<uint8_t> &args){
        platon_event(topic.data(), topic.size(), args.data(), args.size());
    }

    ACTION void platonEcrecover(const std::array<uint8_t, 32> &hash, const std::vector<uint8_t> &sig) {
        uint8_t addr[20] = {};
        platon_ecrecover(hash.data(), sig.data(), sig.size(), addr);
    }

    ACTION void platonRipemd160(const std::vector<uint8_t> &src) {
        uint8_t hash[20] = {};
        platon_ripemd160(src.data(), src.size(), hash);
    }

    ACTION void platonSha256(const std::vector<uint8_t> &src) {
        uint8_t hash[32] = {};
        platon_sha256(src.data(), src.size(), hash);
    }

    ACTION void platonRlpU128Size(uint64_t heigh, uint64_t low) {
        rlp_u128_size(heigh, low);
    }

    ACTION void platonRlpU128(uint64_t heigh, uint64_t low) {
        void *dest = malloc(10);
        platon_rlp_u128(heigh, low, dest);
    }

    ACTION void platonRlpBytesSize(const std::vector<uint8_t> & data) {
        rlp_bytes_size(data.data(), data.size());
    }

    ACTION void platonRlpBytes(const std::vector<uint8_t> & data) {
        void *dest = malloc(10);
        platon_rlp_bytes(data.data(), data.size(), dest);
    }

    ACTION void platonRlpListSize(size_t len) {
        rlp_list_size(len);
    }

    ACTION void platonRlpList(const std::vector<uint8_t> & data) {
        void *dest = malloc(10);
        platon_rlp_list(data.data(), data.size(), dest);
    }

    ACTION void platonContractCodeLength(const std::array<uint8_t, 20> &addr) {
        platon_contract_code_length(addr.data());
    }

    ACTION void platonContractCode(const std::array<uint8_t, 20> &addr) {
        uint8_t code[100] = {};
        platon_contract_code(addr.data(), code, 100);
    }

    ACTION void platonDeploy(const Address & contract_address){
        size_t code_length = platon_contract_code_length(contract_address.data());
        if(0 == code_length) return;
        bytes contract_code(code_length);
        int32_t code_result = platon_contract_code(contract_address.data(), contract_code.data(), code_length);
        if(code_result <= 0) return;

        // value and gas
        bytes value_bytes = value_to_bytes(0U);
        bytes gas_bytes = value_to_bytes(0U);

        // deploy args
        std::vector<byte> magic_number = {0x00, 0x61, 0x73, 0x6d};
        bytes init_rlp = cross_call_args("init");
        RLPSize rlps;
        rlps << contract_code << init_rlp;
        RLPStream stream;
        stream.appendPrefix(magic_number);
        stream.reserve(rlps.size() + 4);
        stream.appendList(2);
        stream << contract_code << init_rlp;
        bytesRef result = stream.out();

        // deploy contract
        Address return_address;
        platon_deploy(return_address.data(), result.data(), result.size(),
                            value_bytes.data(), value_bytes.size(), gas_bytes.data(),
                            gas_bytes.size());
    }

    ACTION void platonClone(const Address & contract_address){
        // value and gas
        bytes value_bytes = value_to_bytes(0U);
        bytes gas_bytes = value_to_bytes(0U);

        // init args
        bytes init_rlp = cross_call_args("init");

        // clone contract
        Address return_address;
        platon_clone(contract_address.data(), return_address.data(), init_rlp.data(),
                        init_rlp.size(), value_bytes.data(), value_bytes.size(),
                        gas_bytes.data(), gas_bytes.size());
    }

};

PLATON_DISPATCH(GasPrice, (init)(platonGasPrice)(platonBlockHash)(platonBlockNumber)(platonGasLimit)
(platonGas)(platonTimestamp)(platonCoinbase)(platonBalance)(platonOrigin)(platonCaller)(platonCallValue)
(platonAddress)(platonSha3)(platonCallerNonce)(platonTransfer)(platonSetState)(platonGetStateLength)
(platonGetState)(platonGetInputLength)(platonGetInput)(platonGetCallOutputLength)(platonGetCallOutput)
(platonReturn)(platonRevert)(platonPanic)(platonDebug)(platonCall)(platonDelegateCall)(platonDestory)
(platonMigrate)(platonMigrateClone)(platonEvent)(platonEcrecover)(platonRipemd160)(platonSha256)
(platonRlpU128Size)(platonRlpU128)(platonRlpBytesSize)(platonRlpBytes)(platonRlpListSize)
(platonRlpList)(platonContractCodeLength)(platonContractCode)(platonDeploy)(platonClone))



