package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
import com.alaya.abi.solidity.datatypes.Function;
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
public class ChainFunction extends Contract {
    private static final String BINARY = "6080604052336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550346001819055506000600260006101000a81548160ff02191690831515021790555061019c806100756000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80637eed92c01461003b5780639f9232f4146100ab575b600080fd5b6100696004803603602081101561005157600080fd5b81019080803515159060200190929190505050610125565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e3600480360360408110156100c157600080fd5b8101908080351515906020019092919080359060200190929190505050610140565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000600115158215151461013857600080fd5b339050919050565b6000600115158315151461015057fe5b600982101561015e57600080fd5b3390509291505056fea265627a7a723158203b510c6eefa5d8925b48c720eac4b221709e391e812b520c54a36f0f03e4c29164736f6c63430005110032";

    public static final String FUNC_DECEASED = "deceased";

    public static final String FUNC_DECEASEDWITHMODIFY = "deceasedWithModify";

    protected ChainFunction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ChainFunction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ChainFunction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId) {
        return deployRemoteCall(ChainFunction.class, web3j, credentials, contractGasProvider, BINARY, "", initialVonValue, chainId);
    }

    public static RemoteCall<ChainFunction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, BigInteger initialVonValue, Long chainId) {
        return deployRemoteCall(ChainFunction.class, web3j, transactionManager, contractGasProvider, BINARY, "", initialVonValue, chainId);
    }

    public RemoteCall<String> deceased(Boolean isDeceased, BigInteger less9) {
        final Function function = new Function(FUNC_DECEASED, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Bool(isDeceased), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(less9)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> deceasedWithModify(Boolean _isDeceased) {
        final Function function = new Function(FUNC_DECEASEDWITHMODIFY, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Bool(_isDeceased)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static ChainFunction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ChainFunction(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ChainFunction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ChainFunction(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
