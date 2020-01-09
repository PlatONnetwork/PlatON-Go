package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

import java.math.BigInteger;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class SumAssembly extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600060019080600181540180825580915050906001820390600052602060002001600090919290919091505550600060029080600181540180825580915050906001820390600052602060002001600090919290919091505550600060039080600181540180825580915050906001820390600052602060002001600090919290919091505550600060049080600181540180825580915050906001820390600052602060002001600090919290919091505550600060059080600181540180825580915050906001820390600052602060002001600090919290919091505550610153806101016000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063853255cc14610030575b600080fd5b61003861004e565b6040518082815260200191505060405180910390f35b600073e70038edb43b693c623bc9d637e75ed38d4358346387fbcc7760006040518263ffffffff1660e01b8152600401808060200182810382528381815481526020019150805480156100c057602002820191906000526020600020905b8154815260200190600101908083116100ac575b50509250505060206040518083038186803b1580156100de57600080fd5b505af41580156100f2573d6000803e3d6000fd5b505050506040513d602081101561010857600080fd5b810190808051906020019092919050505090509056fea265627a7a7231582064b99699beb46ec1ace824714f7b2c8b2d1ea0704a2b6596bd55e0239caf4c3d64736f6c634300050d0032";

    public static final String FUNC_SUM = "sum";

    @Deprecated
    protected SumAssembly(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SumAssembly(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SumAssembly(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SumAssembly(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<SumAssembly> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(SumAssembly.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    public static RemoteCall<SumAssembly> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(SumAssembly.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SumAssembly> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SumAssembly.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SumAssembly> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SumAssembly.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    public RemoteCall<BigInteger> sum() {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    @Deprecated
    public static SumAssembly load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SumAssembly(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SumAssembly load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SumAssembly(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SumAssembly(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SumAssembly load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SumAssembly(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
