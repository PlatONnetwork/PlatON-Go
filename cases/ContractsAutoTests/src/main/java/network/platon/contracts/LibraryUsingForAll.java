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
public class LibraryUsingForAll extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101b0806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e81cf24c14610030575b600080fd5b6100666004803603604081101561004657600080fd5b810190808035906020019092919080359060200190929190505050610068565b005b60008073f42a141fa20381972775b6dd985791f477329f196324fef5c89091856040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b1580156100c357600080fd5b505af41580156100d7573d6000803e3d6000fd5b505050506040513d60208110156100ed57600080fd5b810190808051906020019092919050505090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811415610159576000829080600181540180825580915050906001820390600052602060002001600090919290919091505550610176565b816000828154811061016757fe5b90600052602060002001819055505b50505056fea265627a7a72315820dc9be17c9eca5d102b94f743f9e4f6542c47adf4b5b8f9f4c25c0756433a68aa64736f6c634300050d0032";

    public static final String FUNC_REPLACE = "replace";

    protected LibraryUsingForAll(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected LibraryUsingForAll(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> replace(BigInteger _old, BigInteger _new) {
        final Function function = new Function(
                FUNC_REPLACE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_old), 
                new org.web3j.abi.datatypes.generated.Uint256(_new)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<LibraryUsingForAll> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(LibraryUsingForAll.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new LibraryUsingForAll(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static LibraryUsingForAll load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new LibraryUsingForAll(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
