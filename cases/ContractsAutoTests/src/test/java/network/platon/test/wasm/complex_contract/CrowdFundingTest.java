package network.platon.test.wasm.complex_contract;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CrowdFunding;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @author hudenian
 *
 */
public class CrowdFundingTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.GuessingTest-众筹合约",sourcePrefix = "wasm")
    public void testCrowdFundContract() {

        Long blocks = 30L;//设置截止块高与当前块高为20

        try {
            CrowdFunding crowdFunding = CrowdFunding.deploy(web3j, transactionManager, provider,chainId,Uint64.of("10000"),Uint64.of("10000")).send();
            String contractAddress = crowdFunding.getContractAddress();
            String transactionHash = crowdFunding.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CrowdFunding issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CrowdFunding deploy successfully. gasUsed: " + crowdFunding.getTransactionReceipt().get().getGasUsed().toString());

            //发起众筹
            TransactionReceipt transactionReceipt = crowdFunding.crowdFund(new BigInteger("1000")).send();
            collector.logStepPass("CrowdFunding call transfer hash is:"+transactionReceipt.getTransactionHash());

            //检测众筹目标是否已经达到
            transactionReceipt =crowdFunding.safeWithdrawal().send();
            collector.logStepPass("CrowdFunding call checkGoalReached hash is:"+transactionReceipt.getTransactionHash());

            CrowdFunding.Transfer1EventResponse flg = crowdFunding.getTransfer1Events(transactionReceipt).get(0);
            collector.assertEqual("1000",crowdFunding.getTransfer1Events(transactionReceipt).get(0).arg2.toString());



        } catch (Exception e) {
            collector.logStepFail("CrowdFunding failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
