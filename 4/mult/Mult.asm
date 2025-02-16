// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.
// 
//  pseudo:
//      valueToAdd = RAM[1]
//      n = RAM[0]
//      i = 1
//      sum = 0
//      
//      LOOP:
//          if i > n goto STOP
//          sum = sum + valueToAdd
//          i = i + 1
//          goto LOOP
//      
//      STOP:
//          RAM[2] = sum

// valueToAdd = RAM[1]
@R1
D=M
@valueToAdd
M=D

// n = RAM[0]
@R0
D=M
@n
M=D

// i = 1
@i
M=1

// sum = 0
@sum
M=0

// checking for 0 to pass the test R0=0, R1=0 in 20 clock cycles 
// but if the test had contained 22 clock cycles, the test would have passed
@n
D=M
@STOP
D;JEQ

(LOOP)

// if n - i < 0 goto STOP
@i
D=M
@n
D=M-D
@STOP
D;JLT

// sum = sum + valueToAdd
@valueToAdd
D=M
@sum
M=D+M

// i = i + 1
@i
M=M+1

// goto LOOP
@LOOP
0;JMP

(STOP)
// RAM[2] = sum
@sum
D=M
@R2
M=D

(END)
@END
0;JMP
