package mpacket

import (
	"sync/atomic"
	"testing"
)

// BenchmarkPacketWriteByte benchmarks writing a single byte
func BenchmarkPacketWriteByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteByte(0xFF)
	}
}

// BenchmarkPacketWriteInt16 benchmarks writing int16 values
func BenchmarkPacketWriteInt16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteInt16(12345)
	}
}

// BenchmarkPacketWriteInt32 benchmarks writing int32 values
func BenchmarkPacketWriteInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteInt32(123456789)
	}
}

// BenchmarkPacketWriteInt64 benchmarks writing int64 values
func BenchmarkPacketWriteInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteInt64(1234567890123456789)
	}
}

// BenchmarkPacketWriteString benchmarks writing strings
func BenchmarkPacketWriteString(b *testing.B) {
	testString := "This is a test string for benchmarking"
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteString(testString)
	}
}

// BenchmarkPacketWriteBytes benchmarks writing byte arrays
func BenchmarkPacketWriteBytes(b *testing.B) {
	testBytes := make([]byte, 1024)
	for i := range testBytes {
		testBytes[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		p.WriteBytes(testBytes)
	}
}

// BenchmarkPacketReadByte benchmarks reading a single byte
func BenchmarkPacketReadByte(b *testing.B) {
	p := NewPacket()
	p.WriteByte(0xFF)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadByte()
	}
}

// BenchmarkPacketReadInt16 benchmarks reading int16 values
func BenchmarkPacketReadInt16(b *testing.B) {
	p := NewPacket()
	p.WriteInt16(12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadInt16()
	}
}

// BenchmarkPacketReadInt32 benchmarks reading int32 values
func BenchmarkPacketReadInt32(b *testing.B) {
	p := NewPacket()
	p.WriteInt32(123456789)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadInt32()
	}
}

// BenchmarkPacketReadInt64 benchmarks reading int64 values
func BenchmarkPacketReadInt64(b *testing.B) {
	p := NewPacket()
	p.WriteInt64(1234567890123456789)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadInt64()
	}
}

// BenchmarkPacketReadString benchmarks reading strings
func BenchmarkPacketReadString(b *testing.B) {
	testString := "This is a test string for benchmarking"
	p := NewPacket()
	p.WriteString(testString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadString(int16(len(testString)))
	}
}

// BenchmarkPacketReadBytes benchmarks reading byte arrays
func BenchmarkPacketReadBytes(b *testing.B) {
	testBytes := make([]byte, 1024)
	for i := range testBytes {
		testBytes[i] = byte(i % 256)
	}
	p := NewPacket()
	p.WriteBytes(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		r.ReadBytes(1024)
	}
}

// BenchmarkPacketEncodeDecodeSmall benchmarks a complete encode/decode cycle for a small packet
func BenchmarkPacketEncodeDecodeSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Encode
		p := CreateWithOpcode(0x01)
		p.WriteByte(0x10)
		p.WriteInt16(100)
		p.WriteInt32(1000)

		// Decode
		r := NewReader(&p, 0)
		r.Skip(4) // Skip header
		r.ReadByte()
		r.ReadByte()
		r.ReadInt16()
		r.ReadInt32()
	}
}

// BenchmarkPacketEncodeDecodeMedium benchmarks a complete encode/decode cycle for a medium packet
func BenchmarkPacketEncodeDecodeMedium(b *testing.B) {
	testString := "PlayerName123"
	for i := 0; i < b.N; i++ {
		// Encode
		p := CreateWithOpcode(0x02)
		p.WriteInt32(12345)           // Player ID
		p.WriteString(testString)     // Player name
		p.WriteInt16(50)              // Level
		p.WriteInt32(1000000)         // Experience
		p.WriteInt16(100)             // HP
		p.WriteInt16(100)             // MP
		p.WriteInt16(10)              // Strength
		p.WriteInt16(10)              // Dexterity

		// Decode
		r := NewReader(&p, 0)
		r.Skip(4) // Skip header
		r.ReadByte()
		r.ReadInt32()
		r.ReadString(int16(len(testString)))
		r.ReadInt16()
		r.ReadInt32()
		r.ReadInt16()
		r.ReadInt16()
		r.ReadInt16()
		r.ReadInt16()
	}
}

