package wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeStructContract;
import network.platon.contracts.wasm.ReferenceDataTypeVectorContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试引用类型（Vector）
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeVectorTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeVectorTest验证Vector容器赋值及属性",sourcePrefix = "wasm")
    public void testReferenceDataTypeVector() {
         //部署合约
        ReferenceDataTypeVectorContract referenceDataTypeVectorContract = null;
        try {
            prepare();
            referenceDataTypeVectorContract = ReferenceDataTypeVectorContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeVectorContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeVectorContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeVectorContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeVectorContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：vector类型赋值传参对象
            ReferenceDataTypeVectorContract.Clothes clothes = new ReferenceDataTypeVectorContract.Clothes();
            clothes.color = "yellow";
            TransactionReceipt  transactionReceipt = referenceDataTypeVectorContract.setClothesColorOne(clothes).send();
            collector.logStepPass("referenceDataTypeVectorContract 【验证vector类型赋值传参对象】 successfully hash:" + transactionReceipt.getTransactionHash());
            //2、验证：vector类型赋值传参字符串
            String myColor = "bule";
            TransactionReceipt  transactionReceipt1 = referenceDataTypeVectorContract.setClothesColorTwo(myColor).send();
            collector.logStepPass("referenceDataTypeVectorContract 【验证vector类型赋值传参字符串】 successfully hash:" + transactionReceipt1.getTransactionHash());
            //3、验证：vector类型根据索引取值
            String actualValue = referenceDataTypeVectorContract.getClothesColorIndex().send();
            collector.logStepPass("referenceDataTypeVectorContract 【验证vector类型根据索引取值】 执行getClothesColorIndex() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,"yellow", "checkout execute success.");
            //4、验证：vector类型获取容器大小
            Uint64 expectValue = Uint64.of("2");
            Uint64 actualVectorLength = referenceDataTypeVectorContract.getClothesColorLength().send();
            collector.logStepPass("referenceDataTypeVectorContract 【验证vector类型根据索引取值】 执行getClothesColorLength() successfully actualValue:" + actualValue);
            collector.assertEqual(actualVectorLength,expectValue, "checkout execute success.");
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeVectorContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
