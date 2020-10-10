package network.platon.contracts.evm.v0_6_12;

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
public class ConstructorPublicDeclaraction extends Contract {
    private static final String BINARY = "60806040526000805534801561001457600080fd5b5060405161016c38038061016c8339818101604052602081101561003757600080fd5b8101908080519060200190929190505050806000819055505061010d8061005f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806382ab890a146037578063a87d942c146093575b600080fd5b606060048036036020811015604b57600080fd5b810190808035906020019092919050505060af565b604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b609960ce565b6040518082815260200191505060405180910390f35b6000808260008082825401925050819055503360005491509150915091565b6000805490509056fea26469706673582212209b611323ce297a2be8a0b1694e262222011e280562b069e8ec2c60eb1a7984a564736f6c634300060c0033";

    public static final String FUNC_GETCOUNT = "getCount";

    public static final String FUNC_UPDATE = "update";

    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ConstructorPublicDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<ConstructorPublicDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _count) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(_count)));
        return deployRemoteCall(ConstructorPublicDeclaraction.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
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

    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ConstructorPublicDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ConstructorPublicDeclaraction(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
