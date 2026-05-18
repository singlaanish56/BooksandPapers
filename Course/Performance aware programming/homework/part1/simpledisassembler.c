#include <stdio.h>

const char *reg_table[] = {
    "al", "cl", "dl", "bl", "ah", "ch", "dh", "bh",
    "ax", "cx", "dx", "bx", "sp", "bp", "si", "di"
};

int main(int argc, char *argv[]) {


    FILE *file = fopen(argv[1], "rb");
    if (file == NULL) {
        printf("Failed to open file\n");
        return 1;
    }

    // read the file for the binary instructions on each line and print it
    unsigned char byte1;
    unsigned char byte2;

    printf("%s\n\n", "bits 16");
    while (fread(&byte1, 1, 1, file) && fread(&byte2, 1, 1, file)) {

        unsigned char opcode = (byte1 >> 2);
        unsigned char d = (byte1 >> 1) & 1;
        unsigned char w = (byte1) & 1;

        unsigned char mod = (byte2 >> 6);
        unsigned char reg = (byte2 >> 3) & 0b111;
        unsigned char r_m = (byte2) & 0b111;


        if(mod != 0b11) {
            printf("unsupported mod: %d\n", mod);
            continue;
        }

        if(opcode != 0b100010) {
            printf("unsupported opcode: %d\n", opcode);
            continue;
        }

        unsigned int reg_index = reg + (w * 8);
        unsigned int r_m_index = r_m + (w * 8);
        
        if (d) {
            printf("mov %s, %s\n", reg_table[reg_index], reg_table[r_m_index]);
        } else {
            printf("mov %s, %s\n", reg_table[r_m_index], reg_table[reg_index]);
        }
                
    }
    fclose(file);
    return 0;
}