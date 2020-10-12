package network.platon.contracts.evm.v0_4_26;

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
public class AddressBalance extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060e58061001f6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063761f36d2146044575b600080fd5b348015604f57600080fd5b506082600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506098565b6040518082815260200191505060405180910390f35b60008173ffffffffffffffffffffffffffffffffffffffff163190509190505600a165627a7a723058209052cf1d85626cb1c041b97913e05a42ac5595f0c417f82df782a23a13d817f50029";

    public static final String FUNC_BALANCEOFPLATON = "balanceOfPlatON";

    protected AddressBalance(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AddressBalance(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> balanceOfPlatON(String user) {
        final Function function = new Function(FUNC_BALANCEOFPLATON, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(user)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<AddressBalance> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AddressBalance.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AddressBalance> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AddressBalance.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AddressBalance load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AddressBalance(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AddressBalance load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AddressBalance(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
