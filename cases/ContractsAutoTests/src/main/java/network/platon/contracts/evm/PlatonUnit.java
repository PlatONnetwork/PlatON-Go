package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
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
public class PlatonUnit extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060fd8061001f6000396000f3fe60806040526004361060265760003560e01c806312065fe0146038578063b69ef8a8146060575b6030600054346088565b600081905550005b348015604357600080fd5b50604a60a3565b6040518082815260200191505060405180910390f35b348015606b57600080fd5b50607260c2565b6040518082815260200191505060405180910390f35b600080828401905083811015609957fe5b8091505092915050565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b6000548156fea265627a7a72315820a96106c7f63f399e9b55b5ae9b9a4e1e8c3df95816b4791b551a1200a74a8eed64736f6c634300050d0032";

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
