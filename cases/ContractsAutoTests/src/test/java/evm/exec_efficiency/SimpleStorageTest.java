package evm.exec_efficiency;

import com.alibaba.fastjson.JSONArray;
import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.SimpleStorage;
import network.platon.contracts.SpaceComplexity;
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
import java.util.Arrays;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Semaphore;

public class SimpleStorageTest extends ContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "exec_efficiency.SimpleStorageTest-blockHash", sourcePrefix = "evm")
    public void test() {
        prepare();

        try {
            SimpleStorage simpleStorage = SimpleStorage.deploy(web3j, transactionManager, provider).send();
            contractAddress = simpleStorage.getContractAddress();
            collector.logStepPass("contract deploy successful. contractAddress:" + contractAddress);
            collector.logStepPass("deploy gas used:" + simpleStorage.getTransactionReceipt().get().getGasUsed());


            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "all_addr_and_private_keys_4000_evm.json").toUri().getPath());
            String jsonContent = OneselfFileUtil.readFile(filePath);

            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
            ExecutorService executorService = Executors.newCachedThreadPool();
            // 同时并发执行的线程数
            final Semaphore semaphore = new Semaphore(20);
            // 请求总数
            CountDownLatch countDownLatch = new CountDownLatch(jsonArray.size());

            for (int i = 0; i < jsonArray.size(); i++) {
                int finalI = i;
                executorService.execute(() -> {
                    try {
                        semaphore.acquire();
                        credentials = Credentials.create(jsonArray.getJSONObject(finalI).getString("private_key"));
                        transactionManager = new RawTransactionManager(web3j, credentials, chainId);

                        TransactionReceipt transactionReceipt = SimpleStorage.load(contractAddress, web3j, transactionManager, provider).hello().send();
                        collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash()
                                + ", currentBlockNumber:" + transactionReceipt.getBlockNumber());

                        byte[] hash = SimpleStorage.load(contractAddress, web3j, transactionManager, provider).hash().send();
                        collector.logStepPass("contract load successful, current time:"  + finalI + ", blockHash:" + DataChangeUtil.bytesToHex(hash));
                    } catch (Exception e) {
                        //e.printStackTrace();
                        //collector.logStepFail("call fail.", e.toString());
                    } finally {
                        semaphore.release();
                        countDownLatch.countDown();
                    }
                });

            }
            countDownLatch.await();
            executorService.shutdown();

        } catch (Exception e) {
            collector.logStepFail("The contract fail.", e.toString());
        }
    }
}
