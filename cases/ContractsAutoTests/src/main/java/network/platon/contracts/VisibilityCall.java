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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class VisibilityCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610301806100206000396000f3fe60806040526004361061001e5760003560e01c8063bef55ef314610023575b600080fd5b61002b610048565b604051808381526020018281526020019250505060405180910390f35b60008060006040516100599061019b565b604051809103906000f080158015610075573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1663ca77156f60016040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156100ca57600080fd5b505afa1580156100de573d6000803e3d6000fd5b505050506040513d60208110156100f457600080fd5b810190808051906020019092919050505092508073ffffffffffffffffffffffffffffffffffffffff1663b8b1feb460016040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561015957600080fd5b505afa15801561016d573d6000803e3d6000fd5b505050506040513d602081101561018357600080fd5b81019080805190602001909291905050509150509091565b610124806101a98339019056fe608060405234801561001057600080fd5b50610104806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063b8b1feb4146037578063ca77156f146076575b600080fd5b606060048036036020811015604b57600080fd5b810190808035906020019092919050505060b5565b6040518082815260200191505060405180910390f35b609f60048036036020811015608a57600080fd5b810190808035906020019092919050505060c2565b6040518082815260200191505060405180910390f35b6000600382019050919050565b600060028201905091905056fea265627a7a72315820aa0ee0a873bcf0a3462dc2773922ca5dd9ef169da506b0dddcf367135077c85264736f6c634300050d0032a265627a7a72315820497d9fee5c55ddd0c0990e831b91af5e1be95a3af35061e95ae557164f0cfdae64736f6c634300050d0032";

    public static final String FUNC_READDATA = "readData";

    protected VisibilityCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected VisibilityCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> readData(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_READDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<VisibilityCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(VisibilityCall.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<VisibilityCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(VisibilityCall.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static VisibilityCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new VisibilityCall(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static VisibilityCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new VisibilityCall(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
