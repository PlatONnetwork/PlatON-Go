package network.platon.test;

import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.junit.rules.AssertCollector;
import network.platon.autotest.junit.rules.DriverService;
import org.junit.Rule;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.ECKeyPair;
import org.web3j.crypto.Keys;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.security.InvalidAlgorithmParameterException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;

public class BaseTest {

    @Rule
    public DriverService driverService = new DriverService();
    @Rule
    public AssertCollector collector = new AssertCollector();

    public long chainId;

    public static String getJsonString(Object object){
        return JSONObject.toJSONString(object);
    }

//    public Credentials createCredentials(String privateKey){
//        Credentials credentials;
//        if(this.getParam("secretType").equals("sm2")){
//            credentials = Credentials.createSm2(privateKey);
//        }else {
//            credentials = Credentials.create(privateKey);
//        }
//        return credentials;
//    }
//
//    public String getCredentialsAddress(Credentials credentials){
//        return credentials.getAddress();
//    }
//
//    public ECKeyPair generateEcKeyPair(){
//        ECKeyPair ecKeyPair = null;
//        try {
//            if(this.getParam("secretType").equals("sm2")){
//                ecKeyPair = Keys.createSm2EcKeyPair();
//            }else {
//                ecKeyPair = Keys.createEcKeyPair();
//            }
//        } catch (InvalidAlgorithmParameterException e) {
//            e.printStackTrace();
//        } catch (NoSuchAlgorithmException e) {
//            e.printStackTrace();
//        } catch (NoSuchProviderException e) {
//            e.printStackTrace();
//        }
//        return ecKeyPair;
//    }
//
//    public ECKeyPair generateEcKeyPair(byte[] privateKey){
//        ECKeyPair ecKeyPair;
//        if(this.getParam("secretType").equals("sm2")){
//            ecKeyPair = ECKeyPair.createSm2(privateKey);
//        }else {
//            ecKeyPair = ECKeyPair.create(privateKey);
//        }
//        return ecKeyPair;
//    }

//    public ECKeyPair generateEcKeyPair(String privateKey){
//        return this.generateEcKeyPair(Numeric.hexStringToByteArray(privateKey));
//    }
//
//    public String generateNewWallet() {
//        ECKeyPair ecKeyPair = this.generateEcKeyPair();
//        Credentials credentials = this.createCredentials(ecKeyPair.getPrivateKey().toString(16));
//        return getCredentialsAddress(credentials);
//    }

//    public String generatePrivateKey() {
//        ECKeyPair ecKeyPair = this.generateEcKeyPair();
//        //私钥长度不足64位时需要补0
//        String privateKey = ecKeyPair.getPrivateKey().toString(16);
//        int len = privateKey.length();
//        for(int i=0;i<(64-len);i++){
//            privateKey = "0" + privateKey;
//        }
//        return privateKey;
//    }
//
//    public String generatePublicKeyFromPrivateKey(String privateKey){
//        ECKeyPair ecKeyPair = this.generateEcKeyPair(privateKey);
//        String publicKey = Numeric.toHexStringWithPrefix(ecKeyPair.getPublicKey());
//        //公钥长度不足128位时需要补0
////        int len = publicKey.length();
////        for(int i=0;i<(128-len);i++){
////            publicKey = "0" + publicKey;
////        }
//        return publicKey;
//    }
//
//    public String generatePublicKey(){
//        ECKeyPair ecKeyPair = this.generateEcKeyPair();
//        String publicKey = ecKeyPair.getPublicKey().toString(16);
//        //公钥长度不足128位时需要补0
//        int len = publicKey.length();
//        for(int i=0;i<(128-len);i++){
//            publicKey = "0" + publicKey;
//        }
//        return publicKey;
//    }
//
//    public void transfer(String transferToPrivateKey) throws Exception {
//        Web3j web3j = Web3j.build(new HttpService(this.getParam("nodeUrl")));
//        Credentials credentials = this.createCredentials(this.getParam("privateKey"));
//        collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
//        String transferFrom = getCredentialsAddress(this.createCredentials(driverService.param.get("privateKey")));
//        String transferTo = this.createCredentials(transferToPrivateKey).getAddress();
//        BigInteger transferFromInitialBalance = web3j.platonGetBalance(transferFrom, DefaultBlockParameterName.LATEST).send().getBalance();
//
//        if (transferFromInitialBalance.compareTo(new BigInteger("1000000000000000000")) < 0) {
//            return;
//        }
//
//        BigInteger initialBalance = web3j.platonGetBalance(transferTo, DefaultBlockParameterName.LATEST).send().getBalance();
//
//        TransactionReceipt transactionReceipt = new Transfer(web3j,new RawTransactionManager(web3j, credentials, chainId)).sendFunds(transferTo, new BigDecimal(1), Convert.Unit.LAT).send();
//        BigInteger endBalance = web3j.platonGetBalance(transferTo, DefaultBlockParameterName.LATEST).send().getBalance();
//
//        if (endBalance.compareTo(initialBalance) == 0) {
//            throw new RuntimeException("转账失败");
//        }
//
//        System.out.println("txHash:" + transactionReceipt.getTransactionHash());
//        System.out.println("address:" + transferTo);
//
//    }

    /**
     * 获取链管理者私钥
     *
     * @return
     */
    public String getAdmin() {
        return driverService.param.get("privateKey");
    }

    public String getParam(String name) {
        return driverService.param.get(name);
    }

    public int getIntParam(String name) {
        return Integer.parseInt(getParam(name));
    }

}
