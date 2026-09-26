


#include <cstdint>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <math.h>
#include <sys/stat.h>

#define ArrayCount(Array) (sizeof(Array) / sizeof(Array[0]))

typedef uint8_t u8;
typedef uint64_t u64;
typedef uint32_t u32;
typedef double f64;

#include "repetition_tester.cpp"
#include "read_overhead_tester.cpp"

struct test_function{
    const char* Name;
    read_overhead_function * Func;
};

test_function test_functions[] = {
    { "fread", ReadFRead },
    { "istream", ReadIStream },
    { "posix", ReadPOSIX }
};

int main(int argc, char* argv[]){

    //infinite loop for the the tests
    if(argc >= 2){
        FileContent file_data{};
        file_data.name = argv[1];

        struct stat Stat;
        stat(file_data.name, &Stat);

        

        file_data.size = Stat.st_size;
        file_data.data = (char *) malloc(file_data.size);
            
        if(file_data.data){
            repetition_tester testers[ArrayCount(test_functions)][AllocationType_Count] = {};
            for(;;){
                // loop over the functions of that are supposed to be tested
                for(u32 i = 0; i < ArrayCount(test_functions); ++i) {

                    for(u32 j=0;j<AllocationType_Count;++j){

                        file_data.alloc = (allocation_type)j;
                        test_function func = test_functions[i];
                        repetition_tester *tester = &testers[i][j];
                        
                        printf("\n--- %s%s%s ---\n",
                               GetAllocationType(file_data.alloc),
                               file_data.alloc ? " + " : "",
                               func.Name);
                        InitializeParamForFunc(tester, file_data.size);
                        func.Func(tester, &file_data);
                    }

                }
            }
        }
    }
    return 0;
}