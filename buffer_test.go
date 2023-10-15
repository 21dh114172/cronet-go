package cronet_test

import (
	"bytes"
	"github.com/sagernet/cronet-go"
	"github.com/sagernet/cronet-go/utils"
	"github.com/sagernet/sing/common/buf"
	"sync"
	"testing"
)

func TestBuffer(t *testing.T) {
	randomBytes, err := utils.GetRandomBytes()
	if err != nil {
		t.Fatal(err)
	}
	b := cronet.NewBuffer()
	b.InitWithAlloc(1024)
	if b.Size() != 1024 {
		t.Fatal("bad size")
	}
	data := buf.StackNewSize(1024)
	data.ReadFullFrom(randomBytes, 1024)
	copy(b.DataSlice(), data.Bytes())
	if bytes.Compare(b.DataSlice(), data.Bytes()) != 0 {
		t.Fatal("bad data")
	}
	b.Destroy()
}

func TestManagedBuffer(t *testing.T) {
	randomBytes, err := utils.GetRandomBytes()
	if err != nil {
		t.Fatal(err)
	}
	data := buf.StackNewSize(1024)
	data.ReadFullFrom(randomBytes, 1024)
	dataCopy := utils.CopyBuffer(data)
	b := cronet.NewBuffer()
	var wg sync.WaitGroup
	wg.Add(1)
	b.InitWithDataAndCallback(data.Bytes(), cronet.NewBufferCallback(func(callback cronet.BufferCallback, buffer cronet.Buffer) {
		wg.Done()
	}))
	if b.Size() != 1024 {
		t.Fatal("bad size")
	}
	if bytes.Compare(b.DataSlice(), dataCopy) != 0 {
		t.Fatal("bad data")
	}
	b.Destroy()
	wg.Wait()
}
