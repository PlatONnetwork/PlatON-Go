//package wasm.storage;
//
//import com.alibaba.fastjson.JSONArray;
//import evm.beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.autotest.utils.FileUtil;
//import network.platon.contracts.wasm.WasmStorage;
//import network.platon.utils.OneselfFileUtil;
//import org.junit.Test;
//import org.web3j.crypto.Credentials;
//import org.web3j.protocol.Web3j;
//import org.web3j.protocol.http.HttpService;
//import org.web3j.tx.RawTransactionManager;
//import org.web3j.tx.gas.ContractGasProvider;
//
//import java.math.BigInteger;
//import java.nio.file.Paths;
//import java.util.concurrent.CountDownLatch;
//import java.util.concurrent.ExecutorService;
//import java.util.concurrent.Executors;
//import java.util.concurrent.Semaphore;
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
//            String contractAddress1 = "0xae53d160854dccccb4c4a4c74037c75f1c030b58";
//            WasmStorage wasmStorage3 = WasmStorage.deploy(web3j, transactionManager, provider).send();
//            String contractAddress3 = wasmStorage3.getContractAddress();
//            String transactionHash3 = wasmStorage3.getTransactionReceipt().get().getTransactionHash();
//
//            collector.logStepPass("deploy successfully.contractAddress3:" + contractAddress3
//                    + ", deployTxHash3:" + transactionHash3
//                    + ", gasUsed3:" + wasmStorage3.getTransactionReceipt().get().getGasUsed());
//
//            for (int i = 0; i < 11; i++) {
//                WasmStorage.load(contractAddress3, web3j, transactionManager, provider).action().send();
//                WasmStorage.load(contractAddress3, web3j, transactionManager, provider).debug().send();
//            }
//
//            for (int i = 0; i < 11; i++) {
//                WasmStorage.load(contractAddress1, web3j, transactionManager, provider).action().send();
//                WasmStorage.load(contractAddress1, web3j, transactionManager, provider).debug().send();
//            }
//        } catch (Exception e) {
//            e.printStackTrace();
//            collector.logStepFail("", e.getMessage());
//        }
//    }
//}
