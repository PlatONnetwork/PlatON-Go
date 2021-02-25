package network.platon.test.evm.crossContractCall;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.DelegatecallCallee;
import network.platon.contracts.evm.DelegatecallCaller;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title 0.5.13跨合约调用者,修改的是调用者中的状态变量的值
 * @description:
 * @author: hudenian
 * @create: 2019/12/28
 */
public class DelegatecallCallerTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "DelegatecallCallerTest-跨合约调用者", sourcePrefix = "evm")
    public void crossContractCaller() {
        try {
            //调用者合约地址
            DelegatecallCaller delegatecallCaller = DelegatecallCaller.deploy(web3j, transactionManager, provider, chainId).send();
            String callerContractAddress = delegatecallCaller.getContractAddress();
            TransactionReceipt tx = delegatecallCaller.getTransactionReceipt().get();
            collector.logStepPass("DelegatecallCaller deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + delegatecallCaller.getTransactionReceipt().get().getGasUsed());


            //被调用者合约地址
            DelegatecallCallee delegatecallCallee = DelegatecallCallee.deploy(web3j, transactionManager, provider, chainId).send();
            String calleeContractAddress = delegatecallCallee.getContractAddress();
            TransactionReceipt tx1 = delegatecallCallee.getTransactionReceipt().get();
            collector.logStepPass("DelegatecallCallee deploy successfully.contractAddress:" + calleeContractAddress + ", hash:" + tx1.getTransactionHash());
            collector.logStepPass("deploy gas used:" + delegatecallCallee.getTransactionReceipt().get().getGasUsed());

            //查询调用者x值
            String callerX = delegatecallCaller.getCallerX().send().toString();
            collector.logStepPass("DelegatecallCaller 合约中X的值为："+callerX);

            //查询被调用者x值
            String calleeX = delegatecallCallee.getCalleeX().send().toString();
            collector.logStepPass("DelegatecallCallee 合约中X的值为："+calleeX);


            TransactionReceipt tx2 = delegatecallCaller.inc_delegatecall(calleeContractAddress).send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询调用者x值
            String callerAfterX = delegatecallCaller.getCallerX().send().toString();
            collector.logStepPass("跨合约调用后，DelegatecallCaller 合约中X的值为："+callerAfterX);

            //查询被调用者x值
            String calleeAfterX = delegatecallCallee.getCalleeX().send().toString();
            collector.logStepPass("跨合约调用后，DelegatecallCallee 合约中X的值为："+calleeAfterX);


        } catch (Exception e) {
            collector.logStepFail("DelegatecallCallerTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
