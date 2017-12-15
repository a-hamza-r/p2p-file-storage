package main

import (
	"net"
	"fmt"
	"os"
	"log"
	"strings"
	"flag"
	// "path"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"io"
	"strconv"
)


// functions
func selectMountpoint() string {
	
	mountpoint := *flag.String("mount", "/mnt/fmount", "folder to mount")
	// flag.Parse()
	// if flag.NArg() == 0  {
	// 	log.Printf("Usage of %s:\n", os.Args[0])
	// 	log.Printf("  %s MOUNTPOINT\n", os.Args[0])
	// 	flag.PrintDefaults()

	// 	os.Exit(2)
	// }
	// mountpoint := flag.Arg(1)
	return mountpoint
}


func listen(me Peer) {

	listen, err := net.Listen("tcp", ":" + me.Port)
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %s failed,%s", me.Port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %s", me.Port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}	
		handler(conn, me)
	}
}

func handler(conn net.Conn, me Peer) {

	// s_addr := strings.Split(conn.RemoteAddr().String(),":")
	// otherPeer := Peer{IP: s_addr[0], Port: s_addr[1], Conn: conn}
	// log.Printf("Got request" + " from: %v",otherPeer)
	
	// // Implement a table to keep the IDs and IP addresses
	// // connList = append(connList, &otherPeer)	
	// // // get current working directory
	// // ex, err := os.Executable()		
 // //    if err != nil {
 // //        panic(err)
 // //    }
 // //    exPath := path.Dir(ex)
 // //    _ = exPath
	// // myDir := "/mnt" + "/" + "<" + myport + ">" + otherPeer.String()
	// // // create folder with peer's address
 // //    if _, err := os.Stat(myDir); os.IsNotExist(err) {
	// // 	os.Mkdir(myDir, 0777)
	// // }
	go manageMesseges(conn,me)
}

func connectToServer(me Peer, dst Peer) {
	
	// fix IP and dial
	// myAddr, err := net.ResolveIPAddr("IP", "127.0.0.1")
 //    if err != nil {
 //        panic(err)
 //    }
 //    localTCPAddr := net.TCPAddr{
 //        IP: myAddr.IP,
 //   	    Port: myport}
	// d := &net.Dialer{LocalAddr: &localTCPAddr,Timeout: time.Duration(10)*time.Second}
    
    // use conn istead of _
    conn, err := net.Dial(dst.Network(), dst.String())   	
   	if err != nil {
		log.Fatalln(err)

    } else {
        log.Println("Connected to central server")
	}

	fmt.Print("id: ")
	var id string
    fmt.Scanln(&id)
	myID = id

	// load connList
	// backupDir := "/mnt/"+myID+"_backup/"
	// content, err := ioutil.ReadFile(backupDir + "connList")
    // json.Unmarshal(content, &connList)
    // log.Println(connList)

	msg := message{"login", []byte(""), 0, myID, me.IP+":"+me.Port, ""}
	json.NewEncoder(conn).Encode(&msg)
	// go updatePeers(conn)
}

// func updatePeers(conn net.Conn) {
  	
//   	for {
//     	var msg message
//       	decoder := json.NewDecoder(conn)
//       	err := decoder.Decode(&msg) 
      	
//       	if err == io.EOF {
//       		log.Println("Err: ",err)
//         	conn.Close()
//         	break
//       	}
//       	if msg.Type == "add" {
//         	fmt.Println(msg)
//         	toDialList = append(toDialList, msg.PeerAddr)
//         	ack := central_message{"ack", "", "","",""}
// 			json.NewEncoder(conn).Encode(&ack)
//       	}
// 	}
// } 

func manageMesseges(conn net.Conn, myInfo Peer) {
    	
	var msg message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)    	
	checkFatal(err)
	
	if err == io.EOF {
		conn.Close()
		return
	}
	
	if msg.Type == "send" {
		filename := strconv.Itoa(int(msg.Name))

		myDir := "/mnt/" + (msg.PeerID)
		// create folder with peer's address
	    if _, err := os.Stat(myDir); os.IsNotExist(err) {
			os.Mkdir(myDir, 0777)
		}

		f, err := os.Create(filepath.Join(getPeerDir(msg.PeerID), filename))
		checkFatal(err)
		f.Chmod(0777)
    	b, err := f.WriteString(string(msg.Data))
    	log.Printf("Wrote %d bytes to file: %s \n", b, filename)
    	f.Sync()
    	f.Close()
    	
    	send_ack := message{"send_ack", []byte(""), msg.Name,"","",""}
		json.NewEncoder(conn).Encode(&send_ack)

	} else if msg.Type == "recv" {

		filename := strconv.Itoa(int(msg.Name))
		dat, err := ioutil.ReadFile(filepath.Join(getPeerDir(msg.PeerID), filename))
		if err != nil {
			log.Println(err)
		}
    	log.Printf("Sending file: %s \n", filename)
	    
	    recv_ack := message{"recv_ack", dat, msg.Name,"","",""}
		json.NewEncoder(conn).Encode(&recv_ack)


	} else if msg.Type == "delete" {
		
		filename := strconv.Itoa(int(msg.Name))
    	log.Printf("Removing file: %s \n", filename)
		err := os.Remove(filepath.Join(getPeerDir(msg.PeerID), filename))
		checkError(err)

		del_ack := message{"del_ack", []byte(""), msg.Name,"","",""}
		json.NewEncoder(conn).Encode(&del_ack)


	} else if msg.Type == "add" {
    	
    	fmt.Println(msg)
    	s_addr := strings.Split(msg.PeerAddr,":")
    	otherPeer := Peer{ID: msg.PeerID, IP: s_addr[0], Port: s_addr[1], NetType: "tcp"}

		myDir := "/mnt/" + (otherPeer.ID)
		otherPeer.PathToFiles = myDir
		
		connList = append(connList, &otherPeer)	
    	ack := message{"ack", []byte(""), 0, "", "", ""}
		json.NewEncoder(conn).Encode(&ack)

  	} else if msg.Type == "update" {
    	
    	fmt.Println(msg)
		for _, p := range connList {
			if p.ID == msg.PeerID {
				p_addr := strings.Split(msg.PeerAddr, ":")
				p.IP = p_addr[0]
				p.Port = p_addr[1]
			}
		}
    	ack := message{"ack", []byte(""), 0, "", "", ""}
		json.NewEncoder(conn).Encode(&ack)
  	}
  	conn.Close()

}

func main() {

	myport := flag.Int("port", 9000, "Port to run this node on")
	mountpoint := flag.String("mount", "/mnt/fmount", "folder to mount")
    flag.Parse()

	// interface_addr, _ := net.InterfaceAddrs()
	// local_IP := interface_addr[0].String()
	local_IP := "127.0.0.1"
	me := Peer{IP: local_IP, Port: strconv.Itoa(*myport)}
	fmt.Println("My details:", me)
	
	// connect to central server
	// server_port := 8080
	// server_IP := "127.0.0.1"
	server := Peer{IP: central_IP, Port: strconv.Itoa(central_port), NetType: "tcp"}
	// start listening to incoming connections


	


	go listen(me)

	connectToServer(me, server)

	// mount the FUSE file system
	// mountpoint := selectMountpoint()

	FUSE(*mountpoint)

}

