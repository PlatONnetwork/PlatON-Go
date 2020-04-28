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
public class TimeTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061018a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806328ed13a5146100675780633c35a0c1146100855780637fefad021461008f578063931f8bcd146100ad5780639bd1479a146100cb578063cea52e71146100e9575b600080fd5b61006f610107565b6040518082815260200191505060405180910390f35b61008d610114565b005b610097610121565b6040518082815260200191505060405180910390f35b6100b5610130565b6040518082815260200191505060405180910390f35b6100d361013a565b6040518082815260200191505060405180910390f35b6100f1610147565b6040518082815260200191505060405180910390f35b6000603c60005401905090565b6305f5e100600081905550565b600062093a8060005401905090565b6000424203905090565b6000600160005401905090565b6000610e106000540190509056fea265627a7a72315820498c543657c51c197566458ceff52b3ce46c69241edd731db606a914b60b9d4064736f6c634300050d0032";

    public static final String FUNC_THOURS = "tHours";

    public static final String FUNC_TMINUTES = "tMinutes";

    public static final String FUNC_TSECONDS = "tSeconds";

    public static final String FUNC_TWEEKS = "tWeeks";

    public static final String FUNC_TESTTIME = "testTime";

    public static final String FUNC_TESTIMEDIFF = "testimeDiff";

    @Deprecated
    protected TimeTest(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected TimeTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected TimeTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected TimeTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> tHours() {
        final Function function = new Function(FUNC_THOURS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> tMinutes() {
        final Function function = new Function(FUNC_TMINUTES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> tSeconds() {
        final Function function = new Function(FUNC_TSECONDS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> tWeeks() {
        final Function function = new Function(FUNC_TWEEKS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> testTime() {
        final Function function = new Function(
                FUNC_TESTTIME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testimeDiff() {
        final Function function = new Function(
                FUNC_TESTIMEDIFF, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<TimeTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(TimeTest.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TimeTest> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TimeTest.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<TimeTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(TimeTest.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<TimeTest> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(TimeTest.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static TimeTest load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new TimeTest(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static TimeTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new TimeTest(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static TimeTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new TimeTest(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static TimeTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new TimeTest(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
