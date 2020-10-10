package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class SameNameConstructorDefaultVisibility extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102a8806100206000396000f300608060405260043610610078576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633b016c7e1461007d5780637f14d919146100be5780638d97752a146100ff578063ac84179514610140578063ba1ae46e14610181578063ba91daeb146101ae575b600080fd5b34801561008957600080fd5b506100a8600480360381019080803590602001909291905050506101ef565b6040518082815260200191505060405180910390f35b3480156100ca57600080fd5b506100e960048036038101908080359060200190929190505050610201565b6040518082815260200191505060405180910390f35b34801561010b57600080fd5b5061012a60048036038101908080359060200190929190505050610214565b6040518082815260200191505060405180910390f35b34801561014c57600080fd5b5061016b60048036038101908080359060200190929190505050610227565b6040518082815260200191505060405180910390f35b34801561018d57600080fd5b506101ac6004803603810190808035906020019092919050505061023a565b005b3480156101ba57600080fd5b506101d960048036038101908080359060200190929190505050610244565b6040518082815260200191505060405180910390f35b60006101fa82610256565b9050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b6000816000819055506000549050919050565b8060008190555050565b600061024f82610269565b9050919050565b6000816000819055506000549050919050565b60008160008190555060005490509190505600a165627a7a72305820ac04b63152c728bfe20de3d7db9a5217f73a5ab916320f5edc9f26689daf553e0029";

    public static final String FUNC_PRIVATEVISIBILITYCHECK = "privateVisibilityCheck";

    public static final String FUNC_DEFAULTVISIBILITY = "defaultVisibility";

    public static final String FUNC_PUBLICVISIBILITY = "publicVisibility";

    public static final String FUNC_EXTERNALVISIBILITY = "externalVisibility";

    public static final String FUNC_SAMENAMECONSTRUCTORVISIBILITY = "SameNameConstructorVisibility";

    public static final String FUNC_INTERNALVISIBILITYCHECK = "internalVisibilityCheck";

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected SameNameConstructorDefaultVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> privateVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_PRIVATEVISIBILITYCHECK, 
                Arrays.<Type>asList(new Uint256(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> defaultVisibility(BigInteger param) {
        final Function function = new Function(FUNC_DEFAULTVISIBILITY, 
                Arrays.<Type>asList(new Uint256(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> publicVisibility(BigInteger param) {
        final Function function = new Function(FUNC_PUBLICVISIBILITY, 
                Arrays.<Type>asList(new Uint256(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> externalVisibility(BigInteger param) {
        final Function function = new Function(FUNC_EXTERNALVISIBILITY, 
                Arrays.<Type>asList(new Uint256(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> SameNameConstructorVisibility(BigInteger param) {
        final Function function = new Function(
                FUNC_SAMENAMECONSTRUCTORVISIBILITY, 
                Arrays.<Type>asList(new Uint256(param)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> internalVisibilityCheck(BigInteger param) {
        final Function function = new Function(FUNC_INTERNALVISIBILITYCHECK, 
                Arrays.<Type>asList(new Uint256(param)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<SameNameConstructorDefaultVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(SameNameConstructorDefaultVisibility.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SameNameConstructorDefaultVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SameNameConstructorDefaultVisibility(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
