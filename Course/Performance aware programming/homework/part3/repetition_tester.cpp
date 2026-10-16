
#include "platform_metrics.cpp"

enum tester_state: u32{
    TesterState_Uninitialized,
    TesterState_Running,
    TesterState_Error,
    TesterState_Finished
};

enum allocation_type: u32 {
    AllocationType_None,
    AllocationType_Malloc,
    AllocationType_Count
};

enum tester_result_type{
    TestResult_TestCount,
    TestResult_CPUTimer,
    TestResult_MemPageFault,
    TestResult_ByteCount,

    TestResult_Count
};

struct FileContent{
    
    char* data;
    size_t size;
    char* name;  
    allocation_type alloc;
};

struct tester_value{
    u64 tv[TestResult_Count];
};


struct tester_result{
    tester_value TotalTime;
    tester_value MinTime;
    tester_value MaxTime;
};

struct repetition_tester {

    tester_state state;
 
    u64 CPUFreq;
    u64 TryForTime;
    u64 TimeTestStartedAt;
    
    u64 TimeTrackerStart;
    u64 TimeTrackerEnd;
    
    u64 ByteCount;
    
    tester_value AccumulatorForThisTest;
    tester_result result;
    
};

static char const *GetAllocationType(allocation_type alloc){
    char const *Result;
    switch(alloc){
        case AllocationType_None: Result= ""; break;
        case AllocationType_Malloc: {Result="malloc"; break;}
        default: Result="UNKNOWN"; break;
    }
    return Result;
}

static char*  Allocate(FileContent *file_data){
    switch(file_data->alloc){
        case AllocationType_None: return file_data->data;
        case AllocationType_Malloc: return (char *) malloc(file_data->size);
        default: return nullptr;
    }
}

static void Deallocate(FileContent *file_data, char* buffer){
    switch(file_data->alloc){
        case AllocationType_None: break;
        case AllocationType_Malloc: {free(buffer); break;}
        default: break;
    }
}

static void PrintTime(char const* print_label, tester_value value, u64 cpu_freq, u64 bytes_read){
    u64 testCount=value.tv[TestResult_TestCount];
    f64 divisor = testCount!=0?(f64)testCount:1;

    f64 E[TestResult_Count];

    for(u32 index=0; index<ArrayCount(E);index++){
        E[index]=(f64)value.tv[index]/divisor;
    }
    printf("%s : %.0f", print_label, E[TestResult_CPUTimer]);
    if(cpu_freq){
        f64 time_seconds = (E[TestResult_CPUTimer] / (f64)cpu_freq);
        printf(" (%fms)", 1000.0f*time_seconds);

        if(bytes_read){
            f64 gb = (1024.0f * 1024.0f * 1024.0f);
            f64 bandidth = bytes_read / (gb * time_seconds);
            printf(" %f GB/s", bandidth);
        }
    }

    if(E[TestResult_MemPageFault]>0){
        printf("PF:%0.4f (%0.4fk/fault)",E[TestResult_MemPageFault],(f64)bytes_read/(E[TestResult_MemPageFault]*1024.0));    
    }
}

static void PrintResults(tester_result result, u64 cpu_freq, u64 bytes_read){
    PrintTime("Min", result.MinTime, cpu_freq, bytes_read);
    printf("\n");
    PrintTime("Max", result.MaxTime, cpu_freq, bytes_read);
    printf("\n");
    PrintTime("Avg", result.TotalTime, cpu_freq, bytes_read);
    printf("\n");    
    
}

static void InitializeParamForFunc(repetition_tester *tester, u64 ByteCount) {

    u64 CPUFreq = EstimateCPUTimerFreq();   
    if(tester->state == TesterState_Uninitialized) {
        tester->state = TesterState_Running;
        tester->CPUFreq = CPUFreq;
        tester->result.MinTime.tv[TestResult_CPUTimer] = (u64)-1;
        tester->ByteCount = ByteCount;
    }else if (tester->state == TesterState_Finished) {
        tester->state = TesterState_Running;
    }

    tester->TryForTime=10*CPUFreq;
    tester->TimeTestStartedAt = ReadCPUTimer();
}

static void BeginTime(repetition_tester *tester){
    tester->TimeTrackerStart++;

    tester->AccumulatorForThisTest.tv[TestResult_CPUTimer]-= ReadCPUTimer();
    tester->AccumulatorForThisTest.tv[TestResult_MemPageFault]-=ReadOSPageFault();
}

static void EndTime(repetition_tester *tester){
   
    tester->AccumulatorForThisTest.tv[TestResult_CPUTimer]+= ReadCPUTimer();
    tester->AccumulatorForThisTest.tv[TestResult_MemPageFault]+=ReadOSPageFault();
    tester->TimeTrackerEnd++;
}

static bool isTesting(repetition_tester *tester){

    if(tester->state==TesterState_Running){
        u64 currentTime=ReadCPUTimer();

        tester_value accum = tester->AccumulatorForThisTest;
        if(tester->TimeTrackerStart){

            if(tester->TimeTrackerEnd!=tester->TimeTrackerStart){
                tester->state=TesterState_Error;
                fprintf(stderr, "Error: %s\n", "unbalanced begin and end time");
            }
            if(tester->state==TesterState_Running){
                tester_result* results = &tester->result;
       
                accum.tv[TestResult_TestCount]=1;
                for(u32 index=0; index<ArrayCount(accum.tv);index++){
                    results->TotalTime.tv[index]+=accum.tv[index];
                }

                if(results->MaxTime.tv[TestResult_CPUTimer] < accum.tv[TestResult_CPUTimer]){
                    results->MaxTime = accum;
                }
                if(results->MinTime.tv[TestResult_CPUTimer] > accum.tv[TestResult_CPUTimer]){
                    results->MinTime = accum;

                    tester->TimeTestStartedAt= currentTime;
                    PrintTime("Min", results->MinTime, tester->CPUFreq, tester->ByteCount);
                    printf("\n");
                }

                tester->TimeTrackerStart=0;
                tester->TimeTrackerEnd=0;
                tester->AccumulatorForThisTest={};
            }
        }

        if((currentTime - tester->TimeTestStartedAt) > tester->TryForTime){
            tester->state = TesterState_Finished;
            printf("\n");
            PrintResults(tester->result, tester->CPUFreq, tester->ByteCount);
        }
    }

    return tester->state==TesterState_Running;
}