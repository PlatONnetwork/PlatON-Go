package network.platon.contracts.evm;

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
 * <p>Generated with web3j version 0.13.2.0.
 */
public class InterfaceEnableStructAndenumImpl extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060d68061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80636ff6560314602d575b600080fd5b60336056565b60405180826001811115604257fe5b60ff16815260200191505060405180910390f35b6000605e6084565b6040518060200160405280600180811115607457fe5b8152509050806000015191505090565b604051806020016040528060006001811115609b57fe5b8152509056fea265627a7a72315820c18d4a1cd2d227a590049732070b80625fc8bd4e8b4942eda8fc6ffb83ccface64736f6c63430005110032";

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
