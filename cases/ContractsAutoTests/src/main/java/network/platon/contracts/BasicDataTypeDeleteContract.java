package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Bytes3;
import org.web3j.abi.datatypes.generated.Int256;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tuples.generated.Tuple6;
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
public class BasicDataTypeDeleteContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610b07806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063b2548ecc11610071578063b2548ecc1461017a578063c726046714610290578063d5e350e31461029a578063d8c9d5f5146102be578063e9cafac2146102c8578063f0ebce5a146102d2576100a9565b80630849cc99146100ae5780630860dca9146100cc57806309b1b3f2146100d65780630d35d126146101665780633edf92a814610170575b600080fd5b6100b66102fe565b6040518082815260200191505060405180910390f35b6100d461030b565b005b6100de61031f565b604051808360ff1660ff16815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561012a57808201518184015260208101905061010f565b50505050905090810190601f1680156101575780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61016e6103df565b005b610178610655565b005b61018261067e565b60405180871515151581526020018681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001847cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200180602001838152602001828103825284818151815260200191508051906020019080838360005b83811015610250578082015181840152602081019050610235565b50505050905090810190601f16801561027d5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b610298610782565b005b6102a26107e8565b604051808260ff1660ff16815260200191505060405180910390f35b6102c6610814565b005b6102d061083f565b005b6102da61084f565b604051808260028111156102ea57fe5b60ff16815260200191505060405180910390f35b6000600580549050905090565b600660006101000a81549060ff0219169055565b60006060600760000160009054906101000a900460ff166007600101808054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103d05780601f106103a5576101008083540402835291602001916103d0565b820191906000526020600020905b8154815290600101906020018083116103b357829003601f168201915b50505050509050915091509091565b60016000806101000a81548160ff021916908315150217905550600260018190555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3132330000000000000000000000000000000000000000000000000000000000600260146101000a81548162ffffff021916908360e81c02179055506040518060400160405280600581526020017f68656c6c6f000000000000000000000000000000000000000000000000000000815250600390805190602001906104ca929190610866565b5060056004819055506040518060400160405280600160ff1681526020016040518060400160405280600481526020017f456c6c6100000000000000000000000000000000000000000000000000000000815250815250600760008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160010190805190602001906105609291906108e6565b509050506040518060600160405280600160ff168152602001600260ff168152602001600360ff16815250600590600361059b929190610966565b50600160096000600160ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600260096000600260ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600360096000600360ff16815260200190815260200160002060006101000a81548160ff021916908360ff1602179055506002600660006101000a81548160ff0219169083600281111561064e57fe5b0217905550565b60096000600260ff16815260200190815260200160002060006101000a81549060ff0219169055565b600080600080606060008060009054906101000a900460ff16600154600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260149054906101000a900460e81b6003600454818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107675780601f1061073c57610100808354040283529160200191610767565b820191906000526020600020905b81548152906001019060200180831161074a57829003601f168201915b50505050509150955095509550955095509550909192939495565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600260146101000a81549062ffffff0219169055600360006107e09190610a0d565b600460009055565b600060096000600260ff16815260200190815260200160002060009054906101000a900460ff16905090565b6007600080820160006101000a81549060ff021916905560018201600061083b9190610a0d565b5050565b6005600061084d9190610a55565b565b6000600660009054906101000a900460ff16905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108a757805160ff19168380011785556108d5565b828001600101855582156108d5579182015b828111156108d45782518255916020019190600101906108b9565b5b5090506108e29190610a7d565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061092757805160ff1916838001178555610955565b82800160010185558215610955579182015b82811115610954578251825591602001919060010190610939565b5b5090506109629190610a7d565b5090565b82805482825590600052602060002090601f016020900481019282156109fc5791602002820160005b838211156109cd57835183826101000a81548160ff021916908360ff160217905550926020019260010160208160000104928301926001030261098f565b80156109fa5782816101000a81549060ff02191690556001016020816000010492830192600103026109cd565b505b509050610a099190610aa2565b5090565b50805460018160011615610100020316600290046000825580601f10610a335750610a52565b601f016020900490600052602060002090810190610a519190610a7d565b5b50565b50805460008255601f016020900490600052602060002090810190610a7a9190610a7d565b50565b610a9f91905b80821115610a9b576000816000905550600101610a83565b5090565b90565b610acf91905b80821115610acb57600081816101000a81549060ff021916905550600101610aa8565b5090565b9056fea265627a7a723158208e6621ad0d116538240a4a198d7c51c9984ad3fd611c51b1daf871fa42a9e70a64736f6c634300050d0032";

    public static final String FUNC_DELETEARRAY = "deleteArray";

    public static final String FUNC_DELETEBASICDATA = "deleteBasicData";

    public static final String FUNC_DELETEENUM = "deleteEnum";

    public static final String FUNC_DELETEMAPPING = "deleteMapping";

    public static final String FUNC_DELETESTRUCT = "deleteStruct";

    public static final String FUNC_GETARRAYLENGTH = "getArrayLength";

    public static final String FUNC_GETBASICDATA = "getBasicData";

    public static final String FUNC_GETENUM = "getEnum";

    public static final String FUNC_GETMAPPING = "getMapping";

    public static final String FUNC_GETSTRUCT = "getStruct";

    public static final String FUNC_INITBASICDATA = "initBasicData";

    @Deprecated
    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> deleteArray() {
        final Function function = new Function(
                FUNC_DELETEARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteBasicData() {
        final Function function = new Function(
                FUNC_DELETEBASICDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteEnum() {
        final Function function = new Function(
                FUNC_DELETEENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteMapping() {
        final Function function = new Function(
                FUNC_DELETEMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteStruct() {
        final Function function = new Function(
                FUNC_DELETESTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getArrayLength() {
        final Function function = new Function(FUNC_GETARRAYLENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger>> getBasicData() {
        final Function function = new Function(FUNC_GETBASICDATA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}, new TypeReference<Uint256>() {}, new TypeReference<Address>() {}, new TypeReference<Bytes3>() {}, new TypeReference<Utf8String>() {}, new TypeReference<Int256>() {}));
        return new RemoteCall<Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger>>(
                new Callable<Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger>>() {
                    @Override
                    public Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple6<Boolean, BigInteger, String, byte[], String, BigInteger>(
                                (Boolean) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (String) results.get(2).getValue(), 
                                (byte[]) results.get(3).getValue(), 
                                (String) results.get(4).getValue(), 
                                (BigInteger) results.get(5).getValue());
                    }
                });
    }

    public RemoteCall<BigInteger> getEnum() {
        final Function function = new Function(FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getMapping() {
        final Function function = new Function(FUNC_GETMAPPING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<BigInteger, String>> getStruct() {
        final Function function = new Function(FUNC_GETSTRUCT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}, new TypeReference<Utf8String>() {}));
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

    public RemoteCall<TransactionReceipt> initBasicData() {
        final Function function = new Function(
                FUNC_INITBASICDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
