package network.platon.contracts.evm.v0_7_1;

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
public class CallerOne extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506105c7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c146100465780637b8ed01814610064578063f88bef8c14610082575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b6040518082815260200191505060405180910390f35b61008a61009b565b005b60005481565b60008054905090565b60006040516100a9906101bd565b604051809103906000f0801580156100c5573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff167f371303c051bff726100ad13871cababf50c20dd920fca137e519f98f089a74b4604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610152578051825260208201915060208101905060208303925061012f565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101b2576040519150601f19603f3d011682016040523d82523d6000602084013e6101b7565b606091505b50505050565b6103c7806101cb8339019056fe608060405234801561001057600080fd5b506103a7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630c55699c14610046578063371303c0146100645780635a3617561461006e575b600080fd5b61004e61008c565b6040518082815260200191505060405180910390f35b61006c610092565b005b61007661020b565b6040518082815260200191505060405180910390f35b60005481565b60006040516100a090610214565b604051809103906000f0801580156100bc573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff167f371303c051bff726100ad13871cababf50c20dd920fca137e519f98f089a74b4604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106101495780518252602082019150602081019050602083039250610126565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855af49150503d80600081146101a9576040519150601f19603f3d011682016040523d82523d6000602084013e6101ae565b606091505b5050507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a150565b60008054905090565b610150806102228339019056fe608060405234801561001057600080fd5b50610130806100206000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80630c55699c14604157806317f936fb14605d578063371303c0146079575b600080fd5b60476081565b6040518082815260200191505060405180910390f35b60636087565b6040518082815260200191505060405180910390f35b607f6090565b005b60005481565b60008054905090565b60008081548092919060010191905055507fb0333e0e3a6b99318e4e2e0d7e5e5f93646f9cbf62da1587955a4092bf7df6e733600054604051808373ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a156fea2646970667358221220cad78347f3612ff78c53aa548f2a06dd0a405e38230359614191b0c7081bbd2364736f6c63430007010033a26469706673582212209d23f1a2cbf08731a2aeee1a8e86b1fb1df68ecbb44f19c803a453e5802afd1964736f6c63430007010033a2646970667358221220606f95bea2caa068b41b2b9ed5f29cf45311a153dfb5292d580861aa1f13716664736f6c63430007010033";

    public static final String FUNC_GETCALLERX = "getCallerX";

    public static final String FUNC_INC_DELEGATECALL = "inc_delegatecall";

    public static final String FUNC_X = "x";

    protected CallerOne(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CallerOne(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getCallerX() {
        final Function function = new Function(
                FUNC_GETCALLERX, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> inc_delegatecall() {
        final Function function = new Function(
                FUNC_INC_DELEGATECALL, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> x() {
        final Function function = new Function(
                FUNC_X, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<CallerOne> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerOne.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CallerOne> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CallerOne.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CallerOne(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CallerOne load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CallerOne(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