// BenchmarkPacketEncodeDecodeLarge benchmarks a complete encode/decode cycle for a large packet
func BenchmarkPacketEncodeDecodeLarge(b *testing.B) {
	testString := "VeryLongPlayerNameWithManyCharacters"
	testData := make([]byte, 512)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	for i := 0; i < b.N; i++ {
		// Encode
		p := CreateWithOpcode(0x03)
		p.WriteInt32(12345)
		p.WriteString(testString)
		for j := 0; j < 10; j++ {
			p.WriteInt32(int32(j * 1000))
		}
		p.WriteBytes(testData)
		p.WriteInt64(1234567890123456789)

		// Decode
		r := NewReader(&p, 0)
		r.Skip(4) // Skip header
		r.ReadByte()
		r.ReadInt32()
		r.ReadString(int16(len(testString)))
		for j := 0; j < 10; j++ {
			r.ReadInt32()
		}
		r.ReadBytes(512)
		r.ReadInt64()
	}
}

// BenchmarkPacketCreateWithOpcode benchmarks packet creation with opcode
func BenchmarkPacketCreateWithOpcode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateWithOpcode(0x01)
	}
}

// BenchmarkPacketMultipleWrites benchmarks multiple sequential writes
func BenchmarkPacketMultipleWrites(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		for j := 0; j < 100; j++ {
			p.WriteByte(byte(j))
		}
	}
}

// BenchmarkPacketMultipleReads benchmarks multiple sequential reads
func BenchmarkPacketMultipleReads(b *testing.B) {
	p := NewPacket()
	for j := 0; j < 100; j++ {
		p.WriteByte(byte(j))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		for j := 0; j < 100; j++ {
			r.ReadByte()
		}
	}
}

// BenchmarkPacketParallelWrite benchmarks concurrent writes to different packets
func BenchmarkPacketParallelWrite(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := NewPacket()
			p.WriteByte(0xFF)
			p.WriteInt16(12345)
			p.WriteInt32(123456789)
			p.WriteInt64(1234567890123456789)
			p.WriteString("Test String")
		}
	})
}

// BenchmarkPacketParallelRead benchmarks concurrent reads from shared packets
func BenchmarkPacketParallelRead(b *testing.B) {
	p := NewPacket()
	p.WriteByte(0xFF)
	p.WriteInt16(12345)
	p.WriteInt32(123456789)
	p.WriteInt64(1234567890123456789)
	p.WriteString("Test String")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := NewReader(&p, 0)
			r.ReadByte()
			r.ReadInt16()
			r.ReadInt32()
			r.ReadInt64()
			r.ReadString(11)
		}
	})
}

// BenchmarkPacketParallelEncodeDecodeSmall benchmarks concurrent small packet operations
func BenchmarkPacketParallelEncodeDecodeSmall(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := CreateWithOpcode(0x01)
			p.WriteByte(0x10)
			p.WriteInt16(100)
			p.WriteInt32(1000)

			r := NewReader(&p, 0)
			r.Skip(4)
			r.ReadByte()
			r.ReadByte()
			r.ReadInt16()
			r.ReadInt32()
		}
	})
}

// BenchmarkPacketParallelEncodeDecodeMedium benchmarks concurrent medium packet operations
func BenchmarkPacketParallelEncodeDecodeMedium(b *testing.B) {
	testString := "PlayerName123"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := CreateWithOpcode(0x02)
			p.WriteInt32(12345)
			p.WriteString(testString)
			p.WriteInt16(50)
			p.WriteInt32(1000000)
			p.WriteInt16(100)
			p.WriteInt16(100)
			p.WriteInt16(10)
			p.WriteInt16(10)

			r := NewReader(&p, 0)
			r.Skip(4)
			r.ReadByte()
			r.ReadInt32()
			r.ReadString(int16(len(testString)))
			r.ReadInt16()
			r.ReadInt32()
			r.ReadInt16()
			r.ReadInt16()
			r.ReadInt16()
			r.ReadInt16()
		}
	})
}

// BenchmarkPacketParallelEncodeDecodeLarge benchmarks concurrent large packet operations
func BenchmarkPacketParallelEncodeDecodeLarge(b *testing.B) {
	testString := "VeryLongPlayerNameWithManyCharacters"
	testData := make([]byte, 512)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := CreateWithOpcode(0x03)
			p.WriteInt32(12345)
			p.WriteString(testString)
			for j := 0; j < 10; j++ {
				p.WriteInt32(int32(j * 1000))
			}
			p.WriteBytes(testData)
			p.WriteInt64(1234567890123456789)

			r := NewReader(&p, 0)
			r.Skip(4)
			r.ReadByte()
			r.ReadInt32()
			r.ReadString(int16(len(testString)))
			for j := 0; j < 10; j++ {
				r.ReadInt32()
			}
			r.ReadBytes(512)
			r.ReadInt64()
		}
	})
}

