TEXT ·rdtsc(SB),7,$0
	MOVQ	hi+0(FP), R9
	MOVQ	lo+8(FP), R10
	RDTSC
	MOVL	DX, (R9)
	MOVL	AX, (R10)
	RET
