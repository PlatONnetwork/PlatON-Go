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
//public class Wasm2StorageTest extends ContractPrepareTest {
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "qcxiao", showName = "wasm.storage.WasmStorageTest", sourcePrefix = "wasm")
//    public void test() {
//
//        try {
//            String contractAddress2 = "0x8a4c045bf6c8c58ee2b9b335db3d01300c3c167b";
//            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "all_addr_and_private_keys_4000_evm.json").toUri().getPath());
//            String jsonContent = OneselfFileUtil.readFile(filePath);
//
//            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
//            ExecutorService executorService = Executors.newCachedThreadPool();
//            // 同时并发执行的线程数
//            final Semaphore semaphore = new Semaphore(100);
//            // 请求总数
//            CountDownLatch countDownLatch = new CountDownLatch(150000);
//            for (int j = 0; j < 50; j++) {
//                for (int i = 0; i < 3000; i++) {
//                    int finalI = i;
//                    int finalJ = j;
//                    executorService.execute(() -> {
//                        try {
//                            semaphore.acquire();
//                            chainId = Integer.valueOf(driverService.param.get("chainId"));
//                            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
//                            credentials = Credentials.create(jsonArray.getJSONObject(finalI).getString("private_key"));
//                            provider = new ContractGasProvider(new BigInteger(gasPrice), new BigInteger(gasLimit));
//                            transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//
//                            TransactionReceipt transactionReceipt = WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
//                            collector.logStepPass("transactionHash: " + transactionReceipt.getTransactionHash() + ",this time: " + finalI + ", j: " + finalJ);
//                            //WasmStorage.load(contractAddress2, web3j, transactionManager, provider).debug().send();
//                        } catch (Exception e) {
//                            //e.printStackTrace();
//                            collector.logStepFail("call fail. this time: " + finalI + ", j: " + finalJ, e.toString());
//                        } finally {
//                            semaphore.release();
//                            countDownLatch.countDown();
//                        }
//                    });
//                }
//            }
//
//
//
//            countDownLatch.await();
//            executorService.shutdown();
//
//        } catch (Exception e) {
//            e.printStackTrace();
//            collector.logStepFail("", e.getMessage());
//        }
//    }
//}
