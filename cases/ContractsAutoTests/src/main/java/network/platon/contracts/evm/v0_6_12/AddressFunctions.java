package network.platon.contracts.evm.v0_6_12;

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
public class AddressFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061025a806100206000396000f3fe60806040526004361061003f5760003560e01c80631a695230146100445780633e58c58c14610088578063ecbde5e6146100e2578063f8b2cb4f1461010d575b600080fd5b6100866004803603602081101561005a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610172565b005b6100ca6004803603602081101561009e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506101bc565b60405180821515815260200191505060405180910390f35b3480156100ee57600080fd5b506100f76101fb565b6040518082815260200191505060405180910390f35b34801561011957600080fd5b5061015c6004803603602081101561013057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610203565b6040518082815260200191505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f193505050501580156101b8573d6000803e3d6000fd5b5050565b60008173ffffffffffffffffffffffffffffffffffffffff166108fc60019081150290604051600060405180830381858888f193505050509050919050565b600047905090565b60008173ffffffffffffffffffffffffffffffffffffffff1631905091905056fea264697066735822122088c83159280a4b5fc51f49eb91395a0e9a2d3a3222122abb7adf1033efb041d764736f6c634300060c0033";

    public static final String FUNC_GETBALANCE = "getBalance";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_SEND = "send";

    public static final String FUNC_TRANSFER = "transfer";

    protected AddressFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AddressFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getBalance(String addr) {
        final Function function = new Function(FUNC_GETBALANCE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBalanceOf() {
        final Function function = new Function(FUNC_GETBALANCEOF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> send(String addr, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_SEND, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> transfer(String addr, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TRANSFER, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AddressFunctions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AddressFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AddressFunctions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AddressFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AddressFunctions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AddressFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AddressFunctions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
