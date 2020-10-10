package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Address;
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
    private static final String BINARY = "608060405234801561001057600080fd5b50610542806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806302e9cd8b1461005c57806346df069e146100a6578063a52e2905146100f0578063a7010a661461013a578063c105b57c14610144575b600080fd5b61006461018e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100ae6101b8565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100f86101e1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61014261020b565b005b61014c61040f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600060405161021990610439565b604051809103906000f080158015610235573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060405161028190610445565b604051809103906000f08015801561029d573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040516102ea90610445565b604051809103906000f080158015610306573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60588061045283390190565b6064806104aa8339019056fe6080604052348015600f57600080fd5b50603b80601d6000396000f3fe608060405200fea265627a7a72315820365d0631b9957f63acafe1ae6efa4e43b8ede9b264900652a2b0f0a127764c0b64736f6c634300051100326080604052348015600f57600080fd5b50604780601d6000396000f3fe6080604052348015600f57600080fd5b5000fea265627a7a723158201e2b8da001f5bfa8d514df4902c1fbcf2e91f58798e8a03d0c9fec5a0b9e4c9364736f6c63430005110032a265627a7a72315820aa085e8507579522a580d0418725fe4d6af51793695d5b332314fd496e239b3064736f6c63430005110032";

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

    public RemoteCall<String> getAddressToPayable() {
        final Function function = new Function(FUNC_GETADDRESSTOPAYABLE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getNonalContractAddress() {
        final Function function = new Function(FUNC_GETNONALCONTRACTADDRESS, 
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
