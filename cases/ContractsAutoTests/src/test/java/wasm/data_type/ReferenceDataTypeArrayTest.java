package wasm.data_type;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithParams;
import network.platon.contracts.wasm.ReferenceDataTypeArrayContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型数组(array类型)
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeArrayTest extends WASMContractPrepareTest {

    private String indexNo;
    private String expectIndexValue;
    private String expectArrayLength;
    private String indexByteNo;

    @Before
    public void before() {
        indexNo = driverService.param.get("indexNo");
        expectIndexValue = driverService.param.get("expectIndexValue");
        expectArrayLength = driverService.param.get("expectArrayLength");
        indexByteNo = driverService.param.get("indexByteNo");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeArrayTest验证数组定义及赋值",sourcePrefix = "wasm")
    public void testReferenceDataTypeArray() {

         //部署合约
        ReferenceDataTypeArrayContract referenceDataTypeArrayContract = null;
        try {
            prepare();
            referenceDataTypeArrayContract = ReferenceDataTypeArrayContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeArrayContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeArrayContract.getTransactionReceipt().get();
            collector.logStepPass("ReferenceDataTypeArrayContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：数组初始化赋值
            TransactionReceipt  transactionReceipt = referenceDataTypeArrayContract.setInitArray().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证数组初始化赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：数组取值
            String actualValue = referenceDataTypeArrayContract.getArrayStringIndex(Uint32.of(indexNo)).send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证数组取值】 执行getArrayIndex() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");
            //3、验证：获取数组容器大小
            Uint8 actualArrayLength = referenceDataTypeArrayContract.getArrayUintSize().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证获取数组容器大小】 执行getArraySize() successfully actualArrayLength:" + actualArrayLength);
            collector.assertEqual(actualArrayLength,Uint8.of(expectArrayLength), "checkout  execute success.");
            //4、验证：数组定义person类型
           /*ReferenceDataTypeArrayContract.Person person = new ReferenceDataTypeArrayContract.Person();
            person.name = "lucy";
            person.age = Long.parseLong("20");
            TransactionReceipt  transactionReceipt2 =referenceDataTypeArrayContract.setArrayPerson(person).send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证person类型数组赋值】 successfully hash:" + transactionReceipt2.getTransactionHash());
            //person类型取值
            String actualValueName = referenceDataTypeArrayContract.getArrayPersonNameIndex().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证person类型数组取值】 执行getArrayPersonNameIndex() successfully actualValueName:" + actualValueName);
            collector.assertEqual(actualValueName,"lucy", "checkout  execute success.");*/
           //5、验证：字节数组赋值&取值
            TransactionReceipt receipt = referenceDataTypeArrayContract.setBytesArray().send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证字节数组赋值】 successfully hash:" + receipt.getTransactionHash());
            Uint8 actualByteValue = referenceDataTypeArrayContract.getBytesArrayIndex(Uint32.of(indexByteNo)).send();
            collector.logStepPass("ReferenceDataTypeArrayContract 【验证字节数组取值】 执行getBytesArrayIndex() successfully actualByteValue:" + actualByteValue);
            //collector.assertEqual(actualByteValue,expectByteValue, "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("ReferenceDataTypeArrayContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
