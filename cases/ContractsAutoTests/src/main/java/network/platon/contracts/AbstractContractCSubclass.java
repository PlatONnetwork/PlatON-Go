package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Int256;
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
public class AbstractContractCSubclass extends Contract {
    private static final String BINARY = "60806040526040518060200160405280600081525060009080519060200190610029929190610062565b50604051806020016040528060008152506001908051906020019061004f929190610062565b5034801561005c57600080fd5b50610107565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a357805160ff19168380011785556100d1565b828001600101855582156100d1579182015b828111156100d05782518255916020019190600101906100b5565b5b5090506100de91906100e2565b5090565b61010491905b808211156101005760008160009055506001016100e8565b5090565b90565b6105c6806101166000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806331aa8a6e146100675780633af1a463146100ea5780634fc341131461016d5780639c72890b14610228578063accab56b14610246578063e652e56514610301575b600080fd5b61006f610384565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100af578082015181840152602081019050610094565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100f26103c6565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610132578082015181840152602081019050610117565b50505050905090810190601f16801561015f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102266004803603602081101561018357600080fd5b81019080803590602001906401000000008111156101a057600080fd5b8201836020820111156101b257600080fd5b803590602001918460018302840111640100000000831117156101d457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610468565b005b610230610482565b6040518082815260200191505060405180910390f35b6102ff6004803603602081101561025c57600080fd5b810190808035906020019064010000000081111561027957600080fd5b82018360208201111561028b57600080fd5b803590602001918460018302840111640100000000831117156102ad57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610490565b005b6103096104aa565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561034957808201518184015260208101905061032e565b50505050905090810190601f1680156103765780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6060806040518060400160405280600881526020017f635375624e616d6500000000000000000000000000000000000000000000000081525090508091505090565b606060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561045e5780601f106104335761010080835404028352916020019161045e565b820191906000526020600020905b81548152906001019060200180831161044157829003601f168201915b5050505050905090565b806000908051906020019061047e9291906104ec565b5050565b600080601490508091505090565b80600190805190602001906104a69291906104ec565b5050565b6060806040518060400160405280600a81526020017f706172656e744e616d650000000000000000000000000000000000000000000081525090508091505090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061052d57805160ff191683800117855561055b565b8280016001018555821561055b579182015b8281111561055a57825182559160200191906001019061053f565b5b509050610568919061056c565b5090565b61058e91905b8082111561058a576000816000905550600101610572565b5090565b9056fea265627a7a72315820b26c2972c24ef19a18088e18f83fdba51a181d905071038cae20c607ee8d750f64736f6c634300050d0032";

    public static final String FUNC_ASUBAGE = "aSubAge";

    public static final String FUNC_ASUBNAME = "aSubName";

    public static final String FUNC_CSUBNAME = "cSubName";

    public static final String FUNC_PARENTNAME = "parentName";

    public static final String FUNC_SETASUBNAME = "setASubName";

    public static final String FUNC_SETPARENTNAME = "setParentName";

    protected AbstractContractCSubclass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AbstractContractCSubclass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> aSubAge() {
        final Function function = new Function(FUNC_ASUBAGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> aSubName() {
        final Function function = new Function(FUNC_ASUBNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> cSubName() {
        final Function function = new Function(FUNC_CSUBNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> parentName() {
        final Function function = new Function(FUNC_PARENTNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> setASubName(String v) {
        final Function function = new Function(
                FUNC_SETASUBNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(v)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setParentName(String name) {
        final Function function = new Function(
                FUNC_SETPARENTNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(name)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<AbstractContractCSubclass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractCSubclass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AbstractContractCSubclass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractCSubclass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AbstractContractCSubclass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractCSubclass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AbstractContractCSubclass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractCSubclass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
