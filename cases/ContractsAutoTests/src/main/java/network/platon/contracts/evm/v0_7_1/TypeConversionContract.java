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
public class TypeConversionContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610305806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806399a909621161005b57806399a909621461010f578063a136096714610159578063ad42221214610196578063dcefd42f146101b75761007d565b8063744708f814610082578063853255cc146100cc5780639311ca69146100ed575b600080fd5b61008a6101f7565b604051808363ffffffff168152602001827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b6100d4610219565b604051808260010b815260200191505060405180910390f35b6100f5610233565b604051808261ffff16815260200191505060405180910390f35b61011761024a565b604051808361ffff168152602001827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681526020019250505060405180910390f35b61016161026a565b60405180827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b61019e6102a3565b604051808260000b815260200191505060405180910390f35b6101bf6102b7565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6000806000611234905060008161ffff169050808160e01b9350935050509091565b60008060029050600060649050808260000b019250505090565b600080600a905060008160ff169050809250505090565b6000806000631234567890506000819050808160f01b9350935050509091565b60008061123460f01b90506000817dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19169050809250505090565b600080600190506000819050809250505090565b60008061123460f01b9050600081905080925050509056fea2646970667358221220fd50f9fdae5a0364dc69eec44a79fc55cabddfbad646a6b5555937d674d6d87464736f6c63430007010033";

    public static final String FUNC_CONVERSION = "conversion";

    public static final String FUNC_DISPLAYCONVERSION = "displayConversion";

    public static final String FUNC_DISPLAYCONVERSION1 = "displayConversion1";

    public static final String FUNC_DISPLAYCONVERSION2 = "displayConversion2";

    public static final String FUNC_DISPLAYCONVERSION3 = "displayConversion3";

    public static final String FUNC_DISPLAYCONVERSION4 = "displayConversion4";

    public static final String FUNC_SUM = "sum";

    protected TypeConversionContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected TypeConversionContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> conversion() {
        final Function function = new Function(
                FUNC_CONVERSION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> displayConversion() {
        final Function function = new Function(
                FUNC_DISPLAYCONVERSION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> displayConversion1() {
        final Function function = new Function(
                FUNC_DISPLAYCONVERSION1, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> displayConversion2() {
        final Function function = new Function(
                FUNC_DISPLAYCONVERSION2, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> displayConversion3() {
        final Function function = new Function(
                FUNC_DISPLAYCONVERSION3, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> displayConversion4() {
        final Function function = new Function(
                FUNC_DISPLAYCONVERSION4, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> sum() {
        final Function function = new Function(
                FUNC_SUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<TypeConversionContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeConversionContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<TypeConversionContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeConversionContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static TypeConversionContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new TypeConversionContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static TypeConversionContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new TypeConversionContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
