package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class InsertSort extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061023c806100206000396000f3fe60806040526004361061001e5760003560e01c80631df339cf14610023575b600080fd5b6100e36004803603604081101561003957600080fd5b810190808035906020019064010000000081111561005657600080fd5b82018360208201111561006857600080fd5b8035906020019184602083028401116401000000008311171561008a57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019092919050505061013a565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561012657808201518184015260208101905061010b565b505050509050019250505060405180910390f35b60606000806000600192505b848310156101fb5760008087858151811061015d57fe5b602002602001015190508492505b60018310158015610191575087600184038151811061018657fe5b602002602001015181125b156101d3578760018403815181106101a557fe5b60200260200101518884815181106101b957fe5b60200260200101818152505082806001900393505061016b565b808884815181106101e057fe5b60200260200101818152505050508280600101935050610146565b8593505050509291505056fea265627a7a72315820860e628e451dfe6fb3fd7990c41dddb1607006cc9b42f59e5c8554bbec43348664736f6c634300050d0032";

    public static final String FUNC_OUPUTARRAYS = "OuputArrays";

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
