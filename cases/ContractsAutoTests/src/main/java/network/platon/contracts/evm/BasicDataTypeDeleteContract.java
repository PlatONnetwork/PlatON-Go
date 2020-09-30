package network.platon.contracts.evm;

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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.1.5.
 */
public class BasicDataTypeDeleteContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610bc9806100206000396000f3fe6080604052600436106100a9576000357c0100000000000000000000000000000000000000000000000000000000900480630849cc99146100ae5780630860dca9146100d957806309b1b3f2146100f05780630d35d1261461018d5780633edf92a8146101a4578063b2548ecc146101bb578063c7260467146102de578063d5e350e3146102f5578063d8c9d5f514610326578063e9cafac21461033d578063f0ebce5a14610354575b600080fd5b3480156100ba57600080fd5b506100c361038d565b6040518082815260200191505060405180910390f35b3480156100e557600080fd5b506100ee61039a565b005b3480156100fc57600080fd5b506101056103ae565b604051808360ff1660ff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610151578082015181840152602081019050610136565b50505050905090810190601f16801561017e5780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561019957600080fd5b506101a261046e565b005b3480156101b057600080fd5b506101b9610703565b005b3480156101c757600080fd5b506101d061072c565b60405180871515151581526020018681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001847cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200180602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561029e578082015181840152602081019050610283565b50505050905090810190601f1680156102cb5780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b3480156102ea57600080fd5b506102f361084d565b005b34801561030157600080fd5b5061030a6108b3565b604051808260ff1660ff16815260200191505060405180910390f35b34801561033257600080fd5b5061033b6108df565b005b34801561034957600080fd5b5061035261090a565b005b34801561036057600080fd5b5061036961091a565b6040518082600281111561037957fe5b60ff16815260200191505060405180910390f35b6000600580549050905090565b600660006101000a81549060ff0219169055565b60006060600760000160009054906101000a900460ff166007600101808054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561045f5780601f106104345761010080835404028352916020019161045f565b820191906000526020600020905b81548152906001019060200180831161044257829003601f168201915b50505050509050915091509091565b60016000806101000a81548160ff021916908315150217905550600260018190555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3132330000000000000000000000000000000000000000000000000000000000600260146101000a81548162ffffff02191690837d010000000000000000000000000000000000000000000000000000000000900402179055506040805190810160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525060039080519060200190610577929190610931565b5060056004819055506040805190810160405280600160ff1681526020016040805190810160405280600481526020017f456c6c6100000000000000000000000000000000000000000000000000000000815250815250600760008201518160000160006101000a81548160ff021916908360ff160217905550602082015181600101908051906020019061060d9291906109b1565b50905050606060405190810160405280600160ff168152602001600260ff168152602001600360ff168152506005906003610649929190610a31565b50600160096000600160ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600260096000600260ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600360096000600360ff16815260200190815260200160002060006101000a81548160ff021916908360ff1602179055506002600660006101000a81548160ff021916908360028111156106fc57fe5b0217905550565b60096000600260ff16815260200190815260200160002060006101000a81549060ff0219169055565b600080600080606060008060009054906101000a900460ff16600154600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260149054906101000a90047d010000000000000000000000000000000000000000000000000000000000026003600454818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108325780601f1061080757610100808354040283529160200191610832565b820191906000526020600020905b81548152906001019060200180831161081557829003601f168201915b50505050509150955095509550955095509550909192939495565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600260146101000a81549062ffffff0219169055600360006108ab9190610ad8565b600460009055565b600060096000600260ff16815260200190815260200160002060009054906101000a900460ff16905090565b6007600080820160006101000a81549060ff02191690556001820160006109069190610ad8565b5050565b600560006109189190610b20565b565b6000600660009054906101000a900460ff16905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061097257805160ff19168380011785556109a0565b828001600101855582156109a0579182015b8281111561099f578251825591602001919060010190610984565b5b5090506109ad9190610b48565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106109f257805160ff1916838001178555610a20565b82800160010185558215610a20579182015b82811115610a1f578251825591602001919060010190610a04565b5b509050610a2d9190610b48565b5090565b82805482825590600052602060002090601f01602090048101928215610ac75791602002820160005b83821115610a9857835183826101000a81548160ff021916908360ff1602179055509260200192600101602081600001049283019260010302610a5a565b8015610ac55782816101000a81549060ff0219169055600101602081600001049283019260010302610a98565b505b509050610ad49190610b6d565b5090565b50805460018160011615610100020316600290046000825580601f10610afe5750610b1d565b601f016020900490600052602060002090810190610b1c9190610b48565b5b50565b50805460008255601f016020900490600052602060002090810190610b459190610b48565b50565b610b6a91905b80821115610b66576000816000905550600101610b4e565b5090565b90565b610b9a91905b80821115610b9657600081816101000a81549060ff021916905550600101610b73565b5090565b9056fea165627a7a723058200ce8b59af2ee5b501a1f954ec7c2a2126a54406ef88ec161f087f022ed96dacd0029";

    public static final String FUNC_GETARRAYLENGTH = "getArrayLength";

    public static final String FUNC_DELETEENUM = "deleteEnum";

    public static final String FUNC_GETSTRUCT = "getStruct";

    public static final String FUNC_INITBASICDATA = "initBasicData";

    public static final String FUNC_DELETEMAPPING = "deleteMapping";

    public static final String FUNC_GETBASICDATA = "getBasicData";

    public static final String FUNC_DELETEBASICDATA = "deleteBasicData";

    public static final String FUNC_GETMAPPING = "getMapping";

    public static final String FUNC_DELETESTRUCT = "deleteStruct";

    public static final String FUNC_DELETEARRAY = "deleteArray";

    public static final String FUNC_GETENUM = "getEnum";

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getArrayLength() {
        final Function function = new Function(FUNC_GETARRAYLENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> deleteEnum() {
        final Function function = new Function(
                FUNC_DELETEENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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

    public RemoteCall<TransactionReceipt> deleteMapping() {
        final Function function = new Function(
                FUNC_DELETEMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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

    public RemoteCall<TransactionReceipt> deleteBasicData() {
        final Function function = new Function(
                FUNC_DELETEBASICDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getMapping() {
        final Function function = new Function(FUNC_GETMAPPING, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> deleteStruct() {
        final Function function = new Function(
                FUNC_DELETESTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteArray() {
        final Function function = new Function(
                FUNC_DELETEARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getEnum() {
        final Function function = new Function(FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BasicDataTypeDeleteContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeDeleteContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BasicDataTypeDeleteContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeDeleteContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
