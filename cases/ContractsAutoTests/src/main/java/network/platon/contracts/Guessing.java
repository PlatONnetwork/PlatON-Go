package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class Guessing extends Contract {
    private static final String BINARY = "60806040526000600260006101000a81548160ff021916908315150217905550674563918244f400006003556001600a556000600d5534801561004157600080fd5b5060405161121d38038061121d8339818101604052602081101561006457600080fd5b810190808051906020019092919050505033600c60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060008190555050611150806100cd6000396000f3fe6080604052600436106101405760003560e01c806382d333a0116100b6578063b69ef8a81161006f578063b69ef8a814610618578063c2930f9114610643578063c29bde8c1461066e578063c7f40b6414610699578063ecbde5e6146106f0578063ff37be161461071b57610140565b806382d333a0146104335780638941f2f01461045e57806394696a92146104ad578063a46c363714610519578063ae44794114610544578063b03e0077146105a957610140565b806313eaca431161010857806313eaca43146102c35780631ef9c56f146102f2578063220bc55e1461035e57806327ebd9ab14610368578063629374ab146103935780636e5ab671146103ce57610140565b8063045f9c971461014c578063062d6a98146101c7578063083c6323146101f2578063094cc1ab1461021d5780630b8b850214610248575b61014a3433610746565b005b34801561015857600080fd5b506101856004803603602081101561016f57600080fd5b8101908080359060200190929190505050610909565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101d357600080fd5b506101dc610945565b6040518082815260200191505060405180910390f35b3480156101fe57600080fd5b5061020761094b565b6040518082815260200191505060405180910390f35b34801561022957600080fd5b50610232610951565b6040518082815260200191505060405180910390f35b34801561025457600080fd5b506102816004803603602081101561026b57600080fd5b8101908080359060200190929190505050610957565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102cf57600080fd5b506102d861098a565b604051808215151515815260200191505060405180910390f35b3480156102fe57600080fd5b5061030761099d565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561034a57808201518184015260208101905061032f565b505050509050019250505060405180910390f35b610366610a2b565b005b34801561037457600080fd5b5061037d610a37565b6040518082815260200191505060405180910390f35b34801561039f57600080fd5b506103cc600480360360208110156103b657600080fd5b8101908080359060200190929190505050610a61565b005b3480156103da57600080fd5b5061041d600480360360208110156103f157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610db2565b6040518082815260200191505060405180910390f35b34801561043f57600080fd5b50610448610dca565b6040518082815260200191505060405180910390f35b34801561046a57600080fd5b506104976004803603602081101561048157600080fd5b8101908080359060200190929190505050610dd0565b6040518082815260200191505060405180910390f35b3480156104b957600080fd5b506104c2610df9565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156105055780820151818401526020810190506104ea565b505050509050019250505060405180910390f35b34801561052557600080fd5b5061052e610e8e565b6040518082815260200191505060405180910390f35b34801561055057600080fd5b506105936004803603602081101561056757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610e9b565b6040518082815260200191505060405180910390f35b3480156105b557600080fd5b50610602600480360360408110156105cc57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610eb3565b6040518082815260200191505060405180910390f35b34801561062457600080fd5b5061062d610ee1565b6040518082815260200191505060405180910390f35b34801561064f57600080fd5b50610658610ee7565b6040518082815260200191505060405180910390f35b34801561067a57600080fd5b50610683610eed565b6040518082815260200191505060405180910390f35b3480156106a557600080fd5b506106ae610ef3565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156106fc57600080fd5b50610705610f19565b6040518082815260200191505060405180910390f35b34801561072757600080fd5b50610730610f38565b6040518082815260200191505060405180910390f35b43600054106109055760035482101561075e57600080fd5b600860008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600a5490806001815401808255809150509060018203906000526020600020016000909192909190915055508060076000600a54815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600a6000815480929190600101919050555081600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282540192505081905550816004600082825401925050819055507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf681836001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a15b5050565b600b818154811061091657fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600a5481565b60005481565b600d5481565b60076020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900460ff1681565b6060600b805480602002602001604051908101604052809291908181526020018280548015610a2157602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190600101908083116109d7575b5050505050905090565b610a353433610746565b565b600060011515600260009054906101000a900460ff16151514610a5957600080fd5b600d54905090565b600054431115610daf57600260009054906101000a900460ff16158015610ad557503373ffffffffffffffffffffffffffffffffffffffff16600c60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b8015610ae357506001600a54115b15610dae57600081604051602001808281526020019150506040516020818303038152906040528051906020012060001c90506000600a548281610b2357fe5b0690506064600a541015610b4157610b3c81600a610f42565b610b6b565b612710600a541015610b5d57610b58816064610f42565b610b6a565b610b69816103e8610f42565b5b5b600b8054905060045481610b7b57fe5b04600581905550600080600090505b600b80549050811015610d8757600160096000600b8481548110610baa57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020540160096000600b8481548110610c2257fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600073ffffffffffffffffffffffffffffffffffffffff16600b8281548110610cb057fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610d7a57600b8181548110610d0357fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691508173ffffffffffffffffffffffffffffffffffffffff166108fc6005549081150290604051600060405180830381858888f19350505050158015610d78573d6000803e3d6000fd5b505b8080600101915050610b8a565b506001600260006101000a81548160ff021916908315150217905550836001819055505050505b5b50565b60096020528060005260406000206000915090505481565b60055481565b6000610101430382118015610de757506001430382105b610df057600080fd5b81409050919050565b6060600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020805480602002602001604051908101604052809291908181526020018280548015610e8457602002820191906000526020600020905b815481526020019060010190808311610e70575b5050505050905090565b6000600b80549050905090565b60066020528060005260406000206000915090505481565b60086020528160005260406000208181548110610ecc57fe5b90600052602060002001600091509150505481565b60045481565b60035481565b60015481565b600c60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b6000600154905090565b808281610f4b57fe5b06600d819055506000600d541415611031576000600190505b600a5481101561102b57600082600d54830381610f7d57fe5b06141561101e57600b6007600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b8080600101915050610f64565b50611117565b6000600190505b600a5481101561111557600082828161104d57fe5b06141580156110695750600082600d5483038161106657fe5b06145b1561110857600b6007600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b8080600101915050611038565b505b505056fea265627a7a72315820a1b3f7cf42689a6608ca798e2e4763512eb0243ed6b2c3dbea09d66ba71aa3c964736f6c634300050d0032";

    public static final String FUNC_AVERAGEAMOUNT = "averageAmount";

    public static final String FUNC_BALANCE = "balance";

    public static final String FUNC_BASEUNIT = "baseUnit";

    public static final String FUNC_BLOCK_HASH = "block_hash";

    public static final String FUNC_CREATEADDRESS = "createAddress";

    public static final String FUNC_DRAW = "draw";

    public static final String FUNC_ENDBLOCK = "endBlock";

    public static final String FUNC_GENERATEBLOCKHASH = "generateBlockHash";

    public static final String FUNC_GETBALANCEOF = "getBalanceOf";

    public static final String FUNC_GETENDBLOCKHASH = "getEndBlockHash";

    public static final String FUNC_GETMYGUESSCODES = "getMyGuessCodes";

    public static final String FUNC_GETPOSTFIX = "getPostfix";

    public static final String FUNC_GETWINNERADDRESSES = "getWinnerAddresses";

    public static final String FUNC_GETWINNERCOUNT = "getWinnerCount";

    public static final String FUNC_GUESSINGCLOSED = "guessingClosed";

    public static final String FUNC_GUESSINGWITHLAT = "guessingWithLat";

    public static final String FUNC_GUSSINGERCODES = "gussingerCodes";

    public static final String FUNC_GUSSINGERLAT = "gussingerLat";

    public static final String FUNC_INDEXKEY = "indexKey";

    public static final String FUNC_INDEXOFGUSSINGER = "indexOfgussinger";

    public static final String FUNC_POSTFIX = "postfix";

    public static final String FUNC_WINNERADDRESSES = "winnerAddresses";

    public static final String FUNC_WINNERMAP = "winnerMap";

    public static final Event CURRENTBALANCE_EVENT = new Event("CurrentBalance", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event FUNDTRANSFER_EVENT = new Event("FundTransfer", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}, new TypeReference<Uint256>() {}, new TypeReference<Bool>() {}));
    ;

    protected Guessing(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Guessing(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<Guessing> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger _endBlock) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_endBlock)));
        return deployRemoteCall(Guessing.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public List<CurrentBalanceEventResponse> getCurrentBalanceEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(CURRENTBALANCE_EVENT, transactionReceipt);
        ArrayList<CurrentBalanceEventResponse> responses = new ArrayList<CurrentBalanceEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, CurrentBalanceEventResponse>() {
            @Override
            public CurrentBalanceEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(CURRENTBALANCE_EVENT, log);
                CurrentBalanceEventResponse typedResponse = new CurrentBalanceEventResponse();
                typedResponse.log = log;
                typedResponse._msgSenderAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._balance = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<CurrentBalanceEventResponse> currentBalanceEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(CURRENTBALANCE_EVENT));
        return currentBalanceEventObservable(filter);
    }

    public List<FundTransferEventResponse> getFundTransferEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(FUNDTRANSFER_EVENT, transactionReceipt);
        ArrayList<FundTransferEventResponse> responses = new ArrayList<FundTransferEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            FundTransferEventResponse typedResponse = new FundTransferEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, FundTransferEventResponse>() {
            @Override
            public FundTransferEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(FUNDTRANSFER_EVENT, log);
                FundTransferEventResponse typedResponse = new FundTransferEventResponse();
                typedResponse.log = log;
                typedResponse._backer = (String) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse._amount = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                typedResponse._isSuccess = (Boolean) eventValues.getNonIndexedValues().get(2).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<FundTransferEventResponse> fundTransferEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(FUNDTRANSFER_EVENT));
        return fundTransferEventObservable(filter);
    }

    public RemoteCall<BigInteger> averageAmount() {
        final Function function = new Function(FUNC_AVERAGEAMOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> balance() {
        final Function function = new Function(FUNC_BALANCE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> baseUnit() {
        final Function function = new Function(FUNC_BASEUNIT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> block_hash() {
        final Function function = new Function(FUNC_BLOCK_HASH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> createAddress() {
        final Function function = new Function(FUNC_CREATEADDRESS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> draw(byte[] _block_hash) {
        final Function function = new Function(
                FUNC_DRAW, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Bytes32(_block_hash)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> endBlock() {
        final Function function = new Function(FUNC_ENDBLOCK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> generateBlockHash(BigInteger _blocknumber) {
        final Function function = new Function(FUNC_GENERATEBLOCKHASH, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_blocknumber)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> getBalanceOf() {
        final Function function = new Function(FUNC_GETBALANCEOF, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> getEndBlockHash() {
        final Function function = new Function(FUNC_GETENDBLOCKHASH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<List> getMyGuessCodes() {
        final Function function = new Function(FUNC_GETMYGUESSCODES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<BigInteger> getPostfix() {
        final Function function = new Function(FUNC_GETPOSTFIX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<List> getWinnerAddresses() {
        final Function function = new Function(FUNC_GETWINNERADDRESSES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Address>>() {}));
        return new RemoteCall<List>(
                new Callable<List>() {
                    @Override
                    @SuppressWarnings("unchecked")
                    public List call() throws Exception {
                        List<Type> result = (List<Type>) executeCallSingleValueReturn(function, List.class);
                        return convertToNative(result);
                    }
                });
    }

    public RemoteCall<BigInteger> getWinnerCount() {
        final Function function = new Function(FUNC_GETWINNERCOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> guessingClosed() {
        final Function function = new Function(FUNC_GUESSINGCLOSED, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> guessingWithLat(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GUESSINGWITHLAT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> gussingerCodes(String param0, BigInteger param1) {
        final Function function = new Function(FUNC_GUSSINGERCODES, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(param0), 
                new org.web3j.abi.datatypes.generated.Uint256(param1)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> gussingerLat(String param0) {
        final Function function = new Function(FUNC_GUSSINGERLAT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> indexKey() {
        final Function function = new Function(FUNC_INDEXKEY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> indexOfgussinger(BigInteger param0) {
        final Function function = new Function(FUNC_INDEXOFGUSSINGER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> postfix() {
        final Function function = new Function(FUNC_POSTFIX, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> winnerAddresses(BigInteger param0) {
        final Function function = new Function(FUNC_WINNERADDRESSES, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> winnerMap(String param0) {
        final Function function = new Function(FUNC_WINNERMAP, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(param0)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static Guessing load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Guessing(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Guessing load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Guessing(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class CurrentBalanceEventResponse {
        public Log log;

        public String _msgSenderAddress;

        public BigInteger _balance;
    }

    public static class FundTransferEventResponse {
        public Log log;

        public String _backer;

        public BigInteger _amount;

        public Boolean _isSuccess;
    }
}
