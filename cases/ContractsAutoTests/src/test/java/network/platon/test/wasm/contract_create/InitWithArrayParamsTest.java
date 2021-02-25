package network.platon.test.wasm.contract_create;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithArrayParams;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init入参包含array
 * @description:
 * @author: hudenian
 * @create: 2020/02/27
 */
public class InitWithArrayParamsTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init入参包含array",sourcePrefix = "wasm")
    public void testArrayParams() {

        String[] array = new String[]{"array1","array2","array3","array4","array5","array6","array7","array8","array9","array10"};
        try {
            prepare();
            InitWithArrayParams initWithArrayParams = InitWithArrayParams.deploy(web3j, transactionManager, provider, chainId, array).send();
            String contractAddress = initWithArrayParams.getContractAddress();
            String transactionHash = initWithArrayParams.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithArrayParamsTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithArrayParams.getTransactionReceipt().get().getGasUsed());

            //查询包含array对象
            String[] chainArray = initWithArrayParams.get_array().send();

            collector.assertEqual(chainArray[0].toString(),array[0].toString());
            collector.assertEqual(chainArray[1].toString(),array[1].toString());
            collector.assertEqual(chainArray[2].toString(),array[2].toString());

            //查看数组大小
            Uint8 arrarySize = initWithArrayParams.get_array_size().send();
            collector.assertEqual("10",arrarySize.value.toString());

            //查看数组是否包含指定的元素
           boolean flg = initWithArrayParams.get_array_contain_element("array1").send();
           collector.assertEqual(true,flg);

           flg = initWithArrayParams.get_array_contain_element("array01").send();
           collector.assertEqual(false,flg);

        } catch (Exception e) {
            collector.logStepFail("InitWithArrayParamsTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
