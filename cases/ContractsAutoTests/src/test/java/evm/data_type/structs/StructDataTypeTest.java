package evm.data_type.structs;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.StructDataType;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple6;

import java.math.BigInteger;

/**
 * @title 测试：结构体数据类型
 * @description:
 * @author: hudenian
 * @create: 2020/1/11 16:03
 **/
public class StructDataTypeTest extends ContractPrepareTest {

    @Before
    public void before() {
       this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls",sheetName = "structs", author = "hudenian", showName = "StructDataTypeTest.结构体数据类型",sourcePrefix = "evm")
    public void testTypeConversionContract() {

        StructDataType structDataType = null;
        try {
            //合约部署
            structDataType = StructDataType.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = structDataType.getContractAddress();
            TransactionReceipt tx =  structDataType.getTransactionReceipt().get();
            collector.logStepPass("StructDataType issued successfully.contractAddress:" + contractAddress
                                    + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

            //调用合约
            tx = structDataType.run().send();
            collector.logStepPass("StructDataTypeTest run() successful.transactionHash:" + tx.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + tx.getBlockNumber());

            //查询执行结果
            Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>  tuple6 =  structDataType.getRunValue().send();

            collector.assertEqual("2",tuple6.getValue1().toString());
            collector.assertEqual("2",tuple6.getValue2().toString());
            collector.assertEqual("2",tuple6.getValue3().toString());
            collector.assertEqual("6",tuple6.getValue4().toString());
            collector.assertEqual("9",tuple6.getValue5().toString());
            collector.assertEqual("7",tuple6.getValue6().toString());

        } catch (Exception e) {
            collector.logStepFail("StructDataTypeTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
