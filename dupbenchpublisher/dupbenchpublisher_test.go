package main

import (
	"hash"
	"hash/crc32"
	"testing"
	"time"

	"github.com/evanj/streambench/messages"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// May or may not be necessary; better safe than sorry https://github.com/golang/go/issues/27400
var preventCompilerRemovingBenchmark uint32

var crc32cTable = crc32.MakeTable(crc32.Castagnoli)

type testFixture struct {
	msg    *messages.DuplicateTest
	hasher hash.Hash32
	buf    []byte
}

func makeFixture(b *testing.B) *testFixture {
	ts := timestamppb.New(time.Unix(1574604732, 123456789))
	msg := &messages.DuplicateTest{
		GoroutineId: "00000000000000000000000000000000",
		Sequence:    123456,
		Created:     ts,
	}
	hasher := crc32.New(crc32cTable)
	b.ReportAllocs()

	return &testFixture{msg, hasher, nil}
}

func (f *testFixture) mustMarshalAppend() {
	byteLen := proto.Size(f.msg)
	f.buf = f.buf[:0]
	if cap(f.buf) < byteLen {
		f.buf = make([]byte, 0, 2*byteLen)
	}
	out, err := proto.MarshalOptions{UseCachedSize: true}.MarshalAppend(f.buf, f.msg)
	if err != nil {
		panic(err)
	}
	f.hasher.Write(out)
}

func BenchmarkMarshal(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serialized, err := proto.Marshal(fixture.msg)
		if err != nil {
			b.Fatal(err.Error())
		}
		fixture.hasher.Write(serialized)
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkMarshalAppend(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fixture.mustMarshalAppend()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkTimestampPBNow(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fixture.msg.Created = timestamppb.Now()

		fixture.mustMarshalAppend()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkTimestampPBReuse(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := time.Now()
		fixture.msg.Created.Seconds = t.Unix()
		fixture.msg.Created.Nanos = int32(t.Nanosecond())

		fixture.mustMarshalAppend()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}
