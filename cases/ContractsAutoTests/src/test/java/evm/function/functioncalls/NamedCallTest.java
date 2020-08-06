package evm.function.functioncalls;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.NamedCall;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;


/**
 * @title 验证函数具名调用
 * @description:
 * @author: liweic
 * @create: 2020/01/02 19:11
 **/

public class NamedCallTest extends ContractPrepareTest {

    @Before
    public void before() {

        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.NamedCallTest-函数具名调用测试", sourcePrefix = "evm")
    public void namedcall() {
        try {
            NamedCall namedcall = NamedCall.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = namedcall.getContractAddress();
            TransactionReceipt tx = namedcall.getTransactionReceipt().get();
            collector.logStepPass("namedcall deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("namedcall deploy gasUsed:" + namedcall.getTransactionReceipt().get().getGasUsed());

            //交换传入值的顺序并返回
            Tuple2 result = namedcall.exchange(new BigInteger("1"),new BigInteger("2")).send();
            collector.logStepPass("exchange函数返回值：" + result);

            //任意顺序的通过变量名来指定参数值
            Tuple2 nametesult = namedcall.namecall().send();
            collector.logStepPass("namecall函数返回值：" + nametesult);
            collector.assertEqual("2",result.getValue1().toString());
            collector.assertEqual("1",result.getValue2().toString());


        } catch (Exception e) {
            collector.logStepFail("NamedCallContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}


