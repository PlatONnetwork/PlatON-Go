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
public class CallTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101f7806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063de29278914610030575b600080fd5b61003861004e565b6040518082815260200191505060405180910390f35b60008060405161005d90610103565b604051809103906000f080158015610079573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1663569c5f6d6040518163ffffffff1660e01b815260040160206040518083038186803b1580156100c257600080fd5b505afa1580156100d6573d6000803e3d6000fd5b505050506040513d60208110156100ec57600080fd5b810190808051906020019092919050505091505090565b60b3806101108339019056fe6080604052348015600f57600080fd5b5060958061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063569c5f6d14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b60008060019050600060029050808201925050509056fea265627a7a723158208a3736efc139b4aa2f80754ab8c26633b235325396967b868409ab191afd753264736f6c634300050d0032a265627a7a723158205958500930784d8b1995da651638ef77e56b41175b6da8f71769615ac22091e164736f6c634300050d0032";

    public static final String FUNC_GETRESULT = "getResult";

    @Deprecated
    protected CallTest(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected CallTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected CallTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected CallTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> getResult() {
        final Function function = new Function(
                FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<CallTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(CallTest.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<CallTest> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(CallTest.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<CallTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(CallTest.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<CallTest> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(CallTest.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static CallTest load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new CallTest(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static CallTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new CallTest(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static CallTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new CallTest(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static CallTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new CallTest(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
