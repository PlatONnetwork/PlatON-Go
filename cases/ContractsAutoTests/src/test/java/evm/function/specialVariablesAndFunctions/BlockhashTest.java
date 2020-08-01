package evm.function.specialVariablesAndFunctions;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.BlockTransactionPropertiesFunctions;
import network.platon.contracts.Blockhash;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * @title blockHash函数(只能获取（当前块-256）至（当前块-1）之间的块hash)
 * @description:
 * @author: hudenian
 * @create: 2020/03/09 19:17
 **/

public class BlockhashTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.BlockhashTest-验证以blockhash函数", sourcePrefix = "evm")
    public void blockhashTest() {
        try {
            Blockhash blockhash = Blockhash.deploy(web3j, transactionManager, provider, chainId).send();

            collector.logStepPass(DataChangeUtil.bytesToHex(blockhash.getBlockhashbefore0().send()));
            collector.logStepPass(DataChangeUtil.bytesToHex(blockhash.getBlockhashbefore30().send()));
            collector.logStepPass(DataChangeUtil.bytesToHex(blockhash.getBlockhashbefore255().send()));
            collector.logStepPass(DataChangeUtil.bytesToHex(blockhash.getBlockhashbefore256().send()));
            collector.logStepPass(DataChangeUtil.bytesToHex(blockhash.getBlockhashbefore257().send()));

        } catch (Exception e) {
            collector.logStepFail("BlockTransactionPropertiesFunctionsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
