package network.platon.test.evm.function.paramandreturns;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.PramaAndReturns;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;

/**
 * @title 参数和返回类型
 * @description:
 * @author: liweic
 * @create: 2019/12/30 15:01
 **/


public class PramaAndReturnsTest extends ContractPrepareTest {

    private String a;
    private String b;
    private String c;

    @Before
    public void before() {
        this.prepare();
        a = driverService.param.get("a");
        b = driverService.param.get("b");
        c = driverService.param.get("c");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.PramaAndReturnsTest-参数和返回类型测试", sourcePrefix = "evm")
    public void ParamAndReturn() {
        try {

            PramaAndReturns pramaAndReturns = PramaAndReturns.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = pramaAndReturns.getContractAddress();
            TransactionReceipt tx = pramaAndReturns.getTransactionReceipt().get();
            collector.logStepPass("pramaAndReturnstest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("pramaAndReturnstest deploy gasUsed:" + pramaAndReturns.getTransactionReceipt().get().getGasUsed());

            BigInteger resultA = pramaAndReturns.InputParam(new BigInteger(a)).send();
            collector.logStepPass("InputParam函数返回值：" + resultA);
            collector.assertEqual("100",resultA.toString());

            //调用没有返回值的函数
            TransactionReceipt transactionReceipt = pramaAndReturns.NoOutput(new BigInteger("1"),new BigInteger("1")).send();
            collector.logStepPass("NoOutput NoOutput successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            ////验证函数NoOutput是否调用成功
            BigInteger resultB = pramaAndReturns.getS().send();
            collector.logStepPass("验证的S值为: " + resultB);
            collector.assertEqual(new BigInteger("1"), resultB);

            BigInteger resultC = pramaAndReturns.OmitParam(new BigInteger("10"), new BigInteger("20")).send();
            collector.logStepPass("OmitParam函数返回值为: " + resultC);
            collector.assertEqual(new BigInteger("10"), resultC);

            //入参和出参均为数组
            List<BigInteger> bigIntegerList = new ArrayList<BigInteger>();
            bigIntegerList.add(new BigInteger(a));
            bigIntegerList.add(new BigInteger(b));
            bigIntegerList.add(new BigInteger(c));

            List resList =  pramaAndReturns.IuputArray(bigIntegerList).send();
            collector.logStepPass("IuputArray第一个数: " + resList.get(0));
            collector.assertEqual("100", resList.get(0).toString());
            collector.logStepPass("IuputArray第二个数: " + resList.get(1));
            collector.assertEqual("1", resList.get(1).toString());
            collector.logStepPass("IuputArray第三个数: " + resList.get(2));
            collector.assertEqual("3", resList.get(2).toString());

            //返回值类型是字符串
            String resultD = pramaAndReturns.OuputString().send();
            collector.logStepPass("返回的字符串是: " + resultD);
            collector.assertEqual(resultD, "What's up man");

            //多个返回值且返回类型是数组
            Tuple2 resultE = pramaAndReturns.OuputArrays().send();
            collector.logStepPass("第一个数组是: " + resultE.getValue1());
            collector.assertEqual("[10, 2, 3]", resultE.getValue1().toString());
            collector.logStepPass("第二个数组是: " + resultE.getValue2());
            collector.assertEqual("[10, 2, 3]", resultE.getValue2().toString());

        } catch (Exception e) {
            collector.logStepFail("ParamAndReturnsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

}


