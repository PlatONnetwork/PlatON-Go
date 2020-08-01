package wasm.contract_termination;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Contract_panic;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.tx.gas.ContractGasProvider;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title 合约gas不足终止
 * @description:
 * @author: hudenian
 * @create: 2020/02/19
 */
public class ContractPanicWithSmalllGasTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_termination合约gas不足终止",sourcePrefix = "wasm")
    public void testPanicContractWithSmallGas() {

        String name = "hudenian";
        Long value = 3L;
        try {
            prepare();
            //gas设置过小
//            provider = new ContractGasProvider(BigInteger.valueOf(50L), BigInteger.valueOf(90000000L));
            Contract_panic contractPanic = Contract_panic.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractPanic.getContractAddress();
            String transactionHash = contractPanic.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractPanic issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

        } catch (Exception e) {
            if(e instanceof RuntimeException && e.getMessage().contains("transaction underpriced")){
                collector.logStepPass("gas 不足合约部署终止:"+e.getMessage());
            }else{
                collector.logStepFail("ContractPanicWithSmalllGasTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
