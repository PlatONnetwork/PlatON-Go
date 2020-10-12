package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class FallBack extends Contract {
    private static final String BINARY = "6080604052600160005534801561001557600080fd5b5061012b806100256000396000f3fe6080604052348015600f57600080fd5b506004361060355760003560e01c8063d1f1548f146040578063d46300fd146048576036565b5b6064600081905550005b60466064565b005b604e60ec565b6040518082815260200191505060405180910390f35b3073ffffffffffffffffffffffffffffffffffffffff1660405180807f66756e6374696f6e4e6f744578697374282900000000000000000000000000008152506012019050600060405180830381855af49150503d806000811460e2576040519150601f19603f3d011682016040523d82523d6000602084013e60e7565b606091505b505050565b6000805490509056fea2646970667358221220df5f5c723b5d1e2c6296acc40e8a3d35043c3ab3089bb85c519ac2f277d7426064736f6c63430007010033";

    public static final String FUNC_CALLFUNCTIONNOTEXIST = "CallFunctionNotExist";

    public static final String FUNC_GETA = "getA";

    protected FallBack(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected FallBack(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> CallFunctionNotExist() {
        final Function function = new Function(
                FUNC_CALLFUNCTIONNOTEXIST, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getA() {
        final Function function = new Function(FUNC_GETA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<FallBack> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FallBack.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<FallBack> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FallBack.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static FallBack load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new FallBack(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static FallBack load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new FallBack(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
