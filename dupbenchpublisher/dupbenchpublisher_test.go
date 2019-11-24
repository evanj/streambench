package main

import (
	"hash"
	"hash/crc32"
	"testing"
	"time"

	"github.com/evanj/streambench/messages"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// May or may not be necessary; better safe than sorry https://github.com/golang/go/issues/27400
var preventCompilerRemovingBenchmark uint32

var crc32cTable = crc32.MakeTable(crc32.Castagnoli)

type testFixture struct {
	msg    *messages.DuplicateTest
	hasher hash.Hash32
	buf    *proto.Buffer
}

func makeFixture(b *testing.B) *testFixture {
	ts, err := ptypes.TimestampProto(time.Unix(1574604732, 123456789))
	if err != nil {
		panic(err)
	}
	msg := &messages.DuplicateTest{
		GoroutineId: "00000000000000000000000000000000",
		Sequence:    123456,
		Created:     ts,
	}
	hasher := crc32.New(crc32cTable)
	buf := proto.NewBuffer(nil)
	b.ReportAllocs()

	return &testFixture{msg, hasher, buf}
}

func (f *testFixture) mustMarshalBuffer() {
	err := f.buf.Marshal(f.msg)
	if err != nil {
		panic(err)
	}
	f.hasher.Write(f.buf.Bytes())
	f.buf.Reset()
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

func BenchmarkBuffer(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fixture.mustMarshalBuffer()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkTimeNew(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fixture.msg.Created = ptypes.TimestampNow()

		fixture.mustMarshalBuffer()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkTimeReuse(b *testing.B) {
	fixture := makeFixture(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := time.Now()
		fixture.msg.Created.Seconds = t.Unix()
		fixture.msg.Created.Nanos = int32(t.Nanosecond())

		fixture.mustMarshalBuffer()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}
