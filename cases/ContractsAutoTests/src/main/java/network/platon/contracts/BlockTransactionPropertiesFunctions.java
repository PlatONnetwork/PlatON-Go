package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes4;
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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class BlockTransactionPropertiesFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061053c806100206000396000f3fe6080604052600436106100c25760003560e01c8063796b89b91161007f578063bbe4fd5011610059578063bbe4fd50146102fa578063d12d910214610325578063df1f29ee1461038e578063e9413d38146103e5576100c2565b8063796b89b914610279578063a16963b3146102a4578063ab70fd69146102cf576100c2565b806312e05dd1146100c757806320965255146100f25780632df8e949146101105780633bc5de301461016757806342cbb15c146101f75780635e01eb5a14610222575b600080fd5b3480156100d357600080fd5b506100dc610434565b6040518082815260200191505060405180910390f35b6100fa61043c565b6040518082815260200191505060405180910390f35b34801561011c57600080fd5b50610125610444565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561017357600080fd5b5061017c61044c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101bc5780820151818401526020810190506101a1565b50505050905090810190601f1680156101e95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561020357600080fd5b5061020c610499565b6040518082815260200191505060405180910390f35b34801561022e57600080fd5b506102376104a1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561028557600080fd5b5061028e6104a9565b6040518082815260200191505060405180910390f35b3480156102b057600080fd5b506102b96104b1565b6040518082815260200191505060405180910390f35b3480156102db57600080fd5b506102e46104b9565b6040518082815260200191505060405180910390f35b34801561030657600080fd5b5061030f6104c1565b6040518082815260200191505060405180910390f35b34801561033157600080fd5b5061033a6104c9565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561039a57600080fd5b506103a36104f4565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103f157600080fd5b5061041e6004803603602081101561040857600080fd5b81019080803590602001909291905050506104fc565b6040518082815260200191505060405180910390f35b600044905090565b600034905090565b600041905090565b60606000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050905090565b600043905090565b600033905090565b600042905090565b600045905090565b60003a905090565b600042905090565b600080357fffffffff0000000000000000000000000000000000000000000000000000000016905090565b600032905090565b60008140905091905056fea265627a7a72315820bd6a8ef314b48721498d1ef9e494201a3ba6ee6a032643b5e3f9ff308fbec57c64736f6c634300050d0032";

    public static final String FUNC_GETBLOCKCOINBASE = "getBlockCoinbase";

    public static final String FUNC_GETBLOCKDIFFICULTY = "getBlockDifficulty";

    public static final String FUNC_GETBLOCKNUMBER = "getBlockNumber";

    public static final String FUNC_GETBLOCKTIMESTAMP = "getBlockTimestamp";

    public static final String FUNC_GETBLOCKHASH = "getBlockhash";

    public static final String FUNC_GETDATA = "getData";

    public static final String FUNC_GETGASLIMIT = "getGaslimit";

    public static final String FUNC_GETGASPRICE = "getGasprice";

    public static final String FUNC_GETNOW = "getNow";

    public static final String FUNC_GETORIGIN = "getOrigin";

    public static final String FUNC_GETSENDER = "getSender";

    public static final String FUNC_GETSIG = "getSig";

    public static final String FUNC_GETVALUE = "getValue";

    @Deprecated
    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> getBlockCoinbase() {
        final Function function = new Function(FUNC_GETBLOCKCOINBASE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getBlockDifficulty() {
        final Function function = new Function(FUNC_GETBLOCKDIFFICULTY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBlockNumber() {
        final Function function = new Function(FUNC_GETBLOCKNUMBER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getBlockTimestamp() {
        final Function function = new Function(FUNC_GETBLOCKTIMESTAMP, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> getBlockhash(BigInteger blockNumber) {
        final Function function = new Function(
                FUNC_GETBLOCKHASH, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(blockNumber)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<byte[]> getData() {
        final Function function = new Function(FUNC_GETDATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getGaslimit() {
        final Function function = new Function(FUNC_GETGASLIMIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getGasprice() {
        final Function function = new Function(FUNC_GETGASPRICE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getNow() {
        final Function function = new Function(FUNC_GETNOW, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getOrigin() {
        final Function function = new Function(FUNC_GETORIGIN, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getSender() {
        final Function function = new Function(FUNC_GETSENDER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> getSig() {
        final Function function = new Function(FUNC_GETSIG, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes4>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<TransactionReceipt> getValue(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_GETVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
