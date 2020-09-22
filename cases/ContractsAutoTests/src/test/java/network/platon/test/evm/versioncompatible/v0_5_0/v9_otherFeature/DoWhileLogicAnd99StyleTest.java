package network.platon.test.evm.versioncompatible.v0_5_0.v9_otherFeature;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.DoWhileLogicAnd99Style;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title do...while循环里的continue不再跳转到循环体内, 而是跳转到while处判断循环条件,
 * 若条件为假,就退出循环。这一修改更符合一般编程语言的设计风格
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class DoWhileLogicAnd99StyleTest extends ContractPrepareTest {

    //do...while循环起始值
    private String doWhileStart;

    //for循环起始值
    private String forStart;

    @Before
    public void before() {
        this.prepare();
        doWhileStart = driverService.param.get("doWhileStart");
        forStart = driverService.param.get("forStart");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "DoWhileLogicAnd99StyleTest-c99语法风格", sourcePrefix = "evm")
    public void doWhileAndFor() {
        try {
            DoWhileLogicAnd99Style doWhileLogicAnd99Style = DoWhileLogicAnd99Style.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = doWhileLogicAnd99Style.getContractAddress();
            TransactionReceipt tx = doWhileLogicAnd99Style.getTransactionReceipt().get();

            collector.logStepPass("DoWhileLogicAnd99StyleTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + doWhileLogicAnd99Style.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = doWhileLogicAnd99Style.dowhile(new BigInteger(doWhileStart)).send();

            collector.logStepPass("DoWhileLogicAnd99StyleTest dowhile successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            String chainDoWhileValue = doWhileLogicAnd99Style.getDoWhileSum().send().toString();

            collector.assertEqual(dowhile(Integer.valueOf(doWhileStart)), chainDoWhileValue);

            //测试for循环
            TransactionReceipt transactionReceiptFor = doWhileLogicAnd99Style.forsum(new BigInteger(forStart)).send();

            collector.logStepPass("DoWhileLogicAnd99StyleTest dowhile successful.transactionHash:" + transactionReceiptFor.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceiptFor.getBlockNumber());

            String chainForValue = doWhileLogicAnd99Style.getForSum().send().toString();

            collector.assertEqual(forsum(Integer.valueOf(forStart)), chainForValue);

        } catch (Exception e) {
            collector.logStepFail("DoWhileLogicAnd99StyleTest doWhileAndFor process fail.", e.toString());
            e.printStackTrace();
        }
    }


    /**
     * @title do...while结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String dowhile(int x) {
        int y = x + 10;
        int z = x + 9;
        do {
            x += 1;
            if (x > z) continue;
        } while (x < y);
        return String.valueOf(x);
    }

    /**
     * @title for 循环后的结果值
     * @description:
     * @author: hudenian
     * @create: 2019/12/27
     */
    public static String forsum(int x) {
        int forSum = 0;
        for (int i = 0; i < x; i++) {
            forSum = forSum + i;
        }
        return String.valueOf(forSum);
    }

}
