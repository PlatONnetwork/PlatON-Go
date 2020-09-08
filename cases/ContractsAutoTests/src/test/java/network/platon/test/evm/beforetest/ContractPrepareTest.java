package network.platon.test.evm.beforetest;

import network.platon.test.BaseTest;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;

/**
 * @title 所有合约部署前相关准备工作，需要初始化gas值
 * @author: albedo
 * @create: 2019/12/26 11:27
 **/
public class ContractPrepareTest extends BaseTest {

    protected Web3j web3j = null;
    protected Credentials credentials = null;
    protected RawTransactionManager transactionManager;
    protected ContractGasProvider provider;
    protected String walletAddress;
    protected String gasLimit = "4712388";
    protected String gasPrice = "3000000000000000";

    /**
     * 合约部署相关准备工作
     * @param gasPrice
     * @param gasLimit
     */
    protected void prepare(String gasPrice,String gasLimit){
        chainId = Integer.valueOf(driverService.param.get("chainId"));
        try {
            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
            credentials = Credentials.create(driverService.param.get("privateKey"));
            collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
            walletAddress = driverService.param.get("address");
        } catch (Exception e) {
            collector.logStepFail("The node is unable to connect", e.toString());
            e.printStackTrace();
        }
        provider = new ContractGasProvider(new BigInteger(gasPrice), new BigInteger(gasLimit));
        transactionManager = new RawTransactionManager(web3j, credentials, chainId);
    }

    /**
     * 合约部署相关准备工作
     */
    protected void prepare(){
        this.prepare(gasPrice,gasLimit);
    }
}
