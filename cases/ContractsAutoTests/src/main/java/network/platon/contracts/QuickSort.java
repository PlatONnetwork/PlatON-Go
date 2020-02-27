package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Int256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.6-SNAPSHOT.
 */
public class QuickSort extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b5b5b5b6103b4806100216000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680631703539c1461004957806371e5ee5f146100b5575b600080fd5b341561005457600080fd5b6100b36004808035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919080359060200190919080359060200190919050506100ec565b005b34156100c057600080fd5b6100d6600480803590602001909190505061015d565b6040518082815260200191505060405180910390f35b60006100f9848484610182565b600090505b8351811015610156576000805480600101828161011b9190610337565b916000526020600020900160005b868481518110151561013757fe5b90602001906020020151909190915055505b80806001019150506100fe565b5b50505050565b60008181548110151561016c57fe5b906000526020600020900160005b915090505481565b6000818310156101c1576101978484846101c8565b90506000811415156101b2576101b1848460018403610182565b5b6101c0846001830184610182565b5b5b50505050565b60008060008086868151811015156101dc57fe5b9060200190602002015192508591508490505b808214151561030c575b808210801561021e575082878281518110151561021257fe5b90602001906020020151135b1561022e576001810390506101f9565b8082101561027157868181518110151561024457fe5b90602001906020020151878381518110151561025c57fe5b90602001906020020181815250506001820191505b5b8082108015610297575082878381518110151561028b57fe5b90602001906020020151125b156102a757600182019150610272565b808210156102ea5786828151811015156102bd57fe5b9060200190602002015187828151811015156102d557fe5b90602001906020020181815250506001810390505b8287838151811015156102f957fe5b90602001906020020181815250506101ef565b82878381518110151561031b57fe5b90602001906020020181815250508193505b5050509392505050565b81548183558181151161035e5781836000526020600020918201910161035d9190610363565b5b505050565b61038591905b80821115610381576000816000905550600101610369565b5090565b905600a165627a7a7230582068adb5a14c660b70a15ee00ff3111c50c032455eed8f1f46ce2611e44ee10eda0029";

    public static final String FUNC_SORT = "sort";

    public static final String FUNC_ARR = "arr";

    @Deprecated
    protected QuickSort(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected QuickSort(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected QuickSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected QuickSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> sort(List<BigInteger> _arr, BigInteger low, BigInteger high) {
        final Function function = new Function(
                FUNC_SORT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Int256>(
                        org.web3j.abi.Utils.typeMap(_arr, org.web3j.abi.datatypes.generated.Int256.class)), 
                new org.web3j.abi.datatypes.generated.Uint256(low), 
                new org.web3j.abi.datatypes.generated.Uint256(high)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> arr(BigInteger param0) {
        final Function function = new Function(FUNC_ARR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(QuickSort.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(QuickSort.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<QuickSort> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(QuickSort.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<QuickSort> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(QuickSort.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static QuickSort load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new QuickSort(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static QuickSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new QuickSort(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static QuickSort load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new QuickSort(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static QuickSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new QuickSort(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
