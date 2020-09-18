package network.platon.test.evm.complexcontracts;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.junit.rules.AssertCollector;
import network.platon.autotest.junit.rules.DriverService;
import network.platon.contracts.evm.Ballot;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;

/**
 * @title 投票功能合约测试
 * @description:
 * @author: qcxiao
 * @create: 2019/12/18 15:09
 **/
public class BallotTest {
    @Rule
    public AssertCollector collector = new AssertCollector();

    @Rule
    public DriverService driverService = new DriverService();

    // 底层链ID
    private long chainId;

    @Before
    public void before() {
        chainId = Integer.valueOf(driverService.param.get("chainId"));
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", author = "qcxiao"
            , showName = "complexcontracts.BallotTest-投票功能合约", sourcePrefix = "evm")
    public void testBallot() {
        Web3j web3j = null;
        Credentials credentials = null;
        try {
            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
            credentials = Credentials.create(driverService.param.get("privateKey"));
            collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("The node is unable to connect", e.toString());
            e.printStackTrace();
        }


        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {
            Ballot ballot = Ballot.deploy(web3j, transactionManager, provider, chainId, BigInteger.valueOf(100)).send();

        } catch (Exception e) {
            collector.logStepFail("ballot deploy fail.", e.toString());
            e.printStackTrace();
        }

    }

}
