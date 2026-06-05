#include <complex.h>
#include <stdio.h>
#include <stdint.h>
#include <stdnoreturn.h>
#include <string.h>
const char *reg_table[] = {
    "al", "cl", "dl", "bl", "ah", "ch", "dh", "bh",
    "ax", "cx", "dx", "bx", "sp", "bp", "si", "di"
};

const char *r_m_encoding[] = {
    "bx + si", "bx + di", "bp + si", "bp + di", "si", "di", "bp", "bx"
};

short get_two_displacement_address(FILE* file) {
    unsigned char low;
    unsigned char high;

    if (fread(&low, 1, 1, file) != 1) {
        printf("failed to read displacement\n");
        return 0;
    }

    if (fread(&high, 1, 1, file) != 1) {
        printf("failed to read displacement\n");
        return 0;
    }

    unsigned short address = low | (high << 8);
    return address;
}


void reg_to_reg(const char* mnemonic,unsigned char byte1, FILE* file) {
    
    unsigned char byte2;

    if (fread(&byte2, 1, 1, file) != 1) {
        printf("failed to read byte2\n");
        return;
    }
    
    unsigned char d = (byte1 >> 1) & 1;
    unsigned char w = (byte1) & 1;
    unsigned char mod = (byte2 >> 6);
    unsigned char reg = (byte2 >> 3) & 0b111;
    unsigned char r_m = (byte2) & 0b111;
    unsigned int reg_index=reg_index = reg + (w * 8);
    unsigned int r_m_index = r_m + (w * 8);

    switch (mod){
        case 0b11: {
            if (d) {
                printf("%s %s, %s\n",mnemonic, reg_table[reg_index], reg_table[r_m_index]);
            } else {
                printf("%s %s, %s\n",mnemonic, reg_table[r_m_index], reg_table[reg_index]);
            }
            break;
        }
        case 0b00: {
            if(r_m!=0b110){
                if (d) {
                    printf("%s %s,[%s]\n",mnemonic, reg_table[reg_index], r_m_encoding[r_m]);
                } else {
                    printf("%s [%s], %s\n",mnemonic, r_m_encoding[r_m], reg_table[reg_index]);
                }
            } else {

                unsigned short address = get_two_displacement_address(file);
                if (d) {
                    printf("%s %s, [%d]\n",mnemonic,
                           reg_table[reg_index],
                           address);
                } else {
                    printf("%s [%d], %s\n",mnemonic,
                           address,
                           reg_table[reg_index]);
                }
            }
            break;
        }
        case 0b01: {
                int8_t displ;
                
                if (fread(&displ, 1, 1, file) != 1) {
                    printf("failed to read displacement\n");
                    return;
                }
        
                if (d) {
                    if (displ >= 0) {
                        printf("%s %s, [%s + %d]\n",mnemonic,
                            reg_table[reg_index],
                            r_m_encoding[r_m],
                            displ);
                    } else {
                        printf("%s %s, [%s - %d]\n",mnemonic,
                            reg_table[reg_index],
                            r_m_encoding[r_m],
                            -displ);
                    }
                } else {
                    if (displ >= 0) {
                        printf("%s [%s + %d], %s\n",mnemonic,
                            r_m_encoding[r_m],
                            displ,
                            reg_table[reg_index]);
                    } else {
                        printf("%s [%s - %d], %s\n",mnemonic,
                            r_m_encoding[r_m],
                            -displ,
                            reg_table[reg_index]);
                    }
                }

            break;
        }
        case 0b10: {

                int16_t displ = get_two_displacement_address(file);
        
                if (d) {
                    printf("%s %s, [%s + %d]\n",mnemonic,
                        reg_table[reg_index],
                        r_m_encoding[r_m],
                        displ);
                } else {
                    printf("%s [%s + %d], %s\n",mnemonic,
                        r_m_encoding[r_m],
                        displ,
                        reg_table[reg_index]);
                }

            break;
        }
        default:
            printf("unsupported mod for reg to reg: %d\n", mod);
            return;
    }
    
}

