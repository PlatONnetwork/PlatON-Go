package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.FunctionEncoder;
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
public class NewTargetCreateContract extends Contract {
    private static final String BINARY = "60806040526040516101053803806101058339818101604052602081101561002657600080fd5b810190808051906020019092919050505080600081905550346001819055505060b1806100546000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c55699c146037578063aa8c217c146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b60596075565b6040518082815260200191505060405180910390f35b60005481565b6001548156fea264697066735822122065bde9461d0fa89d433499fcaced7f17b978e6ad8d972119d9c56d2f86d98e8964736f6c63430007010033";

    public static final String FUNC_AMOUNT = "amount";

    public static final String FUNC_X = "x";

    protected NewTargetCreateContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected NewTargetCreateContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<NewTargetCreateContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _a) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_a)));
        return deployRemoteCall(NewTargetCreateContract.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<NewTargetCreateContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _a) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_a)));
        return deployRemoteCall(NewTargetCreateContract.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<TransactionReceipt> amount() {
        final Function function = new Function(
                FUNC_AMOUNT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> x() {
        final Function function = new Function(
                FUNC_X, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static NewTargetCreateContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new NewTargetCreateContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static NewTargetCreateContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new NewTargetCreateContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
