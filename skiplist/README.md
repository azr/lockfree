Implementation of [A Provably Correct Scalable Concurrent Skip
List](https://www.cs.tau.ac.il/~shanir/nir-pubs-web/Papers/OPODIS2006-BA.pdf)

* Searches are lock free.
* Inserts/Deletes will lock locally.

Abstract:
> [...] skip list algorithm distinguished
by a combination of simplicity and scalability. The algorithm employs
optimistic synchronization, searching without acquiring locks, followed
by short lock-based validation before adding or removing nodes. It also
logically removes an item before physically unlinking it. Unlike some
other concurrent skip list algorithms, this algorithm preserves the skip
list properties at all times, which facilitates reasoning about its correctness.
Experimental evidence shows that this algorithm performs as well
as the best previously known algorithm under most circumstances


The race detector will trigger for bool read/writes.
