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
public class InheritContractBParentBase extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b50609d8061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80636c3364ea14602d575b600080fd5b6033605f565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60003390509056fea2646970667358221220fd2d92df976a9837bda8c7396597a8c265ece7ac98864f1bcecf93a23970ccff64736f6c634300060c0033";

    public static final String FUNC_GETADDRESSB = "getAddressB";

    protected InheritContractBParentBase(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractBParentBase(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getAddressB() {
        final Function function = new Function(
                FUNC_GETADDRESSB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritContractBParentBase> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractBParentBase.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractBParentBase> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractBParentBase.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractBParentBase load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBParentBase(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractBParentBase load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBParentBase(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
