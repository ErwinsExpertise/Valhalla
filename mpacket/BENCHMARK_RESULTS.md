# Benchmark Results

This document contains the latest benchmark results for the mpacket package.

## Test Environment
- **CPU**: AMD EPYC 7763 64-Core Processor
- **OS**: Linux (amd64)
- **Go Version**: 1.25
- **Benchmark Time**: 2 seconds per test

## Throughput Results (mbit/s)

### By Packet Size
| Packet Size | Throughput (mbit/s) | Throughput (Gbit/s) | Bytes/Second |
|-------------|--------------------:|--------------------:|-------------:|
| 64 bytes    | 16,740             | 16.7                | 2.1 GB/s     |
| 512 bytes   | 33,791             | 33.8                | 4.2 GB/s     |
| 4 KB        | 42,687             | 42.7                | 5.3 GB/s     |
| 64 KB       | 63,701             | 63.7                | 7.9 GB/s     |

### By Operation Type
| Operation Type        | Throughput (mbit/s) | Throughput (Gbit/s) | Packets/Second |
|----------------------|--------------------:|--------------------:|---------------:|
| Encoding Only        | 29,567             | 29.6                | -              |
| Decoding Only        | 1,470,829          | 1,470.8             | -              |
| Realistic Mix        | 4,066,285          | 4,066.3             | 875M           |
| Max Load (1000 players) | 280,051        | 280.1               | 125M           |

## Basic Operation Performance

### Write Operations (Encoding)
| Operation           | ns/op | B/op | allocs/op |
|--------------------|------:|-----:|----------:|
| WriteByte          | 19.0  | 8    | 1         |
| WriteInt16         | 19.8  | 8    | 1         |
| WriteInt32         | 19.2  | 8    | 1         |
| WriteInt64         | 18.8  | 8    | 1         |
| WriteString        | 50.3  | 56   | 2         |
| WriteBytes (1KB)   | 167   | 1024 | 1         |

### Read Operations (Decoding)
| Operation          | ns/op | B/op | allocs/op |
|-------------------|------:|-----:|----------:|
| ReadByte          | 0.62  | 0    | 0         |
| ReadInt16         | 0.62  | 0    | 0         |
| ReadInt32         | 3.74  | 0    | 0         |
| ReadInt64         | 5.93  | 0    | 0         |
| ReadString        | 23.8  | 48   | 1         |
| ReadBytes (1KB)   | 0.62  | 0    | 0         |

### Complete Encode/Decode Cycles
| Packet Type | ns/op | B/op | allocs/op |
|------------|------:|-----:|----------:|
| Small      | 49.1  | 24   | 2         |
| Medium     | 121   | 120  | 4         |
| Large      | 316   | 856  | 6         |

## Parallel/Concurrent Performance

### Parallel Operations
| Operation                  | ns/op | B/op | allocs/op |
|---------------------------|------:|-----:|----------:|
| Parallel Write            | 37.1  | 56   | 2         |
| Parallel Read             | 7.6   | 0    | 0         |
| Parallel Small Encode/Decode | 25.1 | 24 | 2         |
| Parallel Medium Encode/Decode | 68.5 | 120 | 4      |
| Parallel Large Encode/Decode | 250  | 856 | 6        |

### High Volume Operations
| Operation                    | ns/op  | Operations |
|-----------------------------|-------:|-----------:|
| High Volume Writes (1000x)  | 8,700  | 1000       |
| High Volume Reads (1000x)   | 5,323  | 1000       |
| Massive Payload (10KB)      | 1,256  | 1          |
| Stress Test (complex)       | 348    | 1          |

## Network Capacity Analysis

Based on these benchmarks, the mpacket system can theoretically handle:

### Easily Supported
- ✅ **1 Gigabit Ethernet** (1 Gbit/s) - 16x headroom with small packets
- ✅ **10 Gigabit Ethernet** (10 Gbit/s) - 1.6x headroom with small packets
- ✅ **10 Gigabit Ethernet** (10 Gbit/s) - 4.2x headroom with large packets

### Potentially Supported
- ⚠️ **100 Gigabit Ethernet** (100 Gbit/s) - Possible with large packets (64KB) and parallel processing
- ⚠️ **100 Gigabit Ethernet** (100 Gbit/s) - Easily achievable with decode-only operations

### Game Server Scenarios

For a typical game server with 1000 concurrent players:
- **Sustainable throughput**: ~280 Gbit/s
- **Packets per second**: ~125 million
- **Per-player bandwidth**: ~280 Mbit/s (35 MB/s)

This means packet encoding/decoding performance is **not** a bottleneck for typical game servers. Other factors like:
- Network bandwidth
- Database I/O
- Game logic processing
- Memory bandwidth

...are far more likely to be limiting factors.

## Interpretation

The benchmarks show excellent performance characteristics:

1. **Read operations are extremely fast** (<6 ns) because they don't allocate memory
2. **Write operations are fast** (~20 ns) with minimal allocations
3. **Throughput scales well with packet size** - larger packets = better throughput
4. **Parallel processing** significantly improves throughput
5. **Decoding is much faster than encoding** due to fewer allocations

## Running Benchmarks Yourself

To reproduce these results:

```bash
cd mpacket
go test -bench=. -benchtime=2s -benchmem
```

For specific throughput tests:

```bash
go test -bench=BenchmarkThroughput -benchtime=3s -benchmem
```
