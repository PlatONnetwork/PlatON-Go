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
public class AssertContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060fa8061005f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806314fef936146037578063de292789146080575b600080fd5b606a60048036036040811015604b57600080fd5b810190808035906020019092919080359060200190929190505050609c565b6040518082815260200191505060405180910390f35b608660ba565b6040518082815260200191505060405180910390f35b600081831160a657fe5b818301600181905550600154905092915050565b600060015490509056fea2646970667358221220c9a7ffdd68a686a3fd213f24ad86923eee8e670e0b3bad9891fa1df04f30a7e864736f6c63430007010033";

    public static final String FUNC_GETRESULT = "getResult";

    public static final String FUNC_TOSENDERAMOUNT = "toSenderAmount";

    protected AssertContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AssertContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<AssertContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssertContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AssertContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AssertContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> getResult() {
        final Function function = new Function(
                FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> toSenderAmount(BigInteger frist, BigInteger second) {
        final Function function = new Function(
                FUNC_TOSENDERAMOUNT, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(frist), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(second)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static AssertContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AssertContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AssertContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AssertContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
