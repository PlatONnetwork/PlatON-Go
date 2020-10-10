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
public class ContractAndAddress extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506104fe806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806302e9cd8b1461005c57806346df069e14610090578063a52e2905146100c4578063a7010a66146100f8578063c105b57c14610102575b600080fd5b610064610136565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610098610160565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100cc610189565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101006101b3565b005b61010a6103b6565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60006040516101c1906103e0565b604051809103906000f0801580156101dd573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051610229906103ec565b604051809103906000f080158015610245573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051610292906103ec565b604051809103906000f0801580156102ae573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b606b806103f983390190565b6065806104648339019056fe6080604052348015600f57600080fd5b50604e80601d6000396000f3fe608060405236600a57005b348015601557600080fd5b5000fea26469706673582212202101511779e642022e671470f528eb7cf5638ce6489f60c62443288e41e15bb864736f6c634300070100336080604052348015600f57600080fd5b50604880601d6000396000f3fe6080604052348015600f57600080fd5b5000fea2646970667358221220ffda2f1c06bf9075b9a2f41f788546080985da9627afc1a689a7d4cb5d436d3a64736f6c63430007010033a264697066735822122005376edc9139184139f84aae478c0476a537ab6d3fd5029ebfb4d0101201f5e364736f6c63430007010033";

    public static final String FUNC_GETADDRESSTOPAYABLE = "getAddressToPayable";

    public static final String FUNC_GETNONALCONTRACTADDRESS = "getNonalContractAddress";

    public static final String FUNC_GETNONALPAYABLEADDRESS = "getNonalPayableAddress";

    public static final String FUNC_GETPAYABLETOADDRESS = "getPayableToAddress";

    public static final String FUNC_PAYABLEORNOT = "payableOrNot";

    protected ContractAndAddress(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ContractAndAddress(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getAddressToPayable() {
        final Function function = new Function(
                FUNC_GETADDRESSTOPAYABLE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getNonalContractAddress() {
        final Function function = new Function(
                FUNC_GETNONALCONTRACTADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getNonalPayableAddress() {
        final Function function = new Function(
                FUNC_GETNONALPAYABLEADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getPayableToAddress() {
        final Function function = new Function(
                FUNC_GETPAYABLETOADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> payableOrNot() {
        final Function function = new Function(
                FUNC_PAYABLEORNOT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractAndAddress.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ContractAndAddress> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractAndAddress.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ContractAndAddress load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ContractAndAddress(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ContractAndAddress load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ContractAndAddress(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
