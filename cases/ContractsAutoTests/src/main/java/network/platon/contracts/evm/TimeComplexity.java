package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
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
public class TimeComplexity extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610100806100206000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80637003f6c2146041578063d25f264014606c578063e65284be146097575b600080fd5b606a60048036036020811015605557600080fd5b810190808035906020019092919050505060c2565b005b609560048036036020811015608057600080fd5b810190808035906020019092919050505060c5565b005b60c06004803603602081101560ab57600080fd5b810190808035906020019092919050505060c8565b005b50565b50565b5056fea265627a7a72315820ea1b28c428172d9f5e496e520c6c383142e208a5e3549bfbbb0fbbad4fca207d64736f6c63430005110032";

    public static final String FUNC_LOGNTEST = "logNTest";

    public static final String FUNC_NSQUARETEST = "nSquareTest";

    public static final String FUNC_NTEST = "nTest";

    protected TimeComplexity(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected TimeComplexity(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> logNTest(BigInteger n) {
        final Function function = new Function(
                FUNC_LOGNTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> nSquareTest(BigInteger n) {
        final Function function = new Function(
                FUNC_NSQUARETEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> nTest(BigInteger n) {
        final Function function = new Function(
                FUNC_NTEST, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<TimeComplexity> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TimeComplexity.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<TimeComplexity> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TimeComplexity.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static TimeComplexity load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new TimeComplexity(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static TimeComplexity load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new TimeComplexity(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
