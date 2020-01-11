#define WASM_EXPORT __attribute__((visibility("default")))

#include <stddef.h>
#include <stdint.h>

uint64_t platon_gas_price();
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
int32_t platon_static_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
int32_t platon_destroy();
int32_t platon_migrate(uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
void platon_event(const uint8_t* args, size_t argsLen);
void platon_event1(const uint8_t* topic, size_t topicLen, const uint8_t* args, size_t argsLen);
void platon_event2(const uint8_t* topic1, size_t topic1Len, const uint8_t* topic2, size_t topic2Len, const uint8_t* args, size_t argsLen);
void platon_event3(const uint8_t* topic1, size_t topic1Len, const uint8_t* topic2, size_t topic2Len, const uint8_t* topic3, size_t topic3Len, uint8_t* args, size_t argsLen);
void platon_event4(const uint8_t* topic1, size_t topic1Len, const uint8_t* topic2, size_t topic2Len, const uint8_t* topic3, size_t topic3Len, const uint8_t* topic4, size_t topic4Len, uint8_t* args, size_t argsLen);





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
  platon_call(addr, &data, datalen, &value, 1, &gas, 5);
}

WASM_EXPORT
void platon_delegate_call_contract_test () {
    uint8_t addr[20] = {1, 2, 4}; // don't change it
    uint8_t data[1024];
    size_t datalen = platon_get_input_length();
    platon_get_input(data);
    uint8_t gas = 100000;
    platon_delegate_call(addr, &data, datalen, &gas, 5);

}

WASM_EXPORT
void platon_static_call_contract_test () {
   uint8_t addr[20] = {1, 2, 4}; // don't change it
   uint8_t data[1024];
   size_t datalen = platon_get_input_length();
   platon_get_input(data);
   uint8_t gas = 100000;
   platon_static_call(addr, &data, datalen, &gas, 5);
}

WASM_EXPORT
void platon_destroy_contract_test () {
    platon_destroy();
}

WASM_EXPORT
void platon_migrate_contract_test () {
    uint8_t newAddr[20];
    uint8_t data[1024];
    size_t datalen = platon_get_input_length();
    platon_get_input(data);
    uint8_t gas = 100000;
    uint8_t value = 2;
    platon_migrate(newAddr, &data, datalen, &value, 1, &gas, 5);
    platon_return(newAddr, 20);
}

WASM_EXPORT
void platon_event_test () {
    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);
    platon_event(data, len);
}

WASM_EXPORT
void platon_event1_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    uint8_t *topic = (uint8_t*)"t1";
    size_t topicLen = sizeof(topic);
    platon_event1(topic, topicLen, data, len);
}

WASM_EXPORT
void platon_event2_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    uint8_t *topic1 = (uint8_t*)"t1";
    size_t topic1Len = sizeof(topic1);

    uint8_t *topic2 = (uint8_t*)"t2";
    size_t topic2Len = sizeof(topic2);

    platon_event2(topic1, topic1Len, topic2, topic2Len, data, len);
}


WASM_EXPORT
void platon_event3_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    uint8_t *topic1 = (uint8_t*)"t1";
    size_t topic1Len = sizeof(topic1);

    uint8_t *topic2 = (uint8_t*)"t2";
    size_t topic2Len = sizeof(topic2);

    uint8_t *topic3 = (uint8_t*)"t3";
    size_t topic3Len = sizeof(topic3);

    platon_event3(topic1, topic1Len, topic2, topic2Len, topic3, topic3Len, data, len);
}


WASM_EXPORT
void platon_event4_test () {

    uint8_t data[1024];
    size_t len = platon_get_input_length();
    platon_get_input(data);

    uint8_t *topic1 = (uint8_t*)"t1";
    size_t topic1Len = sizeof(topic1);

    uint8_t *topic2 = (uint8_t*)"t2";
    size_t topic2Len = sizeof(topic2);

    uint8_t *topic3 = (uint8_t*)"t3";
    size_t topic3Len = sizeof(topic3);

    uint8_t *topic4 = (uint8_t*)"t4";
    size_t topic4Len = sizeof(topic4);

    platon_event4(topic1, topic1Len, topic2, topic2Len, topic3, topic3Len, topic4, topic4Len, data, len);
}