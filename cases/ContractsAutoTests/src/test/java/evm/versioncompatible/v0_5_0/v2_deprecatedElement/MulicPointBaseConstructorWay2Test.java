package evm.versioncompatible.v0_5_0.v2_deprecatedElement;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.MulicPointBaseConstructorWay2;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title * 0.5.0不允许在同一继承层次结构中多次指定基类构造函数参数
 * 0.4.x可以同时使用2种方式，但如果2种方式都存在，优先选择修饰符方式
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class MulicPointBaseConstructorWay2Test extends ContractPrepareTest {

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
            author = "hudenian", showName = "version_compatible.0.5.0.constructorDeprecatedElementTest-弃用元素测试", sourcePrefix = "evm")
    public void update() {
        try {

            MulicPointBaseConstructorWay2 mulicPointBaseConstructorWay2 = MulicPointBaseConstructorWay2.deploy(web3j, transactionManager, provider, chainId, new BigInteger(initValue)).send();

            String contractAddress = mulicPointBaseConstructorWay2.getContractAddress();
            TransactionReceipt tx = mulicPointBaseConstructorWay2.getTransactionReceipt().get();

            collector.logStepPass("链上函数的初始值为："+initValue);

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + mulicPointBaseConstructorWay2.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt =mulicPointBaseConstructorWay2.update(new BigInteger(addValue)).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String afterValue = mulicPointBaseConstructorWay2.getA().send().toString();
            collector.logStepPass("链上函数的执行update后的值为："+afterValue);

            collector.assertEqual(String.valueOf(Integer.valueOf(initValue)*Integer.valueOf(initValue)+Integer.valueOf(addValue)),afterValue);
        } catch (Exception e) {
            collector.logStepFail("MulicPointBaseConstructorWay2Test process fail.", e.toString());
            e.printStackTrace();
        }
    }


}
