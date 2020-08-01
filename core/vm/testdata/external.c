#define WASM_EXPORT __attribute__((visibility("default")))

#include <stddef.h>
#include <stdint.h>

uint8_t platon_gas_price(uint8_t gas_price[32]);
void platon_return(const uint8_t *value, size_t len);
void platon_block_hash(int64_t num,  uint8_t hash[32]);
uint64_t platon_block_number();
uint64_t platon_gas_limit();
uint64_t platon_gas();
int64_t platon_timestamp();
void platon_coinbase(uint8_t addr[20]);
uint8_t platon_balance(const uint8_t  addr[20], uint8_t balance[32]);
void platon_origin(uint8_t addr[20]);
void platon_caller(uint8_t addr[20]);
int32_t platon_transfer(const uint8_t to[20], const uint8_t *amount, size_t len);
uint8_t platon_call_value(uint8_t val[32]);
void platon_address(uint8_t addr[20]);
void platon_sha3(const uint8_t *src, size_t srcLen, uint8_t *dest, size_t destLen);
uint64_t platon_caller_nonce();


// c++
size_t platon_get_input_length();
void platon_get_input(const uint8_t *value);

void platon_set_state(const uint8_t* key, size_t klen, const uint8_t *value, size_t vlen);
size_t platon_get_state_length(const uint8_t* key, size_t klen);
size_t platon_get_state(const uint8_t* key, size_t klen, uint8_t *value, size_t vlen);
size_t platon_get_call_output_length();
void platon_get_call_output(const uint8_t *value);
void platon_revert();
void platon_panic();
void platon_debug(uint8_t *dst, size_t len);


