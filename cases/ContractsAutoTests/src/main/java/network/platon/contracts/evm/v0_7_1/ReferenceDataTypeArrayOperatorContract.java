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
public class ReferenceDataTypeArrayOperatorContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506105fe806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806399a1c3691161005b57806399a1c369146101605780639a6e3cb7146101aa578063b9903334146101ee578063fd081ef8146102385761007d565b80631ff0db4014610082578063676cb904146100cc57806386d1171014610116575b600080fd5b61008a610282565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b6100d46102b4565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b61011e610300565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b61016861034c565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b6101b261037e565b60405180861515815260200185151581526020018415158152602001831515815260200182151581526020019550505050505060405180910390f35b6101f661056f565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b6102406105a1565b60405180837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020018260ff1681526020019250505060405180910390f35b6000806000608160f81b90506000608060f81b90506000818316905060008160f81c9050818195509550505050509091565b6000806000608160f81b905060006001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916901b905060008160f81c90508181945094505050509091565b6000806000608160f81b905060006001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916901c905060008160f81c90508181945094505050509091565b6000806000608160f81b90506000608060f81b90506000818318905060008160f81c9050818195509550505050509091565b6000806000806000807f6100000000000000000000000000000000000000000000000000000000000000905060007f620000000000000000000000000000000000000000000000000000000000000090506000606160f81b90506000827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916847effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161090506000837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916857effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161190506000837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916867effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161490506000847effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916867effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916141590506000857effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916887effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161015905084848484849c509c509c509c509c5050505050505050509091929394565b6000806000608160f81b90506000608060f81b90506000818317905060008160f81c9050818195509550505050509091565b6000806000608160f81b905060008119905060008160f81c9050818194509450505050909156fea264697066735822122089f5a84dc14c9c68503c52bc35e16ca86d1d5246a3d7a821bd71148dad0ccc8264736f6c63430007010033";

    public static final String FUNC_ARRAYBITANDOPERATORS = "arrayBitAndOperators";

    public static final String FUNC_ARRAYBITINVERSEOPERATORS = "arrayBitInverseOperators";

    public static final String FUNC_ARRAYBITLEFTSHIFTPERATORS = "arrayBitLeftShiftperators";

    public static final String FUNC_ARRAYBITOROPERATORS = "arrayBitOrOperators";

    public static final String FUNC_ARRAYBITRIGHTSHIFTPERATORS = "arrayBitRightShiftperators";

    public static final String FUNC_ARRAYBITXOROPERATORS = "arrayBitXOROperators";

    public static final String FUNC_ARRAYCOMPARE = "arrayCompare";

    protected ReferenceDataTypeArrayOperatorContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeArrayOperatorContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> arrayBitAndOperators() {
        final Function function = new Function(
                FUNC_ARRAYBITANDOPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayBitInverseOperators() {
        final Function function = new Function(
                FUNC_ARRAYBITINVERSEOPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayBitLeftShiftperators() {
        final Function function = new Function(
                FUNC_ARRAYBITLEFTSHIFTPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayBitOrOperators() {
        final Function function = new Function(
                FUNC_ARRAYBITOROPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayBitRightShiftperators() {
        final Function function = new Function(
                FUNC_ARRAYBITRIGHTSHIFTPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayBitXOROperators() {
        final Function function = new Function(
                FUNC_ARRAYBITXOROPERATORS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> arrayCompare() {
        final Function function = new Function(
                FUNC_ARRAYCOMPARE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ReferenceDataTypeArrayOperatorContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayOperatorContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeArrayOperatorContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayOperatorContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ReferenceDataTypeArrayOperatorContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayOperatorContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeArrayOperatorContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayOperatorContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