// BenchmarkPacketHighVolumeWrites benchmarks heavy sequential writes
func BenchmarkPacketHighVolumeWrites(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewPacket()
		for j := 0; j < 1000; j++ {
			p.WriteByte(byte(j % 256))
			p.WriteInt16(int16(j))
			p.WriteInt32(int32(j * 1000))
		}
	}
}

// BenchmarkPacketHighVolumeReads benchmarks heavy sequential reads
func BenchmarkPacketHighVolumeReads(b *testing.B) {
	p := NewPacket()
	for j := 0; j < 1000; j++ {
		p.WriteByte(byte(j % 256))
		p.WriteInt16(int16(j))
		p.WriteInt32(int32(j * 1000))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewReader(&p, 0)
		for j := 0; j < 1000; j++ {
			r.ReadByte()
			r.ReadInt16()
			r.ReadInt32()
		}
	}
}

// BenchmarkPacketMassivePayload benchmarks very large packet operations
func BenchmarkPacketMassivePayload(b *testing.B) {
	testData := make([]byte, 10240) // 10KB
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := CreateWithOpcode(0x04)
		p.WriteBytes(testData)
		
		r := NewReader(&p, 0)
		r.Skip(5)
		r.ReadBytes(10240)
	}
}

// BenchmarkPacketStressTest simulates realistic high-load scenario
func BenchmarkPacketStressTest(b *testing.B) {
	playerNames := []string{
		"Player1", "Player2", "Player3", "Player4", "Player5",
		"Player6", "Player7", "Player8", "Player9", "Player10",
	}
	
	b.RunParallel(func(pb *testing.PB) {
		count := 0
		for pb.Next() {
			// Simulate complex player data packet
			p := CreateWithOpcode(0x05)
			p.WriteInt32(int32(count))
			p.WriteString(playerNames[count%len(playerNames)])
			p.WriteInt16(int16(count % 200))
			p.WriteInt32(int32(count * 1000))
			
			// Add inventory items
			p.WriteByte(byte(20)) // item count
			for j := 0; j < 20; j++ {
				p.WriteInt32(int32(1000 + j))
				p.WriteInt16(int16(j + 1))
			}
			
			// Add skills
			p.WriteByte(byte(10)) // skill count
			for j := 0; j < 10; j++ {
				p.WriteInt32(int32(2000 + j))
				p.WriteByte(byte(j + 1))
			}
			
			// Add quest data
			p.WriteInt16(int16(5)) // quest count
			for j := 0; j < 5; j++ {
				p.WriteInt16(int16(3000 + j))
				p.WriteByte(byte(j))
			}
			
			// Decode everything
			r := NewReader(&p, 0)
			r.Skip(4)
			r.ReadByte()
			r.ReadInt32()
			nameLen := playerNames[count%len(playerNames)]
			r.ReadString(int16(len(nameLen)))
			r.ReadInt16()
			r.ReadInt32()
			
			itemCount := r.ReadByte()
			for j := 0; j < int(itemCount); j++ {
				r.ReadInt32()
				r.ReadInt16()
			}
			
			skillCount := r.ReadByte()
			for j := 0; j < int(skillCount); j++ {
				r.ReadInt32()
				r.ReadByte()
			}
			
			questCount := r.ReadInt16()
			for j := 0; j < int(questCount); j++ {
				r.ReadInt16()
				r.ReadByte()
			}
			
			count++
		}
	})
}

// BenchmarkPacketThroughput measures packet processing throughput
func BenchmarkPacketThroughput(b *testing.B) {
	b.ReportAllocs()
	
	// Create 100 different packet templates
	packets := make([]Packet, 100)
	for i := range packets {
		p := CreateWithOpcode(byte(i % 256))
		p.WriteInt32(int32(i * 1000))
		p.WriteString("ThroughputTest")
		p.WriteInt16(int16(i))
		p.WriteInt32(int32(i * 10000))
		p.WriteBytes(make([]byte, 128))
		packets[i] = p
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		count := 0
		for pb.Next() {
			// Decode packets in rotation
			p := packets[count%len(packets)]
			r := NewReader(&p, 0)
			r.Skip(4)
			r.ReadByte()
			r.ReadInt32()
			r.ReadString(14)
			r.ReadInt16()
			r.ReadInt32()
			r.ReadBytes(128)
			count++
		}
	})
}

