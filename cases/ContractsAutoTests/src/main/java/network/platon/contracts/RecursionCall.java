package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class RecursionCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060b08061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806357e9813914602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060005410156078576000808154600101919050819055506077816058565b5b5056fea265627a7a72315820866d428f796d6d211292654fa6d470f86cfe14f15462c43b12b5e08885ca1dc964736f6c634300050d0032";

    public static final String FUNC_RECURSIONCALLTEST = "recursionCallTest";

    @Deprecated
    protected RecursionCall(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected RecursionCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected RecursionCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected RecursionCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> recursionCallTest(BigInteger n) {
        final Function function = new Function(
                FUNC_RECURSIONCALLTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<RecursionCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(RecursionCall.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RecursionCall> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RecursionCall.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<RecursionCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(RecursionCall.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<RecursionCall> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(RecursionCall.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RecursionCall load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new RecursionCall(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static RecursionCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new RecursionCall(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static RecursionCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new RecursionCall(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static RecursionCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new RecursionCall(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