int32_t platon_call(const uint8_t to[20], const uint8_t *args, size_t argsLen, const uint8_t *value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
int32_t platon_delegate_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
//int32_t platon_static_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
int32_t platon_destroy(const uint8_t to[20]);
int32_t platon_migrate(uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
void platon_event(const uint8_t* topic, size_t topicLen, const uint8_t* args, size_t argsLen);


void platon_sha256(const uint8_t *input, uint32_t input_len, uint8_t hash[32]);
void platon_ripemd160(const uint8_t *input, uint32_t input_len, uint8_t addr[20]);
int32_t platon_ecrecover(const uint8_t hash[32], const uint8_t* sig, const uint8_t sig_len, uint8_t addr[20]);

// u128
size_t rlp_u128_size(uint64_t heigh, uint64_t low);
void platon_rlp_u128(uint64_t heigh, uint64_t low, void * dest);

// bytes
size_t rlp_bytes_size(const void *data, size_t len);
void platon_rlp_bytes(const void *data, size_t len, void * dest);

// list
size_t rlp_list_size(size_t len);
void platon_rlp_list(const void *data, size_t len, void * dest);

WASM_EXPORT
void platon_gas_price_test() {
    uint8_t gas[32] = {0};
    uint8_t len = platon_gas_price(gas);
    platon_return(gas, len);
}

WASM_EXPORT
void platon_block_hash_test() {
  uint8_t hash[32];
  platon_block_hash(0, hash);
  platon_return(hash, sizeof(hash));
}

WASM_EXPORT
void platon_block_number_test() {
  uint64_t num = platon_block_number();
  platon_return((uint8_t*)&num, sizeof(num));
}
WASM_EXPORT
void platon_gas_limit_test() {
  uint64_t num = platon_gas_limit();
  platon_return((uint8_t*)&num, sizeof(num));
}

WASM_EXPORT
void platon_gas_test() {
  uint64_t num = platon_gas();
  platon_return((uint8_t*)&num, sizeof(num));
}

WASM_EXPORT
void platon_timestamp_test() {
  uint64_t num = platon_timestamp();
  platon_return((uint8_t*)&num, sizeof(num));
}

WASM_EXPORT
void platon_coinbase_test() {
  uint8_t hash[20];
  platon_coinbase(hash);
  platon_return(hash, sizeof(hash));
}

WASM_EXPORT
void platon_balance_test() {
  uint8_t hash[32] = {1};
  uint8_t balance[32] = {0};
  uint8_t len = platon_balance(hash, balance);
  platon_return(balance, len);
}

WASM_EXPORT
void platon_origin_test() {
  uint8_t hash[20];
  platon_origin(hash);
  platon_return(hash, sizeof(hash));
}

WASM_EXPORT
void platon_caller_test() {
  uint8_t hash[20];
  platon_caller(hash);
  platon_return(hash, sizeof(hash));
}

WASM_EXPORT
void platon_call_value_test() {
  uint8_t hash[32];
  uint8_t len = platon_call_value(hash);
  platon_return(hash, len);
}

WASM_EXPORT
void platon_address_test() {
  uint8_t hash[20];
  platon_address(hash);
  platon_return(hash, sizeof(hash));
}

WASM_EXPORT
void platon_sha3_test() {
  uint8_t data[1024];
  size_t len = platon_get_input_length();
  platon_get_input(data);
  uint8_t hash[32];
  platon_sha3(data, len, hash, 32);
  platon_return(hash, sizeof(hash));
}


WASM_EXPORT
void platon_caller_nonce_test() {
  uint64_t num = platon_caller_nonce();
  platon_return((uint8_t*)&num, sizeof(num));
}

WASM_EXPORT
void platon_set_state_test() {
  uint8_t data[1024];
  size_t len = platon_get_input_length();
  platon_get_input(data);
  platon_set_state((uint8_t*)"key", 3, data, len);
}

WASM_EXPORT
void platon_get_state_test() {
  uint8_t data[1024];
  size_t len = platon_get_state_length((uint8_t*)"key", 3);
  platon_get_state((uint8_t*)"key", 3, data, 1024);
  platon_return(data, len);
}

WASM_EXPORT
void platon_get_call_output_test() {
  uint8_t data[1024];
  size_t len = platon_get_call_output_length();
  platon_get_call_output(data);
  platon_return(data, len);
}

WASM_EXPORT
void platon_revert_test() {
  platon_revert();
}

WASM_EXPORT
void platon_panic_test() {
  platon_panic();
}

WASM_EXPORT
void platon_debug_test() {
  uint8_t data[1024];
  size_t len = platon_get_input_length();
  platon_get_input(data);
  platon_debug(data, len);
}

WASM_EXPORT
void platon_transfer_test() {
  uint8_t data[1024];
  size_t len = platon_get_input_length();
  platon_get_input(data);
  uint8_t value = 1;
  platon_transfer(data, &value, 1);
  platon_return(&value, 1);
}

WASM_EXPORT
void platon_call_contract_test() {
  uint8_t addr[20] = {1, 2, 4}; // don't change it
  uint8_t data[1024];
  size_t datalen = platon_get_input_length();
  platon_get_input(data);
  uint8_t gas = 100000;
  uint8_t value = 2;
  platon_call(addr, data, datalen, &value, 1, &gas, 5);
}

WASM_EXPORT
void platon_delegate_call_contract_test () {
    uint8_t addr[20] = {1, 2, 4}; // don't change it
    uint8_t data[1024];
    size_t datalen = platon_get_input_length();
    platon_get_input(data);
    uint8_t gas = 100000;
    platon_delegate_call(addr, data, datalen, &gas, 5);

}

//WASM_EXPORT
//void platon_static_call_contract_test () {
//   uint8_t addr[20] = {1, 2, 4}; // don't change it
//   uint8_t data[1024];
//   size_t datalen = platon_get_input_length();
//   platon_get_input(data);
//   uint8_t gas = 100000;
//   platon_static_call(addr, &data, datalen, &gas, 5);
//}

WASM_EXPORT
void platon_destroy_contract_test () {
    uint8_t addr[20] = {1, 2, 6};
    platon_destroy(addr);
}

WASM_EXPORT
void platon_migrate_contract_test () {
    uint8_t newAddr[20];
    uint8_t data[1024];
    size_t datalen = platon_get_input_length();
    platon_get_input(data);
    uint8_t gas = 100000;
    uint8_t value = 2;
    platon_migrate(newAddr, data, datalen, &value, 1, &gas, 5);
    platon_return(newAddr, 20);
}

WASM_EXPORT
void platon_event0_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    // empty topic
    uint8_t topics[1] = {0};

    platon_event(topics, 0, data, len);
}

WASM_EXPORT
void platon_event3_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    // rlp([topic1, topic2, topic3])
    uint8_t topics[10] = {201, 130, 116, 49, 130, 116, 50, 130, 116, 51};

    platon_event(topics, 10, data, len);
}
void platon_sha256(const uint8_t *input, uint32_t input_len, uint8_t hash[32]);
void platon_ripemd160(const uint8_t *input, uint32_t input_len, uint8_t addr[20]);
int32_t platon_ecrecover(const uint8_t hash[32], const uint8_t* sig, const uint8_t sig_len, uint8_t addr[20]);
WASM_EXPORT
void platon_sha256_test() {
    uint8_t input[3] = {1,2,3};
//    uint8_t hash[32] = {3,144,88,198,242,192,203,73,44,83,59,10,77,20,239,119,204,15,120,171,204,206,213,40,125,132,161,162,1,28,251,129};
    uint8_t res[32] = {0};
    platon_sha256(input, 3, res);
    platon_return(res, 32);
}

