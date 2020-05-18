package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Fibonacci;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class exists for docs.
 */
public class ContractFibonacciTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_fibonacci",sourcePrefix = "wasm")
    public void testFibonacciContract() {

        try {
            // deploy contract.
            Fibonacci contract = Fibonacci.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("Fibonacci deploy successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("Fibonacci deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());
            Long number = Long.valueOf(5);
            TransactionReceipt tr = contract.fibonacci_notify(Uint64.of(number)).send();
            collector.logStepPass("Fibonacci notify successfully, hash: " + tr.getTransactionHash());
            Fibonacci.NotifyEventResponse eventResponse = contract.getNotifyEvents(tr).get(0);

            collector.logStepPass("To invoke fibonacci_notify success, args0: " + eventResponse.arg1
                    + " args2: " + eventResponse.arg2
                    + " args3: " + eventResponse.arg3);

            Uint64 result = contract.fibonacci_call(Uint64.of(number)).send();
            collector.logStepPass("To invoke fibonacci success, result: " + result.toString());

        } catch (Exception e) {
            collector.logStepFail("Fibonacci failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
