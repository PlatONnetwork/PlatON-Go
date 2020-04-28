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
public class VisibilityCallTest extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610305806100206000396000f3fe60806040526004361061001e5760003560e01c8063bef55ef314610023575b600080fd5b61002b610048565b604051808381526020018281526020019250505060405180910390f35b60008060006040516100599061019f565b604051809103906000f080158015610075573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1663ca77156f60016040518263ffffffff1660e01b815260040180828152602001915050602060405180830381600087803b1580156100cc57600080fd5b505af11580156100e0573d6000803e3d6000fd5b505050506040513d60208110156100f657600080fd5b810190808051906020019092919050505092508073ffffffffffffffffffffffffffffffffffffffff1663b8b1feb460016040518263ffffffff1660e01b815260040180828152602001915050602060405180830381600087803b15801561015d57600080fd5b505af1158015610171573d6000803e3d6000fd5b505050506040513d602081101561018757600080fd5b81019080805190602001909291905050509150509091565b610124806101ad8339019056fe608060405234801561001057600080fd5b50610104806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063b8b1feb4146037578063ca77156f146076575b600080fd5b606060048036036020811015604b57600080fd5b810190808035906020019092919050505060b5565b6040518082815260200191505060405180910390f35b609f60048036036020811015608a57600080fd5b810190808035906020019092919050505060c2565b6040518082815260200191505060405180910390f35b6000600382019050919050565b600060028201905091905056fea265627a7a72315820684bcd0a99ef8f6af1e6a109c2b23359af3bbde5cb3dd01257ed0caa2feb8b8464736f6c634300050d0032a265627a7a72315820d65524e612df7bfe1eea3981df70c18b9fdc449c49b17be54b83b9f2d4e3284c64736f6c634300050d0032";

    public static final String FUNC_READDATA = "readData";

    @Deprecated
    protected VisibilityCallTest(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected VisibilityCallTest(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected VisibilityCallTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected VisibilityCallTest(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> readData(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_READDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public static RemoteCall<VisibilityCallTest> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(VisibilityCallTest.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<VisibilityCallTest> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(VisibilityCallTest.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<VisibilityCallTest> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(VisibilityCallTest.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<VisibilityCallTest> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(VisibilityCallTest.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static VisibilityCallTest load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new VisibilityCallTest(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static VisibilityCallTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new VisibilityCallTest(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static VisibilityCallTest load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new VisibilityCallTest(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static VisibilityCallTest load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new VisibilityCallTest(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
