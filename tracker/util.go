package tracker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var Log *log.Logger = log.New(os.Stderr, "", log.LstdFlags)

type debugging int

const (
	debugl debugging = iota
	errorl
)

const debug_flag = false

func (d debugging) Println(args ...interface{}) {
	if d == debugl && !debug_flag {
		return
	}
	Log.Println(args...)
}

func (d debugging) Printf(format string, args ...interface{}) {
	if d == debugl && !debug_flag {
		return
	}
	Log.Printf(format, args...)
}

func RandomBytes(n int) []byte {
	out := make([]byte, n)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		out[i] = byte(rand.Int())
	}
	return out
}

func RandomInt(n int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(n)
}

func to4Bytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func toInt64(b []byte) int64 {
	var i int64
	if len(b) < 8 {
		tmp := make([]byte, 8-len(b))
		tmp = append(tmp, b...)
		b = tmp
	}
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		fmt.Println("Error converting to int: ", err)
	}
	return i
}

func toInt(b []byte) int {
	return int(toInt64(b))
}
