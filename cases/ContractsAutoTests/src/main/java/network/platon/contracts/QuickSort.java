package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.FunctionEncoder;
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
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class QuickSort extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b60405161042a38038061042a833981016040528080518201919050505b60008090505b815181101561008a576000805480600101828161004f9190610092565b916000526020600020900160005b848481518110151561006b57fe5b90602001906020020151909190915055505b8080600101915050610032565b5b50506100e3565b8154818355818115116100b9578183600052602060002091820191016100b891906100be565b5b505050565b6100e091905b808211156100dc5760008160009055506001016100c4565b5090565b90565b610338806100f26000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806371e5ee5f146100545780637b395ec21461008b5780639ae8886a146100b7575b600080fd5b341561005f57600080fd5b61007560048080359060200190919050506100e0565b6040518082815260200191505060405180910390f35b341561009657600080fd5b6100b56004808035906020019091908035906020019091905050610105565b005b34156100c257600080fd5b6100ca610114565b6040518082815260200191505060405180910390f35b6000818154811015156100ef57fe5b906000526020600020900160005b915090505481565b61010f828261011a565b5b5050565b60015481565b60008183101561015d5761012e8383610163565b905060008114151561014857610147836001830361011a565b5b610155600182018361011a565b806001819055505b5b505050565b60008060008060008681548110151561017857fe5b906000526020600020900160005b505492508591508490505b80821415156102dc575b80821080156101c75750826000828154811015156101b557fe5b906000526020600020900160005b5054115b156101d75760018103905061019b565b80821015610227576000818154811015156101ee57fe5b906000526020600020900160005b505460008381548110151561020d57fe5b906000526020600020900160005b50819055506001820191505b5b808210801561025457508260008381548110151561024257fe5b906000526020600020900160005b5054105b1561026457600182019150610228565b808210156102b45760008281548110151561027b57fe5b906000526020600020900160005b505460008281548110151561029a57fe5b906000526020600020900160005b50819055506001810390505b826000838154811015156102c457fe5b906000526020600020900160005b5081905550610191565b826000838154811015156102ec57fe5b906000526020600020900160005b50819055508193505b505050929150505600a165627a7a7230582061979b9a9f012d9562afe7e13531a1a57f9fc4cac4f883085ff0c1ca7523a1c60029";

    public static final String FUNC_ARR = "arr";

    public static final String FUNC_SORT = "sort";

    public static final String FUNC_P = "p";

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

    public RemoteCall<BigInteger> arr(BigInteger param0) {
        final Function function = new Function(FUNC_ARR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> sort(BigInteger low, BigInteger high) {
        final Function function = new Function(
                FUNC_SORT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(low), 
                new org.web3j.abi.datatypes.generated.Uint256(high)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> p() {
        final Function function = new Function(FUNC_P, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, List<BigInteger> _arr) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(_arr, org.web3j.abi.datatypes.generated.Uint256.class))));
        return deployRemoteCall(QuickSort.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor);
    }

    public static RemoteCall<QuickSort> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, List<BigInteger> _arr) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(_arr, org.web3j.abi.datatypes.generated.Uint256.class))));
        return deployRemoteCall(QuickSort.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<QuickSort> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit, List<BigInteger> _arr) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(_arr, org.web3j.abi.datatypes.generated.Uint256.class))));
        return deployRemoteCall(QuickSort.class, web3j, credentials, gasPrice, gasLimit, BINARY, encodedConstructor);
    }

    @Deprecated
    public static RemoteCall<QuickSort> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit, List<BigInteger> _arr) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Uint256>(
                        org.web3j.abi.Utils.typeMap(_arr, org.web3j.abi.datatypes.generated.Uint256.class))));
        return deployRemoteCall(QuickSort.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, encodedConstructor);
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
