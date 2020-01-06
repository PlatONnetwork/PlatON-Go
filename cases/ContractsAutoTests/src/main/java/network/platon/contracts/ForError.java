package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
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
public class ForError extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610120806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063895e3ada146037578063ec56ae5d146057575b600080fd5b603d6077565b604051808215151515815260200191505060405180910390f35b605d609a565b604051808215151515815260200191505060405180910390f35b600080608060bd565b90508060000160009054906101000a900460ff1691505090565b60008060a360c7565b90508060000160009054906101000a900460ff1691505090565b60008090505b60c3565b60005b600190508060000160009054906101000a900460ff161560e85760ca565b9056fea265627a7a72315820ea46a9bc53c92eb0aa6f1bc459a084e0e3e75ad5b6820d6204789fac77962eae64736f6c634300050d0032";

    public static final String FUNC_GETFORCONTROLRES = "getForControlRes";

    public static final String FUNC_GETFORCONTROLRES1 = "getForControlRes1";

    @Deprecated
    protected ForError(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected ForError(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected ForError(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected ForError(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<Boolean> getForControlRes() {
        final Function function = new Function(FUNC_GETFORCONTROLRES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<Boolean> getForControlRes1() {
        final Function function = new Function(FUNC_GETFORCONTROLRES1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(ForError.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ForError> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ForError.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<ForError> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(ForError.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ForError> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ForError.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static ForError load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new ForError(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static ForError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new ForError(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static ForError load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new ForError(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static ForError load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new ForError(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
