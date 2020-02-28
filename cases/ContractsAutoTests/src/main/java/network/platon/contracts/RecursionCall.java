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
 * <p>Generated with web3j version 0.9.1.0-SNAPSHOT.
 */
public class RecursionCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060c28061001f6000396000f3fe608060405260043610601c5760003560e01c806357e98139146021575b600080fd5b604a60048036036020811015603557600080fd5b81019080803590602001909291905050506060565b6040518082815260200191505060405180910390f35b60008160005410156083576000808154600101919050819055506081826060565b505b600054905091905056fea265627a7a72315820b0ab0185b7531690b114598e2e5d5ccb80d31015347c84abe9508dfbd4dc468f64736f6c634300050d0032";

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

    public RemoteCall<TransactionReceipt> recursionCallTest(BigInteger n, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_RECURSIONCALLTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
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
