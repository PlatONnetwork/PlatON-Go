package lib;

import beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.LibraryUsingFor;
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
 * 解释：指令using A for B 可用于附加库函数（从库 A）到任何类型（B）。 这些函数将接收到调用它们的对象作为它们的第一个参数。
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class LibraryUsingForTest extends ContractPrepareTest {
    protected static final BigInteger GAS_LIMIT = BigInteger.valueOf(4700000);
    protected static final BigInteger GAS_PRICE = BigInteger.valueOf(1000000000L);

    public static final int DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH = 40;
    public static final long DEFAULT_POLLING_FREQUENCY = 2 * 1000;
    String LIBRARY_BINARY="610130610026600b82828239805160001a60731461001957fe5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060335760003560e01c80637ae1e058146038575b600080fd5b818015604357600080fd5b50607760048036036040811015605857600080fd5b8101908080359060200190929190803590602001909291905050506091565b604051808215151515815260200191505060405180910390f35b600082600001600083815260200190815260200160002060009054906101000a900460ff161560c2576000905060f5565b600183600001600084815260200190815260200160002060006101000a81548160ff021916908315150217905550600190505b9291505056fea265627a7a72315820ed7d9db6e4f6b48c0a8f8c388d4323361bff73be7e62229ed3e7a582a0ca195764736f6c634300050d0032";
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "lib.LibraryUsingForTest-using A for B")
    public void testRegister() {
        try {
            prepare();
            TransactionReceipt libReceipt =this.deployLib(GAS_PRICE,GAS_LIMIT,LIBRARY_BINARY);
            String libAddress = libReceipt.getContractAddress();
            collector.logStepPass("libReceipt issued successfully.libAddress:" + libAddress + ", hash:" + libReceipt.getTransactionHash());
            libAddress = StringUtils.substringAfter(libAddress,"0x");
            replaceLibAddress(libAddress);
            LibraryUsingFor using = LibraryUsingFor.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("LibraryUsingFor issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = using.register(new BigInteger("12")).send();
            List<LibraryUsingFor.ResultEventResponse> eventData = using.getResultEvents(receipt);
            String data = eventData.get(0).log.getData();
            collector.assertEqual(subHexData(data), subHexData("1"), "checkout using A for B");
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


    private String subHexData(String hexStr) {
        if (StringUtils.isBlank(hexStr)) {
            throw new IllegalArgumentException("string is blank");
        }
        if (StringUtils.startsWith(hexStr, "0x")) {
            hexStr = StringUtils.substringAfter(hexStr, "0x");
        }
        byte[] addi = hexStr.getBytes();
        for (int i = 0; i < addi.length; i++) {
            if (addi[i] != 0) {
                hexStr = StringUtils.substring(hexStr, i - 1);
                break;
            }
        }
        return hexStr;
    }


    private void replaceLibAddress(String address){
        String contractBinary= LibraryUsingFor.BINARY;
        int startIndex=StringUtils.indexOf(contractBinary,"__$");
        int endIndex=StringUtils.indexOf(contractBinary,"$__");
        if(startIndex==0||endIndex==0){
            return;
        }
        String replaceStr =StringUtils.substring(contractBinary,startIndex,endIndex+3);
        LibraryUsingFor.BINARY=StringUtils.replace(contractBinary,replaceStr,address);

    }
}
