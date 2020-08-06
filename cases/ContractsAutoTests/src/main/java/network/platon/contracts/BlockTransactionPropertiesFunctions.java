package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class BlockTransactionPropertiesFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061058a806100206000396000f3fe6080604052600436106100dd5760003560e01c8063a16963b31161007f578063d12d910211610059578063d12d910214610340578063df1f29ee146103a9578063e9413d3814610400578063edb4b8651461044f576100dd565b8063a16963b3146102bf578063ab70fd69146102ea578063bbe4fd5014610315576100dd565b80633bc5de30116100bb5780633bc5de301461018257806342cbb15c146102125780635e01eb5a1461023d578063796b89b914610294576100dd565b806312e05dd1146100e2578063209652551461010d5780632df8e9491461012b575b600080fd5b3480156100ee57600080fd5b506100f761047a565b6040518082815260200191505060405180910390f35b610115610482565b6040518082815260200191505060405180910390f35b34801561013757600080fd5b5061014061048a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561018e57600080fd5b50610197610492565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101d75780820151818401526020810190506101bc565b50505050905090810190601f1680156102045780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561021e57600080fd5b506102276104df565b6040518082815260200191505060405180910390f35b34801561024957600080fd5b506102526104e7565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102a057600080fd5b506102a96104ef565b6040518082815260200191505060405180910390f35b3480156102cb57600080fd5b506102d46104f7565b6040518082815260200191505060405180910390f35b3480156102f657600080fd5b506102ff6104ff565b6040518082815260200191505060405180910390f35b34801561032157600080fd5b5061032a610507565b6040518082815260200191505060405180910390f35b34801561034c57600080fd5b5061035561050f565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156103b557600080fd5b506103be61053a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561040c57600080fd5b506104396004803603602081101561042357600080fd5b8101908080359060200190929190505050610542565b6040518082815260200191505060405180910390f35b34801561045b57600080fd5b5061046461054d565b6040518082815260200191505060405180910390f35b600044905090565b600034905090565b600041905090565b60606000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050905090565b600043905090565b600033905090565b600042905090565b600045905090565b60003a905090565b600042905090565b600080357fffffffff0000000000000000000000000000000000000000000000000000000016905090565b600032905090565b600081409050919050565b60005a90509056fea265627a7a72305820176147ee7465d40bf08511f59fe62908bb1d7c7e2d735041ad6ab61add911da164736f6c63430005090032";

    public static final String FUNC_GETBLOCKDIFFICULTY = "getBlockDifficulty";

    public static final String FUNC_GETVALUE = "getValue";

    public static final String FUNC_GETBLOCKCOINBASE = "getBlockCoinbase";

    public static final String FUNC_GETDATA = "getData";

    public static final String FUNC_GETBLOCKNUMBER = "getBlockNumber";

    public static final String FUNC_GETSENDER = "getSender";

    public static final String FUNC_GETBLOCKTIMESTAMP = "getBlockTimestamp";

    public static final String FUNC_GETGASLIMIT = "getGaslimit";

    public static final String FUNC_GETGASPRICE = "getGasprice";

    public static final String FUNC_GETNOW = "getNow";

    public static final String FUNC_GETSIG = "getSig";

    public static final String FUNC_GETORIGIN = "getOrigin";

    public static final String FUNC_GETBLOCKHASH = "getBlockhash";

    public static final String FUNC_GETGASLEFT = "getGasleft";

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getBlockDifficulty() {
        final Function function = new Function(FUNC_GETBLOCKDIFFICULTY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> getValue(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GETVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<String> getBlockCoinbase() {
        final Function function = new Function(FUNC_GETBLOCKCOINBASE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> getData() {
        final Function function = new Function(FUNC_GETDATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getBlockNumber() {
        final Function function = new Function(FUNC_GETBLOCKNUMBER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getSender() {
        final Function function = new Function(FUNC_GETSENDER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getBlockTimestamp() {
        final Function function = new Function(FUNC_GETBLOCKTIMESTAMP, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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

    public RemoteCall<byte[]> getSig() {
        final Function function = new Function(FUNC_GETSIG, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes4>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> getOrigin() {
        final Function function = new Function(FUNC_GETORIGIN, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> getBlockhash(BigInteger blockNumber) {
        final Function function = new Function(FUNC_GETBLOCKHASH, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(blockNumber)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getGasleft() {
        final Function function = new Function(FUNC_GETGASLEFT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BlockTransactionPropertiesFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BlockTransactionPropertiesFunctions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BlockTransactionPropertiesFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BlockTransactionPropertiesFunctions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
