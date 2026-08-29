
#include <cstdio>
#include <cstdlib>
#include <filesystem>
#include <iostream>
#include <system_error>
#include <algorithm>
#include<string.h>
#include <fstream>
#include <random>
#include <threads.h>
#include <iomanip>

//#include "my_timer_profile.cpp"
#include "casey_recursive_timer.cpp"

/* Casey's Code */
typedef double f64;
typedef uint64_t u64;

static f64 Square(f64 A)
{
    f64 Result = (A*A);
    return Result;
}

static f64 RadiansFromDegrees(f64 Degrees)
{
    f64 Result = 0.01745329251994329577 * Degrees;
    return Result;
}

// NOTE(casey): EarthRadius is generally expected to be 6372.8
static f64 ReferenceHaversine(f64 X0, f64 Y0, f64 X1, f64 Y1, f64 EarthRadius)
{
    /* NOTE(casey): This is not meant to be a "good" way to calculate the Haversine distance.
       Instead, it attempts to follow, as closely as possible, the formula used in the real-world
       question on which these homework exercises are loosely based.
    */

    f64 lat1 = Y0;
    f64 lat2 = Y1;
    f64 lon1 = X0;
    f64 lon2 = X1;

    f64 dLat = RadiansFromDegrees(lat2 - lat1);
    f64 dLon = RadiansFromDegrees(lon2 - lon1);
    lat1 = RadiansFromDegrees(lat1);
    lat2 = RadiansFromDegrees(lat2);

    f64 a = Square(sin(dLat/2.0)) + cos(lat1)*cos(lat2)*Square(sin(dLon/2));
    f64 c = 2.0*asin(sqrt(a));

    f64 Result = EarthRadius * c;

    return Result;
}
/* End Casey's Code */

struct FileContent{
    char *data;
    size_t size;
};

struct Pairs{
    double x0;
    double y0;
    double x1;
    double y1;
};

FileContent readFile(const char* filename) {
    TimeFunction;
    
    FileContent result = {};

    FILE* file = fopen(filename, "rb");
    if(file == NULL) {
        return result;
    }

    fseek(file, 0, SEEK_END);
    result.size = ftell(file);
    rewind(file);

    result.data = (char *)malloc(result.size+1);
    fread(result.data, 1, result.size, file);
    result.data[result.size] = '\0';

    fclose(file);
    return result;
}

void SkipWhitespace(char *&at) {
    while((*at==' ') || (*at=='\t') || (*at=='\n') || (*at=='\r')) ++at;
}

void SkipUntil(char *&at, char c){
    while(*at && *at != c) ++at;
    if(*at == c) ++at;
}

double ParseNumber(char *&at) {

    SkipWhitespace(at);
    char *end;
    double value = strtod(at, &end);
    at = end;
    return value;
}

void processJson(const FileContent& jsonContent, Pairs* pairs, int numberOfPairs, int& actualPairs) {
    //std::cout<<"heel"<<std::endl;
    TimeFunction;
    
    char *at = jsonContent.data;

    {
        //BlockFunction("processJson for loop");
        TimeBlock("processJson for loop");
        SkipUntil(at, '[');
        for(int i=0;i<numberOfPairs;i++)
        {
            if(*at == ']') break;
            SkipUntil(at, ':');
            pairs[i].x0 = ParseNumber(at);
    
            SkipUntil(at, ':');
            pairs[i].y0 = ParseNumber(at);
    
            SkipUntil(at, ':');
            pairs[i].x1 = ParseNumber(at);
    
            SkipUntil(at, ':');
            pairs[i].y1 = ParseNumber(at);
    
    
            actualPairs++;
            SkipUntil(at, '}');
            SkipWhitespace(at);
    
            if(*at==',') ++at;
        }
    }
}

double computeHarvensineAndSum(const Pairs* pairs, int numberOfPairs) {
    //TimerFunction;
    TimeFunction;
    
    double sum = 0;
    double mul = 1/(double)numberOfPairs;
    double earthRadius = 6372.8;
    
    for(int i=0;i<numberOfPairs;i++) {
        double haversineDistance = ReferenceHaversine(pairs[i].x0, pairs[i].y0, pairs[i].x1, pairs[i].y1, earthRadius);
        sum+=(mul*haversineDistance);
    }
    

    return sum;
}

