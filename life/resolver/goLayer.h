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
	void gadget_createPBVar(int64 varAddr);
	unsigned char gadget_createGadget(int64 input0, int64 input1,
    							int64 input2, int64 res, int32 Type);
	void gadget_setVar(int64 varAddr, int64 Val, unsigned char is_unsigned);
	void gadget_setRetIndex(int64 RetAddr);
	void gadget_generateWitness();
	unsigned char GenerateProofAndResult(const char *pPKEY, char *pProof,
										int prSize, char *pResult, int resSize);
	unsigned char Verify(const char *pVKEY, const char *pPoorf,
						const char *pInput, const char *pOutput);
//};



