package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint8;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class InterfaceEnableStructAndenumImpl extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060d48061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80636ff6560314602d575b600080fd5b60336053565b60405180826001811115604257fe5b815260200191505060405180910390f35b6000605b6081565b6040518060200160405280600180811115607157fe5b8152509050806000015191505090565b604051806020016040528060006001811115609857fe5b8152509056fea2646970667358221220ca61b1202322b40988f11a4f16ea08cb3a07991d429c76eac33e52b153b2ded864736f6c63430007010033";

    public static final String FUNC_GETPRODUCTCONDITION = "getProductCondition";

    protected InterfaceEnableStructAndenumImpl(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InterfaceEnableStructAndenumImpl(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getProductCondition() {
        final Function function = new Function(FUNC_GETPRODUCTCONDITION, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InterfaceEnableStructAndenumImpl> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceEnableStructAndenumImpl.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InterfaceEnableStructAndenumImpl> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InterfaceEnableStructAndenumImpl.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InterfaceEnableStructAndenumImpl load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceEnableStructAndenumImpl(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InterfaceEnableStructAndenumImpl load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InterfaceEnableStructAndenumImpl(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
