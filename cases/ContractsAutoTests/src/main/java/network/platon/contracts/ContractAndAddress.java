package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class ContractAndAddress extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610578806100206000396000f3fe60806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806302e9cd8b1461007257806346df069e146100c9578063a52e290514610120578063a7010a6614610177578063c105b57c1461018e575b600080fd5b34801561007e57600080fd5b506100876101e5565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100d557600080fd5b506100de61020f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012c57600080fd5b50610135610238565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561018357600080fd5b5061018c610262565b005b34801561019a57600080fd5b506101a361045a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600061026c610484565b604051809103906000f080158015610288573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506102d0610493565b604051809103906000f0801580156102ec573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610335610493565b604051809103906000f080158015610351573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b604051604f806104a383390190565b604051605b806104f28339019056fe6080604052348015600f57600080fd5b50603280601d6000396000f3fe608060405200fea165627a7a723058208e5333c65765a5873d66c60c2767fafabac18165c3bc0cd8b0c2fc5c522c4ff200296080604052348015600f57600080fd5b50603e80601d6000396000f3fe6080604052348015600f57600080fd5b5000fea165627a7a7230582081536ebe215ff398ac65007b71edc7371bc56a779dbbbbde86c3800ac3356a890029a165627a7a72305820224c19fc049becc45489cf8610271d5ff753a52c64172fb89df04435583665330029";

    public static final String FUNC_GETADDRESSTOPAYABLE = "getAddressToPayable";

    public static final String FUNC_GETNONALPAYABLEADDRESS = "getNonalPayableAddress";

    public static final String FUNC_GETPAYABLETOADDRESS = "getPayableToAddress";

    public static final String FUNC_PAYABLEORNOT = "payableOrNot";

    public static final String FUNC_GETNONALCONTRACTADDRESS = "getNonalContractAddress";

    @Deprecated
    protected ContractAndAddress(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected ContractAndAddress(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected ContractAndAddress(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected ContractAndAddress(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> getAddressToPayable() {
        final Function function = new Function(FUNC_GETADDRESSTOPAYABLE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getNonalPayableAddress() {
        final Function function = new Function(FUNC_GETNONALPAYABLEADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getPayableToAddress() {
        final Function function = new Function(FUNC_GETPAYABLETOADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> payableOrNot() {
        final Function function = new Function(
                FUNC_PAYABLEORNOT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getNonalContractAddress() {
        final Function function = new Function(FUNC_GETNONALCONTRACTADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(ContractAndAddress.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ContractAndAddress.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(ContractAndAddress.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(ContractAndAddress.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static ContractAndAddress load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new ContractAndAddress(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static ContractAndAddress load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new ContractAndAddress(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static ContractAndAddress load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new ContractAndAddress(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static ContractAndAddress load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new ContractAndAddress(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
