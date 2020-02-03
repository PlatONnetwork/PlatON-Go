package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
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

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class LibraryUsingForAll extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101a3806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e81cf24c14610030575b600080fd5b6100666004803603604081101561004657600080fd5b810190808035906020019092919080359060200190929190505050610068565b005b600061007e8360006100fb90919063ffffffff16565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114156100d95760008290806001815401808255809150509060018203906000526020600020016000909192909190915055506100f6565b81600082815481106100e757fe5b90600052602060002001819055505b505050565b600080600090505b8380549050811015610143578284828154811061011c57fe5b906000526020600020015414156101365780915050610168565b8080600101915050610103565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90505b9291505056fea265627a7a723158209e121e8c2679b9278ccac8237eee7517304f07af7867770be0ddc9efb0e5c88d64736f6c634300050d0032";

    public static final String FUNC_REPLACE = "replace";

    @Deprecated
    protected LibraryUsingForAll(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected LibraryUsingForAll(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected LibraryUsingForAll(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected LibraryUsingForAll(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> replace(BigInteger _old, BigInteger _new) {
        final Function function = new Function(
                FUNC_REPLACE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_old), 
                new org.web3j.abi.datatypes.generated.Uint256(_new)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibraryUsingForAll(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibraryUsingForAll(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new LibraryUsingForAll(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new LibraryUsingForAll(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
