package crossContractCall;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.DelegatecallCallee;
import network.platon.contracts.DelegatecallCallee_050;
import network.platon.contracts.DelegatecallCaller;
import network.platon.contracts.DelegatecallCaller_050;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title 0.5.0跨合约调用者,修改的是调用者中的状态变量的值
 * @description:
 * @author: hudenian
 * @create: 2019/12/30
 */
public class DelegatecallCaller050Test extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "DelegatecallCaller050Test-0.5.0跨合约调用者")
    public void caller050Test() {
        try {
            //调用者合约地址
            DelegatecallCaller_050 delegatecallCaller_050 = DelegatecallCaller_050.deploy(web3j, transactionManager, provider).send();
            String callerContractAddress = delegatecallCaller_050.getContractAddress();
            TransactionReceipt tx = delegatecallCaller_050.getTransactionReceipt().get();
            collector.logStepPass("DelegatecallCaller_050 deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());


            //被调用者合约地址
            DelegatecallCallee_050 delegatecallCallee050 = DelegatecallCallee_050.deploy(web3j, transactionManager, provider).send();
            String calleeContractAddress = delegatecallCallee050.getContractAddress();
            TransactionReceipt tx1 = delegatecallCallee050.getTransactionReceipt().get();
            collector.logStepPass("delegatecallCallee050 deploy successfully.contractAddress:" + calleeContractAddress + ", hash:" + tx1.getTransactionHash());

            //查询调用者x值
            String callerX = delegatecallCaller_050.getCallerX().send().toString();
            collector.logStepPass("DelegatecallCaller_050 合约中X的值为："+callerX);

            //查询被调用者x值
            String calleeX = delegatecallCallee050.getCalleeX().send().toString();
            collector.logStepPass("DelegatecallCallee_050 合约中X的值为："+calleeX);


            TransactionReceipt tx2 = delegatecallCaller_050.inc_delegatecall(calleeContractAddress).send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询调用者x值
            String callerAfterX = delegatecallCaller_050.getCallerX().send().toString();
            collector.logStepPass("跨合约调用后，DelegatecallCaller 合约中X的值为："+callerAfterX);

            //查询被调用者x值
            String calleeAfterX = delegatecallCallee050.getCalleeX().send().toString();
            collector.logStepPass("跨合约调用后，DelegatecallCallee_050 合约中X的值为："+calleeAfterX);


        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
