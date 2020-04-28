@echo off
rem echo "Parse params....."
set VERSION=%1
set CONTRACT=%2
set TARGET=%3
rem echo "Choose solc version: "%VERSION%
rem echo "Source solidity contract path: "%CONTRACT%
rem echo "Compiled abi/bytecode file target path: "%TARGET%
rem echo "Enter solc binary dir ....."
rem echo %cd%
cd ..\solc\solc-windows-%VERSION%
rem echo "Run solc.exe command to compile contract ...."
rem echo %cd%
solc.exe -o %TARGET% --bin --abi --overwrite %CONTRACT%
exit