package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
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
public class ConstructorInternalVisibility extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506040516020806100f58339810180604052810190808051906020019092919050505060078060008190555050806001819055505060a2806100536000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806335f646c0146044575b600080fd5b348015604f57600080fd5b506056606c565b6040518082815260200191505060405180910390f35b60006001549050905600a165627a7a7230582010cdc6ea834b38d98025e3fb857eceebddb56c036d6774fc7aafdebc8f0e7bf80029";

    public static final String FUNC_GETOUTI = "getOutI";

    protected ConstructorInternalVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstructorInternalVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getOutI() {
        final Function function = new Function(FUNC_GETOUTI, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<ConstructorInternalVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(ConstructorInternalVisibility.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ConstructorInternalVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _y) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_y)));
        return deployRemoteCall(ConstructorInternalVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static ConstructorInternalVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalVisibility(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstructorInternalVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalVisibility(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
