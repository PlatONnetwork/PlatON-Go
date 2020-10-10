package network.platon.contracts.evm.v0_6_12;

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
public class GetLibraryAddress extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610126806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80636e7c1504146037578063750c193514603f575b600080fd5b603d6071565b005b604560c7565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b73__$e85707310368214d280ebf362fa781bf90$__6000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509056fea26469706673582212206b57b8a8ec70ae8eae87d4dc06b4421c4b4a77a8162fed34aea7a73fde0e689864736f6c634300060c0033\n"
            + "\n"
            + "// $e85707310368214d280ebf362fa781bf90$ -> /home/platon/.jenkins/workspace/contracts_test_alaya/cases/ContractsAutoTests/src/test/resources/contracts/evm/0.6.12/2.version_compatible/0_5_13/8-address_LibraryName/UserLibrary.sol:UserLibrary";

    public static final String FUNC_GETUSERLIBADDRESS = "getUserLibAddress";

    public static final String FUNC_SETUSERLIBADDRESS = "setUserLibAddress";

    protected GetLibraryAddress(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected GetLibraryAddress(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getUserLibAddress() {
        final Function function = new Function(
                FUNC_GETUSERLIBADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setUserLibAddress() {
        final Function function = new Function(
                FUNC_SETUSERLIBADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<GetLibraryAddress> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(GetLibraryAddress.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<GetLibraryAddress> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(GetLibraryAddress.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static GetLibraryAddress load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new GetLibraryAddress(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static GetLibraryAddress load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new GetLibraryAddress(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
