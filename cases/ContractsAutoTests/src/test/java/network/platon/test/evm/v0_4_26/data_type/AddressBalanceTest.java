package network.platon.test.evm.v0_4_26.data_type;

import com.alaya.protocol.core.methods.response.TransactionReceipt;
import network.platon.contracts.evm.v0_4_26.AddressBalance;
import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import org.junit.Before;
import org.junit.Test;
import java.math.BigInteger;

public class AddressBalanceTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1", author = "qcxiao",
            showName = "AddressBalanceTest.查询某地址余额",sourcePrefix = "evm/0.4.26")
    public void test() {

//        String useAddress = "lax10eycqggu2yawpadtmn7d2zdw0vnmscklynzq8x";
        String useAddress = "atx10eycqggu2yawpadtmn7d2zdw0vnmscklcx6a9v";

        AddressBalance addressBalance = null;
        try {
            //合约部署
            addressBalance = AddressBalance.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = addressBalance.getContractAddress();
            TransactionReceipt tx =  addressBalance.getTransactionReceipt().get();
            collector.logStepPass("AddressBalance deploy successfully.contractAddress:" + contractAddress
                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deploy finish currentBlockNumber:" + tx.getBlockNumber());

            //调用合约
            BigInteger balance = addressBalance.balanceOfPlatON(useAddress).send();
            collector.logStepPass("transactionHash:" + tx.getTransactionHash() + ", currentBlockNumber:" + tx.getBlockNumber());

            System.out.println("address:" + useAddress + ", balance:" + balance);
        } catch (Exception e) {
            collector.logStepFail("AddressBalance process fail.", e.toString());
            e.printStackTrace();
        }
    }

}