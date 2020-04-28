package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
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
public class IUserManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_GETACCOUNTSTATE = "getAccountState";

    public static final String FUNC_FINDBYDEPARTMENTIDTREE = "findByDepartmentIdTree";

    public static final String FUNC_FINDBYMOBILE = "findByMobile";

    public static final String FUNC_USEREXISTS = "userExists";

    public static final String FUNC_CHECKUSERROLE = "checkUserRole";

    public static final String FUNC_UPDATEPASSWORDSTATUS = "updatePasswordStatus";

    public static final String FUNC_LOGIN = "login";

    public static final String FUNC_PAGEBYACCOUNTSTATUS = "pageByAccountStatus";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_GETUSERSTATE = "getUserState";

    public static final String FUNC_FINDBYEMAIL = "findByEmail";

    public static final String FUNC_FINDBYLOGINNAME = "findByLoginName";

    public static final String FUNC_UPDATEUSERSTATUS = "updateUserStatus";

    public static final String FUNC_GETUSERCOUNTMAPPINGBYROLEIDS = "getUserCountMappingByRoleIds";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_GETUSERROLEID = "getUserRoleId";

    public static final String FUNC_CHECKEMAILUNIQUENESS = "checkEmailUniqueness";

    public static final String FUNC_FINDBYROLEID = "findByRoleId";

    public static final String FUNC_GETUSERCOUNTBYDEPARTMENTID = "getUserCountByDepartmentId";

    public static final String FUNC_DELETEBYADDRESS = "deleteByAddress";

    public static final String FUNC_ADDUSERROLE = "addUserRole";

    public static final String FUNC_FINDBYADDRESS = "findByAddress";

    public static final String FUNC_FINDBYDEPARTMENTIDTREEANDCONTION = "findByDepartmentIdTreeAndContion";

    public static final String FUNC_FINDBYACCOUNT = "findByAccount";

    public static final String FUNC_ROLEUSED = "roleUsed";

    public static final String FUNC_GETUSERDEPARTMENTID = "getUserDepartmentId";

    public static final String FUNC_CHECKUSERACTION = "checkUserAction";

    public static final String FUNC_GETUSERADDRBYADDR = "getUserAddrByAddr";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_GETOWNERADDRBYADDR = "getOwnerAddrByAddr";

    public static final String FUNC_GETUSERCOUNTBYACTIONID = "getUserCountByActionId";

    public static final String FUNC_RESETPASSWD = "resetPasswd";

    public static final String FUNC_FINDBYDEPARTMENTID = "findByDepartmentId";

    public static final String FUNC_ISREPETITIVE = "isRepetitive";

    public static final String FUNC_FINDBYUUID = "findByUuid";

    public static final String FUNC_CHECKUSERPRIVILEGE = "checkUserPrivilege";

    public static final String FUNC_UPDATEACCOUNTSTATUS = "updateAccountStatus";

    @Deprecated
    protected IUserManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IUserManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IUserManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IUserManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getAccountState(String _account) {
        final Function function = new Function(FUNC_GETACCOUNTSTATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_account)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByDepartmentIdTree(String _departmentId, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_FINDBYDEPARTMENTIDTREE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByMobile(String _mobile) {
        final Function function = new Function(FUNC_FINDBYMOBILE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_mobile)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> userExists(String _userAddr) {
        final Function function = new Function(FUNC_USEREXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> checkUserRole(String _userAddr, String _roleId) {
        final Function function = new Function(FUNC_CHECKUSERROLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> updatePasswordStatus(String _userAddr, BigInteger _status) {
        final Function function = new Function(
                FUNC_UPDATEPASSWORDSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.generated.Uint256(_status)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> login(String _account) {
        final Function function = new Function(
                FUNC_LOGIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_account)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> pageByAccountStatus(BigInteger _accountStatus, BigInteger _pageNo, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYACCOUNTSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_accountStatus), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNo), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> update(String _userJson) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_userJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getUserState(String _userAddr) {
        final Function function = new Function(FUNC_GETUSERSTATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByEmail(String _email) {
        final Function function = new Function(FUNC_FINDBYEMAIL, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_email)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByLoginName(String _name) {
        final Function function = new Function(FUNC_FINDBYLOGINNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> updateUserStatus(String _userAddr, BigInteger _status) {
        final Function function = new Function(
                FUNC_UPDATEUSERSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.generated.Uint256(_status)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getUserCountMappingByRoleIds(String _roleIds) {
        final Function function = new Function(FUNC_GETUSERCOUNTMAPPINGBYROLEIDS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleIds)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getUserRoleId(String _userAddr, BigInteger _index) {
        final Function function = new Function(FUNC_GETUSERROLEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.generated.Uint256(_index)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> checkEmailUniqueness(String _email, String _mobile) {
        final Function function = new Function(FUNC_CHECKEMAILUNIQUENESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_email), 
                new org.web3j.abi.datatypes.Utf8String(_mobile)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByRoleId(String _roleId) {
        final Function function = new Function(FUNC_FINDBYROLEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getUserCountByDepartmentId(String _departmentId) {
        final Function function = new Function(FUNC_GETUSERCOUNTBYDEPARTMENTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> deleteByAddress(String _userAddr) {
        final Function function = new Function(
                FUNC_DELETEBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> addUserRole(String _userAddr, String _roleId) {
        final Function function = new Function(
                FUNC_ADDUSERROLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByAddress(String _userAddr) {
        final Function function = new Function(FUNC_FINDBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByDepartmentIdTreeAndContion(BigInteger _status, String _name, String _departmentId, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_FINDBYDEPARTMENTIDTREEANDCONTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_status), 
                new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.Utf8String(_departmentId), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByAccount(String _account) {
        final Function function = new Function(FUNC_FINDBYACCOUNT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_account)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> roleUsed(String _roleId) {
        final Function function = new Function(FUNC_ROLEUSED, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getUserDepartmentId(String _userAddr) {
        final Function function = new Function(FUNC_GETUSERDEPARTMENTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> checkUserAction(String _userAddr, String _actionId) {
        final Function function = new Function(FUNC_CHECKUSERACTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getUserAddrByAddr(String _userAddr) {
        final Function function = new Function(FUNC_GETUSERADDRBYADDR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> insert(String _userJson) {
        final Function function = new Function(
                FUNC_INSERT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_userJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getOwnerAddrByAddr(String _userAddr) {
        final Function function = new Function(FUNC_GETOWNERADDRBYADDR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getUserCountByActionId(String _actionId) {
        final Function function = new Function(FUNC_GETUSERCOUNTBYACTIONID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> resetPasswd(String _userAddr, String _ownerAddr, String _publilcKey, String _cipherGroupKey, String _uuid) {
        final Function function = new Function(
                FUNC_RESETPASSWD, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Address(_ownerAddr), 
                new org.web3j.abi.datatypes.Utf8String(_publilcKey), 
                new org.web3j.abi.datatypes.Utf8String(_cipherGroupKey), 
                new org.web3j.abi.datatypes.Utf8String(_uuid)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByDepartmentId(String _departmentId) {
        final Function function = new Function(FUNC_FINDBYDEPARTMENTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> isRepetitive(String _mobile, String _email, String _userAddr, String _uuid, String _publicKey, String _account) {
        final Function function = new Function(FUNC_ISREPETITIVE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_mobile), 
                new org.web3j.abi.datatypes.Utf8String(_email), 
                new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Utf8String(_uuid), 
                new org.web3j.abi.datatypes.Utf8String(_publicKey), 
                new org.web3j.abi.datatypes.Utf8String(_account)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByUuid(String _uuid) {
        final Function function = new Function(FUNC_FINDBYUUID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_uuid)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> checkUserPrivilege(String _userAddr, String _contractAddr, String _funcSha3) {
        final Function function = new Function(FUNC_CHECKUSERPRIVILEGE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.Address(_contractAddr), 
                new org.web3j.abi.datatypes.Utf8String(_funcSha3)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> updateAccountStatus(String _userAddr, BigInteger _status) {
        final Function function = new Function(
                FUNC_UPDATEACCOUNTSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr), 
                new org.web3j.abi.datatypes.generated.Uint256(_status)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<IUserManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IUserManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IUserManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IUserManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IUserManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IUserManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IUserManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IUserManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IUserManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IUserManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IUserManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IUserManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IUserManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IUserManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IUserManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IUserManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
