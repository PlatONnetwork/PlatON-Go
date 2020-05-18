package wasm.contract_distory;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDistory;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title 合约销毁后，对合约账户余额进行校验
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class ContractDistoryCheckBalanceTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_distory合约销毁余额校验",sourcePrefix = "wasm")
    public void testDistoryContract() {

        try {
            prepare();
            ContractDistory contractDistory = ContractDistory.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractDistory.getContractAddress();
            String transactionHash = contractDistory.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistory issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractDistory.getTransactionReceipt().get().getGasUsed());

            //合约销毁
            TransactionReceipt transactionReceipt = contractDistory.distory_contract().send();
            collector.logStepPass("ContractDistory distory_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁后余额为0
            BigInteger afterDistoryBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("After distory, contract balance is: " + afterDistoryBalance);
            collector.assertEqual(afterDistoryBalance.toString(),"0");

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ContractDistoryCheckBalanceTest and could not call contract function");
            }else{
                collector.logStepFail("ContractDistoryCheckBalanceTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
