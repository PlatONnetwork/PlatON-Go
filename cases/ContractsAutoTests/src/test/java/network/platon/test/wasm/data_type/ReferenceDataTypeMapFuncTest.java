package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeMapFuncContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @title map属性方法
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeMapFuncTest extends WASMContractPrepareTest {

    private String cycleNum;
    private String cycleMapNum;

    private String deleteIndex;
    private String deleteMapNum;

    private String insertKey;
    private String insertValue;
    private String insertMapNum;
    private String clearMapNum;

    @Before
    public void before() {
        cycleNum = driverService.param.get("cycleNum");
        cycleMapNum = driverService.param.get("cycleMapNum");
        deleteIndex = driverService.param.get("deleteIndex");
        deleteMapNum = driverService.param.get("deleteMapNum");
        insertKey = driverService.param.get("insertKey");
        insertValue = driverService.param.get("insertValue");
        insertMapNum = driverService.param.get("insertMapNum");
        clearMapNum = driverService.param.get("clearMapNum");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeMapTest验证Map集合属性方法",sourcePrefix = "wasm")
    public void testReferenceDataTypeMapFunc() {

         //部署合约
        ReferenceDataTypeMapFuncContract referenceDataTypeMapFuncContract = null;
        try {
            prepare();
            referenceDataTypeMapFuncContract = ReferenceDataTypeMapFuncContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeMapFuncContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeMapFuncContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeMapFuncContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMapFuncContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：map循环新增值
          TransactionReceipt  transactionReceipt = referenceDataTypeMapFuncContract.addMapByUint(Uint8.of(cycleNum)).send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map循环新增值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：循环map容器数量
            Uint8 actualValueSize = referenceDataTypeMapFuncContract.getMapBySize().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证循环map容器数量】 执行getMapBySize() actualValueSize:" + actualValueSize);
            collector.assertEqual(actualValueSize.getValue(),new BigInteger(cycleMapNum), "checkout  execute success.");
            //3、验证：map容器删除指定值
            TransactionReceipt transactionReceipt1 = referenceDataTypeMapFuncContract.deleteMapByIndex(Uint8.of(deleteIndex)).send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map容器删除指定值】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //获取删除容器数量
            Uint8 actualValueSize1 = referenceDataTypeMapFuncContract.getMapBySize().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证删除map容器数量】 执行getMapBySize() actualValueSize1:" + actualValueSize1);
            collector.assertEqual(actualValueSize1.getValue(),new BigInteger(deleteMapNum), "checkout  execute success.");

            //4、验证：map容器插入方法insert()
            TransactionReceipt  transactionReceipt2 = referenceDataTypeMapFuncContract.insertMapUint(Uint8.of(insertKey),insertValue).send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map容器插入方法insert()】 执行insertMap() successfully hash:" + transactionReceipt2.getTransactionHash());
            //数量
            Uint8 actualValueSize2 = referenceDataTypeMapFuncContract.getMapBySize().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map容器插入元素数量】 执行getMapBySize() actualValueSize2:" + actualValueSize2);
            collector.assertEqual(actualValueSize2.getValue(),new BigInteger(insertMapNum), "checkout  execute success.");

            //5、验证map清空方法clear()
            TransactionReceipt  transactionReceipt3 = referenceDataTypeMapFuncContract.clearMapUint().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map清空方法clear()】 执行clearMapString() successfully hash:" + transactionReceipt3.getTransactionHash());
            //数量
            Uint8 actualValueSize3 = referenceDataTypeMapFuncContract.getMapBySize().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map容器清空数量】 执行getMapStringByKeySize() actualValueSize3:" + actualValueSize3);
            collector.assertEqual(actualValueSize3.getValue(),new BigInteger(clearMapNum), "checkout  execute success.");
            //6、验证map容器判断空
            boolean actualValue4 = referenceDataTypeMapFuncContract.getMapIsEmpty().send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map容器判断空】 执行getMapIsEmpty() actualValue4:" + actualValue4);
            collector.assertEqual(actualValue4,true, "checkout  execute success.");


           /* //5、获取插入的值
            String actualValue2 = referenceDataTypeMapFuncContract.getMapStringByKey("01").send();
            collector.logStepPass("referenceDataTypeMapFuncContract 【验证map获取插入的值】 执行getMapStringByKey() actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,"lucy", "checkout  execute success.");*/

          /*  //4、验证：map容器数量
            Long actualPersonSize = referenceDataTypeMapTestContract.getMapByPersonSize().send();
            collector.logStepPass("referenceDataTypeMapContract 【验证map容器数量】 执行getMapByPersonSize() actualPersonSize:" + actualPersonSize);
            collector.assertEqual(actualPersonSize,Long.parseLong("1"), "checkout  execute success.");
            //5、验证：map容器根据key获取值Person
            String actualValueName = referenceDataTypeMapTestContract.getMapByPerson(Byte.valueOf(keyPerson)).send();
            collector.logStepPass("referenceDataTypeMapContract 【验证map容器根据key获取值】 执行getMapByPerson() actualValueName:" + actualValueName);
            //collector.assertEqual(actualValuePerson.name,personName, "checkout  execute success.");*/

        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeMapContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
