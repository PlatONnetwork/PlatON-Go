package wasm.contract_create;

import com.platon.rlp.datatypes.Uint16;
import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithVector;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init包含vector测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class InitWithVectorTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带vector",sourcePrefix = "wasm")
    public void testNewContract() {

        Uint16 age = Uint16.of(20);

        try {
            prepare();
            InitWithVector initWithVector = InitWithVector.deploy(web3j, transactionManager, provider,age).send();
            String contractAddress = initWithVector.getContractAddress();
            String transactionHash = initWithVector.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithVector issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Uint8 idx = Uint8.of(0);
            Uint64 chainAge = initWithVector.get_vector(idx).send();
            collector.assertEqual(chainAge.value,age.value);

        } catch (Exception e) {
            collector.logStepFail("InitWithVectorTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
