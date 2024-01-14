# auto-timetable
Stochastic timetabling algorithm

## Adding support for deadlines and events
GetInput is the function that takes the "file" and converts it to an inputData type, which is exactly what we want

So we need to adapt GetInput so that it will instead take a top-level directory location, and recurse on that to extract any events or deadlines
It should still do the sorting, checking, and so on of data. This is simply manipulation of where it loads data from and how it finds the data

So, we should have the following:
- Migrating the behaviour that selects a periodic event to auto-timetable (50)
- Remaining reporting changes (25)