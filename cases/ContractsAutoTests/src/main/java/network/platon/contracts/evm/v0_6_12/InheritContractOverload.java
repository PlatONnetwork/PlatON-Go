package network.platon.contracts.evm.v0_6_12;

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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class InheritContractOverload extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101aa806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80635873f056146100515780639450268b1461006f578063a46cf4b3146100c5578063cad0899b146100e3575b600080fd5b61005961012f565b6040518082815260200191505060405180910390f35b6100af6004803603606081101561008557600080fd5b81019080803590602001909291908035906020019092919080359060200190929190505050610144565b6040518082815260200191505060405180910390f35b6100cd610154565b6040518082815260200191505060405180910390f35b610119600480360360408110156100f957600080fd5b810190808035906020019092919080359060200190929190505050610167565b6040518082815260200191505060405180910390f35b600061013f600160026003610144565b905090565b6000818385010190509392505050565b600061016260016002610167565b905090565b600081830190509291505056fea2646970667358221220066f2a56637aabc50c01bf89b3af746fb2cf6ef402beee275d2f1c985bcbfd4964736f6c634300060c0033";

    public static final String FUNC_GETDATAA = "getDataA";

    public static final String FUNC_GETDATAB = "getDataB";

    public static final String FUNC_SUM = "sum";

    protected InheritContractOverload(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractOverload(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getDataA() {
        final Function function = new Function(FUNC_GETDATAA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getDataB() {
        final Function function = new Function(FUNC_GETDATAB, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> sum(BigInteger a, BigInteger b, BigInteger c) {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b),
                new Uint256(c)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> sum(BigInteger a, BigInteger b) {
        final Function function = new Function(FUNC_SUM, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverload.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractOverload> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractOverload.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverload(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractOverload load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractOverload(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
