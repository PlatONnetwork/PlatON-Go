package network.platon.test.evm.versioncompatible.v0_5_0.v1_isibility;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ConstructorInternalDeclaractionSub;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 构造函数可见性必须显示声明public或者internal
 * @description:
 * @author: hudenian
 * @create: 2019/12/27
 */
public class ConstructorInternalDeclaractionTest extends ContractPrepareTest {

    //链上函数初始值
    private String initValue;

    //新增值
    private String addValue;

    @Before
    public void before() {
        this.prepare();
        initValue = driverService.param.get("initValue");
        addValue = driverService.param.get("addValue");
    }



    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ConstructorInternalDeclaractionTest-可见性测试", sourcePrefix = "evm")
    public void update() {
        try {
            ConstructorInternalDeclaractionSub constructorInternalDeclaractionSub = ConstructorInternalDeclaractionSub.deploy(web3j, transactionManager, provider, chainId, new BigInteger(initValue)).send();

            String contractAddress = constructorInternalDeclaractionSub.getContractAddress();
            TransactionReceipt tx = constructorInternalDeclaractionSub.getTransactionReceipt().get();

            collector.logStepPass("ConstructorInternalDeclaractionSub deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + constructorInternalDeclaractionSub.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt =constructorInternalDeclaractionSub.update(new BigInteger(addValue)).send();

            collector.logStepPass("ConstructorInternalDeclaractionTest update successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String afterValue = constructorInternalDeclaractionSub.getCount().send().toString();
            collector.logStepPass("链上函数的执行update后的值为："+afterValue);

            collector.assertEqual(String.valueOf(Integer.valueOf(initValue)+Integer.valueOf(addValue)),afterValue);
        } catch (Exception e) {
            collector.logStepFail("ConstructorInternalDeclaractionTest process fail.", e.toString());
            e.printStackTrace();
        }
    }


}
