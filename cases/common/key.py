import os
import sys
from common.abspath import abspath
from eth_keys import (
    keys,
)
from eth_utils.curried import (
    keccak,
    text_if_str,
    to_bytes,
)


def generate_key():
    """
    generate node public private key
    :return:
        privateKey
        publicKey
    """
    extra_entropy = ''
    extra_key_bytes = text_if_str(to_bytes, extra_entropy)
    key_bytes = keccak(os.urandom(32) + extra_key_bytes)
    privatekey = keys.PrivateKey(key_bytes)
    return privatekey.to_hex()[2:], keys.private_key_to_public_key(privatekey).to_hex()[2:]


def generate_blskey():
    """
    generate bls public and private keys
    :return:
        bls_private_key
        bls_public_key
    """
    if sys.platform in "linux,linux2":
        tool_file = abspath("tool/linux/keytool")
        run("chmod +x {}".format(tool_file))
    else:
        tool_file = abspath("tool/win/keytool.exe")
    keypair = run("{} genblskeypair".format(tool_file))
    if not keypair:
        raise Exception("unable to use generate blskey tool")
    lines = keypair.split("\n")
    bls_private_key = ""
    bls_public_key = ""
    for l in lines:
        kv = l.split(":")
        if kv[0] == "PrivateKey":
            bls_private_key = kv[1].strip()
        elif kv[0] == "PublicKey ":
            bls_public_key = kv[1].strip()
    if not bls_private_key or not bls_public_key:
        raise Exception("Blskey cannot be generated")
    return bls_private_key, bls_public_key


def run(cmd):
    """
    The machine executes the cmd command and gets the result
    :param cmd:
    :return:
    """
    r = os.popen(cmd)
    out = r.read()
    r.close()
    return out


def get_pub_key(url, block):
    """
    obtain signature nodes based on block information
    :param url: node url
    :param block: block height
    :return:
    """
    if sys.platform in "linux,linux2":
        tool_file = abspath("tool/linux/get_pubkey_for_blocknum")
        run("chmod +x {}".format(tool_file))
    else:
        tool_file = abspath("tool/win/get_pubkey_for_blocknum.exe")
    output = run("{} -url={} -blockNumber={}".format(tool_file, url, block))
    if not output:
        raise Exception("unable to use get node id tool")
    if "1111" in output or "2222" in output:
        raise Exception("get node id exception：{}".format(output))
    return output.strip("\n")


def mock_duplicate_sign(dtype, sk, blskey, block_number, epoch=0, view_number=0, block_index=0, index=0):
    """
    forged double sign
    :param dtype:
    :param sk:
    :param blskey:
    :param block_number:
    :param epoch:
    :param view_number:
    :param block_index:
    :param index:
    :return:
    """
    if sys.platform in "linux,linux2":
        tool_file = abspath("tool/linux/duplicateSign")
        run("chmod +x {}".format(tool_file))
    else:
        tool_file = abspath("tool/win/duplicateSign.exe")

    output = run("{} -dtype={} -sk={} -blskey={} -blockNumber={} -epoch={} -viewNumber={} -blockIndex={} -vindex={}".format(
        tool_file, dtype, sk, blskey, block_number, epoch, view_number, block_index, index))
    if not output:
        raise Exception("unable to use double sign tool")
    return output.strip("\n")


if __name__ == "__main__":
    print(generate_blskey())
