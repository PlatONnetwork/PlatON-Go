package network.platon.contracts.evm;

import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class GetLibraryAddress extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610143806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636e7c15041461003b578063750c193514610045575b600080fd5b61004361008f565b005b61004d6100e5565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b73154f27f44ecd53bc3fc73234e13d91a7e82004796000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690509056fea265627a7a723158207372afa32ef7a7ab8b4b8c71d7c51f7714832555b04e875d32c328c3eb99cca664736f6c634300050d0032";

    public static final String FUNC_GETUSERLIBADDRESS = "getUserLibAddress";

    public static final String FUNC_SETUSERLIBADDRESS = "setUserLibAddress";

    protected GetLibraryAddress(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected GetLibraryAddress(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<String> getUserLibAddress() {
        final Function function = new Function(FUNC_GETUSERLIBADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
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
