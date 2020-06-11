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
//import org.web3j.protocol.core.DefaultBlockParameterName;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.protocol.http.HttpService;
//import org.web3j.tx.RawTransactionManager;
//import org.web3j.tx.gas.ContractGasProvider;
//
//import java.io.IOException;
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
//            String contractAddress2 = "lax1keq09ca6esdzfs3q2m2sls6l5n5rwz4e5xstw8";
//
//            chainId = Integer.valueOf(driverService.param.get("chainId"));
//            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
//            provider = new ContractGasProvider(new BigInteger("6000000000"), new BigInteger(gasLimit));
//
//            // addr:lax10zgfmzecyevthcc80jcrxx6m954qhxj5fz9ful,private_key:7a6ffe562fcb8644510d92e42a11ada2f7ec64a0cad3b4f6ec543a76dc88eefb
////            credentials = Credentials.create("7a6ffe562fcb8644510d92e42a11ada2f7ec64a0cad3b4f6ec543a76dc88eefb");
////            transactionManager = new RawTransactionManager(web3j, credentials, chainId);
////            BigInteger endBalance1 = web3j.platonGetBalance("lax10zgfmzecyevthcc80jcrxx6m954qhxj5fz9ful", DefaultBlockParameterName.LATEST).send().getBalance();
////            collector.logStepPass("transferTo:lax10zgfmzecyevthcc80jcrxx6m954qhxj5fz9ful,endBalance:" + endBalance1);
////            TransactionReceipt transactionReceipt1 = WasmStorage.load(contractAddress2, web3j, transactionManager, provider, chainId).random_data().send();
////            collector.logStepPass("transactionHash: " + transactionReceipt1.getTransactionHash());
////            collector.logStepPass("gasuse: " + transactionReceipt1.getGasUsed());
////            WasmStorage.load(contractAddress2, web3j, transactionManager, provider, chainId).debug().send();
//
//
//            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "lax_bech32_all_addr_and_private_keys.json").toUri().getPath());
//            String jsonContent = OneselfFileUtil.readFile(filePath);
//            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
//            ExecutorService executorService = Executors.newCachedThreadPool();
//            // 同时并发执行的线程数
//            final Semaphore semaphore = new Semaphore(30);
//            // 请求总数
//            CountDownLatch countDownLatch = new CountDownLatch(jsonArray.size() * 1);
//
////            credentials = Credentials.create(driverService.param.get("privateKey"));
////            transactionManager = new RawTransactionManager(web3j, credentials, chainId);
////            BigInteger endBalance = web3j.platonGetBalance(driverService.param.get("address"), DefaultBlockParameterName.LATEST).send().getBalance();
////            collector.logStepPass("transferTo:" + driverService.param.get("address") + ",endBalance:" + endBalance);
////            TransactionReceipt transactionReceipt1 = WasmStorage.load(contractAddress2, web3j, transactionManager, provider,chainId).random_data().send();
////            collector.logStepPass("transactionHash: " + transactionReceipt1.getTransactionHash());
//
////            credentials = Credentials.create(jsonArray.getJSONObject(0).getString("private_key"));
////            transactionManager = new RawTransactionManager(web3j, credentials, chainId);
////            BigInteger endBalance = web3j.platonGetBalance(jsonArray.getJSONObject(0).getString("address"), DefaultBlockParameterName.LATEST).send().getBalance();
////            collector.logStepPass("transferTo:" + jsonArray.getJSONObject(0).getString("address") + ",endBalance:" + endBalance);
////            TransactionReceipt transactionReceipt1 = WasmStorage.load(contractAddress2, web3j, transactionManager, provider, chainId).random_data().send();
////            collector.logStepPass("transactionHash: " + transactionReceipt1.getTransactionHash());
////            collector.logStepPass("gasuse: " + transactionReceipt1.getGasUsed());
////            WasmStorage.load(contractAddress2, web3j, transactionManager, provider, chainId).debug().send();
//
//
//
//
//            for (int j = 0; j < 1; j++) {
//                int finalJ = j;
//                for (int i = 0; i < jsonArray.size(); i++) {
//                    int finalI = i;
//                    executorService.execute(() -> {
//                        try {
//                            semaphore.acquire();
//                            credentials = Credentials.create(jsonArray.getJSONObject(finalI).getString("private_key"));
//
//                            transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//
//                            TransactionReceipt transactionReceipt = WasmStorage.load(contractAddress2, web3j, transactionManager, provider,chainId).random_data().send();
//                            collector.logStepPass("transactionHash: " + transactionReceipt.getTransactionHash() + ",gasUsed：" + transactionReceipt.getGasUsed() +
//                                    ",this time: " + finalI + ", j: " + finalJ);
//                            //WasmStorage.load(contractAddress2, web3j, transactionManager, provider, chainId).debug().send();
//                        } catch (Exception e) {
//                            String address = jsonArray.getJSONObject(finalI).getString("address");
//                            BigInteger endBalance = null;
//                            try {
//                                endBalance = web3j.platonGetBalance(address, DefaultBlockParameterName.LATEST).send().getBalance();
//                            } catch (IOException ex) {
//                                ex.printStackTrace();
//                            }
//                            collector.logStepFail("call fail. this time: " + finalI + ", j: " + finalJ +
//                                    ", addr:" + address + ",endBalance：" + endBalance, e.toString());
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
