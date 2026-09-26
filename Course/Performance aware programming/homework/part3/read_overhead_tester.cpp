
#include <alloca.h>
#include <cstddef>
#include <filesystem>
#include <fstream>
#include <fcntl.h>
#include <unistd.h>


typedef void read_overhead_function(repetition_tester* tester, FileContent* content);

static void ReadFRead(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        FILE* file = fopen(content->name, "rb");
        if(file) {
            char *buffer = Allocate(content);
            
            BeginTime(tester);
            fread(buffer, 1, content->size, file);
            EndTime(tester);

            Deallocate(content,buffer);
            fclose(file);
        }
        
    }
}

static void ReadIStream(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        std::ifstream file(content->name, std::ios::binary);
        if(file) {
            char* buffer = Allocate(content);

            BeginTime(tester);
            file.read(buffer, content->size);
            EndTime(tester);

            Deallocate(content,buffer);
            file.close();
        }
    }
}

static void ReadPOSIX(repetition_tester* tester, FileContent* content) {

    while(isTesting(tester)){

        int file = open(content->name, O_RDONLY);
        if(file >= 0) {
            char* buffer = Allocate(content);
            
            BeginTime(tester);
            read(file, buffer, content->size);
            EndTime(tester);

            Deallocate(content,buffer);
            close(file);
        }
    }
}
