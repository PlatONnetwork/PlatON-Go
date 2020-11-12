package network.platon.test.wasm.beforetest;

import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.http.HttpService;
import com.alaya.tx.RawTransactionManager;
import com.alaya.tx.gas.ContractGasProvider;
import com.alaya.utils.Numeric;
import network.platon.test.BaseTest;
import java.math.BigInteger;

/**
 * @title 所有合约部署前相关准备工作，需要初始化gas值
 * @author: albedo
 * @create: 2019/12/26 11:27
 **/
public class WASMContractPrepareTest extends BaseTest {

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

    protected String prependHexPrefix(String input){
        if(Numeric.containsHexPrefix(input)){
            return input;
        }
        return Numeric.prependHexPrefix(input);
    }
}
