package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeSetContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约引用类型(set类型)属性及函数
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeSetTest extends WASMContractPrepareTest {

    private String initSetLength;
    private String insertSetValue;
    private String setLength;
    private String findValue;
    private String eraseValue;
    private String deleteValue;
    private String deleteSetLength;
    private String clearSetLength;



    @Before
    public void before() {
        initSetLength = driverService.param.get("initSetLength");
        insertSetValue = driverService.param.get("insertSetValue");
        setLength = driverService.param.get("setLength");
        findValue = driverService.param.get("findValue");
        eraseValue = driverService.param.get("eraseValue");
        deleteValue = driverService.param.get("deleteValue");
        deleteSetLength = driverService.param.get("deleteSetLength");
        clearSetLength = driverService.param.get("clearSetLength");

    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeSetTest(set类型)属性函数",sourcePrefix = "wasm")
    public void testReferenceDataTypeSet() {
         //部署合约
        ReferenceDataTypeSetContract referenceDataTypeSetContract = null;
        try {
            prepare();
            referenceDataTypeSetContract = ReferenceDataTypeSetContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeSetContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeSetContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeSetContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeSetContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：set类型初始化赋值
            TransactionReceipt  transactionReceipt = referenceDataTypeSetContract.init_set().send();
            collector.logStepPass("referenceDataTypeSetContract 【验证(set类型初始化赋值】 successfully hash:" + transactionReceipt.getTransactionHash());
            //set类型获取容器大小
            getSetLength(referenceDataTypeSetContract,initSetLength);

            //2、验证：set类型插入数据insert
            TransactionReceipt  transactionReceipt1 = referenceDataTypeSetContract.insert_set(Uint8.of(insertSetValue)).send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型插入数据insert】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //set类型获取容器大小
            getSetLength(referenceDataTypeSetContract,setLength);
            //3、验证:set类型查找数据find
            Uint8 actualSetValue = referenceDataTypeSetContract.find_set().send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型查找数据find】 执行find_set() successfully actualSetValue:" + actualSetValue);
            collector.assertEqual(actualSetValue,Uint8.of(findValue), "checkout execute success.");
           //4、验证:set类型删除元素erase
            TransactionReceipt  transactionReceipt2 = referenceDataTypeSetContract.erase_set(Uint8.of(eraseValue)).send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型删除元素erase】 successfully hash:" + transactionReceipt2.getTransactionHash());
            getSetLength(referenceDataTypeSetContract,deleteSetLength);
            //5、验证:set类型判断是否为空
            Boolean actualEmptyValue = referenceDataTypeSetContract.get_set_empty().send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型判断是否为空】 执行find_set() successfully actualEmptyValue:" + actualEmptyValue);
            //6、验证:set类型清空clear
            TransactionReceipt  transactionReceipt3 = referenceDataTypeSetContract.clear_set().send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型清空clear】 successfully hash:" + transactionReceipt3.getTransactionHash());
            getSetLength(referenceDataTypeSetContract,clearSetLength);

        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeVectorFuncContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

    public void getSetLength(ReferenceDataTypeSetContract referenceDataTypeSetContract,String setLength){

        try {
            Uint64 actualSetLength = referenceDataTypeSetContract.get_set_size().send();
            collector.logStepPass("referenceDataTypeSetContract 【验证set类型获取容器大小】 执行get_set_size() successfully actualSetLength:" + actualSetLength);
            collector.assertEqual(actualSetLength,Uint64.of(setLength), "checkout execute success.");
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("referenceDataTypeSetContract Calling Method fail.", e.toString());
        }


    }


}
