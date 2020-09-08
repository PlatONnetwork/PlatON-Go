//package network.platon.test.wasm.storage;
//
//import com.alibaba.fastjson.JSONArray;
//import com.platon.rlp.network.platon.test.datatypes.Uint32;
//import network.platon.test.evm.beforetest.ContractPrepareTest;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.autotest.utils.FileUtil;
//import network.platon.contracts.wasm.SolSimulation;
//import network.platon.contracts.wasm.WasmStorage;
//import network.platon.utils.OneselfFileUtil;
//import org.junit.Test;
//import org.web3j.crypto.Credentials;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.tx.RawTransactionManager;
//
//import java.nio.file.Paths;
//import java.util.concurrent.CountDownLatch;
//import java.util.concurrent.ExecutorService;
//import java.util.concurrent.Executors;
//import java.util.concurrent.Semaphore;
//
///**
// * @title SolSimulationTest
// * @description 验证存储
// * @author qcxiao
// * @updateTime 2020/3/16 20:39
// */
//public class SolSimulationTest extends ContractPrepareTest {
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "qcxiao", showName = "wasm.storage.SolSimulationTest", sourcePrefix = "wasm")
//    public void test() {
//
//        try {
//            prepare();
//            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "all_addr_and_private_keys_4000_evm.json").toUri().getPath());
//            String jsonContent = OneselfFileUtil.readFile(filePath);
//
//            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
//            ExecutorService executorService = Executors.newCachedThreadPool();
//            // 同时并发执行的线程数
//            final Semaphore semaphore = new Semaphore(1);
//            // 请求总数
//            CountDownLatch countDownLatch = new CountDownLatch(10);
//            for (int i = 0; i < 10; i++) {
//                int finalI = i;
//                executorService.execute(() -> {
//                    try {
//                        semaphore.acquire();
//                        credentials = Credentials.create(jsonArray.getJSONObject(finalI + 1000).getString("private_key"));
//                        transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//
//                        SolSimulation solSimulation = SolSimulation.deploy(web3j, transactionManager, provider,chainId, Uint32.of(32), Uint32.of(0)).send();
//                        String contractAddress = solSimulation.getContractAddress();
//                        String transactionHash = solSimulation.getTransactionReceipt().get().getTransactionHash();
//                        collector.logStepPass("deploy successfully.contractAddress:" + contractAddress
//                                + ", hash:" + transactionHash
//                                + ", gasUsed:" + solSimulation.getTransactionReceipt().get().getGasUsed());
//
//                        TransactionReceipt transactionReceiptAction = SolSimulation.load(contractAddress, web3j, transactionManager, provider,chainId).action().send();
//
//                        TransactionReceipt transactionReceipt = SolSimulation.load(contractAddress, web3j, transactionManager, provider,chainId).debug().send();
//
//                        collector.logStepPass("action and debug are gas used: " + transactionReceiptAction.getGasUsed() + "," + transactionReceipt.getGasUsed());
//
//                        semaphore.release();
//                    } catch (Exception e) {
//                        e.printStackTrace();
//                        collector.logStepFail("call fail.", e.toString());
//                    }
//                    countDownLatch.countDown();
//                });
//            }
//            countDownLatch.await();
//            executorService.shutdown();
//        } catch (Exception e) {
//            e.printStackTrace();
//            collector.logStepFail("", e.getMessage());
//        }
//    }
//}
