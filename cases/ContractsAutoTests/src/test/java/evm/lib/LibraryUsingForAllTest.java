package evm.lib;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryUsingForAll;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

/**
 * @title 引用using for方式验证
 * 解释：using A for * 的效果是，库 A 中的函数被附加在任意的类型上。
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryUsingForAllTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "lib.LibraryUsingForAllTest-using A for all type", sourcePrefix = "evm")
    public void testReplace() {
        try {
            prepare();
            LibraryUsingForAll using = LibraryUsingForAll.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryUsingForAll issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + using.getTransactionReceipt().get().getGasUsed());
            TransactionReceipt receipt = using.replace(new BigInteger("12"),new BigInteger("14")).send();
            collector.assertEqual(receipt.getStatus(),"0x1" , "checkout using a for * success");
    } catch (Exception e) {
            collector.logStepFail("LibraryUsingForAll testReplace method failure:",e.getMessage());
            e.printStackTrace();
    }
    }
}
