// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.
// 
//  pseudo:
//      COUNT_SCREEN_WORDS = 8192
//      Reset:
//          i=0
//      Loop:   
//          if (COUNT_SCREEN_WORDS - i) <= 0: goto Reset
//          if (RAM[KBD] == 0): goto PrintWhite
//          else: goto PrintBlack
//          
//      PrintWhite:
//          addr = RAM[SCREEN + i]
//          RAM[addr] = 0
//          i++
//          goto Loop
//      PrintBlack:
//          addr = RAM[SCREEN + i]
//          RAM[addr] = -1
//          i++
//          goto Loop

// COUNT_SCREEN_WORDS = 8192 (32columns * 256rows)
@8192
D=A
@COUNT_SCREEN_WORDS
M=D

// RESET: i = 0
(RESET)
@i
M=0

(LOOP)

// if (count_screen_words - i) < 0 goto RESET
@COUNT_SCREEN_WORDS
D=M
@i
D=D-M
@RESET
D;JLE

// if RAM[KBD] == 0: goto PRINT_WHITE
@KBD
D=M
@PRINT_WHITE
D;JEQ

// else: goto: PRINT_BLACK
@PRINT_BLACK
0;JMP

(PRINT_WHITE)
// addr = RAM[SCREEN + i]
@SCREEN
D=A
@i
A=D+M

// fill white
M=0

// i++
@i
M=M+1

@LOOP
0;JMP

(PRINT_BLACK)
// addr = RAM[SCREEN + i]
@SCREEN
D=A
@i
A=D+M

// fill black
M=-1

// i++
@i
M=M+1

@LOOP
0;JMP
