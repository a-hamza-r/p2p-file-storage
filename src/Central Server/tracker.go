package main

import (
  "net"
  "fmt"
  "os"
  "log"
  "strings"
  "encoding/json"
)

var userMap = make(map[string]*User)  // user name password

type User struct {
	IP string
	Port string
	// Conn net.Conn
	ListenAddr string
}



type message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
	BlockName uint64 `json:"name"`
	PeerID string `json:"id"`
	ListenAddr string `json:"listenaddr"`
	PeerAddr string `json:"peeraddr"`
}



func listen(master User) {
  
	listen, err := net.Listen("tcp", ":" + master.Port)
	defer listen.Close()
	
	if err != nil {
		log.Fatalf("Socket listen port %s failed,%s", master.Port, err)
		os.Exit(1)
	}
	
	log.Printf("Begin listen port: %s", master.Port)
	
	for {
		conn, err := listen.Accept()
		
		if err != nil {
			log.Fatalln(err)
			continue
		} 
		
		go handler(conn)
	}
}


func handler(conn net.Conn) { 
  
	s_addr := strings.Split(conn.RemoteAddr().String(), ":")
	user := User{IP: s_addr[0], Port: s_addr[1]}
	log.Printf("Got request" + " from: %v",user)
	var msg message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	if err != nil {
	    fmt.Println(err)
	}
  
	if msg.Type == "login" {
    
	    name := msg.PeerID

	    if _, ok := userMap[name]; ok {
    
			log.Println(name, "logged in!")
			if userMap[name].IP != user.IP || userMap[name].Port != user.Port || userMap[name].ListenAddr != msg.ListenAddr {
        
		        log.Println(name, "changed its IP/Port")
				for key, value := range userMap {
					if key != name {
						updateUser(name, msg.ListenAddr, value)
					}
				}

		        user.ListenAddr = msg.ListenAddr
				
				for key, value := range userMap {
					if key != name {
						addUser((strings.Split(key, ":"))[0], value.ListenAddr, &user)
					}
				}

		        userMap[name] = &user
			}

	    } else {

			for _, value := range userMap {
		        addUser(name, msg.ListenAddr, value)
			}
	      
			user.ListenAddr = msg.ListenAddr
			for key, value := range userMap {
		        addUser((strings.Split(key, ":"))[0], value.ListenAddr, &user)
			}
	      
			userMap[name] = &user

		}
	}
	conn.Close()

}


func addUser(NewUserID string, NewUserAddr string, User *User) {

	conn, err := net.Dial("tcp", User.ListenAddr)    
	if err != nil {
	    log.Println(err)
	    return
	}

	m := message{"add", []byte(""), 0, NewUserID, "", NewUserAddr}
	err = json.NewEncoder(conn).Encode(&m)
	if err != nil {
	    fmt.Println(err)
	}
	var msg message
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&msg)
	if err != nil {
	    fmt.Println(err)
	}    
}

func updateUser(ChangedUserID string, ChangedUserAddr string, User *User) {

	conn, err := net.Dial("tcp", User.ListenAddr)    
	if err != nil {
	    log.Println(err)
	    return
	}

	m := message{"update", []byte(""), 0, ChangedUserID, "", ChangedUserAddr}
	err = json.NewEncoder(conn).Encode(&m)
	if err != nil {
		fmt.Println(err)
	}
	var msg message
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&msg)
	if err != nil {
	    fmt.Println(err)
	}
}

func main() {

	central_port := "8080"
	interface_addr, _ := net.InterfaceAddrs()
	local_IP := interface_addr[0].String()
	central := User{IP: local_IP, Port: central_port}
	fmt.Println("Central server details: ", central)

	listen(central)
}