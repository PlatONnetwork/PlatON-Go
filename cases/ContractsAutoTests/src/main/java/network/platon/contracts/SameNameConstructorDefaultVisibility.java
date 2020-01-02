package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class SameNameConstructorDefaultVisibility extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102a8806100206000396000f300608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633b016c7e1461007d5780637f14d919146100be5780638d97752a146100ff578063ac84179514610140578063ba1ae46e14610181578063ba91daeb146101ae575b600080fd5b34801561008957600080fd5b506100a8600480360381019080803590602001909291905050506101ef565b6040518082815260200191505060405180910390f35b3480156100ca57600080fd5b506100e960048036038101908080359060200190929190505050610201565b6040518082815260200191505060405180910390f35b34801561010b57600080fd5b5061012a60048036038101908080359060200190929190505050610214565b6040518082815260200191505060405180910390f35b34801561014c57600080fd5b5061016b60048036038101908080359060200190929190505050610227565b6040518082815260200191505060405180910390f35b34801561018d57600080fd5b506101ac6004803603810190808035906020019092919050505061023a565b005b3480156101ba57600080fd5b506101d960048036038101908080359060200190929190505050610244565b6040518082815260200191505060405180910390f35b60006101fa82610256565b9050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b8060008190555050565b600061024f82610269565b9050919050565b6000816000819055506000549050919050565b60008160008190555060005490509190505600a165627a7a7230582084805f11221a93547298af997d7f340b7c65d60df0f416f1ec900bb79c9b45180029";

    public static final String FUNC_PRIVATEVISIBILITYCHECK = "privateVisibilityCheck";

    public static final String FUNC_DEFAULTVISIBILITY = "defaultVisibility";

    public static final String FUNC_PUBLICVISIBILITY = "publicVisibility";

    public static final String FUNC_EXTERNALVISIBILITY = "externalVisibility";

    public static final String FUNC_SAMENAMECONSTRUCTORVISIBILITY = "SameNameConstructorVisibility";

    public static final String FUNC_INTERNALVISIBILITYCHECK = "internalVisibilityCheck";

    @Deprecated
    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> privateVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_PRIVATEVISIBILITYCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> defaultVisibility(BigInteger param) {
        final Function function = new Function(FUNC_DEFAULTVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> publicVisibility(BigInteger param) {
        final Function function = new Function(FUNC_PUBLICVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> externalVisibility(BigInteger param) {
        final Function function = new Function(FUNC_EXTERNALVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> SameNameConstructorVisibility(BigInteger param) {
        final Function function = new Function(
                FUNC_SAMENAMECONSTRUCTORVISIBILITY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> internalVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_INTERNALVISIBILITYCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
