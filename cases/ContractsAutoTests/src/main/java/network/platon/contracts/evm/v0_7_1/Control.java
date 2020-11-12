package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.1.
 */
public class Control extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610c18806100206000396000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c806357609889116100b85780638e418fdb1161007c5780638e418fdb146104b2578063a64be0d5146104d0578063b4feac7c146104ee578063b87df0141461050c578063c0e641fc1461052a578063da193c1f1461054857610137565b80635760988914610352578063687615d71461037057806371ee52021461038e57806378aa6155146104115780637e6b0f571461042f57610137565b806344e24ce0116100ff57806344e24ce01461029c57806347808fc3146102ca5780634b8016b9146102f8578063508242dc1461031657806356230cca1461033457610137565b80631f9c9f3c1461013c578063275ec9761461015a57806335432d3114610178578063383d49e5146101fb5780633f9dbcf914610219575b600080fd5b610144610566565b6040518082815260200191505060405180910390f35b61016261056c565b6040518082815260200191505060405180910390f35b6101806105ca565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101c05780820151818401526020810190506101a5565b50505050905090810190601f1680156101ed5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610203610668565b6040518082815260200191505060405180910390f35b61022161066e565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610261578082015181840152602081019050610246565b50505050905090810190601f16801561028e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102c8600480360360208110156102b257600080fd5b810190808035906020019092919050505061070c565b005b6102f6600480360360208110156102e057600080fd5b8101908080359060200190929190505050610811565b005b6103006108a4565b6040518082815260200191505060405180910390f35b61031e6108aa565b6040518082815260200191505060405180910390f35b61033c6108b0565b6040518082815260200191505060405180910390f35b61035a610904565b6040518082815260200191505060405180910390f35b61037861090e565b6040518082815260200191505060405180910390f35b610396610914565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103d65780820151818401526020810190506103bb565b50505050905090810190601f1680156104035780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6104196109b6565b6040518082815260200191505060405180910390f35b6104376109c0565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561047757808201518184015260208101905061045c565b50505050905090810190601f1680156104a45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6104ba610a62565b6040518082815260200191505060405180910390f35b6104d8610a95565b6040518082815260200191505060405180910390f35b6104f6610aec565b6040518082815260200191505060405180910390f35b610514610b27565b6040518082815260200191505060405180910390f35b610532610b31565b6040518082815260200191505060405180910390f35b610550610b3b565b6040518082815260200191505060405180910390f35b60025481565b6000806005819055506000600190505b600a8110156105c05760006005828161059157fe5b0614156105a3576005549150506105c7565b80600560008282540192505081905550808060010191505061057c565b5060055490505b90565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106605780601f1061063557610100808354040283529160200191610660565b820191906000526020600020905b81548152906001019060200180831161064357829003601f168201915b505050505081565b60035481565b60068054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107045780601f106106d957610100808354040283529160200191610704565b820191906000526020600020905b8154815290600101906020018083116106e757829003601f168201915b505050505081565b6014811015610766576040518060400160405280601381526020017f796f7520617265206120796f756e67206d616e0000000000000000000000000081525060009080519060200190610760929190610b45565b5061080e565b603c8110156107c0576040518060400160405280601481526020017f796f75206172652061206d6964646c65206d616e000000000000000000000000815250600090805190602001906107ba929190610b45565b5061080d565b6040518060400160405280601181526020017f796f75206172652061206f6c64206d616e0000000000000000000000000000008152506000908051906020019061080b929190610b45565b505b5b50565b60148113610854576040518060400160405280600c81526020017f6d6f7265207468616e203230000000000000000000000000000000000000000081525061088b565b6040518060400160405280600c81526020017f6c657373207468616e20323000000000000000000000000000000000000000008152505b600690805190602001906108a0929190610b45565b5050565b60045481565b60015481565b60008060048190555060005b600a8110156108fb576000600282816108d157fe5b0614156108dd576108ee565b806004600082825401925050819055505b80806001019150506108bc565b50600454905090565b6000600454905090565b60055481565b606060068054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109ac5780601f10610981576101008083540402835291602001916109ac565b820191906000526020600020905b81548152906001019060200180831161098f57829003601f168201915b5050505050905090565b6000600554905090565b606060008054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a585780601f10610a2d57610100808354040283529160200191610a58565b820191906000526020600020905b815481529060010190602001808311610a3b57829003601f168201915b5050505050905090565b60008060018190555060005b80600160008282540192505081905550806001019050600a8110610a6e5760015491505090565b6000806003819055506000600190505b600a811015610ae357600060028281610aba57fe5b061415610ac657610ae3565b806003600082825401925050819055508080600101915050610aa5565b50600354905090565b60008060028190555060005b600a811015610b1e57806002600082825401925050819055508080600101915050610af8565b50600254905090565b6000600254905090565b6000600354905090565b6000600154905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b8657805160ff1916838001178555610bb4565b82800160010185558215610bb4579182015b82811115610bb3578251825591602001919060010190610b98565b5b509050610bc19190610bc5565b5090565b5b80821115610bde576000816000905550600101610bc6565b509056fea2646970667358221220223435d63c2199af99d1fd2ab1cf1339bedb5c8c871b6ac5da1ef7eb501025c064736f6c63430007010033";

    public static final String FUNC_DOWHILECONTROL = "doWhileControl";

    public static final String FUNC_DOWHILECONTROLRESULT = "doWhileControlResult";

    public static final String FUNC_FORBREAKCONTROL = "forBreakControl";

    public static final String FUNC_FORBREAKCONTROLRESULT = "forBreakControlResult";

    public static final String FUNC_FORCONTINUECONTROL = "forContinueControl";

    public static final String FUNC_FORCONTINUECONTROLRESULT = "forContinueControlResult";

    public static final String FUNC_FORCONTROL = "forControl";

    public static final String FUNC_FORCONTROLRESULT = "forControlResult";

    public static final String FUNC_FORRETURNCONTROL = "forReturnControl";

    public static final String FUNC_FORRETURNCONTROLRESULT = "forReturnControlResult";

    public static final String FUNC_FORTHREECONTROLCONTROL = "forThreeControlControl";

    public static final String FUNC_FORTHREECONTROLCONTROLRESULT = "forThreeControlControlResult";

    public static final String FUNC_GETFORBREAKCONTROLRESULT = "getForBreakControlResult";

    public static final String FUNC_GETFORCONTINUECONTROLRESULT = "getForContinueControlResult";

    public static final String FUNC_GETFORCONTROLRESULT = "getForControlResult";

    public static final String FUNC_GETFORRETURNCONTROLRESULT = "getForReturnControlResult";

    public static final String FUNC_GETFORTHREECONTROLCONTROLRESULT = "getForThreeControlControlResult";

    public static final String FUNC_GETIFCONTROLRESULT = "getIfControlResult";

    public static final String FUNC_GETDOWHILERESULT = "getdoWhileResult";

    public static final String FUNC_IFCONTROL = "ifControl";

    public static final String FUNC_IFCONTROLRESULT = "ifControlResult";

    protected Control(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Control(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> doWhileControl() {
        final Function function = new Function(
                FUNC_DOWHILECONTROL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> doWhileControlResult() {
        final Function function = new Function(FUNC_DOWHILECONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> forBreakControl() {
        final Function function = new Function(
                FUNC_FORBREAKCONTROL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> forBreakControlResult() {
        final Function function = new Function(FUNC_FORBREAKCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> forContinueControl() {
        final Function function = new Function(
                FUNC_FORCONTINUECONTROL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> forContinueControlResult() {
        final Function function = new Function(FUNC_FORCONTINUECONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> forControl() {
        final Function function = new Function(
                FUNC_FORCONTROL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> forControlResult() {
        final Function function = new Function(FUNC_FORCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> forReturnControl() {
        final Function function = new Function(
                FUNC_FORRETURNCONTROL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> forReturnControlResult() {
        final Function function = new Function(FUNC_FORRETURNCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> forThreeControlControl(BigInteger age) {
        final Function function = new Function(
                FUNC_FORTHREECONTROLCONTROL, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(age)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> forThreeControlControlResult() {
        final Function function = new Function(FUNC_FORTHREECONTROLCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getForBreakControlResult() {
        final Function function = new Function(FUNC_GETFORBREAKCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getForContinueControlResult() {
        final Function function = new Function(FUNC_GETFORCONTINUECONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getForControlResult() {
        final Function function = new Function(FUNC_GETFORCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getForReturnControlResult() {
        final Function function = new Function(FUNC_GETFORRETURNCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getForThreeControlControlResult() {
        final Function function = new Function(FUNC_GETFORTHREECONTROLCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getIfControlResult() {
        final Function function = new Function(FUNC_GETIFCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getdoWhileResult() {
        final Function function = new Function(FUNC_GETDOWHILERESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> ifControl(BigInteger age) {
        final Function function = new Function(
                FUNC_IFCONTROL, 
                Arrays.<Type>asList(new Uint256(age)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> ifControlResult() {
        final Function function = new Function(FUNC_IFCONTROLRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<Control> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Control.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Control> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Control.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Control load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Control(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Control load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Control(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
