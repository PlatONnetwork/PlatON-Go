package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class FallBack extends Contract {
    private static final String BINARY = "6080604052600160005534801561001557600080fd5b5061011d806100256000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063d1f1548f14603c578063d46300fd146044575b6064600081905550005b60426060565b005b604a60e8565b6040518082815260200191505060405180910390f35b3073ffffffffffffffffffffffffffffffffffffffff1660405180807f66756e6374696f6e4e6f744578697374282900000000000000000000000000008152506012019050600060405180830381855af49150503d806000811460de576040519150601f19603f3d011682016040523d82523d6000602084013e60e3565b606091505b505050565b6000805490509056fea165627a7a723058205653ceb98dd3d2ac01ba04fccce27ab37438dfb2cc5b61607376e0acfee109e20029";

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
