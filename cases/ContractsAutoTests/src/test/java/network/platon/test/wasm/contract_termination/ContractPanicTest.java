package network.platon.test.wasm.contract_termination;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Contract_panic;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约终止，把用户的全部 gas 用完
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class ContractPanicTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_termination合约panic终止",sourcePrefix = "wasm")
    public void testPanicContract() {

        String name = "hudenian";
        Uint64 value = Uint64.of(3L);
        try {
            prepare();
            Contract_panic contractPanic = Contract_panic.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractPanic.getContractAddress();
            String transactionHash = contractPanic.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractPanic issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractPanic.getTransactionReceipt().get().getGasUsed());

            //调用包含platon_panic的合约
            TransactionReceipt transactionReceipt = contractPanic.panic_contract(name,value).send();
            collector.logStepPass("ContractDistory panic_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //查询调platon_panic之前设置的值
            String chainName = contractPanic.get_string_storage().send();
            collector.logStepPass("ContractDistory get_string_storage successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName,name);

            //查询调platon_panic之后设置的值
            String chainName1 = contractPanic.get_string_storage1().send();
            collector.logStepPass("ContractDistory get_string_storage1 successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName1,name);

        } catch (Exception e) {
            if(e instanceof TransactionException){
                collector.logStepPass("platon_panic会消耗所有的gas:"+e.getMessage());
            }else{
                collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }

        }
    }
}
