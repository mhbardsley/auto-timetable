# auto-timetable
Stochastic timetabling algorithm

## Supported types of slot
auto-timetable is capable of dealing with one-off events, deadlines, and periodics.

One-off events have a prescribed start and end time, and they only exist between those times.

Deadlines have a prescribed end time, and an estimated number of minutes to achieve the deadline. The idea is that they will be scheduled as evenly as possible.

Periodics are events that can happen whenever, but they continue indefinitely.