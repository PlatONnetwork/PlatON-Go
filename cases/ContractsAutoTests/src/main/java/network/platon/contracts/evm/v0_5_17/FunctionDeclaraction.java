package network.platon.contracts.evm.v0_5_17;

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
public class FunctionDeclaraction extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061020c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806312065fe014610046578063ab95edb114610064578063cb533b38146100d9575b600080fd5b61004e61014e565b6040518082815260200191505060405180910390f35b6100906004803603602081101561007a57600080fd5b8101908080359060200190929190505050610157565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610105600480360360208110156100ef57600080fd5b8101908080359060200190929190505050610175565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b60008054905090565b60008061016383610193565b50503360005481915091509150915091565b600080610181836101b5565b50503360005481915091509150915091565b6000808260008082825401925050819055503360005481915091509150915091565b600080826000808282540192505081905550336000548191509150915091509156fea265627a7a72315820b8bfa0f61106f3d392fef947213099e6cb91e1bbc8b641f5db9df6e57919ba9d64736f6c63430005110032";

    public static final String FUNC_GETBALANCE = "getBalance";

    public static final String FUNC_UPDATE_EXTERNAL = "update_external";

    public static final String FUNC_UPDATE_PUBLIC = "update_public";

    protected FunctionDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected FunctionDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getBalance() {
        final Function function = new Function(FUNC_GETBALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> update_external(BigInteger amount_ex) {
        final Function function = new Function(
                FUNC_UPDATE_EXTERNAL, 
                Arrays.<Type>asList(new Uint256(amount_ex)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> update_public(BigInteger amount_pu) {
        final Function function = new Function(
                FUNC_UPDATE_PUBLIC, 
                Arrays.<Type>asList(new Uint256(amount_pu)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<FunctionDeclaraction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FunctionDeclaraction.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<FunctionDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FunctionDeclaraction.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static FunctionDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new FunctionDeclaraction(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static FunctionDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new FunctionDeclaraction(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
