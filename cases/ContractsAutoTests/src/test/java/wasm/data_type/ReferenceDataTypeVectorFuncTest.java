package wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeVectorContract;
import network.platon.contracts.wasm.ReferenceDataTypeVectorFuncContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约引用类型(verctor类型)属性/函数
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeVectorFuncTest extends WASMContractPrepareTest {

    private String insertValue1;
    private String vectorLength1;
    private String insertValue2;
    private String insertNum2;
    private String vectorLength2;
    private String expectValueAtIndex;
    private String expectValueAtValue;
    private String FrontValue;
    private String BackValue;
    private String vectorLength3;
    private String vectorLength4;
    private String vectorLength5;



    @Before
    public void before() {
        insertValue1 = driverService.param.get("insertValue1");
        vectorLength1 = driverService.param.get("vectorLength1");
        insertValue2 = driverService.param.get("insertValue2");
        insertNum2 = driverService.param.get("insertNum2");
        vectorLength2 = driverService.param.get("vectorLength2");
        expectValueAtIndex = driverService.param.get("expectValueAtIndex");
        expectValueAtValue = driverService.param.get("expectValueAtValue");
        FrontValue = driverService.param.get("FrontValue");
        BackValue = driverService.param.get("BackValue");
        vectorLength3 = driverService.param.get("vectorLength3");
        vectorLength4 = driverService.param.get("vectorLength4");
        vectorLength5 = driverService.param.get("vectorLength5");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeVectorFuncTest(verctor类型)属性函数",sourcePrefix = "wasm")
    public void testReferenceDataTypeVectorFunc() {

         //部署合约
        ReferenceDataTypeVectorFuncContract referenceDataTypeVectorFuncContract = null;
        try {
            prepare();
            referenceDataTypeVectorFuncContract = ReferenceDataTypeVectorFuncContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeVectorFuncContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeVectorFuncContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeVectorFuncContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeVectorFuncContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：(增加函数)指定位置插入数据
            TransactionReceipt  transactionReceipt = referenceDataTypeVectorFuncContract.insertVectorValue(insertValue1).send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证(增加函数)将值添加到begin()起始位置之前】 successfully hash:" + transactionReceipt.getTransactionHash());
            //vector类型获取容器大小
            getVectorLength(referenceDataTypeVectorFuncContract,vectorLength1);

            //2、验证：(增加函数)指定位置插入多个相同数据
            TransactionReceipt  transactionReceipt1 = referenceDataTypeVectorFuncContract.insertVectorMangValue(Uint64.of(insertNum2),insertValue2).send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证(增加函数)指定位置插入多个相同数据】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //vector类型获取容器大小
            getVectorLength(referenceDataTypeVectorFuncContract,vectorLength2);

            //3、验证at()函数，返回index位置元素的引用 [two,two,one]
            String actualVectorAtValue = referenceDataTypeVectorFuncContract.findVectorAt(Uint64.of(expectValueAtIndex)).send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证at()函数，返回index位置元素的引用】 执行findVectorAt() successfully actualVectorAtValue:" + actualVectorAtValue);
            collector.assertEqual(actualVectorAtValue,expectValueAtValue, "checkout execute success.");


            //4、验证front():返回首元素的引用
            String actualVectorFrontValue = referenceDataTypeVectorFuncContract.findVectorFront().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证front():返回首元素的引用】 执行findVectorFront() successfully actualVectorFrontValue:" + actualVectorFrontValue);
            collector.assertEqual(actualVectorFrontValue,FrontValue, "checkout execute success.");
            //5、验证 back():返回尾元素的引用
            String actualVectorBackValue = referenceDataTypeVectorFuncContract.findVectorBack().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证 back():返回尾元素的引用】 执行findVectorBack() successfully actualVectorBackValue:" + actualVectorBackValue);
            collector.assertEqual(actualVectorBackValue,BackValue, "checkout execute success.");

            //6、验证pop_back()删除最后一个元素
            TransactionReceipt  transactionReceipt2 = referenceDataTypeVectorFuncContract.deleteVectorPopBack().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证pop_back()删除最后一个元素】 successfully hash:" + transactionReceipt2.getTransactionHash());
            //vector类型获取容器大小
            getVectorLength(referenceDataTypeVectorFuncContract,vectorLength3);
            //7、验证erase()删除指定元素,将起始位置的元素删除
            TransactionReceipt  transactionReceipt3 = referenceDataTypeVectorFuncContract.deleteVectorErase().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证erase()删除指定元素,将起始位置的元素删除】 successfully hash:" + transactionReceipt3.getTransactionHash());
            //vector类型获取容器大小
            getVectorLength(referenceDataTypeVectorFuncContract,vectorLength4);
            //8、clear()清空元素
            TransactionReceipt  transactionReceipt4 = referenceDataTypeVectorFuncContract.deleteVectorClear().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证clear()清空元素】 successfully hash:" + transactionReceipt4.getTransactionHash());
            //vector类型获取容器大小
            getVectorLength(referenceDataTypeVectorFuncContract,vectorLength5);
            //9、empty()判断函数
            Boolean isEmpty = referenceDataTypeVectorFuncContract.findVectorEmpty().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证empty()判断函数】 执行findVectorEmpty() successfully isEmpty:" + isEmpty);
            collector.assertEqual(isEmpty,true, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeVectorFuncContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }

    public void getVectorLength(ReferenceDataTypeVectorFuncContract referenceDataTypeVectorFuncContract,String vectorLength){

        try {
            Uint64 actualVectorLength = referenceDataTypeVectorFuncContract.getVectorLength().send();
            collector.logStepPass("referenceDataTypeVectorFuncContract 【验证vector类型获取容器大小】 执行getVectorLength() successfully actualVectorLength:" + actualVectorLength);
            collector.assertEqual(actualVectorLength,Uint64.of(vectorLength), "checkout execute success.");
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("referenceDataTypeVectorFuncContract Calling Method fail.", e.toString());
        }


    }


}
