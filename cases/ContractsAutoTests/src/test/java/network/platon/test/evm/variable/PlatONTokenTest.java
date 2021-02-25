package network.platon.test.evm.variable;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.PlatONToken;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;



/**
 * @title 对 PlatON  币的几个单位进行测试
 * @description:
 * @author: liweic
 * @create: 2020/01/15 9:50
 **/

public class PlatONTokenTest extends ContractPrepareTest {

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "network.platon.test.evm.function.TimeTest-Token单位测试", sourcePrefix = "evm")
    public void PlatonTokens() {
        try {
            PlatONToken platonToken = PlatONToken.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = platonToken.getContractAddress();
            TransactionReceipt tx = platonToken.getTransactionReceipt().get();
            collector.logStepPass("PlatONToken deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());


            //lat
            BigInteger getlat = platonToken.Plat().send();
            collector.logStepPass("plat返回值：" + getlat);
            collector.assertEqual(new BigInteger("1000000000000000000"),getlat);

            //finney
            BigInteger getfinney = platonToken.Pfinney().send();
            collector.logStepPass("getfinney返回值：" + getfinney);
            collector.assertEqual(new BigInteger("1000000000000000"),getfinney);

            //szabo
            BigInteger getszabo = platonToken.Pszabo().send();
            collector.logStepPass("getszabo返回值：" + getszabo);
            collector.assertEqual(new BigInteger("1000000000000"),getszabo);

            //von
            BigInteger getvon = platonToken.Pvon().send();
            collector.logStepPass("getvon返回值：" + getvon);
            collector.assertEqual(new BigInteger("1"),getvon);


        } catch (Exception e) {
            collector.logStepFail("PlatONTokenContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}



