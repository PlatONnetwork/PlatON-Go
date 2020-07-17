package wasm.contract_termination;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Contract_termination;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约终止，platon_assert 断言失败会退出合约，此时会花费掉实际执行消耗的 gas
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class ContractAssertTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_termination合约assert终止",sourcePrefix = "wasm")
    public void testAssertContract() {

        String nomalName = "nomalName";
        String errorName = "errorName";
        Uint64 nomalValue = Uint64.of(112L);
        Uint64 errorValue = Uint64.of(12L);
        try {
            prepare();
            Contract_termination contractTermination = Contract_termination.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractTermination.getContractAddress();
            String transactionHash = contractTermination.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractTermination issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractTermination.getTransactionReceipt().get().getGasUsed());

            //调用包含platon_assert的合约,传正常的值
            TransactionReceipt transactionReceipt = contractTermination.transfer_assert(nomalName,nomalValue).send();
            collector.logStepPass("ContractTermination transfer_assert successfully hash:" + transactionReceipt.getTransactionHash());

            //查询调platon_panic之前设置的值
            String chainName = contractTermination.get_string_storage().send();
            collector.logStepPass("ContractTermination assert get_string_storage successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName,nomalName);

            //调用包含platon_assert的合约,传触发assert方法
            transactionReceipt = contractTermination.transfer_assert(errorName,errorValue).send();
            collector.logStepPass("ContractTermination transfer_assert  successfully hash:" + transactionReceipt.getTransactionHash());


        } catch (Exception e) {
            if(e instanceof TransactionException){
                collector.logStepPass("platon_assert会消耗实际使用的gas:"+e.getMessage());
            }else{
                collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }

        }
    }
}
