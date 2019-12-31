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
public class LibraryUsingFor extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061028b806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063e81cf24c1461003b578063f207564e14610073575b600080fd5b6100716004803603604081101561005157600080fd5b8101908080359060200190929190803590602001909291905050506100a1565b005b61009f6004803603602081101561008957600080fd5b81019080803590602001909291905050506101b5565b005b6000600173__$e256e1a656e335905948bdec379cedf07e$__6324fef5c89091856040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b1580156100fd57600080fd5b505af4158015610111573d6000803e3d6000fd5b505050506040513d602081101561012757600080fd5b810190808051906020019092919050505090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114156101935760018290806001815401808255809150509060018203906000526020600020016000909192909190915055506101b0565b81600182815481106101a157fe5b90600052602060002001819055505b505050565b600073__$52c055646ffc806dbb00102514e25bd398$__637ae1e0589091836040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b15801561020f57600080fd5b505af4158015610223573d6000803e3d6000fd5b505050506040513d602081101561023957600080fd5b810190808051906020019092919050505061025357600080fd5b5056fea265627a7a72315820c6ae76b1b8dd679e7ea6e997fcf12b770da408f6d3e12d7b3ccfb507d3cbc59e64736f6c634300050d0032\r\n"
            + "\r\n"
            + "// $e256e1a656e335905948bdec379cedf07e$ -> D:/git/ContractsAutoTests/src/test/resources/contracts/9.library/LibraryUsingFor.sol:SearchLibrary\r\n"
            + "// $52c055646ffc806dbb00102514e25bd398$ -> D:/git/ContractsAutoTests/src/test/resources/contracts/9.library/LibraryUsingFor.sol:BaseLibrary";

    public static final String FUNC_REGISTER = "register";

    public static final String FUNC_REPLACE = "replace";

    @Deprecated
    protected LibraryUsingFor(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected LibraryUsingFor(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected LibraryUsingFor(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected LibraryUsingFor(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> register(BigInteger value) {
        final Function function = new Function(
                FUNC_REGISTER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(value)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> replace(BigInteger _old, BigInteger _new) {
        final Function function = new Function(
                FUNC_REPLACE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_old), 
                new org.web3j.abi.datatypes.generated.Uint256(_new)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibraryUsingFor> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibraryUsingFor.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static LibraryUsingFor load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibraryUsingFor(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static LibraryUsingFor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibraryUsingFor(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static LibraryUsingFor load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new LibraryUsingFor(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static LibraryUsingFor load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new LibraryUsingFor(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
