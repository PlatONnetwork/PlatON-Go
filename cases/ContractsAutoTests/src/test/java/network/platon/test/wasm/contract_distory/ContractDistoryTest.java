package network.platon.test.wasm.contract_distory;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDistory;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约销毁
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class ContractDistoryTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_distory合约销毁",sourcePrefix = "wasm")
    public void testDistoryContract() {

        String name = "hudenian";
        try {
            prepare();
            ContractDistory contractDistory = ContractDistory.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractDistory.getContractAddress();
            String transactionHash = contractDistory.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistory issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractDistory.getTransactionReceipt().get().getGasUsed());

            //合约设置值
            TransactionReceipt transactionReceipt = contractDistory.set_string(name).send();
            collector.logStepPass("ContractDistory set_string successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁前查询合约上的数据
            String chainName = contractDistory.get_string().send();
            collector.logStepPass("ContractDistory get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName,name);

            //合约销毁
            transactionReceipt = contractDistory.distory_contract().send();
            collector.logStepPass("ContractDistory distory_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁后查询合约上的数据
            String chainName1 = contractDistory.get_string().send();
            collector.logStepPass("ContractDistory get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName1,name);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ContractDistoryed and could not call contract function");
            }else{
                collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
