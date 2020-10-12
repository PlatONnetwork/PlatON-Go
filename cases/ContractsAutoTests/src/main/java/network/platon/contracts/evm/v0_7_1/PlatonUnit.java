package network.platon.contracts.evm.v0_7_1;

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
public class PlatonUnit extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060fc8061001f6000396000f3fe60806040526004361060295760003560e01c806312065fe014604d578063b69ef8a814607557603f565b36603f57603760005434609d565b600081905550005b348015604a57600080fd5b50005b348015605857600080fd5b50605f60b8565b6040518082815260200191505060405180910390f35b348015608057600080fd5b50608760c0565b6040518082815260200191505060405180910390f35b60008082840190508381101560ae57fe5b8091505092915050565b600047905090565b6000548156fea26469706673582212206f8dbd1ce8598dda4410e3f7a3303aa2f339935e40e7082d7502e0e387da62d264736f6c63430007010033";

    public static final String FUNC_BALANCE = "balance";

    public static final String FUNC_GETBALANCE = "getBalance";

    protected PlatonUnit(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PlatonUnit(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> balance() {
        final Function function = new Function(FUNC_BALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBalance() {
        final Function function = new Function(FUNC_GETBALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<PlatonUnit> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatonUnit.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PlatonUnit> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatonUnit.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PlatonUnit load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PlatonUnit(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PlatonUnit load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PlatonUnit(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
