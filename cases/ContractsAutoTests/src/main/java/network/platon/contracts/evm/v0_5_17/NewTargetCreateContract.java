package network.platon.contracts.evm.v0_5_17;

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
public class NewTargetCreateContract extends Contract {
    private static final String BINARY = "60806040526040516101043803806101048339818101604052602081101561002657600080fd5b810190808051906020019092919050505080600081905550346001819055505060b0806100546000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c55699c146037578063aa8c217c146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596075565b6040518082815260200191505060405180910390f35b60005481565b6001548156fea265627a7a723158202206163c67423ec8f497273d5293928ce55f70265c4296a4f30bf7b736ad554f64736f6c63430005110032";

    public static final String FUNC_AMOUNT = "amount";

    public static final String FUNC_X = "x";

    protected NewTargetCreateContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected NewTargetCreateContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<NewTargetCreateContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId, BigInteger _a) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_a)));
        return deployRemoteCall(NewTargetCreateContract.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, initialVonValue, chainId);
    }

    public static RemoteCall<NewTargetCreateContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId, BigInteger _a) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(_a)));
        return deployRemoteCall(NewTargetCreateContract.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, initialVonValue, chainId);
    }

    public RemoteCall<BigInteger> amount() {
        final Function function = new Function(FUNC_AMOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> x() {
        final Function function = new Function(FUNC_X, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static NewTargetCreateContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new NewTargetCreateContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static NewTargetCreateContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new NewTargetCreateContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
