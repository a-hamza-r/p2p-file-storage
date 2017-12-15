package main

import (
	"testing"
	"net"
	"bazil.org/fuse"
	"golang.org/x/net/context"
	"fmt"
	"bytes"
	"time"
)

func TestInit(t *testing.T) {
	
	master_port := "8000"
	interface_addr, _ := net.InterfaceAddrs()
	local_ip := interface_addr[0].String()
	master := Peer{Ip: local_ip, Port: master_port}
	fmt.Println("Master details: ", master)
	go listen(master)
	for {

		if len(connList) == 2 {
			break
		}

	}
}

func TestFaultTolerance(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.DataNodes = make(map[uint64][]*Peer)
	f.Replicas = 2		// number of replicas under user's control
	f.InitNode()

	my_size := 10
	data := []byte(RandStringBytes(my_size))
	ctx := context.TODO()
	// Write
	req := &fuse.WriteRequest{
		Offset: 0,
		Data:   data,
	}
	resp := &fuse.WriteResponse{}
	err := f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if resp.Size != my_size {
       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, my_size)
    }
    
    time.Sleep(10)
    f.DataNodes[0][0].Conn.Close()

	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: my_size,
	}
	rresp := &fuse.ReadResponse{}
	err = f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, data) {
		t.Errorf("Data not equal")
	}
}

/*

func TestWriteRead(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.InitNode()
	my_size := 4107
	data := []byte(RandStringBytes(my_size))
	ctx := context.TODO()
	// Write
	req := &fuse.WriteRequest{
		Offset: 0,
		Data:   data,
	}
	resp := &fuse.WriteResponse{}
	err := f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if resp.Size != my_size {
       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, my_size)
    }
	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: my_size,
	}
	rresp := &fuse.ReadResponse{}
	err = f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, data) {
		t.Errorf("Data not equal")
	}
}

func TestMultipleWrite(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.InitNode()
	my_size := 8192
	data := []byte(RandStringBytes(my_size))
	ctx := context.TODO()
	// Write
	for i := int64(0); i < 8192; i += 512 {
			
		req := &fuse.WriteRequest{
			Offset: i,
			Data:   data[i : i+512],
		}
		resp := &fuse.WriteResponse{}
		err := f.Write(ctx, req, resp)
		if err != nil {
	       t.Errorf("Error occurred: %s", err)
		}
		if resp.Size != 512 {
	       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, 512)
	    }
	}
	if f.attributes.Size != uint64(my_size) {
       t.Errorf("Size was incorrect, got: %d, want: %d.", f.attributes.Size, my_size)
    }
    if f.attributes.Blocks != uint64(16) {
       t.Errorf("Size was incorrect, got: %d, want: %d.", f.attributes.Blocks, 16)
    }
	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: my_size,
	}
	rresp := &fuse.ReadResponse{}
	err := f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, data) {
		t.Errorf("Data not equal")
	}
}



func TestTruncateFile(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.InitNode()
	my_size := 4107
	data := []byte(RandStringBytes(my_size))
	ctx := context.TODO()
	// Write
	req := &fuse.WriteRequest{
		Offset: 0,
		Data:   data,
	}
	resp := &fuse.WriteResponse{}
	err := f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if resp.Size != my_size {
       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, my_size)
    }
    // Truncate
	treq := &fuse.SetattrRequest{Size: 1000, Valid: fuse.SetattrSize}
	tresp := &fuse.SetattrResponse{Attr: f.attributes}
	err = f.Setattr(ctx, treq, tresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if f.attributes.Size != uint64(1000) {
       t.Errorf("Size was incorrect, got: %d, want: %d.", f.attributes.Size, 1000)
    }
    if f.attributes.Blocks != uint64(2) {
       t.Errorf("Size was incorrect, got: %d, want: %d.", f.attributes.Blocks, 2)
    }
	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: my_size,
	}
	rresp := &fuse.ReadResponse{}
	err = f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, data[:1000]) {
		t.Errorf("Data not equal")
	}
}	


func TestOverWrite(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.InitNode()
	my_size := 4107
	data := []byte(RandStringBytes(my_size))
	newData := []byte(RandStringBytes(my_size))
	ctx := context.TODO()
	// Write
	req := &fuse.WriteRequest{
		Offset: 0,
		Data:   data,
	}
	resp := &fuse.WriteResponse{}
	err := f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if resp.Size != my_size {
       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, my_size)
    }

    req2 := &fuse.WriteRequest{
		Offset: 0,
		Data:   newData,
	}
	resp2 := &fuse.WriteResponse{}
	err = f.Write(ctx, req2, resp2)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if resp.Size != my_size {
       t.Errorf("Size was incorrect, got: %d, want: %d.", resp.Size, my_size)
    }
	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: my_size,
	}
	rresp := &fuse.ReadResponse{}
	err = f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, newData) {
		t.Errorf("Data not equal")
	}
}

func TestWriteOffset(t *testing.T) {
	
	// Initialise
	f := new(File)
	f.InitNode()
	data := []byte("the cat in the hat sat on the bat")
	ctx := context.TODO()
	// Write
	req := &fuse.WriteRequest{
		Offset: 0,
		Data:   data,
	}
	resp := &fuse.WriteResponse{}
	err := f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
    // Write by Offset
    req = &fuse.WriteRequest{
		Offset: 19,
		Data:   []byte("ran across the mat until he was very tired"),
	}
	resp = &fuse.WriteResponse{}
	err = f.Write(ctx, req, resp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	//Read
 	rreq := &fuse.ReadRequest{
		Offset: 0, Size: 61,
	}
	rresp := &fuse.ReadResponse{}
	err = f.Read(ctx, rreq, rresp)
	if err != nil {
       t.Errorf("Error occurred: %s", err)
	}
	if !bytes.Equal(rresp.Data, []byte("the cat in the hat ran across the mat until he was very tired")) {
		t.Errorf("Data not equal")
	}
	if len(rresp.Data) != 61 {
		t.Errorf("Size was incorrect, got: %d, want: %d.", len(rresp.Data), 61)
	}
}

*/