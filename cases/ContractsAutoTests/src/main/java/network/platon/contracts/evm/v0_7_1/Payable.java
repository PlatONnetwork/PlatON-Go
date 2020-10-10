package network.platon.contracts.evm.v0_7_1;

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
public class Payable extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610178806100206000396000f3fe6080604052600436106100295760003560e01c80631a6952301461002e578063c84aae1714610072575b600080fd5b6100706004803603602081101561004457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506100d7565b005b34801561007e57600080fd5b506100c16004803603602081101561009557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610121565b6040518082815260200191505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f1935050505015801561011d573d6000803e3d6000fd5b5050565b60008173ffffffffffffffffffffffffffffffffffffffff1631905091905056fea26469706673582212203184d163c85f10004fbbb118c4c0dd1179332dcf7b5ffc68327b5eb897578c9064736f6c63430007010033";

    public static final String FUNC_GETBALANCES = "getBalances";

    public static final String FUNC_TRANSFER = "transfer";

    protected Payable(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Payable(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getBalances(String addr) {
        final Function function = new Function(
                FUNC_GETBALANCES, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> transfer(String addr) {
        final Function function = new Function(
                FUNC_TRANSFER, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Payable> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Payable.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Payable> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Payable.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Payable load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Payable(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Payable load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Payable(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
