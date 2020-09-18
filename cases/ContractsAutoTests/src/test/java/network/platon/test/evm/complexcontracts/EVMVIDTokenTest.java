package network.platon.test.evm.complexcontracts;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.VIDToken;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;

public class EVMVIDTokenTest extends ContractPrepareTest {
    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "complexcontracts.EVMVIDTokenTest", sourcePrefix = "evm")
    public void test() {
        try {
            VIDToken token = VIDToken.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = token.getContractAddress();
            TransactionReceipt tx = token.getTransactionReceipt().get();

            collector.logStepPass(
                    "Token issued successfully.contractAddress:" + contractAddress +
                            ", hash:" + tx.getTransactionHash() +
                            ", tokenName:" + token.name().send() +
                            ", symbol:" + token.symbol().send());
            collector.logStepPass("deploy gas used:" + token.getTransactionReceipt().get().getGasUsed());

            BigInteger balanceOf = VIDToken.load(contractAddress, web3j, transactionManager, provider, chainId).balanceOf(contractAddress).send();
            collector.logStepPass("balanceOf:" + balanceOf);

        } catch (Exception e) {
            e.printStackTrace();
        }

    }
}
