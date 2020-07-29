package evm.crossContractCall;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.*;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;


/**
 * @title 0.5.13三个合约间跨合约调用 delegatecall只会修改第一个合约中的状态变量
 * @description:
 * @author: hudenian
 * @create: 2020/1/9
 */
public class ThreeContractCallTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ThreeContractCallTest-三个合约间跨合约调用者", sourcePrefix = "evm")
    public void threeContractCaller() {
        try {
            //第一个合约
            CallerOne callerOne = CallerOne.deploy(web3j, transactionManager, provider, chainId).send();
            String callerContractAddress = callerOne.getContractAddress();
            TransactionReceipt tx = callerOne.getTransactionReceipt().get();
            collector.logStepPass("CallerOne deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + callerOne.getTransactionReceipt().get().getGasUsed());


            //第二个合约
            CallerTwo callerTwo = CallerTwo.deploy(web3j, transactionManager, provider, chainId).send();
            String callerTwoContractAddress = callerTwo.getContractAddress();
            tx = callerTwo.getTransactionReceipt().get();
            collector.logStepPass("CallerTwo deploy successfully.contractAddress:" + callerTwoContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + callerTwo.getTransactionReceipt().get().getGasUsed());

            //第三个合约
            CallerThree callerThree = CallerThree.deploy(web3j, transactionManager, provider, chainId).send();
            String callerThreeContractAddress = callerThree.getContractAddress();
            tx = callerThree.getTransactionReceipt().get();
            collector.logStepPass("DelegatecallCallee deploy successfully.contractAddress:" + callerThreeContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + callerThree.getTransactionReceipt().get().getGasUsed());



            //查询第一个合约x值
            String callerOneX = callerOne.getCallerX().send().toString();
            collector.logStepPass("CallerOne 合约中X的值为："+callerOneX);

            //查询第二个合约x值
            String callerTwoX = callerTwo.getCalleeX().send().toString();
            collector.logStepPass("CallerTwo 合约中X的值为："+callerTwoX);

            //查询第三个合约x值
            String callerThreeX = callerThree.getCalleeThreeX().send().toString();
            collector.logStepPass("CallerThree 合约中X的值为："+callerThreeX);


            TransactionReceipt tx2 = callerOne.inc_delegatecall().send();
            collector.logStepPass("执行跨合约调用后，hash:" + tx2.getTransactionHash());

            //查询第一个合约x值
            callerOneX = callerOne.getCallerX().send().toString();
            collector.logStepPass("CallerOne 合约中X的值为："+callerOneX);
            collector.assertEqual("1",callerOneX);

            //查询第二个合约x值
            callerTwoX = callerTwo.getCalleeX().send().toString();
            collector.logStepPass("CallerTwo 合约中X的值为："+callerTwoX);
            collector.assertEqual("0",callerTwoX);

            //查询第三个合约x值
            callerThreeX = callerThree.getCalleeThreeX().send().toString();
            collector.logStepPass("CallerThree 合约中X的值为："+callerThreeX);
            collector.assertEqual("0",callerThreeX);

        } catch (Exception e) {
            collector.logStepFail("ThreeContractCallTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
