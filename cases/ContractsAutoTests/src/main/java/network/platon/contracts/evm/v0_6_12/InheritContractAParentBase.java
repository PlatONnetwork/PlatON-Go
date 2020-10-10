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
public class InheritContractAParentBase extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060fd8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80631416d3471460375780638da5cb5b146069575b600080fd5b603d609b565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b606f60a3565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600033905090565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea26469706673582212201dbe11df42d84a176054a0003093a9436dc6722d1e800a12acac5e3be6bbea7264736f6c634300060c0033";

    public static final String FUNC_GETADDRESSA = "getAddressA";

    public static final String FUNC_OWNER = "owner";

    protected InheritContractAParentBase(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractAParentBase(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getAddressA() {
        final Function function = new Function(
                FUNC_GETADDRESSA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> owner() {
        final Function function = new Function(
                FUNC_OWNER, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritContractAParentBase> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractAParentBase.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractAParentBase> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractAParentBase.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractAParentBase load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractAParentBase(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractAParentBase load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractAParentBase(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
