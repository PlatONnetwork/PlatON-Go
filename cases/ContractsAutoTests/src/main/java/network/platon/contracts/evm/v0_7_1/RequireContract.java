package network.platon.contracts.evm.v0_7_1;

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
public class RequireContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060df8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806314fef936146037578063de29278914606c575b600080fd5b606a60048036036040811015604b57600080fd5b8101908080359060200190929190803590602001909291905050506088565b005b607260a0565b6040518082815260200191505060405180910390f35b808211609357600080fd5b8082036000819055505050565b6000805490509056fea2646970667358221220753079ebd816a05be05b12e9066319b1f9783e880c847e7b6c37f87d9bb81ec764736f6c63430007010033";

    public static final String FUNC_GETRESULT = "getResult";

    public static final String FUNC_TOSENDERAMOUNT = "toSenderAmount";

    protected RequireContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RequireContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getResult() {
        final Function function = new Function(FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> toSenderAmount(BigInteger frist, BigInteger second) {
        final Function function = new Function(
                FUNC_TOSENDERAMOUNT, 
                Arrays.<Type>asList(new Uint256(frist),
                new Uint256(second)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<RequireContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RequireContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RequireContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RequireContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static RequireContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RequireContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RequireContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RequireContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
