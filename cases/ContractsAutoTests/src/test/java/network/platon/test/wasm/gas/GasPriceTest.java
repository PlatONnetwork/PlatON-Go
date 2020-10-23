package network.platon.test.wasm.gas;

import com.alibaba.fastjson.JSONObject;
import com.platon.rlp.datatypes.*;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.*;
import network.platon.utils.DataChangeUtil;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.rlp.RlpEncoder;
import org.web3j.rlp.RlpList;
import org.web3j.rlp.RlpString;
import org.web3j.rlp.RlpType;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;
import org.web3j.utils.Numeric;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * @title GasTest
 * @description Gas测试
 * @author qcxiao
 */
public class GasPriceTest extends WASMContractPrepareTest {

    private static final String GASMAPSTR = "{\"platonGasPrice\": 2, \"platonBlockHash\": 20, \"platonBlockNumber\": 2," +
            " \"platonGasLimit\": 2, \"platonGas\": 2, \"platonTimestamp\": 2, \"platonCoinbase\": 2, " +
            "\"platonBalance\": 400, \"platonOrigin\": 2, \"platonCaller\": 2, \"platonCallValue\": 2, " +
            "\"platonAddress\": 2, \"platonSha3\": 150, \"platonCallerNonce\": 2, \"platonTransfer\": 7444, " +
            "\"platonGetStateLength\": 200, \"platonGetState\": 200, \"platonGetInputLength\": 2, " +
            "\"platonGetInput\": 2, \"platonGetCallOutputLength\": 2, \"platonGetCallOutput\": 2, \"platonReturn\": 302, " +
            "\"platonRevert\": 0, \"platonPanic\": 0, \"platonDebug\": 310, \"platonCall\": 17429, \"platonDelegateCall\": 17410," +
            " \"platonDestory\": 1293, \"platonMigrate\": 153200, \"platonMigrateClone\": 1293, \"platonEvent\": 1293," +
            "\"platonEcrecover\": 3000, \"platonRipemd160\": 840, \"platonSha256\": 84, \"platonRlpU128Size\": 2, " +
            "\"platonRlpU128\": 17, \"platonRlpBytesSize\": 2, \"platonRlpBytes\": 35," +
            " \"platonRlpListSize\": 2, \"platonRlpList\": 35, \"platonContractCodeLength\": 200, " +
            "\"platonContractCode\": 200, \"platonDeploy\": 110000, \"platonClone\": 110000}";
    private static final JSONObject GASMAP = JSONObject.parseObject(GASMAPSTR);
    private BigInteger getGasValue(TransactionReceipt transactionReceipt, GasPrice gasPrice, String method){
        collector.logStepPass("transactionReceipt: " + JSONObject.toJSONString(transactionReceipt));
        List<GasPrice.GasUsedEventResponse> eventList = gasPrice.getGasUsedEvents(transactionReceipt);
        collector.logStepPass("eventList: " + JSONObject.toJSONString(eventList));

        //topic encoding
        byte[] realFunctionNameBytes = new byte[32];
        byte[] functionNameBytes = method.getBytes();
        int len = functionNameBytes.length;
        if(len < 32){
            for(int i=len-1,j=31;i>=0;i--,j--){
                realFunctionNameBytes[j] = functionNameBytes[i];
            }
        }else{
            System.arraycopy(functionNameBytes, 0, realFunctionNameBytes, 0, 32);
        }

        //filter topic
        String topicFuncName = Numeric.toHexString(realFunctionNameBytes);
        for(GasPrice.GasUsedEventResponse event:eventList){
            if(event.topic.equals(topicFuncName)){
                return event.arg1.value;
            }
        }
        return BigInteger.ZERO;
    }

