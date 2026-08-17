#include <cstdint>
#include "platform_metrics.cpp"
#include <iostream>

typedef uint64_t u64 ;
typedef double f64;

struct ProfilingResult{
    const char* name;
    u64 duration;
};

static ProfilingResult results[1024];
static int resultCount = 0;
static u64 ProfileStartTime;


static void StartProfile(){ 
    ProfileStartTime = ReadCPUTimer();
    resultCount=0;
};
static void StopProfile(){
    u64 endTime = ReadCPUTimer();
    u64 totalTime = endTime - ProfileStartTime;
    u64 guessedCPUFreq = EstimateCPUTimerFreq();
    if(guessedCPUFreq)
    {
        printf("\nTotal time: %0.4fms (CPU freq %llu)\n", 1000.0 * (f64)totalTime / (f64)guessedCPUFreq, guessedCPUFreq);
    }

    for(int i = 0; i < resultCount; i++)
    {
        f64 percent = 100.0 * (f64)results[i].duration / (f64)totalTime;
        printf("%s: %llu  (%0.4f%%)\n", results[i].name, results[i].duration, percent);
    }
    
};

struct Profiler{
    const char* name;
    u64 startTime;
    
    Profiler(const char* name) : name(name), startTime(ReadCPUTimer()) {}
    ~Profiler(){
        u64 endTime = ReadCPUTimer();
        u64 duration = endTime - startTime;
        results[resultCount++] = {name, duration};
    }
};

#define TimerFunction \
    Profiler __functionProfiler(__func__);

#define BlockFunction(name) \
    Profiler __blockProfiler(name);
