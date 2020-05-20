package evm.versioncompatible.v0_5_0.v3_deprecatedType;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.HexLiteralsChangeByte;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * 16进制数值如果长度与 bytesX不相等，也不能直接转换成 bytesX类型
 * @description:
 * @author: hudenian
 * @create: 2019/12/26
 */
public class HexLiteralsChangeByteTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "version_compatible.0.5.0.HexLiteralsChangeByteTest-弃用类型转换", sourcePrefix = "evm")
    public void update() {
        try {

            HexLiteralsChangeByte hexLiteralsChangeByte = HexLiteralsChangeByte.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = hexLiteralsChangeByte.getContractAddress();
            TransactionReceipt tx = hexLiteralsChangeByte.getTransactionReceipt().get();

            collector.logStepPass("FunctionDeclaraction deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + hexLiteralsChangeByte.getTransactionReceipt().get().getGasUsed());

            TransactionReceipt transactionReceipt = hexLiteralsChangeByte.testChange(new BigInteger("1")).send();

            collector.logStepPass("FunctionDeclaraction update_public successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            byte[] afterValueByte = hexLiteralsChangeByte.getY().send(); //f1

            String afterValue = DataChangeUtil.bytesToHex(afterValueByte).toLowerCase();

//            collector.logStepPass(initValue+"转成bytes4后的值为："+afterValue);

            collector.assertEqual("f1",afterValue);
        } catch (Exception e) {
            collector.logStepFail("HexLiteralsChangeByteTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
