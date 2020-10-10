package network.platon.contracts.evm.v0_6_12;

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
public class BasicDataTypeDeleteContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610a37806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063b2548ecc11610071578063b2548ecc14610177578063c726046714610255578063d5e350e31461025f578063d8c9d5f514610280578063e9cafac21461028a578063f0ebce5a14610294576100a9565b80630849cc99146100ae5780630860dca9146100cc57806309b1b3f2146100d65780630d35d126146101635780633edf92a81461016d575b600080fd5b6100b66102bd565b6040518082815260200191505060405180910390f35b6100d46102ca565b005b6100de6102de565b604051808360ff16815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561012757808201518184015260208101905061010c565b50505050905090810190601f1680156101545780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b61016b61039e565b005b610175610614565b005b61017f61063d565b6040518087151581526020018681526020018573ffffffffffffffffffffffffffffffffffffffff168152602001847cffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200180602001838152602001828103825284818151815260200191508051906020019080838360005b838110156102155780820151818401526020810190506101fa565b50505050905090810190601f1680156102425780820380516001836020036101000a031916815260200191505b5097505050505050505060405180910390f35b61025d610741565b005b6102676107a7565b604051808260ff16815260200191505060405180910390f35b6102886107d3565b005b6102926107fe565b005b61029c61080e565b604051808260028111156102ac57fe5b815260200191505060405180910390f35b6000600580549050905090565b600660006101000a81549060ff0219169055565b60006060600760000160009054906101000a900460ff166007600101808054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561038f5780601f106103645761010080835404028352916020019161038f565b820191906000526020600020905b81548152906001019060200180831161037257829003601f168201915b50505050509050915091509091565b60016000806101000a81548160ff021916908315150217905550600260018190555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3132330000000000000000000000000000000000000000000000000000000000600260146101000a81548162ffffff021916908360e81c02179055506040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525060039080519060200190610489929190610825565b5060056004819055506040518060400160405280600160ff1681526020016040518060400160405280600481526020017f456c6c6100000000000000000000000000000000000000000000000000000000815250815250600760008201518160000160006101000a81548160ff021916908360ff160217905550602082015181600101908051906020019061051f929190610825565b509050506040518060600160405280600160ff168152602001600260ff168152602001600360ff16815250600590600361055a9291906108a5565b50600160096000600160ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600260096000600260ff16815260200190815260200160002060006101000a81548160ff021916908360ff160217905550600360096000600360ff16815260200190815260200160002060006101000a81548160ff021916908360ff1602179055506002600660006101000a81548160ff0219169083600281111561060d57fe5b0217905550565b60096000600260ff16815260200190815260200160002060006101000a81549060ff0219169055565b600080600080606060008060009054906101000a900460ff16600154600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260149054906101000a900460e81b6003600454818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107265780601f106106fb57610100808354040283529160200191610726565b820191906000526020600020905b81548152906001019060200180831161070957829003601f168201915b50505050509150955095509550955095509550909192939495565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600260146101000a81549062ffffff02191690556003600061079f919061094c565b600460009055565b600060096000600260ff16815260200190815260200160002060009054906101000a900460ff16905090565b6007600080820160006101000a81549060ff02191690556001820160006107fa919061094c565b5050565b6005600061080c9190610994565b565b6000600660009054906101000a900460ff16905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061086657805160ff1916838001178555610894565b82800160010185558215610894579182015b82811115610893578251825591602001919060010190610878565b5b5090506108a191906109bc565b5090565b82805482825590600052602060002090601f0160209004810192821561093b5791602002820160005b8382111561090c57835183826101000a81548160ff021916908360ff16021790555092602001926001016020816000010492830192600103026108ce565b80156109395782816101000a81549060ff021916905560010160208160000104928301926001030261090c565b505b50905061094891906109d9565b5090565b50805460018160011615610100020316600290046000825580601f106109725750610991565b601f01602090049060005260206000209081019061099091906109bc565b5b50565b50805460008255601f0160209004906000526020600020908101906109b991906109bc565b50565b5b808211156109d55760008160009055506001016109bd565b5090565b5b808211156109fd57600081816101000a81549060ff0219169055506001016109da565b509056fea26469706673582212203a54a9ec2cf7a0b5f64712ac763776cca9a5b8824b109c3ed12a5377a98126f664736f6c634300060c0033";

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

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BasicDataTypeDeleteContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
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

    public RemoteCall<TransactionReceipt> getArrayLength() {
        final Function function = new Function(
                FUNC_GETARRAYLENGTH, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBasicData() {
        final Function function = new Function(
                FUNC_GETBASICDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getEnum() {
        final Function function = new Function(
                FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getMapping() {
        final Function function = new Function(
                FUNC_GETMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getStruct() {
        final Function function = new Function(
                FUNC_GETSTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> initBasicData() {
        final Function function = new Function(
                FUNC_INITBASICDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
