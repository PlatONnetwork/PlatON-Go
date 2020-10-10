package network.platon.contracts.evm.v0_7_1;

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
public class Time extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061018b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806328ed13a5146100675780633c35a0c1146100855780637fefad021461008f578063931f8bcd146100ad5780639bd1479a146100cb578063cea52e71146100e9575b600080fd5b61006f610107565b6040518082815260200191505060405180910390f35b61008d610114565b005b610097610121565b6040518082815260200191505060405180910390f35b6100b5610130565b6040518082815260200191505060405180910390f35b6100d361013a565b6040518082815260200191505060405180910390f35b6100f1610147565b6040518082815260200191505060405180910390f35b6000603c60005401905090565b6305f5e100600081905550565b600062093a8060005401905090565b6000424203905090565b6000600160005401905090565b6000610e106000540190509056fea264697066735822122054631a4db972d76d10dbdc1435c822cc25b21e49bfcd5ca65fb48031b0fc621d64736f6c63430007010033";

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

    public RemoteCall<TransactionReceipt> tHours() {
        final Function function = new Function(
                FUNC_THOURS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> tMinutes() {
        final Function function = new Function(
                FUNC_TMINUTES, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> tSeconds() {
        final Function function = new Function(
                FUNC_TSECONDS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> tWeeks() {
        final Function function = new Function(
                FUNC_TWEEKS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
