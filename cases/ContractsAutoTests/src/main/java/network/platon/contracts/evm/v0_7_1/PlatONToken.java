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
public class PlatONToken extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061010a806100206000396000f3fe6080604052348015600f57600080fd5b506004361060465760003560e01c8063249bb37314604b578063c2412676146067578063c951fdf614606f578063d87698ae14608b575b600080fd5b605160a7565b6040518082815260200191505060405180910390f35b606d60b9565b005b607560c3565b6040518082815260200191505060405180910390f35b609160ce565b6040518082815260200191505060405180910390f35b6000670de0b6b3a76400008101905090565b6001600081905550565b600060018101905090565b6000548156fea26469706673582212201c65c83cd890ad5a58ce0fbac81c686190f6c1ef2b53863019a10d22761c16a664736f6c63430007010033";

    public static final String FUNC_PLAT = "Plat";

    public static final String FUNC_PVON = "Pvon";

    public static final String FUNC_TOKEN = "Token";

    public static final String FUNC_PLATONTOKEN = "platontoken";

    protected PlatONToken(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PlatONToken(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> Plat() {
        final Function function = new Function(FUNC_PLAT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> Pvon() {
        final Function function = new Function(FUNC_PVON, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> Token() {
        final Function function = new Function(
                FUNC_TOKEN, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> platontoken() {
        final Function function = new Function(FUNC_PLATONTOKEN, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatONToken.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatONToken.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PlatONToken(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PlatONToken(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
