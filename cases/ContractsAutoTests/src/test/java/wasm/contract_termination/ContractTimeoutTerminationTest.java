package wasm.contract_termination;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Contract_timeout_termination;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约死循环终止
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class ContractTimeoutTerminationTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_termination合约timeout终止",sourcePrefix = "wasm")
    public void testTimeOutContract() {

        String nomalName = "nomalName";
        String errorName = "errorName";
        Uint64 nomalValue = Uint64.of(12L);
        Uint64 errorValue = Uint64.of(112L);
        try {
            prepare();
            Contract_timeout_termination contractTimeoutTermination = Contract_timeout_termination.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractTimeoutTermination.getContractAddress();
            String transactionHash = contractTimeoutTermination.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractTermination issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractTimeoutTermination.getTransactionReceipt().get().getGasUsed());

            //传正常的值
            TransactionReceipt transactionReceipt = contractTimeoutTermination.forfunction(nomalName,nomalValue).send();
            collector.logStepPass("ContractTermination call forfunction set nomal value successfully hash:" + transactionReceipt.getTransactionHash());

            //查询结果
            String chainName = contractTimeoutTermination.get_string_storage().send();
            collector.logStepPass("ContractTimeoutTerminationTest nomal get_string_storage successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName,nomalName);

            //传异常值
            transactionReceipt = contractTimeoutTermination.forfunction(errorName,errorValue).send();
            collector.logStepPass("ContractTermination call forfunction set error value successfully hash:" + transactionReceipt.getTransactionHash());


        } catch (Exception e) {
            if(e instanceof TransactionException){
                collector.logStepPass("死循环产生gas消耗完:"+e.getMessage());
            }else{
                collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }

        }
    }
}
