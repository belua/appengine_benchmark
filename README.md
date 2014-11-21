# Benchmarking Scratchpad for App Engine

Because we are about to make a number of important architectural decisions based on expected performance of the app engine datastore we should establish some real baselines for those characteristics.

Performance test 1. Low Kind populations

When there are few entities of a particular kind both the entity writes and index writes will end up on the same big-table tablet (same physical machine). Because most benchmarks will start with an empty datastore if we do not control this factor our tests won't be representative. So we need to establish a minumum number of entities of a given kind that are needed before this effect is lost.

Performance test 2. Independent Entity Writes

At what rate can we expect to write entities which are completely independent. Does the number of indexed properties affect this? By how much?
