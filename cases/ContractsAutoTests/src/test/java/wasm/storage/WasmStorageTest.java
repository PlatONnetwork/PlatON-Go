package wasm.storage;

import com.alibaba.fastjson.JSONArray;
import com.platon.rlp.datatypes.Uint32;
import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.SimpleStorage;
import network.platon.contracts.wasm.SolSimulation;
import network.platon.contracts.wasm.WasmStorage;
import network.platon.utils.DataChangeUtil;
import network.platon.utils.OneselfFileUtil;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;
import java.nio.file.Paths;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Semaphore;

/**
 * @title WasmStorageTest
 * @description 验证存储
 * @author qcxiao
 * @updateTime 2020/3/19 20:39
 */
public class WasmStorageTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "wasm.storage.WasmStorageTest", sourcePrefix = "wasm")
    public void test() {

        try {
            prepare();
            WasmStorage wasmStorage1 = WasmStorage.deploy(web3j, transactionManager, provider).send();
            String contractAddress1 = wasmStorage1.getContractAddress();
            String transactionHash1 = wasmStorage1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("deploy successfully.contractAddress1:" + contractAddress1
                    + ", deployTxHash1:" + transactionHash1
                    + ", gasUsed1:" + wasmStorage1.getTransactionReceipt().get().getGasUsed());

            WasmStorage.load(contractAddress1, web3j, transactionManager, provider).random_data().send();

            WasmStorage wasmStorage2 = WasmStorage.deploy(web3j, transactionManager, provider).send();
            String contractAddress2 = wasmStorage2.getContractAddress();
            String transactionHash2 = wasmStorage2.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("deploy successfully.contractAddress2:" + contractAddress2
                    + ", deployTxHash2:" + transactionHash2
                    + ", gasUsed2:" + wasmStorage2.getTransactionReceipt().get().getGasUsed());


            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "all_addr_and_private_keys_4000_evm.json").toUri().getPath());
            String jsonContent = OneselfFileUtil.readFile(filePath);

            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
            ExecutorService executorService = Executors.newCachedThreadPool();
            // 同时并发执行的线程数
            final Semaphore semaphore = new Semaphore(10);
            // 请求总数
            CountDownLatch countDownLatch = new CountDownLatch(3000);
            for (int i = 0; i < 3000; i++) {
                int finalI = i;
                executorService.execute(() -> {
                    try {
                        semaphore.acquire();
                        chainId = Integer.valueOf(driverService.param.get("chainId"));
                        web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
                        credentials = Credentials.create(jsonArray.getJSONObject(finalI).getString("private_key"));
                        provider = new ContractGasProvider(new BigInteger(gasPrice), new BigInteger(gasLimit));
                        transactionManager = new RawTransactionManager(web3j, credentials, chainId);

                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).random_data().send();
                        WasmStorage.load(contractAddress2, web3j, transactionManager, provider).debug().send();

                        collector.logStepPass("this time:" + finalI);
                        semaphore.release();
                    } catch (Exception e) {
                        //e.printStackTrace();
                        collector.logStepFail("call fail. this time:" + finalI, e.toString());
                    }
                    countDownLatch.countDown();
                });
            }
            countDownLatch.await();
            executorService.shutdown();

            WasmStorage wasmStorage3 = WasmStorage.deploy(web3j, transactionManager, provider).send();
            String contractAddress3 = wasmStorage3.getContractAddress();
            String transactionHash3 = wasmStorage3.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("deploy successfully.contractAddress3:" + contractAddress3
                    + ", deployTxHash3:" + transactionHash3
                    + ", gasUsed3:" + wasmStorage3.getTransactionReceipt().get().getGasUsed());

            for (int i = 0; i < 11; i++) {
                WasmStorage.load(contractAddress3, web3j, transactionManager, provider).action().send();
                WasmStorage.load(contractAddress3, web3j, transactionManager, provider).debug().send();
            }

            for (int i = 0; i < 11; i++) {
                WasmStorage.load(contractAddress1, web3j, transactionManager, provider).action().send();
                WasmStorage.load(contractAddress1, web3j, transactionManager, provider).debug().send();
            }


        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("", e.getMessage());
        }
    }
}
