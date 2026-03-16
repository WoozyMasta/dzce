// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import "testing"

const (
	benchmarkLayerSize = 512
	benchmarkUsageBits = 32
	benchmarkValueBits = 8
)

// benchmarkAreaFlagsFixture stores reusable benchmark payload and decoded form.
type benchmarkAreaFlagsFixture struct {
	File *AreaFlagsMapFile
	Data []byte
}

// buildBenchmarkAreaFlagsFixture creates deterministic binary benchmark data.
func buildBenchmarkAreaFlagsFixture(b *testing.B) benchmarkAreaFlagsFixture {
	b.Helper()

	file := &AreaFlagsMapFile{
		Header: AreaFlagsHeader{
			LayerWidth:  benchmarkLayerSize,
			LayerHeight: benchmarkLayerSize,
			WorldWidth:  15360,
			WorldHeight: 15360,
			UsageBits:   benchmarkUsageBits,
		},
		Layers: make([]AreaFlagsLayer, benchmarkUsageBits+benchmarkValueBits),
	}

	layerByteSize, err := file.Header.LayerByteSize()
	if err != nil {
		b.Fatalf("LayerByteSize() error: %v", err)
	}

	for layerIndex := range file.Layers {
		bits := make([]byte, layerByteSize)

		// Fill bytes with deterministic pseudo-random pattern.
		seed := uint64(0x9E3779B97F4A7C15) + uint64(layerIndex+1)*0xBF58476D1CE4E5B9
		for i := range bits {
			seed ^= seed >> 30
			seed *= 0xBF58476D1CE4E5B9
			seed ^= seed >> 27
			seed *= 0x94D049BB133111EB
			seed ^= seed >> 31

			bits[i] = byte(seed)
		}

		file.Layers[layerIndex].Bits = bits
	}

	data, err := EncodeAreaFlagsMap(file)
	if err != nil {
		b.Fatalf("EncodeAreaFlagsMap() setup error: %v", err)
	}

	return benchmarkAreaFlagsFixture{
		File: file,
		Data: data,
	}
}

func BenchmarkAreaFlagsDecode(b *testing.B) {
	fixture := buildBenchmarkAreaFlagsFixture(b)
	b.SetBytes(int64(len(fixture.Data)))

	b.ResetTimer()
	for range b.N {
		if _, err := DecodeAreaFlagsMap(fixture.Data); err != nil {
			b.Fatalf("DecodeAreaFlagsMap() error: %v", err)
		}
	}
}

func BenchmarkAreaFlagsEncode(b *testing.B) {
	fixture := buildBenchmarkAreaFlagsFixture(b)
	b.SetBytes(int64(len(fixture.Data)))

	b.ResetTimer()
	for range b.N {
		if _, err := EncodeAreaFlagsMap(fixture.File); err != nil {
			b.Fatalf("EncodeAreaFlagsMap() error: %v", err)
		}
	}
}

func BenchmarkAreaFlagsDecodeEncodeRoundtrip(b *testing.B) {
	fixture := buildBenchmarkAreaFlagsFixture(b)
	b.SetBytes(int64(len(fixture.Data)))

	b.ResetTimer()
	for range b.N {
		file, err := DecodeAreaFlagsMap(fixture.Data)
		if err != nil {
			b.Fatalf("DecodeAreaFlagsMap() error: %v", err)
		}

		if _, err = EncodeAreaFlagsMap(file); err != nil {
			b.Fatalf("EncodeAreaFlagsMap() error: %v", err)
		}
	}
}
