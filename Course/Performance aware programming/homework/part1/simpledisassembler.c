#include <stdio.h>
#include <stdint.h>

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


void reg_to_reg(unsigned char byte1, FILE* file) {
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
                printf("mov %s, %s\n", reg_table[reg_index], reg_table[r_m_index]);
            } else {
                printf("mov %s, %s\n", reg_table[r_m_index], reg_table[reg_index]);
            }
            break;
        }
        case 0b00: {
            if(r_m!=0b110){
                if (d) {
                    printf("mov %s,[%s]\n", reg_table[reg_index], r_m_encoding[r_m]);
                } else {
                    printf("mov [%s], %s\n", r_m_encoding[r_m], reg_table[reg_index]);
                }
            } else {

                unsigned short address = get_two_displacement_address(file);
                if (d) {
                    printf("mov %s, [%d]\n",
                           reg_table[reg_index],
                           address);
                } else {
                    printf("mov [%d], %s\n",
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
                    printf("mov %s, [%s + %d]\n",
                        reg_table[reg_index],
                        r_m_encoding[r_m],
                        displ);
                } else {
                    printf("mov [%s + %d], %s\n",
                        r_m_encoding[r_m],
                        displ,
                        reg_table[reg_index]);
                }

            break;
        }
        case 0b10: {

                int16_t displ = get_two_displacement_address(file);
        
                if (d) {
                    printf("mov %s, [%s + %d]\n",
                        reg_table[reg_index],
                        r_m_encoding[r_m],
                        displ);
                } else {
                    printf("mov [%s + %d], %s\n",
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

void imd_to_reg(unsigned char byte1,  FILE* file) {
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
    
    printf("mov %s, %d\n", reg_table[reg_index], imd);
}

void reg_to_mem(unsigned char byte1, FILE* file) {
    
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

        unsigned char opcode = (byte1 >> 4);

        switch (opcode) {
            case 0b1000:
                reg_to_reg(byte1, file);
                break;
            case 0b1011:
                imd_to_reg(byte1, file);
                break;
            case 0b1100:
                reg_to_mem(byte1, file);
                break;
            default:
                printf("unsupported opcode: %d\n", opcode);
                continue;
        }
                            
    }
    fclose(file);
    return 0;
}