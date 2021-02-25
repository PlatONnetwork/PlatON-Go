package network.platon.test.evm.variable;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.Time;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 验证时间单位
 * @description:
 * @author: liweic
 * @create: 2020/01/02 19:50
 **/

public class TimeTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "network.platon.test.evm.function.TimeTest-时间单位测试", sourcePrefix = "evm")
    public void time() {
        try {
            Time time = Time.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = time.getContractAddress();
            TransactionReceipt tx = time.getTransactionReceipt().get();
            collector.logStepPass("time deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证时间单位及换算
            //时间戳获取
            BigInteger result = time.testimeDiff().send();
            collector.logStepPass("testimeDiff返回值：" + result);
            collector.assertEqual(new BigInteger("0"),result);

            //秒
            BigInteger second = time.tSeconds().send();
            collector.logStepPass("tSeconds返回值：" + second);
            collector.assertEqual(new BigInteger("1"),second);

            //分
            BigInteger minute = time.tMinutes().send();
            collector.logStepPass("tMinutes返回值：" + minute);
            collector.assertEqual(new BigInteger("60"),minute);

            //时
            BigInteger hour = time.tHours().send();
            collector.logStepPass("tHours返回值：" + hour);
            collector.assertEqual(new BigInteger("3600"),hour);

            //周
            BigInteger week = time.tWeeks().send();
            collector.logStepPass("tWeeks返回值：" + week);
            collector.assertEqual(new BigInteger("604800"),week);


        } catch (Exception e) {
            collector.logStepFail("TimeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


