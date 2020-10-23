//package network.platon.test.wasm.storage;
//
//import com.alibaba.fastjson.JSONArray;
//import com.platon.rlp.network.platon.test.datatypes.Uint32;
//import network.platon.test.evm.beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.autotest.utils.FileUtil;
//import network.platon.contracts.evm.SimpleStorage;
//import network.platon.contracts.wasm.SolSimulation;
//import network.platon.contracts.wasm.WasmStorage;
//import network.platon.utils.DataChangeUtil;
//import network.platon.utils.OneselfFileUtil;
//import org.junit.Test;
//import org.web3j.crypto.Credentials;
//import org.web3j.protocol.Web3j;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
//public class Wasm1StorageTest extends ContractPrepareTest {
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "qcxiao", showName = "wasm.storage.WasmStorageTest", sourcePrefix = "wasm")
//    public void test() {
//
//        try {
//            prepare();
//            WasmStorage wasmStorage1 = WasmStorage.deploy(web3j, transactionManager, provider,chainId).send();
//            String contractAddress1 = wasmStorage1.getContractAddress();
//            String transactionHash1 = wasmStorage1.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("deploy successfully.contractAddress1:" + contractAddress1
//                    + ", deployTxHash1:" + transactionHash1
//                    + ", gasUsed1:" + wasmStorage1.getTransactionReceipt().get().getGasUsed());
//
//            WasmStorage wasmStorage2 = WasmStorage.deploy(web3j, transactionManager, provider,chainId).send();
//            String contractAddress2 = wasmStorage2.getContractAddress();
//            String transactionHash2 = wasmStorage2.getTransactionReceipt().get().getTransactionHash();
//            collector.logStepPass("deploy successfully.contractAddress2:" + contractAddress2
//                    + ", deployTxHash2:" + transactionHash2
//                    + ", gasUsed2:" + wasmStorage2.getTransactionReceipt().get().getGasUsed());
//
//        } catch (Exception e) {
//            e.printStackTrace();
//            collector.logStepFail("", e.getMessage());
//        }
//    }
//}
