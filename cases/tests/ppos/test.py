import rlp
from client_sdk_python import HTTPProvider, Web3, WebsocketProvider, ppos
from client_sdk_python.middleware import geth_poa_middleware
from client_sdk_python.ppos import Ppos


def connect_web3(url, chain_id=101):
    if "ws" in url:
        w3 = Web3(WebsocketProvider(url), chain_id=chain_id)
    else:
        w3 = Web3(HTTPProvider(url), chain_id=chain_id)
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    return w3


def createRestrictingPlan(url, account, plan, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.createRestrictingPlan(account, plan, pri_key)
    print(result)


def delegate(url, typ, node_id, amount, pri_key):
    web3 = connect_web3(url)
    ppos = Ppos(web3)
    result = ppos.delegate(typ, node_id, amount, pri_key)
    print(result)


if __name__ == '__main__':
    url = 'http://192.168.120.146:6789'
    account = 'lax184zj2xdms82dvg5ypacsk48qw3ch0q9rtfrmp3'
    epoch = 100
    amount = Web3.toWei(1000, 'ether')
    # amount = ''
    plan = [{'Epoch': epoch, 'Amount': amount}]
    pri_key = 'd162b28e2ed3c4c0b991c69585bcec362746b86b1666178d7324a3ca56bd4591'
    nodeid = '01027ec8d9ea3c6f334486f88b41f7bfccfaf4aa9412a6cd88e837013b2235b9dba49108b12cc795bba905f5a66e69f6d2fe809f6f048f3fcc0c217360dbc0b2'
    createRestrictingPlan(url, account, plan, pri_key)
    # delegate(url, 0, nodeid, amount, pri_key)
