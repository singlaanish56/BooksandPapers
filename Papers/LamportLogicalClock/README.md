
This is the basic implementation using the the logical clock Implementation Rules defined by Dr. Leslie Lamport
in the Paper

[Time, Clocks, and the Ordering of Events in a Distributed System](https://lamport.azurewebsites.net/pubs/time-clocks.pdf)

this is partial ordering implementation for the logical clocks, which for two given events helps us define the a->b


The relation "->" on a set of events of a system satisfies the conditions
1. if a and b are events of the same process, and a comes before b then a -> b
2. if a is the sending of a message from a process , and b is recieving of the same message on another process then a-> b
3. if a -> b and b -> c then a -> c
4. two events are concurrent  a /-> b and b/-> a


Because logical clocks are used, The correctness cant be measured by the correctness defined for the physical clocks
these clocks are defined correct based on the ordering of the events

The paper defines set of <b>Clock Conditions</b> to be true for this effect. remember C\<a\> is the number clock for the process assigns the event

<b>C1</b>, if a and b are events in a process Pi and a comes before b, then <b>Ci\<a\> < Ci\<b\> </b>
<b>C2</b>, if a is sending of the message from Pi and b is recieving that message on Pj then <b> Ci\<a\> < Cj\<b\> </b>



To satisfy the above condition come the <b>implementation rules</b>, which ensure them

<b>IR1</b>, each process Pi incremenents Ci between any two successive events
<b>IR2, (a)</b> if a event is sending of a message by process Pi, then message m constaind the timestamp Tm = Ci\<a\>
<b>(b)</b> upon recieivng the message m the process Pj sets the time in its clock to be the max of (currentClock, Tm)




Note the same implementation rules can be used for the <b>total ordering of the events</b>, in this case of the tie.
we break even by the process Ids so if <b> Ci\<a\> == Cj\<b\> then Pi < Pj. </b>

The total ordering id defined through the =>.

so the updated condition for total ordering and the clcok condition ensure implies that if <b> a ->b then a=> b </b>


