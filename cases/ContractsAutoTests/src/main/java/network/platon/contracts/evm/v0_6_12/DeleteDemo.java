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
public class DeleteDemo extends Contract {
    private static final String BINARY = "608060405260016000806101000a81548160ff0219169083151502179055506001805533600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040518060400160405280600381526020017f313233000000000000000000000000000000000000000000000000000000000081525060039080519060200190620000b092919062000138565b506040518060400160405280600381526020017f616263000000000000000000000000000000000000000000000000000000000081525060049080519060200190620000fe929190620001bf565b506001600560006101000a81548160ff021916908360028111156200011f57fe5b02179055503480156200013157600080fd5b5062000265565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200017b57805160ff1916838001178555620001ac565b82800160010185558215620001ac579182015b82811115620001ab5782518255916020019190600101906200018e565b5b509050620001bb919062000246565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200020257805160ff191683800117855562000233565b8280016001018555821562000233579182015b828111156200023257825182559160200191906001019062000215565b5b50905062000242919062000246565b5090565b5b808211156200026157600081600090555060010162000247565b5090565b610e3680620002756000396000f3fe6080604052600436106101095760003560e01c8063767800de11610095578063c15bae8411610064578063c15bae84146104e3578063cf08fed514610573578063d1bdda41146105a9578063e5aa3d5814610633578063f02997491461065e57610109565b8063767800de1461035757806393e1ed8314610398578063a1a984e514610428578063ab5170b21461045357610109565b806327c58232116100dc57806327c582321461024e57806332d057c9146102655780633ab0698c1461026f5780634df7e3d01461029a5780635d743b5d146102c757610109565b806305be2c121461010e57806313a5a8af146101a55780631acddabe146101db578063252bd4d31461020d575b600080fd5b34801561011a57600080fd5b5061012361068b565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561016957808201518184015260208101905061014e565b50505050905090810190601f1680156101965780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b3480156101b157600080fd5b506101ba61073e565b604051808260028111156101ca57fe5b815260200191505060405180910390f35b3480156101e757600080fd5b506101f0610755565b604051808381526020018281526020019250505060405180910390f35b34801561021957600080fd5b506102226107a9565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561025a57600080fd5b506102636107d3565b005b61026d61083f565b005b34801561027b57600080fd5b506102846108b7565b6040518082815260200191505060405180910390f35b3480156102a657600080fd5b506102af610947565b60405180821515815260200191505060405180910390f35b3480156102d357600080fd5b506102dc610958565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561031c578082015181840152602081019050610301565b50505050905090810190601f1680156103495780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561036357600080fd5b5061036c6109fa565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103a457600080fd5b506103ad610a20565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103ed5780820151818401526020810190506103d2565b50505050905090810190601f16801561041a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561043457600080fd5b5061043d610abe565b6040518082815260200191505060405180910390f35b34801561045f57600080fd5b50610468610ac8565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104a857808201518184015260208101905061048d565b50505050905090810190601f1680156104d55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156104ef57600080fd5b506104f8610b6a565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561053857808201518184015260208101905061051d565b50505050905090810190601f1680156105655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561057f57600080fd5b50610588610c08565b6040518082600281111561059857fe5b815260200191505060405180910390f35b6105b1610c1b565b6040518083815260200180602001828103825283818151815260200191508051906020019080838360005b838110156105f75780820151818401526020810190506105dc565b50505050905090810190601f1680156106245780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b34801561063f57600080fd5b50610648610d37565b6040518082815260200191505060405180910390f35b34801561066a57600080fd5b50610673610d3d565b60405180821515815260200191505060405180910390f35b600060606006600001546006600101808054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561072f5780601f106107045761010080835404028352916020019161072f565b820191906000526020600020905b81548152906001019060200180831161071257829003601f168201915b50505050509050915091509091565b6000600560009054906101000a900460ff16905090565b600080600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600860010154915091509091565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000806101000a81549060ff0219169055600160009055600260006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003600061081d9190610d53565b6004600061082b9190610d9b565b600560006101000a81549060ff0219169055565b604051806020016040528060c88152506008600082015181600101559050506107d0600860000160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506008600060018201600090555050565b60006060600767ffffffffffffffff811180156108d357600080fd5b506040519080825280602002602001820160405280156109025781602001602082028036833780820191505090505b50905060648160008151811061091457fe5b60200260200101818152505060c88160018151811061092f57fe5b60200260200101818152505060609050805191505090565b60008054906101000a900460ff1681565b606060038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109f05780601f106109c5576101008083540402835291602001916109f0565b820191906000526020600020905b8154815290600101906020018083116109d357829003601f168201915b5050505050905090565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ab65780601f10610a8b57610100808354040283529160200191610ab6565b820191906000526020600020905b815481529060010190602001808311610a9957829003601f168201915b505050505081565b6000600154905090565b606060048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b605780601f10610b3557610100808354040283529160200191610b60565b820191906000526020600020905b815481529060010190602001808311610b4357829003601f168201915b5050505050905090565b60048054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c005780601f10610bd557610100808354040283529160200191610c00565b820191906000526020600020905b815481529060010190602001808311610be357829003601f168201915b505050505081565b600560009054906101000a900460ff1681565b600060606040518060400160405280600a81526020016040518060400160405280600381526020017f6162630000000000000000000000000000000000000000000000000000000000815250815250506006600080820160009055600182016000610c869190610d9b565b50506006600001546006600101808054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d285780601f10610cfd57610100808354040283529160200191610d28565b820191906000526020600020905b815481529060010190602001808311610d0b57829003601f168201915b50505050509050915091509091565b60015481565b60008060009054906101000a900460ff16905090565b50805460018160011615610100020316600290046000825580601f10610d795750610d98565b601f016020900490600052602060002090810190610d979190610de3565b5b50565b50805460018160011615610100020316600290046000825580601f10610dc15750610de0565b601f016020900490600052602060002090810190610ddf9190610de3565b5b50565b5b80821115610dfc576000816000905550600101610de4565b509056fea26469706673582212200565ee5210a110fc3d121925a80b7329ee78f0d4b07d962f8ef0bf36a647343664736f6c634300060c0033";

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

    protected DeleteDemo(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DeleteDemo(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> addr() {
        final Function function = new Function(
                FUNC_ADDR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> b() {
        final Function function = new Function(
                FUNC_B, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> color() {
        final Function function = new Function(
                FUNC_COLOR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> delDynamicArray() {
        final Function function = new Function(
                FUNC_DELDYNAMICARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> delMapping() {
        final Function function = new Function(
                FUNC_DELMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> delStruct() {
        final Function function = new Function(
                FUNC_DELSTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> deleteAttr() {
        final Function function = new Function(
                FUNC_DELETEATTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getaddress() {
        final Function function = new Function(
                FUNC_GETADDRESS, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getbool() {
        final Function function = new Function(
                FUNC_GETBOOL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getbytes() {
        final Function function = new Function(
                FUNC_GETBYTES, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getdelMapping() {
        final Function function = new Function(
                FUNC_GETDELMAPPING, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getenum() {
        final Function function = new Function(
                FUNC_GETENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getstr() {
        final Function function = new Function(
                FUNC_GETSTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getstruct() {
        final Function function = new Function(
                FUNC_GETSTRUCT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getunit() {
        final Function function = new Function(
                FUNC_GETUNIT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> i() {
        final Function function = new Function(
                FUNC_I, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> str() {
        final Function function = new Function(
                FUNC_STR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> varByte() {
        final Function function = new Function(
                FUNC_VARBYTE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
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