// BenchmarkThroughputSmallPackets measures throughput in mbit/s for small packets (64 bytes)
func BenchmarkThroughputSmallPackets(b *testing.B) {
	packetSize := 64
	testData := make([]byte, packetSize-5) // Account for header
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			p := CreateWithOpcode(0x01)
			p.WriteBytes(testData)
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadBytes(len(testData))
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputMediumPackets measures throughput in mbit/s for medium packets (512 bytes)
func BenchmarkThroughputMediumPackets(b *testing.B) {
	packetSize := 512
	testData := make([]byte, packetSize-5)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			p := CreateWithOpcode(0x02)
			p.WriteBytes(testData)
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadBytes(len(testData))
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputLargePackets measures throughput in mbit/s for large packets (4KB)
func BenchmarkThroughputLargePackets(b *testing.B) {
	packetSize := 4096
	testData := make([]byte, packetSize-5)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			p := CreateWithOpcode(0x03)
			p.WriteBytes(testData)
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadBytes(len(testData))
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputMassivePackets measures throughput in mbit/s for massive packets (64KB)
func BenchmarkThroughputMassivePackets(b *testing.B) {
	packetSize := 65536
	testData := make([]byte, packetSize-5)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			p := CreateWithOpcode(0x04)
			p.WriteBytes(testData)
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadBytes(len(testData))
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputRealisticMix measures throughput with realistic packet size distribution
func BenchmarkThroughputRealisticMix(b *testing.B) {
	// Create packets of various realistic sizes
	packets := []Packet{
		createPacketWithSize(32),   // Small control packet
		createPacketWithSize(64),   // Small data packet
		createPacketWithSize(128),  // Movement packet
		createPacketWithSize(256),  // Chat/inventory packet
		createPacketWithSize(512),  // Player info packet
		createPacketWithSize(1024), // Map data packet
		createPacketWithSize(2048), // Large data packet
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		count := 0
		for pb.Next() {
			p := packets[count%len(packets)]
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.GetRestAsBytes()
			count++
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
		b.ReportMetric(float64(b.N)/elapsed, "packets/s")
	}
}

// BenchmarkThroughputEncodeOnly measures encoding throughput in mbit/s
func BenchmarkThroughputEncodeOnly(b *testing.B) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			p := CreateWithOpcode(0x05)
			p.WriteInt32(123456)
			p.WriteString("PlayerName")
			p.WriteInt16(100)
			p.WriteBytes(testData)
			localBytes += int64(len(p))
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputDecodeOnly measures decoding throughput in mbit/s
func BenchmarkThroughputDecodeOnly(b *testing.B) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}
	
	p := CreateWithOpcode(0x06)
	p.WriteInt32(123456)
	p.WriteString("PlayerName")
	p.WriteInt16(100)
	p.WriteBytes(testData)
	
	packetSize := int64(len(p))
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		for pb.Next() {
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadInt32()
			r.ReadString(10)
			r.ReadInt16()
			r.ReadBytes(1024)
			localBytes += packetSize
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
	}
}

// BenchmarkThroughputMaxLoad simulates maximum sustainable load
func BenchmarkThroughputMaxLoad(b *testing.B) {
	// Simulate 1000 concurrent players sending packets
	playerPackets := make([]Packet, 1000)
	for i := range playerPackets {
		p := CreateWithOpcode(byte(i % 256))
		p.WriteInt32(int32(i))
		p.WriteString("Player")
		p.WriteInt16(int16(i % 200))
		p.WriteInt32(int32(i * 1000))
		p.WriteBytes(make([]byte, 256))
		playerPackets[i] = p
	}
	
	var totalBytes int64
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		var localBytes int64
		count := 0
		for pb.Next() {
			p := playerPackets[count%len(playerPackets)]
			localBytes += int64(len(p))
			
			r := NewReader(&p, 0)
			r.Skip(5)
			r.ReadInt32()
			r.ReadString(6)
			r.ReadInt16()
			r.ReadInt32()
			r.ReadBytes(256)
			count++
		}
		atomic.AddInt64(&totalBytes, localBytes)
	})
	
	elapsed := b.Elapsed().Seconds()
	if elapsed > 0 {
		total := atomic.LoadInt64(&totalBytes)
		mbits := float64(total*8) / (1000000 * elapsed)
		b.ReportMetric(mbits, "mbit/s")
		b.ReportMetric(float64(total)/elapsed, "bytes/s")
		b.ReportMetric(float64(b.N)/elapsed, "packets/s")
	}
}

// Helper function to create packets of specific size
func createPacketWithSize(size int) Packet {
	p := CreateWithOpcode(0xFF)
	if size > 5 {
		data := make([]byte, size-5)
		for i := range data {
			data[i] = byte(i % 256)
		}
		p.WriteBytes(data)
	}
	return p
}
