package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.datatypes.Type;
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
public class BaseDefault extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b50604051602080607b83398101806040528101908080519060200190929190505050806000819055505060358060466000396000f3006080604052600080fd00a165627a7a723058201af6de3fe1d29388fe564ca3744b8f96104bf88fbb9c6fb0acce68e208ec99970029";

    protected BaseDefault(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BaseDefault(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<BaseDefault> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _x) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_x)));
        return deployRemoteCall(BaseDefault.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<BaseDefault> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _x) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_x)));
        return deployRemoteCall(BaseDefault.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static BaseDefault load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BaseDefault(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BaseDefault load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BaseDefault(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
