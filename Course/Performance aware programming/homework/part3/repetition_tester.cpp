
#include "platform_metrics.cpp"

enum tester_state: u32{
    TesterState_Uninitialized,
    TesterState_Running,
    TesterState_Error,
    TesterState_Finished
};

struct tester_result{
    u64 Repetitions;
    u64 TotalTime;
    u64 MinTime;
    u64 MaxTime;
};

struct repetition_tester {

    tester_state state;
    tester_result result;
    u64 CPUFreq;
    u64 TryForTime;
    u64 TimeTestStartedAt;
    u64 TimeForThisTest;
    u64 TimeTrackerStart;
    u64 TimeTrackerEnd;
    u64 ByteCount;
};

static void PrintTime(char const* print_label, f64 cpu_time, u64 cpu_freq, u64 bytes_read){

    printf("%s : %.0f", print_label, cpu_time);
    if(cpu_freq){
        f64 time_seconds = (cpu_time / (f64)cpu_freq);
        printf(" (%fms)", 1000.0f*time_seconds);

        if(bytes_read){
            f64 gb = (1024.0f * 1024.0f * 1024.0f);
            f64 bandidth = bytes_read / (gb * time_seconds);
            printf(" %fgb/s", bandidth);
        }
    }
}
static void PrintTime(char const* print_label, u64 cpu_time, u64 cpu_freq, u64 bytes_read){
     PrintTime(print_label, (f64)cpu_time, cpu_freq, bytes_read);
}

static void PrintResults(tester_result result, u64 cpu_freq, u64 bytes_read){
    PrintTime("Min", result.MinTime, cpu_freq, bytes_read);
    printf("\n");
    PrintTime("Max", result.MaxTime, cpu_freq, bytes_read);
    printf("\n");
    if(result.TotalTime){
        PrintTime("Avg", (f64)result.TotalTime/(f64)result.Repetitions, cpu_freq, bytes_read);
        printf("\n");
    }
    
}

static void InitializeParamForFunc(repetition_tester *tester, u64 ByteCount) {

    u64 CPUFreq = EstimateCPUTimerFreq();   
    if(tester->state == TesterState_Uninitialized) {
        tester->state = TesterState_Running;
        tester->CPUFreq = CPUFreq;
        tester->result.MinTime = (u64)-1;
        tester->ByteCount = ByteCount;
    }else if (tester->state == TesterState_Finished) {
        tester->state = TesterState_Running;
    }

    tester->TryForTime=10*CPUFreq;
    tester->TimeTestStartedAt = ReadCPUTimer();
}

static void BeginTime(repetition_tester *tester){
    tester->TimeTrackerStart++;
    tester->TimeForThisTest -= ReadCPUTimer();
}

static void EndTime(repetition_tester *tester){
    tester->TimeTrackerEnd++;
    tester->TimeForThisTest += ReadCPUTimer();
}

static bool isTesting(repetition_tester *tester){

    if(tester->state==TesterState_Running){
        u64 currentTime=ReadCPUTimer();

        if(tester->TimeTrackerStart){

            if(tester->TimeTrackerEnd!=tester->TimeTrackerStart){
                tester->state=TesterState_Error;
                fprintf(stderr, "Error: %s\n", "unbalanced begin and end time");
            }
            if(tester->state==TesterState_Running){
                tester_result* results = &tester->result;
                u64 elapsedTime = tester->TimeForThisTest;
                
                results->Repetitions++;
                results->TotalTime+=elapsedTime;

                if(results->MaxTime < elapsedTime){
                    results->MaxTime = elapsedTime;
                }
                if(results->MinTime > elapsedTime){
                    results->MinTime = elapsedTime;

                    tester->TimeTestStartedAt= currentTime;
                    PrintTime("Min", results->MinTime, tester->CPUFreq, tester->ByteCount);
                    printf("\n");
                }

                tester->TimeTrackerStart=0;
                tester->TimeTrackerEnd=0;
                tester->TimeForThisTest=0;
            }
        }

        if((currentTime - tester->TimeTestStartedAt) > tester->TryForTime){
            tester->state = TesterState_Finished;
            
            PrintResults(tester->result, tester->CPUFreq, tester->ByteCount);
        }
    }

    return tester->state==TesterState_Running;
}