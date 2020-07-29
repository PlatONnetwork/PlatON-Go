package wasm.contract_distory;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDistory;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title wasm同一个合约之间局部状态是否会相互影响验证
 * @description:
 * @author: hudenian
 * @create: 2020/03/05
 */
public class ContractCheckStatusTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.ContractCheckStatusTest同一个合约之间局部状态是否会相互影响验证",sourcePrefix = "wasm")
    public void testDistoryContract() {

        String nameOne = "valueOne";
        String nameTwo = "valueTwo";
        String nameThree = "nameThree";
        try {
            prepare();
            //合约第一次部署
            ContractDistory contractDistory = ContractDistory.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractDistory.getContractAddress();
            String transactionHash = contractDistory.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistory first issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("ContractDistory first deploy gas used:" + contractDistory.getTransactionReceipt().get().getGasUsed());

            //合约第二次部署
            ContractDistory contractDistoryTwo = ContractDistory.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddressTwo = contractDistoryTwo.getContractAddress();
            transactionHash = contractDistoryTwo.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistory second issued successfully.contractAddress:" + contractAddressTwo + ", hash:" + transactionHash);
            collector.logStepPass("ContractDistory second deploy gas used:" + contractDistoryTwo.getTransactionReceipt().get().getGasUsed());


            //contractDistory合约设置值
            TransactionReceipt transactionReceipt = contractDistory.set_string(nameOne).send();
            collector.logStepPass("contractDistory set_string successfully hash:" + transactionReceipt.getTransactionHash());
            //contractDistoryTwo合约设置值
            transactionReceipt = contractDistoryTwo.set_string(nameTwo).send();
            collector.logStepPass("contractDistoryTwo set_string successfully hash:" + transactionReceipt.getTransactionHash());


            String chainNameOne = contractDistory.get_string().send();
            collector.logStepPass("contractDistory get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainNameOne,nameOne);

            String chainNameTwo = contractDistoryTwo.get_string().send();
            collector.logStepPass("contractDistoryTwo get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainNameTwo,nameTwo);

            //第一个合约销毁看对第二个合约是否有影响
            transactionReceipt = contractDistory.distory_contract().send();
            collector.logStepPass("ContractDistory distory_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //验证第一个合约销毁对第二个合约没有影响
            transactionReceipt = contractDistoryTwo.set_string(nameThree).send();
            collector.logStepPass("contractDistoryTwo set_string successfully hash:" + transactionReceipt.getTransactionHash());


            String chainNameThree = contractDistoryTwo.get_string().send();
            collector.logStepPass("contractDistoryTwo get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainNameThree,nameThree);

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
