package drafts

// Consumer assignment refers to these basic tasks:
// - Which broker serves the data
// - Which consumer gets the data
// - How we communicate the assignment or offsets: currently we send a picture of the consumers and the offsets

// |  C0 |  C1 |  C2 |  C3 |  C4 |  C5 |     |     |     |     |     |     |     |     |     |     |
// +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
// |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
// T0 0  |  1  |  2  |  3  T6 0  |  1  |  2  |  3  T3 0  |  1  |  2  |  3  T7 0  |  1  |  2  |  3  T1         Gen v3
// |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
// +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
// |           |           |           |           |           |           |           |           |
// T0    0     |     1     |     2     |     3     T3    0     |     1     |     2     |     3     |          Gen v2
// |           |           |           |           |           |           |           |           |
// +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
// |                       |                       |                       |                       |
// T0          0           |           1           |           2           |           3           |          Gen v1
// |                       |                       |                       |                       |
// +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
//
// These questions are especially relevant when we consider past generations.
//
// Example 1:
// - Cluster of 3 (v1), scaled to 6 (v2) and to 12 brokers (v3).
// - Consumer group A haven't yet read T0/1 v1
// - C5 is assigned to T6/1 and C1 is assigned to T0/1
// - Consumers are assigned to the most recent snapshot of the topology
// - [Q] Who should serve T0/1 v1?
// ----> B6 is leader of range T6 -> T3, B6 should serve T0/1 v1 because it starts at the same position.
// ----> When the current generation has a larger or equal cluster size than the previous one, use the broker assigned
//       to the current range starting at the same starting point (C4).
// ----> When the current generation has a smaller cluster size than the previous one, use the broker assigned
//       to the range containing that portion.
// - [Q] Who is the consumer that should be assigned to T0/1 v1?
// ----> When the current generation has a larger or equal cluster size than the previous one, use the consumer assigned
//       to the current range starting at the same starting point (C4).
// ----> When the current generation has a smaller cluster size than the previous one, use the consumer assigned
//       to the range containing that portion.
// - [Q] Should we store the offset on numerical ranges?
// ----> yes
// - [Q] How we should move the offset when completed
// ----> When the reader identifies that T0/1 v1 is completed, it marks it as completed. The offset state then performs
//       a search for the next generations (T0/2 v1 and T0/3) and marks it as zero.
// - [Q] How should we fix the offset "source"
// ----> Source composed of also a timestamp
// - [Q] Should we replicate the offset on numerical ranges
// - [Q] How Get() and Set() operations and the return value of `Get()` would look like. What's the value struct
// ----> As defined in the draft
// - [Q] How to get the ranges of each generation
// ----> Cluster size gets stored in each generation

// Possible logic:
// - When C5 wants to consume T6/1, B6 (leader of T6/1) navigates through offset _on token numerical ranges_.
// - It does a search and finds that for group A, within that range, it is still in T0/1 v1

// Possible structures:
// - OffsetByPartition
// -- By group and by topic
// -- A ordered range of partitions with int64 values
//
