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
public class ReferenceDataTypeStructRecursiveContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610174806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806397b93b1a14610030575b600080fd5b61003861005c565b60405180848152602001838152602001828152602001935050505060405180910390f35b60008060006100696100c4565b6000610074906100d7565b9050806000015151816000015160008151811061008d57fe5b6020026020010151600001515182600001516001815181106100ab57fe5b6020026020010151600001515193509350935050909192565b6040518060200160405280606081525090565b60405180602001604052908160008201805480602002602001604051908101604052809291908181526020016000905b8282101561013357838290600052602060002001610124906100d7565b81526020019060010190610107565b50505050815250509056fea26469706673582212205eaa0b839d2dd28cc0b411e9d6b46b1fa467ac0dab725e1f21db57d6f68e791d64736f6c634300060c0033";

    public static final String FUNC_GETSTRUCTPERSONLENGTH = "getStructPersonLength";

    protected ReferenceDataTypeStructRecursiveContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeStructRecursiveContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getStructPersonLength() {
        final Function function = new Function(
                FUNC_GETSTRUCTPERSONLENGTH, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ReferenceDataTypeStructRecursiveContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructRecursiveContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructRecursiveContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructRecursiveContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ReferenceDataTypeStructRecursiveContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructRecursiveContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeStructRecursiveContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructRecursiveContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