void imd_to_reg(const char * mnemonic,unsigned char byte1,  FILE* file) {
    unsigned char byte2;

    if (fread(&byte2, 1, 1, file) != 1) {
        printf("failed to read byte2\n");
        return;
    }
    
    unsigned char w = (byte1 >> 3) & 1;
    unsigned char reg = (byte1) & 0b111;

    unsigned int reg_index = reg + (w * 8);
    unsigned short imd = byte2;

    if (w) {
        unsigned char byte3;
    
        if (fread(&byte3, 1, 1, file) != 1) {
            printf("failed to read byte3\n");
            return;
        }   
        imd |= (byte3 << 8);
    }
    
    printf("%s %s, %d\n", mnemonic, reg_table[reg_index], imd);
}

void imd_to_mem(const char * mnemonic,unsigned char byte1, unsigned char byte2, FILE* file) {

    unsigned char w = byte1 & 1;

    unsigned char mod = (byte2 >> 6);
    unsigned char r_m = byte2 & 0b111;

    int displacement = 0;

    switch (mod) {

        case 0b00:

            if (r_m == 0b110) {
                displacement = get_two_displacement_address(file);
            }

            break;

        case 0b01: {

            int8_t displ;

            if (fread(&displ, 1, 1, file) != 1) {
                printf("failed to read displacement\n");
                return;
            }

            displacement = displ;
            break;
        }

        case 0b10:

            displacement = get_two_displacement_address(file);
            break;
        case 0b11: {
            
            int use_sign_extension =
                strcmp(mnemonic, "add") == 0 ||
                strcmp(mnemonic, "sub") == 0 ||
                strcmp(mnemonic, "cmp") == 0;
            
            unsigned char s = use_sign_extension ? ((byte1 >> 1) & 1) : 0;
                unsigned char reg_index = r_m + (w * 8);
            
                int immediate;
            
                if (w == 0) {
            
                    int8_t imd;
            
                    if (fread(&imd, 1, 1, file) != 1) {
                        printf("failed to read immediate\n");
                        return;
                    }
            
                    immediate = imd;
            
                } else {
            
                    if (s == 1) {
            
                        int8_t imd;
            
                        if (fread(&imd, 1, 1, file) != 1) {
                            printf("failed to read immediate\n");
                            return;
                        }
            
                        immediate = imd;
            
                    } else {
            
                        int16_t imd = get_two_displacement_address(file);
                        immediate = imd;
                    }
                }
            
                printf("%s %s, %d\n",
                       mnemonic,
                       reg_table[reg_index],
                       immediate);
            
                return;
            }
        default:
            printf("unsupported mod in immediate-to-memory: %d\n", mod);
            return;
    }

    int use_sign_extension =
        strcmp(mnemonic, "add") == 0 ||
        strcmp(mnemonic, "sub") == 0 ||
        strcmp(mnemonic, "cmp") == 0;
    
    unsigned char s = use_sign_extension ? ((byte1 >> 1) & 1) : 0;
    
    int immediate;
    
    if (w == 0) {
    
        int8_t imd;
    
        if (fread(&imd, 1, 1, file) != 1) {
            printf("failed to read immediate\n");
            return;
        }
    
        immediate = imd;
    
    } else {
    
        if (s == 1) {
    
            int8_t imd;
    
            if (fread(&imd, 1, 1, file) != 1) {
                printf("failed to read immediate\n");
                return;
            }
    
            immediate = imd;
    
        } else {
    
            int16_t imd = get_two_displacement_address(file);
            immediate = imd;
        }
    }

    switch (mod) {

        case 0b00:

            if (r_m == 0b110) {

                printf("%s [%d], %s %d\n",mnemonic,
                       displacement,
                       w ? "word" : "byte",
                       immediate);

            } else {

                printf("%s [%s], %s %d\n",mnemonic,
                       r_m_encoding[r_m],
                       w ? "word" : "byte",
                       immediate);
            }

            break;

        case 0b01:

            if (displacement >= 0) {

                printf("%s [%s + %d], %s %d\n",mnemonic,
                       r_m_encoding[r_m],
                       displacement,
                       w ? "word" : "byte",
                       immediate);

            } else {

                printf("%s [%s - %d], %s %d\n",mnemonic,
                       r_m_encoding[r_m],
                       -displacement,
                       w ? "word" : "byte",
                       immediate);
            }

            break;

        case 0b10:

            if (displacement >= 0) {

                printf("%s [%s + %d], %s %d\n",mnemonic,
                       r_m_encoding[r_m],
                       displacement,
                       w ? "word" : "byte",
                       immediate);

            } else {

                printf("%s [%s - %d], %s %d\n",mnemonic,
                       r_m_encoding[r_m],
                       -displacement,
                       w ? "word" : "byte",
                       immediate);
            }

            break;
    }
}