double parseAndVerifyResult(const FileContent &answersContent, int numberOfPairs) {

    //TimerFunction;
    TimeFunction;
    
    char* at = answersContent.data;
    double sum = 0;
    double mul = 1/(double)numberOfPairs;

    for(int i=0;i<numberOfPairs;i++) {
        SkipWhitespace(at);
        double haversineDistance = ParseNumber(at);
        sum+=(mul*haversineDistance);
        SkipWhitespace(at);
    }

    return sum;

}

void printTimeProf(const char* name, u64 totaltime, u64 start, u64 end) {
    u64 duration = end - start;
    f64 percent = 100.0 * (f64)duration / (f64)totaltime;
    printf("%s: %llu  (%0.4f%%)\n", name, duration, percent);
}

int main(int argc, char* argv[]){

    //StartProfile();
    BeginProfile();
    // u64 ProfBegin=0;
    // u64 ProfRead=0;
    // u64 ProfWholeFileJson=0;
    // u64 ProfWholeFileValues=0;
    // u64 ProfParseJson=0;
    // u64 ProfComputeHarversine=0;
    // u64 ProfParseAnswers=0;
    // u64 ProfEnd=0;

   // ProfBegin = ReadCPUTimer();

    if(argc <2 || argc >3) {
        std::cerr << "Usage: " << argv[0] << " <Json File> <Optional: Answers File>" << std::endl;
        return 1;
    }

    std::string json_file = argv[1];
    //std::string answers_file = (argc == 3) ? argv[2] : "";

//ProfRead = ReadCPUTimer();
    FileContent jsonContent = readFile(json_file.c_str());
    if(jsonContent.data == NULL) {
        std::cerr << "Error: Could not read JSON file " << json_file << std::endl;
        return 1;
    }
   // ProfWholeFileJson = ReadCPUTimer();

    constexpr size_t MinBytesPerPair = sizeof(Pairs); // 32 bytes
    int estimatedPairs = jsonContent.size / MinBytesPerPair;
    int actualPairs = 0;
    Pairs* pairs = new Pairs[estimatedPairs];
    processJson(jsonContent, pairs, estimatedPairs, actualPairs);
   // ProfParseJson = ReadCPUTimer();
    double sum = computeHarvensineAndSum(pairs, actualPairs);
   // ProfComputeHarversine = ReadCPUTimer();

    //std::cout<<"Input Size: "<<jsonContent.size<<"\n"<<"Pair Count: "<<actualPairs<<"\n"<<"Sum: "<<sum<<"\n";

    free(jsonContent.data);
    jsonContent={};

    //reference points file , to read it and then post the reference results
    if(argc == 3) {
        std::string answers_file = argv[2];
        FileContent answersContent = readFile(answers_file.c_str());
        if(answersContent.data == NULL) {
            std::cerr << "Error: Could not read answers file " << answers_file << std::endl;
            return 1;
        }
       // ProfWholeFileValues = ReadCPUTimer();

        double actualSum = parseAndVerifyResult(answersContent, actualPairs);
       // ProfParseAnswers = ReadCPUTimer();
        //std::cout << "Validation\n" << "Reference Sum: " << actualSum << "\nDifference Sum: " << actualSum - sum << "\n";
        free(answersContent.data);
        answersContent = {};
    }
    // ProfEnd = ReadCPUTimer();

    // u64 totaltime = ProfEnd - ProfBegin;
    // u64 guessedCPUFreq = EstimateCPUTimerFreq();
    // if(guessedCPUFreq)
    // {
    //     printf("\nTotal time: %0.4fms (CPU freq %llu)\n", 1000.0 * (f64)totaltime / (f64)guessedCPUFreq, guessedCPUFreq);
    // }

    // printTimeProf("Read Arg", totaltime, ProfBegin, ProfRead);
    // printTimeProf("Read Json File", totaltime, ProfRead, ProfWholeFileJson);
    // printTimeProf("Parse Json", totaltime, ProfWholeFileJson, ProfParseJson);
    // printTimeProf("Compute Harversine", totaltime, ProfParseJson, ProfComputeHarversine);
    // printTimeProf("Read Values File", totaltime, ProfComputeHarversine, ProfWholeFileValues);
    // printTimeProf("Parse Verify  Answers", totaltime, ProfWholeFileValues, ProfParseAnswers);

    //End();
    EndAndPrintProfile();
}
