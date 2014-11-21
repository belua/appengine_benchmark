# Benchmarking Scratchpad for App Engine

Because we are about to make a number of important architectural decisions based on expected performance of the app engine datastore we should establish some real baselines for those characteristics.

Performance test 1. Low Kind populations

When there are few entities of a particular kind both the entity writes and index writes will end up on the same big-table tablet (same physical machine). Because most benchmarks will start with an empty datastore if we do not control this factor our tests won't be representative. So we need to establish a minumum number of entities of a given kind that are needed before this effect is lost.

The result is that this effect is very modest. We can see an initially slow set of puts for the first 100 entities. As long as a performance test creates more than 1000 this effect shouldn't be too pronounced. Each put ranges from just under 200 ms for the first 100 puts settling down to 30-40ms with occasional outliers around 100-200ms.

Performance test 2. Independent Entity Writes

At what rate can we expect to write entities which are completely independent. Does the number of indexed properties affect this? By how much?
