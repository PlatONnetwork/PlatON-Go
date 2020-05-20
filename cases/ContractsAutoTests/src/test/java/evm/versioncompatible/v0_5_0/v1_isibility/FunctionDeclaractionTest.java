package evm.versioncompatible.v0_5_0.v1_isibility;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.FunctionDeclaraction;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 函数可见性必须显式声明合约验证
 * @description:
 * @author: hudenian
 * @create: 2019/12/25
 */
public class FunctionDeclaractionTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    //链上函数初始值
    private String initValue;

    //新增值
    private String addValue="20";



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test1.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "evm.version_compatible.0.5.0.FunctionDeclaractionTest-update_public-可见性测试", sourcePrefix = "evm")
    public void update_public() {
        try {

            FunctionDeclaraction functionDeclaraction = FunctionDeclaraction.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = functionDeclaraction.getContractAddress();
            TransactionReceipt tx = functionDeclaraction.getTransactionReceipt().get();

            initValue = functionDeclaraction.getBalance().send().toString();
            collector.logStepPass("链上函数的初始值为："+initValue);

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + functionDeclaraction.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt =functionDeclaraction.update_public(new BigInteger(addValue)).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String afterValue = functionDeclaraction.getBalance().send().toString();
            collector.logStepPass("链上函数的执行update后的值为："+afterValue);

            collector.assertEqual(String.valueOf(Integer.valueOf(initValue)+Integer.valueOf(addValue)),afterValue);
        } catch (Exception e) {
            collector.logStepFail("FunctionDeclaractionTest update_public process fail.", e.toString());
            e.printStackTrace();
        }
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test2.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "version_compatible.0.5.0.FunctionDeclaractionTest-update_external-可见性测试", sourcePrefix = "evm")
    public void update_external() {
        try {

            FunctionDeclaraction functionDeclaraction = FunctionDeclaraction.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = functionDeclaraction.getContractAddress();
            TransactionReceipt tx = functionDeclaraction.getTransactionReceipt().get();

            initValue = functionDeclaraction.getBalance().send().toString();
            collector.logStepPass("执行update_external前balance值为："+initValue);

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + functionDeclaraction.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt =functionDeclaraction.update_external(new BigInteger(addValue)).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String afterValue = functionDeclaraction.getBalance().send().toString();
            collector.logStepPass("执行update_external后balance的值为："+afterValue);

            collector.assertEqual(String.valueOf(Integer.valueOf(initValue)+Integer.valueOf(addValue)),afterValue);
        } catch (Exception e) {
            collector.logStepFail("FunctionDeclaractionTest update_external process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
