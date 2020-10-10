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
 * <p>Generated with web3j version 0.13.2.0.
 */
public class BasicDataTypeConstantContract extends Contract {
    private static final String BINARY = "60806040526040518060400160405280600581526020017f68656c6c6f0000000000000000000000000000000000000000000000000000008152506000908051906020019061004f9291906100ae565b506040518060400160405280600581526020017f776f726c640000000000000000000000000000000000000000000000000000008152506001908051906020019061009b9291906100ae565b503480156100a857600080fd5b5061014b565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ef57805160ff191683800117855561011d565b8280016001018555821561011d579182015b8281111561011c578251825591602001919060010190610101565b5b50905061012a919061012e565b5090565b5b8082111561014757600081600090555060010161012f565b5090565b610a718061015a6000396000f3fe6080604052600436106100fe5760003560e01c80636273899811610095578063a574971011610064578063a5749710146104cd578063a650b683146104f8578063ccb3441514610588578063ea07cdfa146105d5578063f8b2cb4f14610600576100fe565b8063627389981461031357806369c769171461033e5780637c66959d146103ce578063833f17d51461046c576100fe565b806338cc4831116100d157806338cc48311461022557806343e33562146102665780634dbe1d0f146102b2578063515899f0146102e8576100fe565b8063166aa6e614610103578063209652551461016057806332c7a283146101ab57806338023fb9146101e1575b600080fd5b34801561010f57600080fd5b5061013f6004803603602081101561012657600080fd5b81019080803560ff169060200190929190505050610665565b6040518082600381111561014f57fe5b815260200191505060405180910390f35b34801561016c57600080fd5b5061017561066f565b60405180846fffffffffffffffffffffffffffffffff168152602001838152602001828152602001935050505060405180910390f35b3480156101b757600080fd5b506101c06106a0565b604051808260038111156101d057fe5b815260200191505060405180910390f35b610223600480360360208110156101f757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506106b8565b005b34801561023157600080fd5b5061023a610702565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561027257600080fd5b5061027b610723565b60405180827dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156102be57600080fd5b506102c7610750565b604051808260038111156102d757fe5b815260200191505060405180910390f35b3480156102f457600080fd5b506102fd610761565b6040518082815260200191505060405180910390f35b34801561031f57600080fd5b5061032861077e565b6040518082815260200191505060405180910390f35b34801561034a57600080fd5b5061035361078c565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610393578082015181840152602081019050610378565b50505050905090810190601f1680156103c05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156103da57600080fd5b506103e3610890565b60405180847dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001935050505060405180910390f35b6104ae6004803603602081101561048257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108e8565b6040518083815260200182151581526020019250505060405180910390f35b3480156104d957600080fd5b506104e261092a565b6040518082815260200191505060405180910390f35b34801561050457600080fd5b5061050d610932565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561054d578082015181840152602081019050610532565b50505050905090810190601f16801561057a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561059457600080fd5b5061059d6109d4565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b3480156105e157600080fd5b506105ea610a01565b6040518082815260200191505060405180910390f35b34801561060c57600080fd5b5061064f6004803603602081101561062357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a1a565b6040518082815260200191505060405180910390f35b6000819050919050565b6000806000806003905060006404a817c80090506000660aa87bee5380009050828282955095509550505050909192565b6000806003808111156106af57fe5b90508091505090565b8073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f193505050501580156106fe573d6000803e3d6000fd5b5050565b6000807372ad2b713faa14c2c4cd2d7affe5d8f538968f5a90508091505090565b6000807f01f400000000000000000000000000000000000000000000000000000000000090508091505090565b600061075c6001610665565b905090565b600080805460018160011615610100020316600290049050905090565b600080600690508091505090565b60608060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108255780601f106107fa57610100808354040283529160200191610825565b820191906000526020600020905b81548152906001019060200180831161080857829003601f168201915b505050505090507f61000000000000000000000000000000000000000000000000000000000000008160008151811061085a57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508091505090565b6000806000807f01f4000000000000000000000000000000000000000000000000000000000000905080816000600281106108c757fe5b1a60f81b826001600281106108d857fe5b1a60f81b93509350935050909192565b600080348373ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f1935050505091509150915091565b600047905090565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109ca5780601f1061099f576101008083540402835291602001916109ca565b820191906000526020600020905b8154815290600101906020018083116109ad57829003601f168201915b5050505050905090565b6000807fc80000000000000000000000000000000000000000000000000000000000000090508091505090565b60008060006003811115610a1157fe5b90508091505090565b60008173ffffffffffffffffffffffffffffffffffffffff1631905091905056fea264697066735822122012ac30a9eb40bd6fafb83df4746465e95869a120770d84ea131a0d92a6d5b86464736f6c634300060c0033";

