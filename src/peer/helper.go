package main 

import (
	"encoding/json"
	"time"
	"sort"
	"log"
	"bazil.org/fuse"
	"net"
	"os"
)	



// global variables
var connList []*Peer
var dataBlockSize uint64 = 512
var blockIdentifier uint64 = 0
var inode uint64 = 0
var myID string

var central_IP = "127.0.0.1"
var central_port = 8080



// structures
type Node struct {
	Name string
	Attributes fuse.Attr
}

type Meta struct {
	Name string
	Attributes fuse.Attr
	DataNodes map[uint64][]*OneBlockInfo
	Replicas int
}

type OneBlockInfo struct {
	Name uint64
	PeerInfo *Peer
	Used bool
}

type Peer struct {
	ID string
	IP string
	Port string
	// Conn net.Conn
	DataSent uint64
	ResponseTime time.Duration
	NetType string
	PathToFiles string
}

type message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
	Name uint64 `json:"name"`
	PeerID string `json:"id"`
  	ListenAddr string `json:"listenaddr"`
  	PeerAddr string `json:"peeraddr"`
}

// type central_message struct {
  
//   Type string `json:"type"`
//   Id string `json:"id"`
//   Password string `json:"password"`
//   ListenAddr string `json:"listenaddr"`
//   PeerAddr string `json:"peeraddr"`

// }



// functions
func (n *Node) InitNode() {
	
	t := time.Now()
	n.Attributes.Inode = NewInode()      
    n.Attributes.Size = 0      			
    n.Attributes.Blocks = 0      		
	n.Attributes.Atime = t
	n.Attributes.Mtime = t
	n.Attributes.Ctime = t
	n.Attributes.Crtime = t
	n.Attributes.Mode = 0644 
	n.Attributes.Nlink = 0
	n.Attributes.Uid = 0
	n.Attributes.Gid = 0
	n.Attributes.Rdev = 0
	n.Attributes.BlockSize = uint32(dataBlockSize) // block size

}

func (p *Peer) String() string {
	return p.IP + ":" + p.Port
}

func (p *Peer) Network() string {
	return p.NetType
}

func checkError(e error){
	
	if e != nil {
		log.Println(e)
	}
}

func checkFatal(e error){
	
	if e != nil {
		log.Fatalln(e)
	}	
}

