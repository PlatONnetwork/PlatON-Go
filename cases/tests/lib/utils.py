# -*- coding: utf-8 -*-
import json, rlp
import time
import random
import string
from decimal import Decimal

from hexbytes import HexBytes
from environment.node import Node
from common.log import log
from common.key import get_pub_key
from typing import List


def decorator_sleep(func):
    def wrap():
        result = func()
        if result is None:
            time.sleep(5)
            result = func()
        return result
    return wrap


def find_proposal(proposal_list, block_number):
    for proposal in proposal_list:
        if proposal_effective(proposal, block_number):
            return proposal


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
        blocknumber = node.block_number
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


def get_pledge_list(func, nodeid=None) -> list:
    """
    View the list of specified node IDs
    :param func: Query method, 1. List of current pledge nodes 2,
     the current consensus node list 3, real-time certifier list
    :return:
    """
    validator_info = func().get('Ret')
    if validator_info == "Getting verifierList is failed:The validator is not exist":
        time.sleep(10)
        validator_info = func().get('Ret')
    if validator_info == "Getting candidateList is failed:CandidateList info is not found":
        time.sleep(10)
        validator_info == func().get('Ret')
    if not nodeid:
        validator_list = []
        for info in validator_info:
            validator_list.append(info.get('NodeId'))
        return validator_list
    else:
        for info in validator_info:
            if nodeid == info.get('NodeId'):
                return info.get('RewardPer'), info.get('NextRewardPer')
        raise Exception('Nodeid {} not in the list'.format(nodeid))


def check_node_in_list(nodeid, func) -> bool:
    """
    Check if the node is in the specified list
    :param nodeid: Node id
    :param func: Query method, 1. List of current pledge nodes 2,
     the current consensus node list 3, real-time certifier list
    :return:
    """
    data_dict = func()
    for data in data_dict["Ret"]:
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
    elif isinstance(data, str):
        data = json.loads(data)
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
    elif isinstance(data, str):
        jsoninfo = json.loads(data)
        if key3 is None:
            jsoninfo[key1][key2] = value
        else:
            jsoninfo[key1][key2][key3] = value
        jsondata = json.dumps(jsoninfo)
        return jsondata
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
    if 0 < block - current_block <= 10:
        timeout = 10 + int(time.time()) + 50
    elif block - current_block > 10:
        timeout = int((block - current_block) * interval * 1.5) + int(time.time()) + 50
    else:
        log.info('current block {} is greater than block {}'.format(node.block_number, block))
        return
    print_t = 0
    while int(time.time()) < timeout:
        print_t += 1
        if print_t == 10:
            # Print once every 10 seconds to avoid printing too often
            log.info('The current block height is {}, waiting until {}'.format(node.block_number, block))
            print_t = 0
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
    term_list = []
    nodeid_list = []
    for i in msg["Ret"]:
        term_list.append(i["ValidatorTerm"])
        nodeid_list.append(i["NodeId"])

    max_term = (max(term_list))
    log.info("Maximum tenure{}".format(max_term))
    term_nodeid_dict = dict(zip(nodeid_list, term_list))

    max_term_nodeid = []

    for key in term_nodeid_dict:
        value = term_nodeid_dict[key]
        if value == max_term:
            max_term_nodeid.append(key)
    return max_term_nodeid


def get_max_staking_tx_index(node):
    """
    Get the nodeID of the largest transaction index
    """
    msg = node.ppos.getValidatorList()
    staking_tx_index_list = []
    nodeid = []
    for i in msg["Ret"]:
        staking_tx_index_list.append(i["StakingTxIndex"])
        nodeid.append(i["NodeId"])
    max_staking_tx_index = (max(staking_tx_index_list))
    term_nodeid_dict = dict(zip(staking_tx_index_list, nodeid))
    return term_nodeid_dict[max_staking_tx_index]


def get_block_count_number(node, number):
    """
    Get the number of verifier blocks
    :param url: node url
    :param cycle: Consensus cycle
    :return:
    """
    current_block = node.block_number
    count = 0
    for i in range(number - 1):
        nodeid = node.eth.ecrecover(current_block)
        current_block = current_block - 1
        if nodeid == node.node_id:
            count = count + 1
    return count


def random_string(length=10) -> str:
    """
    Randomly generate a string of letters and numbers of a specified length
    :param length:
    :return:
    """
    return ''.join(
        random.choice(
            string.ascii_lowercase
            + string.ascii_uppercase
            + string.digits
        ) for _ in range(length)
    )


def assert_code(result, code):
    '''
    assert the ErrorCode
    :param result:
    :param code:
    :return:
    '''
    if isinstance(result, int):
        assert result == code, "code error，expect：{}，actually:{}".format(code, result)
    else:
        assert result.get('Code') == code, "code error，expect：{}，actually:{}".format(code, result)


def von_amount(amonut, base):
    """
    Get von amount
    :param amonut:
    :param base:
    :return:
    """
    return int(Decimal(str(amonut)) * Decimal(str(base)))


def get_governable_parameter_value(client_obj, parameter, flag=None):
    """
    Get governable parameter value
    :return:
    """
    # Get governable parameters
    govern_param = client_obj.pip.pip.listGovernParam()
    parameter_information = govern_param['Ret']
    for i in parameter_information:
        if i['ParamItem']['Name'] == parameter:
            log.info("{} ParamValue: {}".format(parameter, i['ParamValue']['Value']))
            log.info("{} Param old Value: {}".format(parameter, i['ParamValue']['StaleValue']))
            if not flag:
                return i['ParamValue']['Value']
            else:
                return int(i['ParamValue']['Value']), int(i['ParamValue']['StaleValue'])


def get_the_dynamic_parameter_gas_fee(data):
    """
    Get the dynamic parameter gas consumption cost
    :return:
    """
    zero_number = 0
    byte_group_length = len(data)
    for i in data:
        if i == 0:
            zero_number = zero_number + 1
    non_zero_number = byte_group_length - zero_number
    dynamic_gas = non_zero_number * 68 + zero_number * 4
    return dynamic_gas


def get_getDelegateReward_gas_fee(client, staking_num, uncalcwheels, gasprice=None):
    data = rlp.encode([rlp.encode(int(5000))])
    if gasprice is None:
        gasprice = client.node.eth.gasPrice
    gas = get_the_dynamic_parameter_gas_fee(data) + 8000 + 3000 + 21000 + staking_num * 1000 + uncalcwheels * 100
    return gas * gasprice