    public static final String FUNC_GETADDRESS = "getAddress";

    public static final String FUNC_GETBALANCE = "getBalance";

    public static final String FUNC_GETCURRENTBALANCE = "getCurrentBalance";

    public static final String FUNC_GETHEXLITERAA = "getHexLiteraA";

    public static final String FUNC_GETHEXLITERAB = "getHexLiteraB";

    public static final String FUNC_GETHEXLITERAC = "getHexLiteraC";

    public static final String FUNC_GETINT = "getInt";

    public static final String FUNC_GETSEASONA = "getSeasonA";

    public static final String FUNC_GETSEASONB = "getSeasonB";

    public static final String FUNC_GETSEASONINDEX = "getSeasonIndex";

    public static final String FUNC_GETSTRA = "getStrA";

    public static final String FUNC_GETSTRALENGTH = "getStrALength";

    public static final String FUNC_GETVALUE = "getValue";

    public static final String FUNC_GOSEND = "goSend";

    public static final String FUNC_GOTRANSFER = "goTransfer";

    public static final String FUNC_PRINTSEASON = "printSeason";

    public static final String FUNC_SETSTRA = "setStrA";

    protected BasicDataTypeConstantContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected BasicDataTypeConstantContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getAddress() {
        final Function function = new Function(
                FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getBalance(String addr) {
        final Function function = new Function(
                FUNC_GETBALANCE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getCurrentBalance() {
        final Function function = new Function(
                FUNC_GETCURRENTBALANCE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getHexLiteraA() {
        final Function function = new Function(
                FUNC_GETHEXLITERAA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getHexLiteraB() {
        final Function function = new Function(
                FUNC_GETHEXLITERAB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getHexLiteraC() {
        final Function function = new Function(
                FUNC_GETHEXLITERAC, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getInt() {
        final Function function = new Function(
                FUNC_GETINT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSeasonA() {
        final Function function = new Function(
                FUNC_GETSEASONA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSeasonB() {
        final Function function = new Function(
                FUNC_GETSEASONB, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSeasonIndex() {
        final Function function = new Function(
                FUNC_GETSEASONINDEX, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getStrA() {
        final Function function = new Function(
                FUNC_GETSTRA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getStrALength() {
        final Function function = new Function(
                FUNC_GETSTRALENGTH, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getValue() {
        final Function function = new Function(
                FUNC_GETVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> goSend(String addr) {
        final Function function = new Function(
                FUNC_GOSEND, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> goTransfer(String addr) {
        final Function function = new Function(
                FUNC_GOTRANSFER, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(addr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> printSeason(BigInteger s) {
        final Function function = new Function(
                FUNC_PRINTSEASON, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint8(s)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setStrA() {
        final Function function = new Function(
                FUNC_SETSTRA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<BasicDataTypeConstantContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeConstantContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<BasicDataTypeConstantContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(BasicDataTypeConstantContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static BasicDataTypeConstantContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeConstantContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static BasicDataTypeConstantContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new BasicDataTypeConstantContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
