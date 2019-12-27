#define WASM_EXPORT __attribute__((visibility("default")))

#include <stddef.h>
#include <stdint.h>

uint64_t platon_gas_price();
void platon_return(const uint8_t *value, size_t len);
void platon_block_hash(int64_t num,  uint8_t hash[32]);
uint64_t platon_block_number();
uint64_t platon_gas_limit();
int64_t platon_timestamp();
void platon_coinbase(uint8_t hash[20]);
uint8_t platon_balance(uint8_t addr[32], uint8_t balance[32]);
void platon_origin(uint8_t hash[20]);
void platon_caller(uint8_t hash[20]);
int32_t platon_transfer(const uint8_t* to, size_t toLen, uint8_t *amount, size_t len);
uint8_t platon_call_value(uint8_t val[32]);
void platon_address(uint8_t hash[20]);
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

WASM_EXPORT
void platon_gas_price_test() {
    uint64_t gas = platon_gas_price();
    platon_return((uint8_t*)&gas, sizeof(gas));
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
void platon_transfer_test() {
  uint8_t data[1024];
  size_t len = platon_get_input_length();
  platon_get_input(data);
  uint8_t value = 1;
  platon_transfer(data, len, &value, 1);
  platon_return(&value, 1);
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