WASM_EXPORT
void platon_ripemd160_test() {
    uint8_t input[3] = {1,2,3};
//    uint8_t addr[20] = {121,249,1,218,38,9,240,32,173,173,191,46,95,104,161,108,140,63,125,87};
    uint8_t res[20] = {0};
    platon_ripemd160(input, 3, res);
    platon_return(res, 20);
}

WASM_EXPORT
void platon_ecrecover_test() {
    uint8_t hash[32] = {65,177,160,100,151,82,175,27,40,179,220,41,161,85,110,238,120,30,74,76,58,31,127,83,249,15,168,52,222,9,140,77};
    uint8_t sig[65] = {209,85,233,67,5,175,126,7,221,140,50,135,62,92,3,203,149,201,224,89,96,239,133,190,156,7,246,113,218,88,199,55,24,193,154,220,57,122,33,26,169,232,126,81,158,32,56,197,163,182,88,97,141,179,53,247,79,128,11,142,12,254,239,68,1};
//    uint8_t addr[20] = {151,14,129,40,171,131,78,142,172,23,171,142,56,18,240,16,103,140,247,145};
    uint8_t res[20] = {0};
    platon_ecrecover(hash, sig, 65, res);
    platon_return(res, 20);
}

WASM_EXPORT
void rlp_u128_size_test(){
  uint64_t heigh = 0x0123456789abcdefULL;
  uint64_t low = 0xfedcba9876543210ULL;
  size_t append_length = rlp_u128_size(heigh, low);
  uint8_t res[8] = {0};
  for(int i = 0; i < 8; i++){
    res[i] = append_length >> (i * 8);
  }
  platon_return(res, 8);
}

WASM_EXPORT
void platon_rlp_u128_test(){
  uint64_t heigh = 0x0123456789abcdefULL;
  uint64_t low = 0xfedcba9876543210ULL;
  uint8_t res[17] = {0};
  platon_rlp_u128(heigh, low, res);
  platon_return(res, 17);
}

WASM_EXPORT
void rlp_bytes_size_test(){
  uint8_t data[16] = {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f};
  size_t append_length = rlp_bytes_size(data, 16);
  uint8_t res[8] = {0};
  for(int i = 0; i < 8; i++){
    res[i] = append_length >> (i * 8);
  }
  platon_return(res, 8);
}

WASM_EXPORT
void platon_rlp_bytes_test(){
  uint8_t data[16] = {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f};
  uint8_t res[17] = {0};
  platon_rlp_bytes(data, 16, res);
  platon_return(res, 17);
}

WASM_EXPORT
void rlp_list_size_test(){
  uint8_t data[16] = {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f};
  size_t append_length = rlp_list_size(sizeof(data));
  uint8_t res[8] = {0};
  for(int i = 0; i < 8; i++){
    res[i] = append_length >> (i * 8);
  }
  platon_return(res, 8);
}

WASM_EXPORT
void platon_rlp_list_test(){
  uint8_t data[16] = {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f};
  uint8_t res[17] = {0};
  platon_rlp_list(data, 16, res);
  platon_return(res, 17);
}