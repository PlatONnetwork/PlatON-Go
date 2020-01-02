package lib;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryUsingForAll;
import org.apache.commons.lang.StringUtils;
import org.junit.Test;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.PlatonGetTransactionCount;
import org.web3j.protocol.core.methods.response.PlatonSendTransaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.utils.Numeric;

import java.io.IOException;
import java.math.BigInteger;
import java.util.List;

/**
 * @title 引用using for方式验证
 * 解释：using A for * 的效果是，库 A 中的函数被附加在任意的类型上。
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryUsingForAllTest extends ContractPrepareTest {
    protected static final BigInteger GAS_LIMIT = BigInteger.valueOf(4700000);
    protected static final BigInteger GAS_PRICE = BigInteger.valueOf(1000000000L);

    public static final int DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH = 40;
    public static final long DEFAULT_POLLING_FREQUENCY = 2 * 1000;
    String LIBRARY_BINARY="610124610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060335760003560e01c806324fef5c8146038575b600080fd5b606b60048036036040811015604c57600080fd5b8101908080359060200190929190803590602001909291905050506081565b6040518082815260200191505060405180910390f35b600080600090505b838054905081101560c4578284828154811060a057fe5b9060005260206000200154141560b8578091505060e9565b80806001019150506089565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90505b9291505056fea265627a7a72315820b38b2a095812dd6f7e05f7afb39c8e72d1fca7a015307728a577672d339abcb864736f6c634300050d0032";
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "lib.LibraryUsingForAllTest-using A for all type")
    public void testReplace() {
        try {
            prepare();
            TransactionReceipt libReceipt =this.deployLib(GAS_PRICE,GAS_LIMIT,LIBRARY_BINARY);
            String libAddress = libReceipt.getContractAddress();
            collector.logStepPass("libReceipt issued successfully.libAddress:" + libAddress + ", hash:" + libReceipt.getTransactionHash());
            libAddress = StringUtils.substringAfter(libAddress,"0x");
            replaceLibAddress(libAddress);
            LibraryUsingForAll using = LibraryUsingForAll.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryUsingForAll issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = using.replace(new BigInteger("12"),new BigInteger("14")).send();
            collector.assertEqual(receipt.getStatus(),"0x1" , "checkout using a for * success");
    } catch (Exception e) {
        e.printStackTrace();
    }
    }

    private TransactionReceipt deployLib(BigInteger gasPrice, BigInteger gasLimit, String data) throws Exception {
        PlatonGetTransactionCount platonGetTransactionCount = web3j
                .platonGetTransactionCount(credentials.getAddress(), DefaultBlockParameterName.LATEST).send();
        BigInteger nonce = platonGetTransactionCount.getTransactionCount();

        String to = "";
        BigInteger value = BigInteger.valueOf(0L);
        RawTransaction rawTransaction = RawTransaction.createTransaction(nonce, gasPrice, gasLimit, to, value, data);

        byte[] signedMessage = TransactionEncoder.signMessage(rawTransaction, chainId, credentials);
        String hexValue = Numeric.toHexString(signedMessage);
        PlatonSendTransaction platonSendTransaction = web3j.platonSendRawTransaction(hexValue).send();

        return processResponse(platonSendTransaction);
    }

    private TransactionReceipt processResponse(PlatonSendTransaction transactionResponse) throws IOException, TransactionException {
        if (transactionResponse.hasError()) {
            throw new RuntimeException("Error processing transaction request: " + transactionResponse.getError().getMessage());
        }

        String transactionHash = transactionResponse.getTransactionHash();

        return new PollingTransactionReceiptProcessor(web3j, DEFAULT_POLLING_FREQUENCY, DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH)
                .waitForTransactionReceipt(transactionHash);
    }


    private void replaceLibAddress(String address){
        String contractBinary= LibraryUsingForAll.BINARY;
        int startIndex=StringUtils.indexOf(contractBinary,"__$");
        int endIndex=StringUtils.indexOf(contractBinary,"$__");
        if(startIndex==0||endIndex==0){
            return;
        }
        String replaceStr =StringUtils.substring(contractBinary,startIndex,endIndex+3);
        LibraryUsingForAll.BINARY=StringUtils.replace(contractBinary,replaceStr,address);

    }
}
