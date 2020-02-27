package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
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
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class InsertSort extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610394806100206000396000f3fe6080604052600436106100295760003560e01c80631df339cf1461002e578063970f17bd14610145575b600080fd5b6100ee6004803603604081101561004457600080fd5b810190808035906020019064010000000081111561006157600080fd5b82018360208201111561007357600080fd5b8035906020019184602083028401116401000000008311171561009557600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291905050506101b1565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610131578082015181840152602081019050610116565b505050509050019250505060405180910390f35b34801561015157600080fd5b5061015a610295565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561019d578082015181840152602081019050610182565b505050509050019250505060405180910390f35b60606000806000600192505b84831015610272576000808785815181106101d457fe5b602002602001015190508492505b6001831015801561020857508760018403815181106101fd57fe5b602002602001015181125b1561024a5787600184038151811061021c57fe5b602002602001015188848151811061023057fe5b6020026020010181815250508280600190039350506101e2565b8088848151811061025757fe5b602002602001018181525050505082806001019350506101bd565b85600090805190602001906102889291906102ed565b5085935050505092915050565b606060008054806020026020016040519081016040528092919081815260200182805480156102e357602002820191906000526020600020905b8154815260200190600101908083116102cf575b5050505050905090565b828054828255906000526020600020908101928215610329579160200282015b8281111561032857825182559160200191906001019061030d565b5b509050610336919061033a565b5090565b61035c91905b80821115610358576000816000905550600101610340565b5090565b9056fea265627a7a72315820e69b19a17a61bcc279ca8b8b68a7e18e809ab259b8b9b6e10088604f2c5f76bf64736f6c634300050d0032";

    public static final String FUNC_OUPUTARRAYS = "OuputArrays";

    public static final String FUNC_GET_ARR = "get_arr";

    @Deprecated
    protected InsertSort(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected InsertSort(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected InsertSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected InsertSort(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> OuputArrays(List<BigInteger> arr, BigInteger n, BigInteger weiValue) {
        final Function function = new Function(
                FUNC_OUPUTARRAYS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicArray<org.web3j.abi.datatypes.generated.Int256>(
                        org.web3j.abi.Utils.typeMap(arr, org.web3j.abi.datatypes.generated.Int256.class)), 
                new org.web3j.abi.datatypes.generated.Uint256(n)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<List> get_arr() {
        final Function function = new Function(FUNC_GET_ARR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Int256>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public static RemoteCall<InsertSort> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(InsertSort.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InsertSort> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InsertSort.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<InsertSort> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(InsertSort.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<InsertSort> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(InsertSort.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static InsertSort load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new InsertSort(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static InsertSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new InsertSort(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static InsertSort load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new InsertSort(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static InsertSort load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new InsertSort(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
