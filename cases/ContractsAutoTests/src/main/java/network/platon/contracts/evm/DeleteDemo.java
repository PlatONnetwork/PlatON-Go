package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
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
public class DeleteDemo extends Contract {
    private static final String BINARY = "608060405260016000806101000a81548160ff0219169083151502179055506001805533600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040805190810160405280600381526020017f313233000000000000000000000000000000000000000000000000000000000081525060039080519060200190620000b092919062000138565b506040805190810160405280600381526020017f616263000000000000000000000000000000000000000000000000000000000081525060049080519060200190620000fe929190620001bf565b506001600560006101000a81548160ff021916908360028111156200011f57fe5b02179055503480156200013157600080fd5b506200026e565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200017b57805160ff1916838001178555620001ac565b82800160010185558215620001ac579182015b82811115620001ab5782518255916020019190600101906200018e565b5b509050620001bb919062000246565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200020257805160ff191683800117855562000233565b8280016001018555821562000233579182015b828111156200023257825182559160200191906001019062000215565b5b50905062000242919062000246565b5090565b6200026b91905b80821115620002675760008160009055506001016200024d565b5090565b90565b610e7a806200027e6000396000f3fe608060405260043610610126576000357c010000000000000000000000000000000000000000000000000000000090048063767800de116100b2578063c15bae8411610081578063c15bae8414610531578063cf08fed5146105c1578063d1bdda41146105fa578063e5aa3d5814610684578063f0299749146106af57610126565b8063767800de1461038f57806393e1ed83146103e6578063a1a984e514610476578063ab5170b2146104a157610126565b806327c58232116100f957806327c582321461028457806332d057c91461029b5780633ab0698c146102a55780634df7e3d0146102d05780635d743b5d146102ff57610126565b806305be2c121461012b57806313a5a8af146101c25780631acddabe146101fb578063252bd4d31461022d575b600080fd5b34801561013757600080fd5b506101406106de565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561018657808201518184015260208101905061016b565b50505050905090810190601f1680156101b35780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156101ce57600080fd5b506101d7610791565b604051808260028111156101e757fe5b60ff16815260200191505060405180910390f35b34801561020757600080fd5b506102106107a8565b604051808381526020018281526020019250505060405180910390f35b34801561023957600080fd5b506102426107fc565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561029057600080fd5b50610299610826565b005b6102a3610892565b005b3480156102b157600080fd5b506102ba61090b565b6040518082815260200191505060405180910390f35b3480156102dc57600080fd5b506102e561098c565b604051808215151515815260200191505060405180910390f35b34801561030b57600080fd5b5061031461099e565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610354578082015181840152602081019050610339565b50505050905090810190601f1680156103815780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561039b57600080fd5b506103a4610a40565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103f257600080fd5b506103fb610a66565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561043b578082015181840152602081019050610420565b50505050905090810190601f1680156104685780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561048257600080fd5b5061048b610b04565b6040518082815260200191505060405180910390f35b3480156104ad57600080fd5b506104b6610b0e565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104f65780820151818401526020810190506104db565b50505050905090810190601f1680156105235780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561053d57600080fd5b50610546610bb0565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561058657808201518184015260208101905061056b565b50505050905090810190601f1680156105b35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156105cd57600080fd5b506105d6610c4e565b604051808260028111156105e657fe5b60ff16815260200191505060405180910390f35b610602610c61565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561064857808201518184015260208101905061062d565b50505050905090810190601f1680156106755780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561069057600080fd5b50610699610d7d565b6040518082815260200191505060405180910390f35b3480156106bb57600080fd5b506106c4610d83565b604051808215151515815260200191505060405180910390f35b600060606006600001546006600101808054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107825780601f1061075757610100808354040283529160200191610782565b820191906000526020600020905b81548152906001019060200180831161076557829003601f168201915b50505050509050915091509091565b6000600560009054906101000a900460ff16905090565b600080600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600860010154915091509091565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600360006108709190610d99565b6004600061087e9190610de1565b600560006101000a81549060ff0219169055565b60206040519081016040528060c88152506008600082015181600101559050506107d0600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506008600060018201600090555050565b60006060600760405190808252806020026020018201604052801561093f5781602001602082028038833980820191505090505b509050606481600081518110151561095357fe5b906020019060200201818152505060c881600181518110151561097257fe5b906020019060200201818152505060609050805191505090565b6000809054906101000a900460ff1681565b606060038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a365780601f10610a0b57610100808354040283529160200191610a36565b820191906000526020600020905b815481529060010190602001808311610a1957829003601f168201915b5050505050905090565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610afc5780601f10610ad157610100808354040283529160200191610afc565b820191906000526020600020905b815481529060010190602001808311610adf57829003601f168201915b505050505081565b6000600154905090565b606060048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ba65780601f10610b7b57610100808354040283529160200191610ba6565b820191906000526020600020905b815481529060010190602001808311610b8957829003601f168201915b5050505050905090565b60048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c465780601f10610c1b57610100808354040283529160200191610c46565b820191906000526020600020905b815481529060010190602001808311610c2957829003601f168201915b505050505081565b600560009054906101000a900460ff1681565b600060606040805190810160405280600a81526020016040805190810160405280600381526020017f6162630000000000000000000000000000000000000000000000000000000000815250815250506006600080820160009055600182016000610ccc9190610de1565b50506006600001546006600101808054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d6e5780601f10610d4357610100808354040283529160200191610d6e565b820191906000526020600020905b815481529060010190602001808311610d5157829003601f168201915b50505050509050915091509091565b60015481565b60008060009054906101000a900460ff16905090565b50805460018160011615610100020316600290046000825580601f10610dbf5750610dde565b601f016020900490600052602060002090810190610ddd9190610e29565b5b50565b50805460018160011615610100020316600290046000825580601f10610e075750610e26565b601f016020900490600052602060002090810190610e259190610e29565b5b50565b610e4b91905b80821115610e47576000816000905550600101610e2f565b5090565b9056fea165627a7a72305820c09e42f9056c0da44419a5d14f78c24ec9c206ef611a936800d20e586cc58b590029";

