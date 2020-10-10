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
import java.math.BigInteger;
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
public class LibraryUsingForAll extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101ae806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e81cf24c14610030575b600080fd5b6100666004803603604081101561004657600080fd5b810190808035906020019092919080359060200190929190505050610068565b005b60008073__$8cb58328c2bc207af2d98cd9c31b792a46$__6324fef5c89091856040518363ffffffff1660e01b8152600401808381526020018281526020019250505060206040518083038186803b1580156100c357600080fd5b505af41580156100d7573d6000803e3d6000fd5b505050506040513d60208110156100ed57600080fd5b810190808051906020019092919050505090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811415610156576000829080600181540180825580915050600190039060005260206000200160009091909190915055610173565b816000828154811061016457fe5b90600052602060002001819055505b50505056fea26469706673582212202dd95455324b1df49fa4b17d7d5b2360e9cb0be066d8152f1bf752ccabe80e1a64736f6c63430007010033\n"
            + "\n"
            + "// $8cb58328c2bc207af2d98cd9c31b792a46$ -> /home/platon/.jenkins/workspace/contracts_test_alaya/cases/ContractsAutoTests/src/test/resources/contracts/evm/0.7.1/9.library/LibraryUserForAll.sol:SearchLibrary";

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
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_old), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(_new)), 
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
