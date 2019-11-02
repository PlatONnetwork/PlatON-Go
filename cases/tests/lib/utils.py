# -*- coding: utf-8 -*-
import time
import random
import string
from hexbytes import HexBytes
from environment.node import Node
# from pip import *
from common.log import log
from typing import List
from .client import Client
from .pip import Pip


def decorator_sleep(func):
    def wrap():
        result = func()
        if result is None:
            time.sleep(5)
            result = func()
        return result
    return wrap


def proposal_list_effective(proposal_list, block_number):
    """
    Determine if there is a proposal in the voting period
    :param proposal_list:
    :param block_number:
    :return:
    """
    for proposal in proposal_list:
        if proposal_effective(proposal, block_number):
            return True
    return False


def proposal_effective(proposal, block_number):
    """
    Determine if the proposal is in the voting period
    :param proposal:
    :param block_number:
    :return:
    """
    if proposal["EndVotingBlock"] > block_number:
        return True
    return False


def upload_platon(node: Node, platon_bin):
    """
    Upload a binary file to the specified node
    :param node:
    :param platon_bin:
    :return:
    """
    node.run_ssh("rm -rf {}".format(node.remote_bin_file))
    node.upload_file(platon_bin, node.remote_bin_file)
    node.run_ssh("chmod +x {}".format(node.remote_bin_file))
    node.run_ssh("mkdir zlp")


def get_blockhash(node, blocknumber=None):
    """
    Get block hash based on block height
    :param node:
    :param blocknumber:
    :return:
    """
    if not blocknumber:
        blocknumber = node.blockNumber
    blockinfo = node.eth.getBlock(blocknumber)
    blockhash = blockinfo.get('hash')
    blockhash = HexBytes(blockhash).hex()
    return blockhash


def int_to_bytes(value):
    return int(value).to_bytes(length=4, byteorder='big', signed=False)


def int16_to_bytes(value):
    return int(value).to_bytes(length=1, byteorder='big', signed=False)


def bytes_to_int(value):
    return int.from_bytes(value, byteorder='big', signed=False)


def compare_two_dict(dict1, dict2, key_list=None):
    """
    Compare two dictionary values
    :param dict1:
    :param dict2:
    :param key_list: Align dictionary key list
    :return:
    """
    if key_list is None:
        key_list = ['blockNumber', 'amount']
    flag = True
    keys1 = dict1.keys()
    keys2 = dict2.keys()
    if len(key_list) != 0:
        for key in key_list:
            if key in keys1 and key in keys2:
                if dict1[key] == dict2[key]:
                    flag = True
                else:
                    flag = False
            else:
                raise Exception('key_list contains error key')
    else:
        raise Exception('key_list is null')
    return flag


def get_no_pledge_node(node_list) -> Node:
    """
    Get nodes that are not pledged
    :param node_list: Node list
    :return:
    """
    for node in node_list:
        result = node.ppos.getCandidateInfo(node.node_id)
        if result['Code'] == 301204:
            return node
    return None


def get_no_pledge_node_list(node_list: List[Node]) -> List[Node]:
    """
    Get all the nodes that can be pledged
    :param node_list: Node list
    :return:
    """
    no_pledge_node_list = []
    for node in node_list:
        result = node.ppos.getCandidateInfo(node.node_id)
        if result['Code'] == 301204:
            no_pledge_node_list.append(node)
    return no_pledge_node_list


def get_pledge_list(func) -> list:
    """
    View the list of specified node IDs
    :param func: Query method, 1. List of current pledge nodes 2,
     the current consensus node list 3, real-time certifier list
    :return:
    """
    validator_info = func().get('Data')
    if not validator_info:
        time.sleep(10)
        validator_info = func().get('Data')
    validator_list = []
    for info in validator_info:
        validator_list.append(info.get('NodeId'))
    return validator_list


def check_node_in_list(nodeid, func) -> bool:
    """
    Check if the node is in the specified list
    :param nodeid: Node id
    :param func: Query method, 1. List of current pledge nodes 2,
     the current consensus node list 3, real-time certifier list
    :return:
    """
    data_dict = func()
    for data in data_dict["Data"]:
        if data["NodeId"] == nodeid:
            return True
    return False


