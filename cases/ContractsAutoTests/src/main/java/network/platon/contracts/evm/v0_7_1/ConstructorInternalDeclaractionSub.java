package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.FunctionEncoder;
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
public class ConstructorInternalDeclaractionSub extends Contract {
    private static final String BINARY = "6080604052600a60005534801561001557600080fd5b506040516101a73803806101a78339818101604052602081101561003857600080fd5b810190808051906020019092919050505080806000819055505050610145806100626000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806306661abd1461004657806382ab890a14610064578063a87d942c146100c3575b600080fd5b61004e6100e1565b6040518082815260200191505060405180910390f35b6100906004803603602081101561007a57600080fd5b81019080803590602001909291905050506100e7565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b6100cb610106565b6040518082815260200191505060405180910390f35b60005481565b6000808260008082825401925050819055503360005491509150915091565b6000805490509056fea264697066735822122058c44994e40d27daee26833f27853dcf97839d56a058cf752fe0816f2672eafd64736f6c63430007010033";

    public static final String FUNC_COUNT = "count";

    public static final String FUNC_GETCOUNT = "getCount";

    public static final String FUNC_UPDATE = "update";

    protected ConstructorInternalDeclaractionSub(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstructorInternalDeclaractionSub(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ConstructorInternalDeclaractionSub> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorInternalDeclaractionSub.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ConstructorInternalDeclaractionSub> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorInternalDeclaractionSub.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public RemoteCall<TransactionReceipt> count() {
        final Function function = new Function(
                FUNC_COUNT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getCount() {
        final Function function = new Function(
                FUNC_GETCOUNT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(amount)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static ConstructorInternalDeclaractionSub load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalDeclaractionSub(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstructorInternalDeclaractionSub load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorInternalDeclaractionSub(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
