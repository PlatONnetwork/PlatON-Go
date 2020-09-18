package network.platon.test.evm.data_type.structs;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RecursiveStorageMemoryComplex;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.List;

/**
 * @title 复杂结构体数据类递归调用验证
 * @description:
 * @author: hudenian
 * @create: 2020/1/13 10:03
 **/
public class RecursiveStorageMemoryComplexTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls",sheetName = "structs", author = "hudenian", showName = "StructDataTypeTest.结构体数据类型",sourcePrefix = "evm")
    public void testRecursiveStorageMemoryComplex() {

        RecursiveStorageMemoryComplex recursiveStorageMemoryComplex = null;
        try {
            //合约部署
            recursiveStorageMemoryComplex = RecursiveStorageMemoryComplex.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = recursiveStorageMemoryComplex.getContractAddress();
            TransactionReceipt tx =  recursiveStorageMemoryComplex.getTransactionReceipt().get();
            collector.logStepPass("RecursiveStorageMemoryComplex issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

            //调用合约
            tx = recursiveStorageMemoryComplex.run().send();
            collector.logStepPass("RecursiveStorageMemoryComplexTest run() successful.transactionHash:" + tx.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + tx.getBlockNumber());

            //查询执行结果
            List resultList =  recursiveStorageMemoryComplex.getRunResult().send();

            collector.assertEqual("66",resultList.get(0).toString());
            collector.assertEqual("16896",resultList.get(1).toString());
            collector.assertEqual("4325376",resultList.get(2).toString());
            collector.assertEqual("4325377",resultList.get(3).toString());
            collector.assertEqual("4325378",resultList.get(4).toString());
            collector.assertEqual("16897",resultList.get(5).toString());
            collector.assertEqual("4325632",resultList.get(6).toString());
            collector.assertEqual("4325633",resultList.get(7).toString());
            collector.assertEqual("4325634",resultList.get(8).toString());
            collector.assertEqual("4325635",resultList.get(9).toString());

        } catch (Exception e) {
            collector.logStepFail("RecursiveStorageMemoryComplexTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