    private void checkGas(BigInteger gasValue, String method){
        collector.logStepPass("gas of " + method + ": " + gasValue);
        collector.assertTrue((Math.abs(GASMAP.getIntValue(method) - gasValue.intValue())) < 100);
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.GasTest",sourcePrefix = "wasm")
    public void test() {

        try {
            prepare();
            GasPrice gasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = gasPrice.getContractAddress();
            String transactionHash = gasPrice.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("GasPrice deploy successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + gasPrice.getTransactionReceipt().get().getGasUsed());

            gasPrice = GasPrice.load(contractAddress, web3j, transactionManager, provider,chainId);

            TransactionReceipt transactionReceipt;
            BigInteger gas;
            switch (Integer.valueOf(driverService.param.get("seq"))){
                case 1:{
                    //查询当前交易的 gas price
                    transactionReceipt = gasPrice.platonGasPrice().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_gas_price");
                    this.checkGas(gas, "platonGasPrice");
                    break;
                }
                case 2:{
                    //根据blockNumber查询blockHash
                    BigInteger currentBlockNumber = web3j.platonBlockNumber().send().getBlockNumber();
//                    transactionReceipt = gasPrice.platonBlockHash(Int64.of(BigInteger.TEN)).send();
                    transactionReceipt = gasPrice.platonBlockHash(Int64.of(currentBlockNumber)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_block_hash");
                    this.checkGas(gas, "platonBlockHash");
                    break;
                }
                case 3:{
                    //查询当前blockNumber
                    transactionReceipt = gasPrice.platonBlockNumber().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_block_number");
                    this.checkGas(gas, "platonBlockNumber");
                    break;
                }
                case 4:{
                    //查询当前tx的gasLimit
                    transactionReceipt = gasPrice.platonGasLimit().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_gas_limit");
                    this.checkGas(gas, "platonGasLimit");
                    break;
                }
                case 5:{
                    //查询当前tx的gas
                    transactionReceipt = gasPrice.platonGas().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_gas");
                    this.checkGas(gas, "platonGas");
                    break;
                }
                case 6:{
                    //查询当前block的时间戳
                    transactionReceipt = gasPrice.platonTimestamp().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_timestamp");
                    this.checkGas(gas, "platonTimestamp");
                    break;
                }
                case 7:{
                    //查询当前block的coinbase
                    transactionReceipt = gasPrice.platonCoinbase().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_coinbase");
                    this.checkGas(gas, "platonCoinbase");
                    break;
                }
                case 8:{
                    //根据addr查询addr的余额
                    Uint8[] addr = new Uint8[20];
                    transactionReceipt = gasPrice.platonBalance(addr).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_balance");
                    this.checkGas(gas, "platonBalance");
                    break;
                }
                case 9:{
                    //查询tx非原始发送者
                    transactionReceipt = gasPrice.platonOrigin().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_origin");
                    this.checkGas(gas, "platonOrigin");
                    break;
                }
                case 10:{
                    //查询合约的上一级调用者账户地址
                    transactionReceipt = gasPrice.platonCaller().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_caller");
                    this.checkGas(gas, "platonCaller");
                    break;
                }
                case 11:{
                    //查询合约上一级调用者的余额
                    transactionReceipt = gasPrice.platonCallValue().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_call_value");
                    this.checkGas(gas, "platonCallValue");
                    break;
                }
                case 12:{
                    //查询当前合约的地址
                    transactionReceipt = gasPrice.platonAddress().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_address");
                    this.checkGas(gas, "platonAddress");
                    break;
                }
                case 13:{
                    //对输入的内容做 sha3
                    transactionReceipt = gasPrice.platonSha3("helloworldworldhello".getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_sha3");
                    this.checkGas(gas, "platonSha3");
                    break;
                }
                case 14:{
                    //查询当前合约的调用方账户的nonce
                    transactionReceipt = gasPrice.platonCallerNonce().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_caller_nonce");
                    this.checkGas(gas, "platonCallerNonce");
                    break;
                }
                case 15:{
                    //转账 (金额转移)
                    Uint8[] addr = new Uint8[20];
                    transactionReceipt = gasPrice.platonTransfer(addr).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_transfer");
                    this.checkGas(gas, "platonTransfer");
                    break;
                }
                case 16:{
                    //存储单个账户的store (调用SetState()), 新增数据
                    String key = RandomStringUtils.randomAlphanumeric(20);
                    String value = RandomStringUtils.randomAlphanumeric(20);
                    transactionReceipt = gasPrice.platonSetState(key.getBytes(), value.getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_set_state");
                    collector.logStepPass("gas of platonSetState: " + gas);
                    collector.assertTrue((Math.abs(40000 - gas.intValue())) < 100);

                    //存储单个账户的store (调用SetState()), 修改数据
                    value = RandomStringUtils.randomAlphanumeric(20);
                    transactionReceipt = gasPrice.platonSetState(key.getBytes(), value.getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_set_state");
                    collector.logStepPass("gas of platonSetState: " + gas);
                    collector.assertTrue((Math.abs(5000 - gas.intValue())) < 100);

                    //存储单个账户的store (调用SetState()), 删除数据
                    value = "";
                    transactionReceipt = gasPrice.platonSetState(key.getBytes(), value.getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_set_state");
                    collector.logStepPass("gas of platonSetState: " + gas);
                    collector.assertTrue((Math.abs(10000 - gas.intValue())) < 100);
                    break;
                }
                case 17:{
                    //根据key获取 store 中value 的长度
                    transactionReceipt = gasPrice.platonGetStateLength("helloworldworldhello".getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_state_length");
                    this.checkGas(gas, "platonGetStateLength");
                    break;
                }
                case 18:{
                    //获取单个账户的store (调用GetState())
                    transactionReceipt = gasPrice.platonGetState("helloworldworldhello".getBytes(), Uint32.of(10)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_state");
                    this.checkGas(gas, "platonGetState");
                    break;
                }
                case 19:{
                    //获取输入参数长度
                    transactionReceipt = gasPrice.platonGetInputLength().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_input_length");
                    this.checkGas(gas, "platonGetInputLength");
                    break;
                }
                case 20:{
                    //获取输入参数
                    transactionReceipt = gasPrice.platonGetInput().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_input");
                    this.checkGas(gas, "platonGetInput");
                    break;
                }
                case 21:{
                    //获取 跨合约调用返回 output数据长度
                    transactionReceipt = gasPrice.platonGetCallOutputLength().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_call_output_length");
                    this.checkGas(gas, "platonGetCallOutputLength");
                    break;
                }
                case 22:{
                    //获取 跨合约调用返回 output数据
                    transactionReceipt = gasPrice.platonGetCallOutput().send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_get_call_output");
                    this.checkGas(gas, "platonGetCallOutput");
                    break;
                }
                case 23:{
                    //获取当前合约返回值
                    transactionReceipt = gasPrice.platonReturn(Uint32.of(100)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_return");
                    this.checkGas(gas, "platonReturn");
                    break;
                }
                case 24:{
                    //合约终止指令
                    try{
                        transactionReceipt = gasPrice.platonRevert().send();
                        collector.logStepFail("合约终止调用需要报错", "合约终止调用未报错");
                    }catch (Exception e){

                    }
                    break;
                }
                case 25:{
                    //合约异常恐慌中断指令
                    try{
                        transactionReceipt = gasPrice.platonPanic().send();
                        collector.logStepFail("合约中断调用需要报错", "合约中断调用未报错");
                    }catch (Exception e){

                    }
                    break;
                }
                case 26:{
                    //打印合约调试信息
                    transactionReceipt = gasPrice.platonDebug(Uint32.of(100)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_debug");
                    this.checkGas(gas, "platonDebug");
                    break;
                }
                case 27:{
                    //跨合约普通调用
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonCall(wasmAddress, "platonTimestamp").send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_call");
                    collector.logStepPass("gas of platonCall: " + gas);
                    collector.assertTrue((Math.abs(GASMAP.getIntValue("platonCall") - gas.intValue())) < 10000);
                    break;
                }
                case 28:{
                    //跨合约代理调用
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonDelegateCall(wasmAddress, "platonTimestamp").send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_delegate_call");
                    collector.logStepPass("gas of platonDelegateCall: " + gas);
                    collector.assertTrue((Math.abs(GASMAP.getIntValue("platonDelegateCall") - gas.intValue())) < 10000);
                    break;
                }
                case 29:{
                    //合约销毁，地址使用过
                    WasmAddress wasmAddress = new WasmAddress(gasPrice.getContractAddress());
                    gasPrice.platonBlockNumber();
                    transactionReceipt = gasPrice.platonDestory(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_destroy");
                    collector.logStepPass("gas of platonDestory: " + gas);
                    collector.assertTrue((Math.abs(5000 - gas.intValue())) < 100);

                    //合约销毁，地址未使用过
                    gasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    wasmAddress = new WasmAddress("lax1uqug0zq7rcxddndleq4ux2ft3tv6dqljphydrl");
                    transactionReceipt = gasPrice.platonDestory(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_destroy");
                    collector.logStepPass("gas of platonDestory: " + gas);
                    collector.assertTrue((Math.abs(5000 - gas.intValue())) < 100);
                    break;
                }
                case 30:{
                    //合约迁移
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonMigrate(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_migrate");
                    collector.logStepPass("gas of platonMigrate: " + gas);
                    collector.assertTrue((Math.abs(153200 - gas.intValue())) < 10000);
                    break;
                }
                case 31:{
                    //合约克隆迁移
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonMigrateClone(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_clone_migrate");
                    collector.logStepPass("gas of platonMigrateClone: " + gas);
                    collector.assertTrue((Math.abs(153200 - gas.intValue())) < 10000);
                    break;
                }
                case 32:{
                    //合约事件
                    RlpString rlpString = RlpString.create(RandomStringUtils.randomAlphanumeric(20).getBytes());
                    byte[] topic = RlpEncoder.encode(rlpString);
                    List<String> args = Arrays.asList("hello", "world");
                    RlpType[] values = new RlpType[args.size()];
                    for(int i=0;i<args.size();i++){
                        values[i] = RlpString.create(args.get(i));
                    }
                    byte[] endcodedArgs = RlpEncoder.encode(new RlpList(values));
                    transactionReceipt = gasPrice.platonEvent(endcodedArgs, topic).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_event");
                    this.checkGas(gas, "platonEvent");
                    break;
                }
                case 33:{
                    //根据Hash和sig解出对应的addr
                    String hexHash = "4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45";
                    String hexSignature = "f4128988cbe7df8315440adde412a8955f7f5ff9a5468a791433727f82717a6753bd71882079522207060b681fbd3f5623ee7ed66e33fc8e581f442acbcf6ab800";
                    byte[] signature = DataChangeUtil.hexToByteArray(hexSignature);
                    Uint8[] hash = new Uint8[32];
                    for(int i=0;i<32;i++){
                        hash[i] = Uint8.of(Integer.parseInt(hexHash.substring(2*i, 2*i+1),16));
                    }
                    transactionReceipt = gasPrice.platonEcrecover(hash, signature).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_ecrecover");
                    this.checkGas(gas, "platonEcrecover");
                    break;
                }
                case 34:{
                    //ripemd160算法求Hash
                    transactionReceipt = gasPrice.platonRipemd160(RandomStringUtils.randomAlphanumeric(33).getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_ripemd160");
                    this.checkGas(gas, "platonRipemd160");
                    break;
                }
                case 35:{
                    //sha256算法求Hash
                    transactionReceipt = gasPrice.platonSha256(RandomStringUtils.randomAlphanumeric(33).getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_sha256");
                    this.checkGas(gas, "platonSha256");
                    break;
                }
                case 36:{
                    //计算u128数据在rlp之后的数据长度
                    transactionReceipt = gasPrice.platonRlpU128Size(Uint64.of(100000), Uint64.of(100000)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "rlp_u128_size");
                    this.checkGas(gas, "platonRlpU128Size");
                    break;
                }
                case 37:{
                    //计算u128数据在rlp之后的数据值
                    transactionReceipt = gasPrice.platonRlpU128(Uint64.of(100000), Uint64.of(100000)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_rlp_u128");
                    this.checkGas(gas, "platonRlpU128");
                    break;
                }
                case 38:{
                    //计算bytes数据在rlp之后的数据长度
                    transactionReceipt = gasPrice.platonRlpBytesSize("helloworld".getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "rlp_bytes_size");
                    this.checkGas(gas, "platonRlpBytesSize");
                    break;
                }
                case 39:{
                    //计算bytes数据在rlp之后的数据值
                    transactionReceipt = gasPrice.platonRlpBytes("helloworld".getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_rlp_bytes");
                    this.checkGas(gas, "platonRlpBytes");
                    break;
                }
                case 40:{
                    //计算list数据在rlp之后的数据长度
                    transactionReceipt = gasPrice.platonRlpListSize(Uint32.of(10000)).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "rlp_list_size");
                    this.checkGas(gas, "platonRlpListSize");
                    break;
                }
                case 41:{
                    //计算list数据在rlp之后的数据值
                    transactionReceipt = gasPrice.platonRlpList("helloworld".getBytes()).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_rlp_list");
                    this.checkGas(gas, "platonRlpList");
                    break;
                }
                case 42:{
                    //获取合约代码长度
                    Uint8[] addr = new Uint8[20];
                    transactionReceipt = gasPrice.platonContractCodeLength(addr).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_contract_code_length");
                    this.checkGas(gas, "platonContractCodeLength");
                    break;
                }
                case 43:{
                    //获取合约代码
                    Uint8[] addr = new Uint8[20];
                    transactionReceipt = gasPrice.platonContractCode(addr).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_contract_code");
                    this.checkGas(gas, "platonContractCode");
                    break;
                }
                case 44:{
                    //创建合约
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonDeploy(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_deploy");
                    collector.logStepPass("gas of platonDeploy: " + gas);
                    collector.assertTrue((Math.abs(110000 - gas.intValue())) < 10000);
                    break;
                }
                case 45:{
                    //克隆合约
                    GasPrice newGasPrice = GasPrice.deploy(web3j, transactionManager, provider,chainId).send();
                    WasmAddress wasmAddress = new WasmAddress(newGasPrice.getContractAddress());

                    transactionReceipt = gasPrice.platonClone(wasmAddress).send();
                    gas = this.getGasValue(transactionReceipt, gasPrice, "platon_clone");
                    collector.logStepPass("gas of platonClone: " + gas);
                    collector.assertTrue((Math.abs(110000 - gas.intValue())) < 10000);
                    break;
                }
            }

        } catch (Exception e) {
            collector.logStepFail("Gas price test fail.", e.getMessage());
            e.printStackTrace();
        }
    }

}
