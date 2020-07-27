package network.platon.contracts;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.DynamicArray;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class ContractArray extends Contract {
    private static final String BINARY = "608060405260036040519080825280602002602001820160405280156100345781602001602082028038833980820191505090505b506000908051906020019061004a92919061005d565b5034801561005757600080fd5b5061012a565b8280548282559060005260206000209081019282156100d6579160200282015b828111156100d55782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509160200191906001019061007d565b5b5090506100e391906100e7565b5090565b61012791905b8082111561012357600081816101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055506001016100ed565b5090565b90565b610319806101396000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806326121ff014610046578063807b4c3314610050578063e276c799146100af575b600080fd5b61004e6100f9565b005b6100586101fd565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561009b578082015181840152602081019050610080565b505050509050019250505060405180910390f35b6100b761028b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101016102c2565b60003090806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505030600160006003811061017657fe5b0160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555030816000600381106101c357fe5b602002019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505050565b6060600080548060200260200160405190810160405280929190818152602001828054801561028157602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610237575b5050505050905090565b6000600160006003811061029b57fe5b0160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b604051806060016040528060039060208202803883398082019150509050509056fea265627a7a72315820b8af178b05ba84ef3d64c785742261662a105e31498aec0b3533949ab04de67c64736f6c634300050d0032";

    public static final String FUNC_F = "f";

    public static final String FUNC_GETX = "getx";

    public static final String FUNC_GETY = "gety";

    protected ContractArray(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ContractArray(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> f() {
        final Function function = new Function(
                FUNC_F, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getx() {
        final Function function = new Function(FUNC_GETX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<List> gety() {
        final Function function = new Function(FUNC_GETY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public static RemoteCall<ContractArray> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractArray.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ContractArray> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractArray.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ContractArray load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ContractArray(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ContractArray load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ContractArray(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
