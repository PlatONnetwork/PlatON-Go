package network.platon.test.wasm.contract_object_oriented;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractInterface;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约面向对象接口特性的函数测试
 * @description:
 * @author: xuwen
 * @create: 2020/02/07
 */
public class ContractInterfaceTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "xuwen", showName = "wasm.contract_interface接口特性测试",sourcePrefix = "wasm")
    public void testNewContract() {

        String name = "xuwen";
        try {
            prepare();
            ContractInterface contractInterface = ContractInterface.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractInterface.getContractAddress();
            String transactionHash = contractInterface.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractInterface issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Uint64 number = Uint64.of(10);
            TransactionReceipt transactionReceipt = contractInterface.setCount(number).send();
            collector.logStepPass("ContractInterface setCount successfully hash:" + transactionReceipt.getTransactionHash());

            Uint64 count = contractInterface.getCount().send();
            collector.assertEqual(count,number);
        } catch (Exception e) {
            collector.logStepFail("ContractInterfaceTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
