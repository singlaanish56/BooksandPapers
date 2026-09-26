
#include <cstddef>
#include <fstream>
#include <fcntl.h>
#include <unistd.h>
struct FileContent{
    char* data;
    size_t size;
    char* name;  
};

typedef void read_overhead_function(repetition_tester* tester, FileContent* content);

static void ReadFRead(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        FILE* file = fopen(content->name, "rb");
        if(file) {

            BeginTime(tester);
            fread(content->data, 1, content->size, file);
            EndTime(tester);
            
            fclose(file);
        }
        
    }
}
static void ReadIStream(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        std::ifstream file(content->name, std::ios::binary);
        if(file) {

            BeginTime(tester);
            file.read(content->data, content->size);
            EndTime(tester);
        }
    }
}

static void ReadPOSIX(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        int file = open(content->name, O_RDONLY);
        if(file >= 0) {

            BeginTime(tester);
            read(file, content->data, content->size);
            EndTime(tester);

            close(file);
        }
    }
}
