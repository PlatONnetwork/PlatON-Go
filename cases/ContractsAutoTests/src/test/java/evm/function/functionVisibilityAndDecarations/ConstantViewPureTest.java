package evm.function.functionVisibilityAndDecarations;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.ConstantViewPure;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 验证constant,view,pure
 * @description:
 * @author: liweic
 * @create: 2020/01/02 14:01
 **/

public class ConstantViewPureTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.ConstantViewPureTest-函数声明方式测试", sourcePrefix = "evm")
    public void constantviewPure() {
        try {
            ConstantViewPure constantviewpure = ConstantViewPure.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = constantviewpure.getContractAddress();
            TransactionReceipt tx = constantviewpure.getTransactionReceipt().get();
            collector.logStepPass("ConstantViewPure deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("ConstantViewPure deploy gasUsed:" + constantviewpure.getTransactionReceipt().get().getGasUsed());


            TransactionReceipt age = constantviewpure.constantViewPure().send();
            collector.logStepPass("age交易Hash" + age.getTransactionHash());
            //验证constant声明
            BigInteger constantage = constantviewpure.getAgeByConstant().send();
            collector.logStepPass("constant声明函数后返回值：" + constantage);
            collector.assertEqual(new BigInteger("20"),constantage);

            //验证view声明
            BigInteger viewage = constantviewpure.getAgeByView().send();
            collector.logStepPass("view声明函数后返回值：" + viewage);
            collector.assertEqual(new BigInteger("20"),viewage);

            //验证pure声明
            BigInteger pureage = constantviewpure.getAgeByPure().send();
            collector.logStepPass("pure声明函数后返回值：" + pureage);
            collector.assertEqual(new BigInteger("1"),pureage);


        } catch (Exception e) {
            collector.logStepFail("ConstantViewPureContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}

