package network.platon.contracts.evm.v0_6_12;

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
public class AssertHandle extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610170806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80638c671e0a14610067578063ad92212f14610071578063afcd320e1461007b578063cda0a5eb146100a9578063f25e0471146100f9578063f81cf6db14610103575b600080fd5b61006f61010d565b005b61007961010f565b005b6100a76004803603602081101561009157600080fd5b8101908080359060200190929190505050610111565b005b6100d8600480360360208110156100bf57600080fd5b81019080803560000b906020019092919050505061011e565b604051808260038111156100e857fe5b815260200191505060405180910390f35b610101610136565b005b61010b610138565b005b565b565b600a811061011b57fe5b50565b60008160000b600381111561012f57fe5b9050919050565b565b56fea2646970667358221220b78a6c1a3ea91bde04ec447a23646e8cc9beeeb89755e0bde412441cdf0ea88164736f6c634300060c0033";

    public static final String FUNC_BINARYMOVEMINUSEXCEPTION = "binaryMoveMinusException";

    public static final String FUNC_DIVIDENDZEROEXCEPTION = "dividendZeroException";

    public static final String FUNC_INTCHANGEEXCEPTION = "intChangeException";

    public static final String FUNC_NOOUTOFBOUNDSEXCEPTION = "noOutOfBoundsException";

    public static final String FUNC_OUTOFBOUNDSEXCEPTION = "outOfBoundsException";

    public static final String FUNC_PARAMEXCEPTION = "paramException";

    protected AssertHandle(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AssertHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> binaryMoveMinusException() {
        final Function function = new Function(
                FUNC_BINARYMOVEMINUSEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> dividendZeroException() {
        final Function function = new Function(
                FUNC_DIVIDENDZEROEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> intChangeException(BigInteger param) {
        final Function function = new Function(
                FUNC_INTCHANGEEXCEPTION, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int8(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> noOutOfBoundsException() {
        final Function function = new Function(
                FUNC_NOOUTOFBOUNDSEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> outOfBoundsException() {
        final Function function = new Function(
                FUNC_OUTOFBOUNDSEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> paramException(BigInteger param) {
        final Function function = new Function(
                FUNC_PARAMEXCEPTION, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<AssertHandle> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssertHandle.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AssertHandle> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssertHandle.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AssertHandle load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AssertHandle(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AssertHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AssertHandle(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
