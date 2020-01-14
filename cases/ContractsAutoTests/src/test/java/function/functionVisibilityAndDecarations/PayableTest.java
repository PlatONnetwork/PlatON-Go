package function.functionVisibilityAndDecarations;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Payable;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title 验证Payable声明
 * @description:
 * @author: liweic
 * @create: 2020/01/02 16:01
 **/

public class PayableTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.PayableTest-函数声明方式Payable测试")
    public void payable() {
        try {
            Payable payable = Payable.deploy(web3j, transactionManager, provider).send();

            String contractAddress = payable.getContractAddress();
            TransactionReceipt tx = payable.getTransactionReceipt().get();
            collector.logStepPass("paybale deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证payable声明
            BigInteger payablepremoney = payable.getBalances("0x8a9B36694F1eeeb500c84A19bB34137B05162EC4").send();
            collector.logStepPass("转账前余额：" + payablepremoney);
            TransactionReceipt result = payable.transfer("0x8a9B36694F1eeeb500c84A19bB34137B05162EC4",new BigInteger("100")).send();
            BigInteger payableaftermoney = payable.getBalances("0x8a9B36694F1eeeb500c84A19bB34137B05162EC4").send();
            collector.logStepPass("转账后余额：" + payableaftermoney);
            int a = Integer.valueOf(payableaftermoney.toString());
            int b = Integer.valueOf(payablepremoney.toString());
            int payablecount = a - b;
            collector.assertEqual(100,payablecount);


        } catch (Exception e) {
            collector.logStepFail("PayableContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}