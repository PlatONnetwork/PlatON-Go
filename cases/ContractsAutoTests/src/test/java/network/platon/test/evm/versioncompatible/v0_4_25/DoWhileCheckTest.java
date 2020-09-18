package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.DoWhileCheck;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;

import java.math.BigDecimal;
import java.math.BigInteger;
/**
 * @title 0.4.25版本重大bug测试
 * 1. 0.4.25版本do...while循环里的continue跳转到循环体内，可能会产生死循环验证
 * 2. 0.4.25版本局部变量上级作用域生效验证
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class DoWhileCheckTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "network.platon.test.evm.DoWhileCheckTest-重大bug验证", sourcePrefix = "evm")
    public void testDoWhileCheck() {
        try {
            prepare();
            DoWhileCheck doWhileCheck = DoWhileCheck.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = doWhileCheck.getContractAddress();
            String transactionHash = doWhileCheck.getTransactionReceipt().get().getTransactionHash();
            Transfer.sendFunds(web3j,credentials, chainId,contractAddress, BigDecimal.valueOf(300.00), Convert.Unit.GLAT);
            collector.logStepPass("DoWhileCheck issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + doWhileCheck.getTransactionReceipt().get().getGasUsed());
            Tuple2<BigInteger, BigInteger> result = doWhileCheck.doWhileCheck().send();
            Tuple2<BigInteger, BigInteger> expect =new Tuple2(new BigInteger("21"),new BigInteger("14")) ;
            collector.assertEqual(JSONObject.toJSONString(result), JSONObject.toJSONString(expect), "checkout continue bug and scope is not clear");
        } catch (Exception e) {
            collector.logStepFail("DoWhileCheckTest testDoWhileCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
