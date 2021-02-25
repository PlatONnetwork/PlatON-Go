package network.platon.contracts.evm;

import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class AbstractContractFSubclass extends Contract {
    private static final String BINARY = "6080604052604051806020016040528060008152506000908051906020019061002992919061003c565b5034801561003657600080fd5b506100e1565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061007d57805160ff19168380011785556100ab565b828001600101855582156100ab579182015b828111156100aa57825182559160200191906001019061008f565b5b5090506100b891906100bc565b5090565b6100de91905b808211156100da5760008160009055506001016100c2565b5090565b90565b610481806100f06000396000f3fe608060405234801561001057600080fd5b506004361061004b5760003560e01c8062af095414610050578063accab56b1461010b578063e16ae640146101c6578063e652e56514610249575b600080fd5b6101096004803603602081101561006657600080fd5b810190808035906020019064010000000081111561008357600080fd5b82018360208201111561009557600080fd5b803590602001918460018302840111640100000000831117156100b757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506102cc565b005b6101c46004803603602081101561012157600080fd5b810190808035906020019064010000000081111561013e57600080fd5b82018360208201111561015057600080fd5b8035906020019184600183028401116401000000008311171561017257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506102e6565b005b6101ce610300565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561020e5780820151818401526020810190506101f3565b50505050905090810190601f16801561023b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610251610305565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610291578082015181840152602081019050610276565b50505050905090810190601f1680156102be5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b80600090805190602001906102e29291906103a7565b5050565b80600090805190602001906102fc9291906103a7565b5050565b606090565b606060008054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561039d5780601f106103725761010080835404028352916020019161039d565b820191906000526020600020905b81548152906001019060200180831161038057829003601f168201915b5050505050905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106103e857805160ff1916838001178555610416565b82800160010185558215610416579182015b828111156104155782518255916020019190600101906103fa565b5b5090506104239190610427565b5090565b61044991905b8082111561044557600081600090555060010161042d565b5090565b9056fea265627a7a7231582079c589d9001e32cbe327b5c3d0eb5d88cda8a0bfdb1457721d90ad8f4970c17264736f6c634300050d0032";

    public static final String FUNC_DSUBCLASSNAME = "dSubClassName";

    public static final String FUNC_PARENTNAME = "parentName";

    public static final String FUNC_SETPARENTNAME = "setParentName";

    public static final String FUNC_SETPARENTNAMED = "setParentNameD";

    protected AbstractContractFSubclass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AbstractContractFSubclass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<String> dSubClassName() {
        final Function function = new Function(FUNC_DSUBCLASSNAME, 
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

    public RemoteCall<TransactionReceipt> setParentName(String name) {
        final Function function = new Function(
                FUNC_SETPARENTNAME, 
                Arrays.<Type>asList(new Utf8String(name)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setParentNameD(String name) {
        final Function function = new Function(
                FUNC_SETPARENTNAMED, 
                Arrays.<Type>asList(new Utf8String(name)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<AbstractContractFSubclass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractFSubclass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AbstractContractFSubclass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractFSubclass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AbstractContractFSubclass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractFSubclass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AbstractContractFSubclass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractFSubclass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