    public static final String FUNC_GETSTRUCT = "getstruct";

    public static final String FUNC_GETENUM = "getenum";

    public static final String FUNC_GETDELMAPPING = "getdelMapping";

    public static final String FUNC_GETADDRESS = "getaddress";

    public static final String FUNC_DELETEATTR = "deleteAttr";

    public static final String FUNC_DELMAPPING = "delMapping";

    public static final String FUNC_DELDYNAMICARRAY = "delDynamicArray";

    public static final String FUNC_B = "b";

    public static final String FUNC_GETBYTES = "getbytes";

    public static final String FUNC_ADDR = "addr";

    public static final String FUNC_VARBYTE = "varByte";

    public static final String FUNC_GETUNIT = "getunit";

    public static final String FUNC_GETSTR = "getstr";

    public static final String FUNC_STR = "str";

    public static final String FUNC_COLOR = "color";

    public static final String FUNC_DELSTRUCT = "delStruct";

    public static final String FUNC_I = "i";

    public static final String FUNC_GETBOOL = "getbool";

    protected DeleteDemo(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DeleteDemo(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple2<BigInteger, String>> getstruct() {
        final Function function = new Function(FUNC_GETSTRUCT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Utf8String>() {}));
        return new RemoteCall<Tuple2<BigInteger, String>>(
                new Callable<Tuple2<BigInteger, String>>() {
                    @Override
                    public Tuple2<BigInteger, String> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, String>(
                                (BigInteger) results.get(0).getValue(), 
                                (String) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<BigInteger> getenum() {
        final Function function = new Function(FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> getdelMapping() {
        final Function function = new Function(FUNC_GETDELMAPPING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<String> getaddress() {
        final Function function = new Function(FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> deleteAttr() {
        final Function function = new Function(
                FUNC_DELETEATTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> delMapping(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_DELMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> delDynamicArray() {
        final Function function = new Function(FUNC_DELDYNAMICARRAY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> b() {
        final Function function = new Function(FUNC_B, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<byte[]> getbytes() {
        final Function function = new Function(FUNC_GETBYTES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> addr() {
        final Function function = new Function(FUNC_ADDR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> varByte() {
        final Function function = new Function(FUNC_VARBYTE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getunit() {
        final Function function = new Function(FUNC_GETUNIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getstr() {
        final Function function = new Function(FUNC_GETSTR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> str() {
        final Function function = new Function(FUNC_STR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> color() {
        final Function function = new Function(FUNC_COLOR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> delStruct(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_DELSTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> i() {
        final Function function = new Function(FUNC_I, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> getbool() {
        final Function function = new Function(FUNC_GETBOOL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DeleteDemo.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DeleteDemo.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DeleteDemo load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DeleteDemo(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DeleteDemo load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DeleteDemo(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