void acc_and_mem(const char * mnemonic,unsigned char byte1, FILE* file) {

    unsigned char d = (byte1 >> 1) &1;
    unsigned char w = byte1 & 1;
    unsigned char byte2;
    
    if (fread(&byte2, 1, 1, file) != 1) {
        printf("failed to read byte2\n");
        return;
    }

    unsigned short memory = byte2;
    
    if (w) {
        unsigned char byte3;
    
        if (fread(&byte3, 1, 1, file) != 1) {
            printf("failed to read byte3\n");
            return;
        }   
        memory |= (byte3 << 8);
    }
    
    if(d){
        printf("%s [%d], %s\n",mnemonic, memory, reg_table[w*8]);
    } else {
        printf("%s %s, [%d]\n",mnemonic, reg_table[w*8], memory);
    }
}

void acc_immediate(const char *mnemonic, unsigned char byte1, FILE *file) {

    unsigned char w = byte1 & 1;

    int immediate;

    if (w) {

        int16_t imd = get_two_displacement_address(file);
        immediate = imd;

    } else {

        int8_t imd;

        if (fread(&imd, 1, 1, file) != 1) {
            printf("failed to read immediate\n");
            return;
        }

        immediate = imd;
    }

    printf("%s %s, %d\n",
           mnemonic,
           w ? "ax" : "al",
           immediate);
}

int main(int argc, char *argv[]) {


    FILE *file = fopen(argv[1], "rb");
    if (file == NULL) {
        printf("Failed to open file\n");
        return 1;
    }

    // read the file for the binary instructions on each line and print it
    unsigned char byte1;


    printf("%s\n\n", "bits 16");
    while (fread(&byte1, 1, 1, file)) {

        if ((byte1 >> 1) == 0b0000010) {
            acc_immediate("add", byte1, file);
            continue;
        }
        
        else if ((byte1 >> 1) == 0b0010110) {
            acc_immediate("sub", byte1, file);
            continue;
        }
        
        else if ((byte1 >> 1) == 0b0011110) {
            acc_immediate("cmp", byte1, file);
            continue;
        }
        else if ((byte1 >> 2) == 0b100010){
            reg_to_reg("mov",byte1, file);
            continue;
        }else if ((byte1 >>2)== 0b000000){
            reg_to_reg("add",byte1, file);
            continue;
        }
        else if ((byte1 >> 2)== 0b001010) {
            reg_to_reg("sub",byte1, file);
            continue;
        }else if ((byte1 >> 2) == 0b001110) {
            reg_to_reg("cmp",byte1, file);
            continue;
        }else if ((byte1 >> 2) == 0b100000) {
            unsigned char byte2;
            if (fread(&byte2, 1, 1, file) != 1) {
                printf("failed to read modrm\n");
                return 0;
            }
            unsigned char reg = (byte2>>3) & 0b111;
            imd_to_mem(reg==0b000?"add":(reg==0b111?"cmp":"sub"),byte1,byte2, file);
            continue;
        }

        unsigned char opcode = (byte1 >> 4);
        
        switch (opcode) {
            case 0b1011:
                imd_to_reg("mov",byte1, file);
                break;
            case 0b1100:
                unsigned char byte2;
                if (fread(&byte2, 1, 1, file) != 1) {
                    printf("failed to read modrm\n");
                    return 0;
                }
                imd_to_mem("mov",byte1,byte2, file);
                break;
            case 0b1010:
                acc_and_mem("mov",byte1, file);
                break;
            default:
                printf("unsupported opcode: %d\n", opcode);
                continue;
        }
                            
    }
    fclose(file);
    return 0;
}