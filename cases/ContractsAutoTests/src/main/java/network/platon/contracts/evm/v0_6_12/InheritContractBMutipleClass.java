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
public class InheritContractBMutipleClass extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060c88061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063430fe9c114603757806354b39533146053575b600080fd5b603d606f565b6040518082815260200191505060405180910390f35b6059607c565b6040518082815260200191505060405180910390f35b600060776089565b905090565b60006084606f565b905090565b6000600190509056fea2646970667358221220a6c71d806c25be00e643d709a752cbfec651f80d94f663fba06fc73afe95fe0864736f6c634300060c0033";

    public static final String FUNC_CALLGETDATEB = "callGetDateB";

    public static final String FUNC_GETDATE = "getDate";

    protected InheritContractBMutipleClass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractBMutipleClass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> callGetDateB() {
        final Function function = new Function(
                FUNC_CALLGETDATEB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getDate() {
        final Function function = new Function(
                FUNC_GETDATE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritContractBMutipleClass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractBMutipleClass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractBMutipleClass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractBMutipleClass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractBMutipleClass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBMutipleClass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractBMutipleClass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractBMutipleClass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
