package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeLinkedlistContract;
import network.platon.contracts.wasm.ReferenceDataTypeMapContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

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
            referenceDataTypeLinkedlistContract = ReferenceDataTypeLinkedlistContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = referenceDataTypeLinkedlistContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeLinkedlistContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeLinkedlistContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeLinkedlistContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：定义单向链表
            TransactionReceipt  transactionReceipt = referenceDataTypeLinkedlistContract.setListNode().send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证定义单向链表】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：链表新增结点
            TransactionReceipt  transactionReceipt1 = referenceDataTypeLinkedlistContract.addListNode().send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证链表新增结点】 执行addListNode() successfully hash:" + transactionReceipt1.getTransactionHash());
            //3、验证：获取结点指针地址
            Long actualAdrr = referenceDataTypeLinkedlistContract.getListNode().send();
            collector.logStepPass("referenceDataTypeLinkedlistContract 【验证获取结点指针地址】 执行getListNode() successfully actualAdrr:" + actualAdrr);
            //collector.assertEqual(actualMapLength,expectMapLength, "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeLinkedlistContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
