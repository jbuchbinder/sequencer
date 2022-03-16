# SEQUENCER

Unique ID generator, capable of distributed operation. Based off of [Twitter Snowflake](https://github.com/twitter-archive/snowflake/tree/snowflake-2010) and one of its [succesors](https://www.callicoder.com/distributed-unique-id-sequence-number-generator/).

Inherited advantages:

* No DB syncing required
* Supports massive concurrency and parallel installs

Additional advantages:

* Simple REST interface
* Written in Go, so much leaner than the Java alternatives
 
