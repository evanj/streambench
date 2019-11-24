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
}

func makeFixture() *testFixture {
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
	return &testFixture{msg, hasher}
}

func BenchmarkMarshal(b *testing.B) {
	fixture := makeFixture()
	b.ReportAllocs()

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
	fixture := makeFixture()
	buf := proto.NewBuffer(make([]byte, proto.Size(fixture.msg)))
	buf.Reset()
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := buf.Marshal(fixture.msg)
		if err != nil {
			b.Fatal(err.Error())
		}
		fixture.hasher.Write(buf.Bytes())
		buf.Reset()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkBufferTimeNew(b *testing.B) {
	fixture := makeFixture()
	buf := proto.NewBuffer(make([]byte, proto.Size(fixture.msg)))
	buf.Reset()
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fixture.msg.Created = ptypes.TimestampNow()

		err := buf.Marshal(fixture.msg)
		if err != nil {
			b.Fatal(err.Error())
		}
		fixture.hasher.Write(buf.Bytes())
		buf.Reset()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}

func BenchmarkBufferTimeReuse(b *testing.B) {
	fixture := makeFixture()
	buf := proto.NewBuffer(make([]byte, proto.Size(fixture.msg)))
	buf.Reset()
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t := time.Now()
		fixture.msg.Created.Seconds = t.Unix()
		fixture.msg.Created.Nanos = int32(t.Nanosecond())

		err := buf.Marshal(fixture.msg)
		if err != nil {
			b.Fatal(err.Error())
		}
		fixture.hasher.Write(buf.Bytes())
		buf.Reset()
	}
	preventCompilerRemovingBenchmark = fixture.hasher.Sum32()
}
