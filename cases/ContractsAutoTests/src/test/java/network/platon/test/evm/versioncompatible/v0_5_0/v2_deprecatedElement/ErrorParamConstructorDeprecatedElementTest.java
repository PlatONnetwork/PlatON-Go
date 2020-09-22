package network.platon.test.evm.versioncompatible.v0_5_0.v2_deprecatedElement;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import org.junit.Before;
import org.junit.Test;
import network.platon.contracts.evm.ErrorParamConstructor;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 不允许调用带参数但具有错误参数计数的构造函数。
 * 如果只想在不提供参数的情况下指定继承关系，请不要提供括号
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class ErrorParamConstructorDeprecatedElementTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    //链上函数初始值
    private String initValue = "10";

    //新增值
    private String addValue="20";



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "network.platon.test.evm.version_compatible.0.5.0.errorParamConstructorDeprecatedElementTest-弃用元素测试", sourcePrefix = "evm")
    public void update() {
        try {

            ErrorParamConstructor errorParamConstructor = ErrorParamConstructor.deploy(web3j, transactionManager, provider, chainId, new BigInteger(initValue)).send();

            String contractAddress = errorParamConstructor.getContractAddress();
            TransactionReceipt tx = errorParamConstructor.getTransactionReceipt().get();

            collector.logStepPass("链上函数的初始值为："+initValue);

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + errorParamConstructor.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt =errorParamConstructor.update(new BigInteger(addValue)).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String aValue = errorParamConstructor.getA().send().toString();
            collector.logStepPass("链上函数基类a值为："+aValue);
            collector.assertEqual("10",aValue);

            String bValue = errorParamConstructor.getB().send().toString();
            collector.logStepPass("链上函数子类b值为："+bValue);

            collector.assertEqual(String.valueOf(Integer.valueOf(initValue)+Integer.valueOf(addValue)),bValue);
        } catch (Exception e) {
            collector.logStepFail("ErrorParamConstructorDeprecatedElementTest process fail.", e.toString());
            e.printStackTrace();
        }
    }


}
