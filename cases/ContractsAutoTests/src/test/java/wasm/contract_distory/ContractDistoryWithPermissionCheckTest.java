package wasm.contract_distory;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDistory;
import network.platon.contracts.wasm.ContractDistoryWithPermissionCheck;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约销毁,带用户权限校验
 * @description:
 * @author: hudenian
 * @create: 2020/02/19
 */
public class ContractDistoryWithPermissionCheckTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_distory合约销毁带用户权限校验",sourcePrefix = "wasm")
    public void testDistoryContractCheckPermission() {

        String name = "hudenian";
        try {
            prepare();
            ContractDistoryWithPermissionCheck ontractDistoryWithPermissionCheck = ContractDistoryWithPermissionCheck.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = ontractDistoryWithPermissionCheck.getContractAddress();
            String transactionHash = ontractDistoryWithPermissionCheck.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistoryWithPermissionCheck issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + ontractDistoryWithPermissionCheck.getTransactionReceipt().get().getGasUsed());

            //合约设置值
            TransactionReceipt transactionReceipt = ontractDistoryWithPermissionCheck.set_string(name).send();
            collector.logStepPass("ContractDistoryWithPermissionCheck set_string successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁前查询合约上的数据
            String chainName = ontractDistoryWithPermissionCheck.get_string().send();
            collector.logStepPass("ContractDistoryWithPermissionCheck get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName,name);

            //合约销毁
            transactionReceipt = ontractDistoryWithPermissionCheck.distory_contract().send();
            collector.logStepPass("ContractDistoryWithPermissionCheck distory_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁后查询合约上的数据
            String chainName1 = ontractDistoryWithPermissionCheck.get_string().send();
            collector.logStepPass("ContractDistoryWithPermissionCheck get_string successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainName1,name);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ContractDistoryWithPermissionCheckTest and could not call contract function");
            }else{
                collector.logStepFail("ContractDistoryWithPermissionCheckTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
