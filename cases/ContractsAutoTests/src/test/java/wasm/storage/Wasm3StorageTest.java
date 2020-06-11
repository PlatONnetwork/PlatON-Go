//package wasm.storage;
//
//import evm.beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.contracts.wasm.WasmStorage;
//import org.junit.Test;
//
///**
// * @title WasmStorageTest
// * @description 验证存储
// * @author qcxiao
// * @updateTime 2020/3/19 20:39
// */
//public class Wasm3StorageTest extends ContractPrepareTest {
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "qcxiao", showName = "wasm.storage.WasmStorageTest", sourcePrefix = "wasm")
//    public void test() {
//
//        try {
//            prepare();
//            String contractAddress1 = "lax1hahqha2mmujx0wktwh3t37s956hxxk5574m2c4";
//            WasmStorage wasmStorage3 = WasmStorage.deploy(web3j, transactionManager, provider,chainId).send();
//            String contractAddress3 = wasmStorage3.getContractAddress();
//            String transactionHash3 = wasmStorage3.getTransactionReceipt().get().getTransactionHash();
//
//            collector.logStepPass("deploy successfully.contractAddress3:" + contractAddress3
//                    + ", deployTxHash3:" + transactionHash3
//                    + ", gasUsed3:" + wasmStorage3.getTransactionReceipt().get().getGasUsed());
//
//            for (int i = 0; i < 11; i++) {
//                WasmStorage.load(contractAddress3, web3j, transactionManager, provider,chainId).action().send();
//                WasmStorage.load(contractAddress3, web3j, transactionManager, provider,chainId).debug().send();
//            }
//
//            for (int i = 0; i < 11; i++) {
//                WasmStorage.load(contractAddress1, web3j, transactionManager, provider,chainId).action().send();
//                WasmStorage.load(contractAddress1, web3j, transactionManager, provider,chainId).debug().send();
//            }
//        } catch (Exception e) {
//            e.printStackTrace();
//            collector.logStepFail("", e.getMessage());
//        }
//    }
//}
