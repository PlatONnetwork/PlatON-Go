package wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeStructContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型结构体（Struct）
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeStructTest extends WASMContractPrepareTest {

    private String name;
    private String age;


    @Before
    public void before() {
        name = driverService.param.get("name");
        age = driverService.param.get("age");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeStructTest验证结构体赋值&取值",sourcePrefix = "wasm")
    public void testReferenceDataTypeStruct() {

         //部署合约
        ReferenceDataTypeStructContract referenceDataTypeStructContract = null;
        try {
            prepare();
            referenceDataTypeStructContract = ReferenceDataTypeStructContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeStructContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeStructContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeStructContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：定义struct类型并赋值
            TransactionReceipt  transactionReceipt = referenceDataTypeStructContract.setStructPersonA(name,Uint64.of(age)).send();
            collector.logStepPass("referenceDataTypeStructContract 【验证定义struct类型并赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：struct取值
            String actualValue = referenceDataTypeStructContract.getPersonName().send();
            collector.logStepPass("referenceDataTypeStructContract 【验证struct取值】 执行getPersonName() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,name, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeStructContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
