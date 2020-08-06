package evm.createcontract;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.CreateContract;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple2;

import java.math.BigInteger;

/**
 * @title new关键字创建合约测试
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class CreateContractTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "createcontract.CreateContractTest-new创建合约", sourcePrefix = "evm")
    public void testNewContract() {
        try {
            prepare();
            CreateContract createContract = CreateContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = createContract.getContractAddress();
            String transactionHash = createContract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CreateContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + createContract.getTransactionReceipt().get().getGasUsed());
            Tuple2<BigInteger, BigInteger> result = createContract.getTargetCreateContractData().send();
            Tuple2<BigInteger, BigInteger> expect = new Tuple2<>(new BigInteger("1000"),new BigInteger("0"));
            collector.assertEqual(result, expect, "checkout new contract param");
        } catch (Exception e) {
            collector.logStepFail("CreateContractTest testNewContract failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
