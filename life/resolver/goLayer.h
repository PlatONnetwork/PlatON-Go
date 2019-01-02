#pragma once

#include <stdlib.h>

typedef int	int32;
typedef unsigned int uint32;
typedef long long	int64;
typedef unsigned long long	uint64;

//extern "C" 
//{
    void gadget_initEnv();
	void gadget_uninitEnv();
	void gadget_createPBVar(const char *varName);
	unsigned char gadget_createGadget(const char *input0Name, const char *input1Name,
    							const char *input2Name, const char *resName, int32 Type);
	void gadget_setVar(const char *varName, uint64 Val);
	void gadget_generateWitness();
	unsigned char GenerateProofAndResult(const char *pPKEY, char *pProof,
										int prSize, char *pResult, int resSize);
	unsigned char Verify(const char *pVKEY, const char *pPoorf,
						const char *pInput, const char *pOutput);
//};



