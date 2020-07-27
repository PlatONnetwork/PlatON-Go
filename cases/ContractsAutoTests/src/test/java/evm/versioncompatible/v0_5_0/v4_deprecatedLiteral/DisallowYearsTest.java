package evm.versioncompatible.v0_5_0.v4_deprecatedLiteral;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DisallowYears;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title  07- 存储区域
 * 1-结构体(struct)，数组(array)，
 * 映射(mapping)类型的变量必须显式声明存储区域( storage， memeory， calldata)，
 * 包括函数参数和返回值变量都必须显式声明
 * 2-external 的函数的引用参数和映射参数需显式声明为 calldata
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class DisallowYearsTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "DisallowYearsTest-弃用字面量及后缀", sourcePrefix = "evm")
    public void update() {
        try {

            DisallowYears disallowYears = DisallowYears.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = disallowYears.getContractAddress();
            TransactionReceipt tx = disallowYears.getTransactionReceipt().get();

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + disallowYears.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = disallowYears.testyear(new BigInteger("10"),new BigInteger("1")).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String time1 = disallowYears.getTime1().send().toString();
            //与以太坊的测试网数据进行比较
            collector.assertEqual("31536000",time1);

            String etherValue = disallowYears.getEtherValue().send().toString();
            //与以太坊的测试网数据进行比较
            collector.assertEqual("255000000000000000000",etherValue);

            String hexValue = disallowYears.getHexValue().send().toString();
            //与以太坊的测试网数据进行比较
            collector.assertEqual("255",hexValue);

            String hexComValue = disallowYears.getHexComValue().send().toString();
            //与以太坊的测试网数据进行比较
            collector.assertEqual("255000000000000000000",hexComValue);

        } catch (Exception e) {
            collector.logStepFail("DisallowYearsTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
