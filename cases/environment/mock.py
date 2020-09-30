#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#   @Time    : 2019/12/30 20:19
#   @Author  : PlatON-Developer
#   @Site    : https://github.com/PlatONnetwork/

from common.log import log

class MockSftp:
    def put(self, a, b):
        pass

    def get(self, a, b):
        pass


class MockSsh:
    def exec_command(self, cmd):
        return MockStdin(), MockStdout(cmd), MockStderr()


class MockStdin:

    def write(self, cmd):
        pass


class MockStdout:
    def __init__(self, cmd):
        self.cmd = cmd

    def readlines(self):
        if "init" in self.cmd:
            return []
        return [""]


class MockStderr:
    pass


class MockT:
    def close(self):
        pass


def mock_connect_linux():
    log.info("mock server connect linux")
    return MockSsh(), MockSftp(), MockT()