def get_param_by_dict(data, *args):
    """
    Query parameter values​based on json data
    :param data: dictionary
    :param args: Key
    :return:
    """
    i = 0
    if isinstance(data, dict):
        for key in args:
            data = data.get(key)
            i = i + 1
            if isinstance(data, dict) and i > len(args):
                raise Exception("The parameters entered are incorrect.")
        return data

    raise Exception("Data format error")


def update_param_by_dict(data, key1, key2, key3, value):
    """
    Modify the json parameter
    :param data:
    :param key1:
    :param key2:
    :param key3:
    :param value:
    :return:
    """
    if isinstance(data, dict):
        if key3 is None:
            data[key1][key2] = value
        else:
            data[key1][key2][key3] = value
        return data
    raise Exception("Data format error")


def wait_block_number(node, block, interval=1):
    """
    Waiting until the specified block height
    :param node: Node
    :param block: Block height
    :param interval: Block interval, default is 1s
    :return:
    """
    current_block = node.block_number
    timeout = int((block - current_block) * interval * 1.5) + int(time.time())
    while int(time.time()) < timeout:
        log.info('The current block height is {}, waiting until {}'.format(node.block_number, block))
        if node.block_number > block:
            return
        time.sleep(1)
    raise Exception("Unable to pop out the block normally, the "
                    "current block height is: {}, the target block height is: {}".format(node.block_number, block))


def get_validator_term(node):
    """
    Get the nodeID with the highest term
    """
    msg = node.ppos.getValidatorList()
    term = []
    nodeid = []
    for i in msg["Data"]:
        term.append(i["ValidatorTerm"])
        nodeid.append(i["NodeId"])
    max_term = (max(term))
    term_nodeid_dict = dict(zip(term, nodeid))
    return term_nodeid_dict[max_term]


def get_max_staking_tx_index(node):
    """
    Get the nodeID of the largest transaction index
    """
    msg = node.ppos.getValidatorList()
    staking_tx_index_list = []
    nodeid = []
    for i in msg["Data"]:
        staking_tx_index_list.append(i["StakingTxIndex"])
        nodeid.append(i["NodeId"])
    max_staking_tx_index = (max(staking_tx_index_list))
    term_nodeid_dict = dict(zip(staking_tx_index_list, nodeid))
    return term_nodeid_dict[max_staking_tx_index]


def random_string(length=10) -> str:
    """
    Randomly generate a string of letters and numbers of a specified length
    :param length:
    :return:
    """
    return ''.join(
        random.choice(
            string.ascii_lowercase +
            string.ascii_uppercase +
            string.digits
        ) for _ in range(length)
    )


def get_pip_obj(nodeid, pip_obj_list: List[Pip]) -> Pip:
    """
    Get the pip object according to the node id
    :param nodeid:
    :param pip_obj_list:
    :return:
    """
    for pip_obj in pip_obj_list:
        if nodeid == pip_obj.node.node_id:
            return pip_obj


def get_pip_obj_list(nodeid_list, pip_obj_list: List[Pip]) -> List[Pip]:
    """
    Get a list of pip objects based on the node id list
    :param nodeid_list:
    :param pip_obj_list:
    :return:
    """
    new_pip_obj_list = []
    for nodeid in nodeid_list:
        new_pip_obj_list.append(get_pip_obj(nodeid, pip_obj_list))
    return new_pip_obj_list


def get_client_obj(nodeid, client_obj_list: List[Client]) -> Client:
    """
    Get the client object according to the node id
    :param nodeid:
    :param client_obj_list:
    :return:
    """
    for client_obj in client_obj_list:
        if nodeid == client_obj.node.node_id:
            return client_obj


def get_client_obj_list(nodeid_list, client_obj_list: List[Client]) -> List[Client]:
    """
    Get the client object list according to the node id list
    :param nodeid_list:
    :param client_obj_list:
    :return:
    """
    new_client_obj_list = []
    for nodeid in nodeid_list:
        new_client_obj_list.append(get_client_obj(nodeid, client_obj_list))
    return new_client_obj_list

def assert_code(result, code):
    assert result.get('Code') == code, "状态码错误，预期状态码：{}，实际状态码:{}".format(code, result.get("Code"))