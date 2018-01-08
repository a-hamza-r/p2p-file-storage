# p2p-file-system
A peer-to-peer file storage system implemented in GOLANG

## Initial Setup:
Copy the repository to your $GOPATH and run the following commands:
- "go install peer"
- "go install tracker"  

Both the executables will be created in the 'bin' folder of your repository.
In order to avoid temporary files made by vim while editing any file, add the following lines in your .vimrc file (vim configuration file on unix based OS) present in you $HOME directory:
set nobackup
set nowritebackup
set noswapfile

## How To Run:
The tracker runs on port 8080. The port number and the mount point of peer can be given from the command line. To run the tracker: "./tracker" from the bin directory. To run the peer: "./peer -mount /path/to/mountpoint/ -port some_random_port" with admin rights. Remember not to run different peers on same port or same mountpoint in case testing on the single machine. You have to enter an 'ID' manually after running the peer which serves as a unique identity for the peer in the network. 

## Description:
The peer module runs its own FUSE filesystem. It contacts the tracker to get the address of other peers in the system. When writing anything to a file, the data gets divided into block 512 bytes each. These blocks are sent to other peers by initiating a TCP connection with them. For load balancing, the peer can either send the next block to the peer with least amount of data stored previously or it can choose the peer with the least response time. For fault tolerance, the peer simply keeps multiple copies of its blocks on separate peers. If one peer fails to respond it asks the other peer holding that same data.

The tracker follows a flood-fill algorithm. When a new node arrives, it sends it address to all the peers in the node. Similarly, when a particular node changes its address it informs all the other nodes.

#### Functions implemented:
- Division of file into blocks
- Distributed writeAll function
- Distributed readAll function
- Write by offset
- Read by offset
- Seperate working peers
- Tracker to share IPs
- Fault tolerance
- Load balancing
