package wasm.contract_cross_call;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.*;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

public class ContractCrossCallFlagByRetTest extends WASMContractPrepareTest {

    // 跨合约调用有返回值的合约
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_flag_ByRet",sourcePrefix = "wasm")
    public void testCrossCallContractByRet() {

        try {
            prepare();

            // deploy the target contract which the name is `receiver_byret`, first
            ContractReceiverByRet receiver  = ContractReceiverByRet.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy receiver_byret contract:" + receiver.getTransactionReceipt().get().getGasUsed());

            String receiverAddr = receiver.getContractAddress();
            String receiverTxHash = receiver.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("receiver_byret deployed sucessfully, contractAddress:" + receiverAddr + ", txHash:" + receiverTxHash);


            // deploy the cross_caller_byret  contract second
             ContractCallerByRet caller = ContractCallerByRet.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_caller_byret contract:" + caller.getTransactionReceipt().get().getGasUsed());

            String callerAddr = caller.getContractAddress();
            String callerTxHash = caller.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_caller_byret deployed sucessfully, contractAddress:" + callerAddr + ", txHash:" + callerTxHash);


            // 给予receiver的gas足够时
            TransactionReceipt receipt = caller.callFeed(receiverAddr, Uint64.of(508651l)).send();
            collector.logStepPass("cross_caller_byret callFeed by gas 98651l successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            Uint64 status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 0l);

            // 给予receiver的gas 不够时
            receipt = caller.callFeed(receiverAddr, Uint64.of(10000l)).send();
            collector.logStepPass("cross_caller_byret callFeed by gas 10000 successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 1l);


            // 直接给予 0gas, 就会使用默认的 剩余gas
            receipt = caller.callFeed(receiverAddr, Uint64.of(0l)).send();
            collector.logStepPass("cross_caller_byret callFeed by gas 0 successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 0l);


        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_caller_byret Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }



    // 跨合约调用无返回值的合约
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xujiacan", showName = "wasm.contract_cross_call_flag_ByNoRet",sourcePrefix = "wasm")
    public void testCrossCallContractByNoRet() {

        try {
            prepare();

            // deploy the target contract which the name is `receiver_noret`, first
            ContractReceiverNoRet receiver  = ContractReceiverNoRet.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy receiver_noret contract:" + receiver.getTransactionReceipt().get().getGasUsed());

            String receiverAddr = receiver.getContractAddress();
            String receiverTxHash = receiver.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("receiver_noret deployed sucessfully, contractAddress:" + receiverAddr + ", txHash:" + receiverTxHash);


            // deploy the cross_caller_noret  contract second
            ContractCallerNoRet caller = ContractCallerNoRet.deploy(web3j, transactionManager, provider, chainId).send();
            collector.logStepPass("gas used after deploy cross_caller_noret contract:" + caller.getTransactionReceipt().get().getGasUsed());

            String callerAddr = caller.getContractAddress();
            String callerTxHash = caller.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("cross_caller_noret deployed sucessfully, contractAddress:" + callerAddr + ", txHash:" + callerTxHash);


            // 给予receiver的gas足够时
            TransactionReceipt receipt = caller.callFeed(receiverAddr, Uint64.of(98651l)).send();
            collector.logStepPass("cross_caller_noret callFeed by gas 98651l successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            Uint64 status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 0l);

            // 给予receiver的gas 不够时
            receipt = caller.callFeed(receiverAddr, Uint64.of(10000l)).send();
            collector.logStepPass("cross_caller_noret callFeed by gas 10000 successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 1l);


            // 直接给予 0gas, 就会使用默认的 剩余gas
            receipt = caller.callFeed(receiverAddr, Uint64.of(0l)).send();
            collector.logStepPass("cross_caller_noret callFeed by gas 0 successfully txHash:" + receipt.getTransactionHash());

            // 检查下调用receiver 成功与否
            status = caller.get_status().send();
            collector.assertEqual(status.getValue().longValue(), 0l);

        } catch (Exception e) {
            collector.logStepFail("Failed to call cross_caller_noret Contract,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
