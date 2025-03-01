This is a simple implementation for th vector clocks in the paper

[Timestamps is message passing systems](https://fileadmin.cs.lth.se/cs/Personal/Amr_Ergawy/dist-algos-papers/4.pdf)

One of the important aspects of concurrent programming in a distributes system is determining the order in which the events occured in a specific process.

Now the lamport ties to solve this problem by determing a specific total/partial-order within the events occured across the processes.
But it only yields one of the many possible , valid orderings that exist for the distributed system.

This paper allows to have all possible global states to compute the order version other processes as well , without adhering to specific ordering as suggested by the Lamport's Paper.

This is solved by preserving the start or timestamp for each process to know when it last communicated with every other process in the system.

These are the rules to be followed while attempting to establish such system without adding on to communnication chain.

## Steps for Asynchronous Communication

1. Each process maintains an array for each process in the system, or start with 0 of it encounters a message from a new process
2. The local clock value was incremented at least once before sending an event.
3. The current value of the timestamps is piggybacked onto the comunication to the other procecss
4. Upon recieving the signal, the recieving process checks two things so there are two arrays to think about localArray and the recievedArray
5. where q is the id of the sender,

```
if localArray[q] <= recievedArray[q]
    localArray[q]=1+receivedArray[q]

for i:=1 to n do
    localArray[i]=max(localArray[i], recievedArray[i])
```
6. the sender's local timestamp is +1 to take array for the delay in communication

## Steps for Synchronous Commnucation

even lamports algo can be adapted to ensure the processes updated the local timstamps to their max valueas and send back the timestamp
i.e exchanging the timestamp for each event

1. Each process maintains an array for each process in the system, or start with 0 of it encounters a message from a new process
2. The local clock value was incremented at least once before sending an event.
3. Upon recieving the signal, the recieving process checks two things so there are two arrays to think about localArray and the recievedArray
4. processes exchange the timestamps for each event, and the localArray to set the max value of localValue vs the recieved value

```
for i:=1 to n do
    localArray[i]=max(localArray[i], recievedArray[i])
```

Hence, we can derive 

```ep -> fq iff Tep[p] <= Tfq[p] ^ Tep[q] < Tfq[q]
```
what does this mean?

the first half of the conjunction ensures that the timestamp of clock recieved by process q from process p is as recent as execution of event ep, hence ensuring that fq happens after ep

the seconf half of the conjuction ensures that process p does not have up to date information of q. hence ep!=fq


## Sample Run from the Code

```
[Send] P0 -> P1| Time: map[0:3] | Content: This is a message from process 0

[Recieve] P1 from P0 | Updated Timestamp: map[0:4 1:0] | Content: This is a message from process 0

[Send] P1 -> P2| Time: map[0:4 1:1] | Content: This is a message from process 1

[Recieve] P2 from P1 | Updated Timestamp: map[0:4 1:2 2:0] | Content: This is a message from process 1

[Send] P0 -> P1| Time: map[0:6] | Content: This is a message from process 0

[Recieve] P1 from P0 | Updated Timestamp: map[0:7 1:1] | Content: This is a message from process 0

[Send] P1 -> P0| Time: map[0:7 1:2] | Content: This is a message from process 1

[Recieve] P0 from P1 | Updated Timestamp: map[0:7 1:3] | Content: This is a message from process 1

[Send] P1 -> P2| Time: map[0:7 1:5] | Content: This is a message from process 1

[Recieve] P2 from P1 | Updated Timestamp: map[0:7 1:6 2:0] | Content: This is a message from process 1

the vector Clocks for various process
P0 | T0 map[0:7 1:3]
P1 | T1 map[0:7 1:5]
P2 | T2 map[0:7 1:6 2:0]

```