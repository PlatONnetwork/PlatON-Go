package network.platon.contracts.evm.v0_6_12;

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
 * <p>Generated with web3j version 0.13.2.0.
 */
public class BlockTransactionPropertiesFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061052a806100206000396000f3fe6080604052600436106100dd5760003560e01c8063a16963b31161007f578063d12d910211610059578063d12d910214610314578063df1f29ee1461035e578063e9413d381461039f578063edb4b865146103ee576100dd565b8063a16963b314610293578063ab70fd69146102be578063bbe4fd50146102e9576100dd565b80633bc5de30116100bb5780633bc5de301461016c57806342cbb15c146101fc5780635e01eb5a14610227578063796b89b914610268576100dd565b806312e05dd1146100e2578063209652551461010d5780632df8e9491461012b575b600080fd5b3480156100ee57600080fd5b506100f7610419565b6040518082815260200191505060405180910390f35b610115610421565b6040518082815260200191505060405180910390f35b34801561013757600080fd5b50610140610429565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561017857600080fd5b50610181610431565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101c15780820151818401526020810190506101a6565b50505050905090810190601f1680156101ee5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561020857600080fd5b5061021161047e565b6040518082815260200191505060405180910390f35b34801561023357600080fd5b5061023c610486565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561027457600080fd5b5061027d61048e565b6040518082815260200191505060405180910390f35b34801561029f57600080fd5b506102a8610496565b6040518082815260200191505060405180910390f35b3480156102ca57600080fd5b506102d361049e565b6040518082815260200191505060405180910390f35b3480156102f557600080fd5b506102fe6104a6565b6040518082815260200191505060405180910390f35b34801561032057600080fd5b506103296104ae565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561036a57600080fd5b506103736104d9565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103ab57600080fd5b506103d8600480360360208110156103c257600080fd5b81019080803590602001909291905050506104e1565b6040518082815260200191505060405180910390f35b3480156103fa57600080fd5b506104036104ec565b6040518082815260200191505060405180910390f35b600044905090565b600034905090565b600041905090565b60606000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050905090565b600043905090565b600033905090565b600042905090565b600045905090565b60003a905090565b600042905090565b600080357fffffffff0000000000000000000000000000000000000000000000000000000016905090565b600032905090565b600081409050919050565b60005a90509056fea26469706673582212205a609d38a4324262095dfbcebc4d451f3ed7255f3e4f18087a2776818683439464736f6c634300060c0033";

    public static final String FUNC_GETBLOCKCOINBASE = "getBlockCoinbase";

    public static final String FUNC_GETBLOCKDIFFICULTY = "getBlockDifficulty";

    public static final String FUNC_GETBLOCKNUMBER = "getBlockNumber";

    public static final String FUNC_GETBLOCKTIMESTAMP = "getBlockTimestamp";

    public static final String FUNC_GETBLOCKHASH = "getBlockhash";

    public static final String FUNC_GETDATA = "getData";

    public static final String FUNC_GETGASLEFT = "getGasleft";

    public static final String FUNC_GETGASLIMIT = "getGaslimit";

    public static final String FUNC_GETGASPRICE = "getGasprice";

    public static final String FUNC_GETNOW = "getNow";

    public static final String FUNC_GETORIGIN = "getOrigin";

    public static final String FUNC_GETSENDER = "getSender";

    public static final String FUNC_GETSIG = "getSig";

    public static final String FUNC_GETVALUE = "getValue";

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BlockTransactionPropertiesFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getBlockCoinbase() {
        final Function function = new Function(
                FUNC_GETBLOCKCOINBASE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBlockDifficulty() {
        final Function function = new Function(
                FUNC_GETBLOCKDIFFICULTY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBlockNumber() {
        final Function function = new Function(
                FUNC_GETBLOCKNUMBER, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBlockTimestamp() {
        final Function function = new Function(
                FUNC_GETBLOCKTIMESTAMP, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBlockhash(BigInteger blockNumber) {
        final Function function = new Function(
                FUNC_GETBLOCKHASH, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(blockNumber)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getData() {
        final Function function = new Function(
                FUNC_GETDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getGasleft() {
        final Function function = new Function(
                FUNC_GETGASLEFT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getGaslimit() {
        final Function function = new Function(
                FUNC_GETGASLIMIT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getGasprice() {
        final Function function = new Function(
                FUNC_GETGASPRICE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getNow() {
        final Function function = new Function(
                FUNC_GETNOW, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getOrigin() {
        final Function function = new Function(
                FUNC_GETORIGIN, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSender() {
        final Function function = new Function(
                FUNC_GETSENDER, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSig() {
        final Function function = new Function(
                FUNC_GETSIG, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getValue() {
        final Function function = new Function(
                FUNC_GETVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
