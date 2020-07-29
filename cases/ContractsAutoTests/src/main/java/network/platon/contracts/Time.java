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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class Time extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061018a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806328ed13a5146100675780633c35a0c1146100855780637fefad021461008f578063931f8bcd146100ad5780639bd1479a146100cb578063cea52e71146100e9575b600080fd5b61006f610107565b6040518082815260200191505060405180910390f35b61008d610114565b005b610097610121565b6040518082815260200191505060405180910390f35b6100b5610130565b6040518082815260200191505060405180910390f35b6100d361013a565b6040518082815260200191505060405180910390f35b6100f1610147565b6040518082815260200191505060405180910390f35b6000603c60005401905090565b6305f5e100600081905550565b600062093a8060005401905090565b6000424203905090565b6000600160005401905090565b6000610e106000540190509056fea265627a7a7231582007fd97e3e5341b7aafed64bad6d6f917c84f308f6b0c0f6d65952951ce1cb1ec64736f6c634300050d0032";

    public static final String FUNC_THOURS = "tHours";

    public static final String FUNC_TMINUTES = "tMinutes";

    public static final String FUNC_TSECONDS = "tSeconds";

    public static final String FUNC_TWEEKS = "tWeeks";

    public static final String FUNC_TESTTIME = "testTime";

    public static final String FUNC_TESTIMEDIFF = "testimeDiff";

    protected Time(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Time(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
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

    public RemoteCall<BigInteger> testimeDiff() {
        final Function function = new Function(FUNC_TESTIMEDIFF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<Time> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Time.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Time> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Time.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Time load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Time(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Time load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Time(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
