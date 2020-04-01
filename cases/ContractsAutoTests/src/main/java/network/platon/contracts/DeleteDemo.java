package network.platon.contracts;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.*;
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

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class DeleteDemo extends Contract {
    private static final String BINARY = "608060405260016000806101000a81548160ff0219169083151502179055506001805533600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040518060400160405280600381526020017f313233000000000000000000000000000000000000000000000000000000000081525060039080519060200190620000b092919062000138565b506040518060400160405280600381526020017f616263000000000000000000000000000000000000000000000000000000000081525060049080519060200190620000fe929190620001bf565b506001600560006101000a81548160ff021916908360028111156200011f57fe5b02179055503480156200013157600080fd5b506200026e565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200017b57805160ff1916838001178555620001ac565b82800160010185558215620001ac579182015b82811115620001ab5782518255916020019190600101906200018e565b5b509050620001bb919062000246565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200020257805160ff191683800117855562000233565b8280016001018555821562000233579182015b828111156200023257825182559160200191906001019062000215565b5b50905062000242919062000246565b5090565b6200026b91905b80821115620002675760008160009055506001016200024d565b5090565b90565b610e5d806200027e6000396000f3fe6080604052600436106101095760003560e01c8063767800de11610095578063c15bae8411610064578063c15bae8414610514578063cf08fed5146105a4578063d1bdda41146105dd578063e5aa3d5814610667578063f02997491461069257610109565b8063767800de1461037257806393e1ed83146103c9578063a1a984e514610459578063ab5170b21461048457610109565b806327c58232116100dc57806327c582321461026757806332d057c91461027e5780633ab0698c146102885780634df7e3d0146102b35780635d743b5d146102e257610109565b806305be2c121461010e57806313a5a8af146101a55780631acddabe146101de578063252bd4d314610210575b600080fd5b34801561011a57600080fd5b506101236106c1565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561016957808201518184015260208101905061014e565b50505050905090810190601f1680156101965780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156101b157600080fd5b506101ba610774565b604051808260028111156101ca57fe5b60ff16815260200191505060405180910390f35b3480156101ea57600080fd5b506101f361078b565b604051808381526020018281526020019250505060405180910390f35b34801561021c57600080fd5b506102256107df565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561027357600080fd5b5061027c610809565b005b610286610875565b005b34801561029457600080fd5b5061029d6108ed565b6040518082815260200191505060405180910390f35b3480156102bf57600080fd5b506102c8610966565b604051808215151515815260200191505060405180910390f35b3480156102ee57600080fd5b506102f7610978565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561033757808201518184015260208101905061031c565b50505050905090810190601f1680156103645780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561037e57600080fd5b50610387610a1a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103d557600080fd5b506103de610a40565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561041e578082015181840152602081019050610403565b50505050905090810190601f16801561044b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561046557600080fd5b5061046e610ade565b6040518082815260200191505060405180910390f35b34801561049057600080fd5b50610499610ae8565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104d95780820151818401526020810190506104be565b50505050905090810190601f1680156105065780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561052057600080fd5b50610529610b8a565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561056957808201518184015260208101905061054e565b50505050905090810190601f1680156105965780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156105b057600080fd5b506105b9610c28565b604051808260028111156105c957fe5b60ff16815260200191505060405180910390f35b6105e5610c3b565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561062b578082015181840152602081019050610610565b50505050905090810190601f1680156106585780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561067357600080fd5b5061067c610d57565b6040518082815260200191505060405180910390f35b34801561069e57600080fd5b506106a7610d5d565b604051808215151515815260200191505060405180910390f35b600060606006600001546006600101808054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107655780601f1061073a57610100808354040283529160200191610765565b820191906000526020600020905b81548152906001019060200180831161074857829003601f168201915b50505050509050915091509091565b6000600560009054906101000a900460ff16905090565b600080600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600860010154915091509091565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600360006108539190610d73565b600460006108619190610dbb565b600560006101000a81549060ff0219169055565b604051806020016040528060c88152506008600082015181600101559050506107d0600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506008600060018201600090555050565b6000606060076040519080825280602002602001820160405280156109215781602001602082028038833980820191505090505b50905060648160008151811061093357fe5b60200260200101818152505060c88160018151811061094e57fe5b60200260200101818152505060609050805191505090565b6000809054906101000a900460ff1681565b606060038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a105780601f106109e557610100808354040283529160200191610a10565b820191906000526020600020905b8154815290600101906020018083116109f357829003601f168201915b5050505050905090565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ad65780601f10610aab57610100808354040283529160200191610ad6565b820191906000526020600020905b815481529060010190602001808311610ab957829003601f168201915b505050505081565b6000600154905090565b606060048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b805780601f10610b5557610100808354040283529160200191610b80565b820191906000526020600020905b815481529060010190602001808311610b6357829003601f168201915b5050505050905090565b60048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c205780601f10610bf557610100808354040283529160200191610c20565b820191906000526020600020905b815481529060010190602001808311610c0357829003601f168201915b505050505081565b600560009054906101000a900460ff1681565b600060606040518060400160405280600a81526020016040518060400160405280600381526020017f6162630000000000000000000000000000000000000000000000000000000000815250815250506006600080820160009055600182016000610ca69190610dbb565b50506006600001546006600101808054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d485780601f10610d1d57610100808354040283529160200191610d48565b820191906000526020600020905b815481529060010190602001808311610d2b57829003601f168201915b50505050509050915091509091565b60015481565b60008060009054906101000a900460ff16905090565b50805460018160011615610100020316600290046000825580601f10610d995750610db8565b601f016020900490600052602060002090810190610db79190610e03565b5b50565b50805460018160011615610100020316600290046000825580601f10610de15750610e00565b601f016020900490600052602060002090810190610dff9190610e03565b5b50565b610e2591905b80821115610e21576000816000905550600101610e09565b5090565b9056fea265627a7a72315820a3c30cf8578fb5aa065e040716de2b212c44fc2ee549d527be70966ac0392c1364736f6c634300050d0032";

    public static final String FUNC_ADDR = "addr";

    public static final String FUNC_B = "b";

    public static final String FUNC_COLOR = "color";

    public static final String FUNC_DELDYNAMICARRAY = "delDynamicArray";

    public static final String FUNC_DELMAPPING = "delMapping";

    public static final String FUNC_DELSTRUCT = "delStruct";

    public static final String FUNC_DELETEATTR = "deleteAttr";

    public static final String FUNC_GETADDRESS = "getaddress";

    public static final String FUNC_GETBOOL = "getbool";

    public static final String FUNC_GETBYTES = "getbytes";

    public static final String FUNC_GETDELMAPPING = "getdelMapping";

    public static final String FUNC_GETENUM = "getenum";

    public static final String FUNC_GETSTR = "getstr";

    public static final String FUNC_GETSTRUCT = "getstruct";

    public static final String FUNC_GETUNIT = "getunit";

    public static final String FUNC_I = "i";

    public static final String FUNC_STR = "str";

    public static final String FUNC_VARBYTE = "varByte";

    @Deprecated
    protected DeleteDemo(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected DeleteDemo(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected DeleteDemo(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected DeleteDemo(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> addr() {
        final Function function = new Function(FUNC_ADDR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Boolean> b() {
        final Function function = new Function(FUNC_B, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<BigInteger> color() {
        final Function function = new Function(FUNC_COLOR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> delDynamicArray() {
        final Function function = new Function(FUNC_DELDYNAMICARRAY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> delMapping(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_DELMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<TransactionReceipt> delStruct(BigInteger weiValue) {
        final Function function = new Function(
                FUNC_DELSTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, weiValue);
    }

    public RemoteCall<TransactionReceipt> deleteAttr() {
        final Function function = new Function(
                FUNC_DELETEATTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getaddress() {
        final Function function = new Function(FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Boolean> getbool() {
        final Function function = new Function(FUNC_GETBOOL, 
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

    public RemoteCall<BigInteger> getenum() {
        final Function function = new Function(FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getstr() {
        final Function function = new Function(FUNC_GETSTR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
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

    public RemoteCall<BigInteger> getunit() {
        final Function function = new Function(FUNC_GETUNIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> i() {
        final Function function = new Function(FUNC_I, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> str() {
        final Function function = new Function(FUNC_STR, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> varByte() {
        final Function function = new Function(FUNC_VARBYTE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(DeleteDemo.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(DeleteDemo.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(DeleteDemo.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<DeleteDemo> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(DeleteDemo.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static DeleteDemo load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new DeleteDemo(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static DeleteDemo load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new DeleteDemo(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static DeleteDemo load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new DeleteDemo(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static DeleteDemo load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new DeleteDemo(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
