package main

import (
    "bytes"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    cadence := sync.NewCond(&sync.Mutex{})
    go func() {
        for range time.Tick(1 * time.Millisecond) {
            cadence.Broadcast()
        }
    }()

    takeStep := func() {
        cadence.L.Lock()
        cadence.Wait()
        cadence.L.Unlock()
    }

    // try to move in a certain direction with the count variable denote which direction
    tryDirection := func(dirName string, dir *int32, out *bytes.Buffer) bool {
        fmt.Fprintf(out, " %v", dirName)
        atomic.AddInt32(dir, 1)
        takeStep()
        if atomic.LoadInt32(dir) == 1 {
            fmt.Fprintf(out, ". Success!")
            return true
        }

        takeStep()
        atomic.AddInt32(dir, -1)
        return false
    }

    var left, right int32
    tryleft := func(out *bytes.Buffer) bool { return tryDirection("left", &left, out) }
    tryRight := func(out *bytes.Buffer) bool { return tryDirection("right", &right, out) }

    walk := func(walking *sync.WaitGroup, name string) {
        var out bytes.Buffer
        defer func() { fmt.Println(out.String()) }()
        defer walking.Done()

        fmt.Fprintf(&out, "%v is trying to move:", name)
        for i := 0; i < 5; i++ {
            if tryleft(&out) || tryRight(&out) {
                return
            }
        }

        fmt.Fprintf(&out, "\n%v is just lost in its own world.", name)
    }

    var peopleInHall sync.WaitGroup
    peopleInHall.Add(2)

    time.Sleep(10 * time.Millisecond) // Initial delay to let the broadcasting start

    go walk(&peopleInHall, "Anish")
    go walk(&peopleInHall, "Rijul")

    peopleInHall.Wait()
}

// Rijul is trying to move: left right left right left right left right left right
// Rijul is just lost in its own world.
// Anish is trying to move: left right left right left right left right left right
// Anish is just lost in its own world.

// result no one is able to move in this thing altho the threads are always busy unlike deadlock