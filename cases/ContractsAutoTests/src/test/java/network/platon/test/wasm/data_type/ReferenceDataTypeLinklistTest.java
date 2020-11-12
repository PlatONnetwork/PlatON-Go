package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeLinkedlistContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型链表
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeLinklistTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeLinklistTest验证引用类型链表",sourcePrefix = "wasm")
    public void testReferenceDataTypeLinklist() {

         //部署合约
        ReferenceDataTypeLinkedlistContract referenceDataTypeLinkedlistContract = null;
        try {
            prepare();
            referenceDataTypeLinkedlistContract = ReferenceDataTypeLinkedlistContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeLinkedlistContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeLinkedlistContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeLinkedlistContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeLinkedlistContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:链表新增
            String nodeData = "one";
            String nodeData1 = "two";
            TransactionReceipt  transactionReceipt = referenceDataTypeLinkedlistContract.insertNodeElement(nodeData).send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证链表新增】 successfully hash:" + transactionReceipt.getTransactionHash());
            TransactionReceipt  transactionReceipt2 = referenceDataTypeLinkedlistContract.insertNodeElement(nodeData1).send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证链表新增】 successfully hash:" + transactionReceipt2.getTransactionHash());
           //2、验证:获取结点数据
            Uint8 a = Uint8.of(0);
            Uint8 b = Uint8.of(0);
            Uint8 expectValue = Uint8.of(1);
            Uint8 actualValueIndex = referenceDataTypeLinkedlistContract.getNodeElementIndex(a,b).send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证获取结点数据】 执行getNodeElementIndex() successfully actualValueIndex:" + actualValueIndex);
            collector.assertEqual(actualValueIndex,expectValue, "checkout  execute success.");
           //3、验证:遍历结点链表数据
            Uint8 index = Uint8.of(0);
            String[]  arrayStr = referenceDataTypeLinkedlistContract.findNodeElement(index).send();
            for (int i = 0; i < arrayStr.length; i++) {
                collector.logStepPass("referenceDataTypeLinkedlistContract 【验证获取结点数据】 执行findNodeElement() successfully nodeName:" +  arrayStr[i]);
            }
            //4、验证:链表删除数据
            TransactionReceipt  transactionReceipt3 = referenceDataTypeLinkedlistContract.clearNodeElement().send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证链表删除数据】 successfully hash:" + transactionReceipt3.getTransactionHash());
            //清空后查询数据
            Uint8 index1 = Uint8.of(0);
            String[]  arrayStr1 = referenceDataTypeLinkedlistContract.findNodeElement(index1).send();
            collector.assertEqual(arrayStr1.length,0, "checkout  execute success.");


        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeLinkedlistContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
