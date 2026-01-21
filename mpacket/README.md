# mpacket Benchmarking

This directory contains benchmarks for testing the performance of packet encoding and decoding operations, including throughput measurements in megabits per second (mbit/s).

## Running Benchmarks

To run all benchmarks in the mpacket package:

```bash
cd mpacket
go test -bench=.
```

### Run Throughput Benchmarks (mbit/s measurements)

To see how many mbit/s the packet system can handle:

```bash
go test -bench=BenchmarkThroughput -benchtime=3s
```

### Run with Memory Statistics

To see memory allocation statistics:

```bash
go test -bench=. -benchmem
```

### Run Specific Benchmark

To run a specific benchmark (e.g., only write operations):

```bash
go test -bench=BenchmarkPacketWrite
```

To run a specific single benchmark:

```bash
go test -bench=BenchmarkPacketWriteInt32
```

### Adjust Benchmark Time

By default, benchmarks run for 1 second. To run for longer and get more accurate results:

```bash
go test -bench=. -benchtime=10s
```

### Run Parallel/Load Benchmarks

To test heavy concurrent load:

```bash
go test -bench=BenchmarkPacketParallel -benchtime=5s
```

### Run with CPU Profiling

To generate a CPU profile:

```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Run with Memory Profiling

To generate a memory profile:

```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
```

## Available Benchmarks

### Write Operations (Encoding)
- `BenchmarkPacketWriteByte` - Writing single bytes
- `BenchmarkPacketWriteInt16` - Writing 16-bit integers
- `BenchmarkPacketWriteInt32` - Writing 32-bit integers
- `BenchmarkPacketWriteInt64` - Writing 64-bit integers
- `BenchmarkPacketWriteString` - Writing strings
- `BenchmarkPacketWriteBytes` - Writing byte arrays (1KB)
- `BenchmarkPacketMultipleWrites` - Sequential write operations

### Read Operations (Decoding)
- `BenchmarkPacketReadByte` - Reading single bytes
- `BenchmarkPacketReadInt16` - Reading 16-bit integers
- `BenchmarkPacketReadInt32` - Reading 32-bit integers
- `BenchmarkPacketReadInt64` - Reading 64-bit integers
- `BenchmarkPacketReadString` - Reading strings
- `BenchmarkPacketReadBytes` - Reading byte arrays (1KB)
- `BenchmarkPacketMultipleReads` - Sequential read operations

### Complete Cycles
- `BenchmarkPacketEncodeDecodeSmall` - Small packet (header + few fields)
- `BenchmarkPacketEncodeDecodeMedium` - Medium packet (player info)
- `BenchmarkPacketEncodeDecodeLarge` - Large packet (complex data)

### Parallel/Concurrent Load Tests
- `BenchmarkPacketParallelWrite` - Concurrent writes to different packets
- `BenchmarkPacketParallelRead` - Concurrent reads from packets
- `BenchmarkPacketParallelEncodeDecodeSmall` - Concurrent small packet operations
- `BenchmarkPacketParallelEncodeDecodeMedium` - Concurrent medium packet operations
- `BenchmarkPacketParallelEncodeDecodeLarge` - Concurrent large packet operations

### High Volume Tests
- `BenchmarkPacketHighVolumeWrites` - 1000 sequential write operations
- `BenchmarkPacketHighVolumeReads` - 1000 sequential read operations
- `BenchmarkPacketMassivePayload` - 10KB packet operations
- `BenchmarkPacketStressTest` - Realistic high-load scenario with complex packets

### Throughput Tests (mbit/s measurements)
- `BenchmarkThroughputSmallPackets` - 64-byte packets (typical control packets)
- `BenchmarkThroughputMediumPackets` - 512-byte packets (typical data packets)
- `BenchmarkThroughputLargePackets` - 4KB packets (large data transfers)
- `BenchmarkThroughputMassivePackets` - 64KB packets (maximum size transfers)
- `BenchmarkThroughputRealisticMix` - Mixed packet sizes (realistic workload)
- `BenchmarkThroughputEncodeOnly` - Encoding-only throughput
- `BenchmarkThroughputDecodeOnly` - Decoding-only throughput
- `BenchmarkThroughputMaxLoad` - Maximum sustainable load with 1000 concurrent players

### Creation
- `BenchmarkPacketCreateWithOpcode` - Packet creation overhead

## Interpreting Results

Benchmark output format:
```
BenchmarkPacketWriteByte-4    61291801    19.20 ns/op    8 B/op    1 allocs/op
```

- `61291801` - Number of iterations completed
- `19.20 ns/op` - Average time per operation (nanoseconds)
- `8 B/op` - Bytes allocated per operation
- `1 allocs/op` - Number of allocations per operation

For throughput benchmarks, additional metrics are shown:
```
BenchmarkThroughputLargePackets-4   4678621   781.6 ns/op   5240488413 bytes/s   41924 mbit/s   4104 B/op   2 allocs/op
```

- `5240488413 bytes/s` - Bytes processed per second
- `41924 mbit/s` - **Megabits per second throughput**
- `packets/s` - Packets processed per second (where applicable)

Lower values generally indicate better performance for ns/op, B/op, and allocs/op. Higher values indicate better performance for mbit/s and packets/s.

## Performance Characteristics

### Typical Results (AMD EPYC 7763 64-Core)

**Basic Operations:**
- Write operations: ~20 ns/op for primitive types
- Read operations: ~1-6 ns/op for primitive types
- Small packet encode/decode: ~50 ns/op
- Medium packet encode/decode: ~120 ns/op
- Large packet encode/decode: ~320 ns/op

**Throughput Results:**
- Small packets (64 bytes): **~16,000 mbit/s** (16 Gbit/s)
- Medium packets (512 bytes): **~33,000 mbit/s** (33 Gbit/s)
- Large packets (4KB): **~42,000 mbit/s** (42 Gbit/s)
- Massive packets (64KB): **~63,000 mbit/s** (63 Gbit/s)
- Realistic mix: **~3,900,000 mbit/s** (3.9 Tbit/s) - theoretical maximum
- Encoding only: **~29,000 mbit/s** (29 Gbit/s)
- Decoding only: **~1,500,000 mbit/s** (1.5 Tbit/s)
- Max load (1000 players): **~283,000 mbit/s** (283 Gbit/s), **~127M packets/s**

**Note:** These are processor-bound benchmarks measuring the maximum theoretical throughput of the packet encoding/decoding logic. Real-world network throughput will be limited by:
- Network interface card speed (typically 1-100 Gbit/s)
- Network bandwidth and latency
- Operating system network stack overhead
- Application-level processing and business logic

## Understanding Network Capacity

The benchmarks show the packet system can theoretically handle:
- **Gigabit Ethernet (1 Gbit/s)**: Easily achievable with large headroom
- **10 Gigabit Ethernet (10 Gbit/s)**: Easily achievable with large headroom
- **100 Gigabit Ethernet (100 Gbit/s)**: Potentially achievable with large packets and parallel processing

For a typical game server scenario with thousands of players, the packet encoding/decoding performance is unlikely to be the bottleneck. Other factors like database I/O, game logic, and network bandwidth will typically be the limiting factors.