func Split(buf []byte, lim int) [][]byte {
	
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func NewInode() uint64 {
	
	inode += 1
	return inode
}


func Blocks(value uint64) uint64 { // Blocks returns the number of 512 byte blocks required
	
	if value == 0 {
		return 0
	}
	blocks := value / dataBlockSize
	if value%dataBlockSize > 0 {
		return blocks + 1
	}
	return blocks
}

func Offset2Block (value uint64) uint64 {
	
	return (value / dataBlockSize)
}

func BlockCheck(offsetBlock uint64, dataNodes *map[uint64][]*OneBlockInfo, buffer []byte, numReplicas *int) {
	
	if *numReplicas > len(connList) {
		*numReplicas = len(connList)
	}
	if offsetBlock < uint64(len(*dataNodes)) {
		if (*dataNodes)[offsetBlock][0].Used {
			for j := 0; j < len((*dataNodes)[offsetBlock]); j++ {
				sendBlock((*dataNodes)[offsetBlock][j].PeerInfo, buffer, (*dataNodes)[offsetBlock][j].Name)
			}
		} else {
			sortPeers1("data", connList)
			name := (*dataNodes)[offsetBlock][0].Name
			(*dataNodes)[offsetBlock] = make([]*OneBlockInfo, 0, *numReplicas)
			for j := 0; j < *numReplicas; j++ {
				singleBlock := &OneBlockInfo{name, connList[j], true}
				(*dataNodes)[offsetBlock] = append((*dataNodes)[offsetBlock], singleBlock)
				sendBlock(connList[j], buffer, (*singleBlock).Name)
			}
		}
	} else {
		sortPeers1("data", connList)
		for j := 0; j < *numReplicas; j++ {
			singleBlock := &OneBlockInfo{blockIdentifier, connList[j], true}
			if _, ok := (*dataNodes)[offsetBlock]; !ok {
				(*dataNodes)[offsetBlock] = make([]*OneBlockInfo, 0, *numReplicas)
			}
			(*dataNodes)[offsetBlock] = append((*dataNodes)[offsetBlock], singleBlock)
			sendBlock(connList[j], buffer, (*singleBlock).Name)
		}
		blockIdentifier++
		log.Println(blockIdentifier)
	}
}


func sendBlock(peer *Peer, data []byte, block uint64) {

	conn, err := net.Dial(peer.Network(), peer.String())   	
   	if err != nil {
		log.Fatalln(err)

    }
	// peer.Conn = getConn(peer).Conn
	// conn := peer.Conn	
	encrypted_data := Encrypt(data, []byte("123"), int64(len(data)))
	m := message{"send", encrypted_data, block, myID, "", ""}	
	json.NewEncoder(conn).Encode(&m)
	start_time := time.Now()
	var ack message
	decoder := json.NewDecoder(conn)
	peer.ResponseTime = time.Since(start_time)
	peer.DataSent += uint64(len(data)) 
 	err = decoder.Decode(&ack)
 	checkError(err)

 	conn.Close()

}

func recvBlock(peer *Peer, block uint64) ([]byte, error) {

	conn, err := net.Dial(peer.Network(), peer.String())   	
   	if err != nil {
		log.Fatalln(err)

    }
	// peer.Conn = getConn(peer).Conn
	// conn := peer.Conn
	m := message{"recv", []byte(""), block, myID, "", ""}
	err = json.NewEncoder(conn).Encode(&m)
	start_time := time.Now()	
	var ack message
	decoder := json.NewDecoder(conn)
 	peer.ResponseTime = time.Since(start_time)
 	err = decoder.Decode(&ack)
 	decrypted_data := Decrypt(ack.Data, []byte("123"))

 	conn.Close()
	return decrypted_data, err
}

func deleteBlock(peer *Peer, block uint64) {
	
	conn, err := net.Dial(peer.Network(), peer.String())   	
   	if err != nil {
		log.Fatalln(err)

    }

	// peer.C/onn = getConn(peer).Conn
	// conn := peer.Conn
	m := message{"delete", []byte(""), block, myID, "", ""}
	json.NewEncoder(conn).Encode(&m)
	start_time := time.Now()
	var ack message
	decoder := json.NewDecoder(conn)
 	peer.ResponseTime = time.Since(start_time)
 	err = decoder.Decode(&ack)
 	checkError(err)

 	conn.Close()

}

func sortPeers1(loadType string, peerArray []*Peer) {

	if loadType == "data" {
		
		sort.Slice(peerArray[:], func(i, j int) bool {
    		return peerArray[i].DataSent < peerArray[j].DataSent
		})

	} else if loadType == "time" {

		sort.Slice(peerArray[:], func(i, j int) bool {
    		return peerArray[i].ResponseTime < peerArray[j].ResponseTime
		})

	}
}

func sortPeers(loadType string, peerArray []*OneBlockInfo) {

	if loadType == "data" {
		
		sort.Slice(peerArray[:], func(i, j int) bool {
    		return peerArray[i].PeerInfo.DataSent < peerArray[j].PeerInfo.DataSent
		})

	} else if loadType == "time" {

		sort.Slice(peerArray[:], func(i, j int) bool {
    		return peerArray[i].PeerInfo.ResponseTime < peerArray[j].PeerInfo.ResponseTime
		})

	}
}


// func getConn(addr *Peer) (*Peer) {
	
// 	s_addr := addr.IP + ":" + addr.Port
// 	var conn *Peer
// 	for _, p := range connList {
// 		if p.Conn.RemoteAddr().String() == s_addr {
// 			conn = p
// 		}
// 	}
// 	return conn
// }

func getPeerDir(id string) (string) {
	
	var requiredPeer *Peer
	for _, p := range connList {
		if (*p).ID == id {
			requiredPeer = p
		}
	}
	return requiredPeer.PathToFiles
}


func SaveConnList() {

    j, err := json.Marshal(connList)
    if err != nil {
        log.Println("Error converting connList to json ",err)
        return
    }
	handle, err := os.Create("/mnt/"+myID+"_backup/"+"connList")
	defer handle.Close()
	if err != nil {
	    log.Println("Error creating connList file ",err)
	    return
	}
	handle.Chmod(0777)
	_, err = handle.WriteString(string(j))
	if err != nil {
	    log.Println("Error saving connList file ",err)
	    return
	}
	handle.Sync()
	log.Println("Saving connList file")

}

