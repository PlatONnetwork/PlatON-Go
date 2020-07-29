package wasm.contract_create;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitOverload;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约带空init函数测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/07
 */
public class InitOverloadTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约",sourcePrefix = "wasm")
    public void testNewContract() {

        String name = "hudenian";
        try {
            prepare();
            InitOverload initOverload = InitOverload.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = initOverload.getContractAddress();
            String transactionHash = initOverload.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitOverload issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initOverload.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = initOverload.add_vector(name).send();
            collector.logStepPass("InitOverload add_vector successfully hash:" + transactionReceipt.getTransactionHash());

            Uint8 idx = Uint8.of(0);
            String chainName = initOverload.get_vector(idx).send();
            collector.assertEqual(chainName,name);

            Uint64 size = initOverload.get_vector_size().send();
            collector.logStepPass("vector size is:"+size);
        } catch (Exception e) {
            collector.logStepFail("InitOverloadTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
