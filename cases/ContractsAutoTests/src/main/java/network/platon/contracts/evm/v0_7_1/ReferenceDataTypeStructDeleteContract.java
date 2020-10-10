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
public class ReferenceDataTypeStructDeleteContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506005600081905550600a600160020181905550600180600301600060018152602001908152602001600020819055506002600160030160006002815260200190815260200160002081905550600360016000016000018190555060018060000160010160006001815260200190815260200160002060006101000a81548160ff02191690831515021790555060018060000160010160006002815260200190815260200160002060006101000a81548160ff021916908315150217905550600160008082016000808201600090555050600282016000905550506000805561019e806100fe6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806311977c5c1461005c5780631268893e1461007a5780635ff76c8a1461009857806379e44a38146100b6578063d587919c146100d4575b600080fd5b6100646100f4565b6040518082815260200191505060405180910390f35b6100826100fd565b6040518082815260200191505060405180910390f35b6100a061011c565b6040518082815260200191505060405180910390f35b6100be61012c565b6040518082815260200191505060405180910390f35b6100dc610139565b60405180821515815260200191505060405180910390f35b60008054905090565b6000600160030160006001815260200190815260200160002054905090565b6000600160000160000154905090565b6000600160020154905090565b6000600160000160010160006001815260200190815260200160002060009054906101000a900460ff1690509056fea26469706673582212201c976f47960fc90a00b71038bff527155beea276b6a8adbb48b05a56e0f9c73364736f6c63430007010033";

    public static final String FUNC_GETNESTEDMAPPING = "getNestedMapping";

    public static final String FUNC_GETNESTEDVALUE = "getNestedValue";

    public static final String FUNC_GETTODELETEINT = "getToDeleteInt";

    public static final String FUNC_GETTOPMAPPING = "getTopMapping";

    public static final String FUNC_GETTOPVALUE = "getTopValue";

    protected ReferenceDataTypeStructDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeStructDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructDeleteContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructDeleteContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeStructDeleteContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeStructDeleteContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> getNestedMapping() {
        final Function function = new Function(
                FUNC_GETNESTEDMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getNestedValue() {
        final Function function = new Function(
                FUNC_GETNESTEDVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getToDeleteInt() {
        final Function function = new Function(
                FUNC_GETTODELETEINT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getTopMapping() {
        final Function function = new Function(
                FUNC_GETTOPMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getTopValue() {
        final Function function = new Function(
                FUNC_GETTOPVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static ReferenceDataTypeStructDeleteContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructDeleteContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeStructDeleteContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeStructDeleteContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
