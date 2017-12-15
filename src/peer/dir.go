package main

import (
	"log"
	"os"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context" // need this cause bazil lib doesn't use syslib context lib
)

// Dir implements both Node and Handle for the root directory.
type Dir struct{
	Node
	files       *[]*File
	directories *[]*Dir
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	
	a.Inode = d.Attributes.Inode
	a.Mode = os.ModeDir | 0555
	return nil
}

func (d *Dir) Lookup(ctx context.Context, Name string) (fs.Node, error) { //** find command **//
	
	log.Println("Requested lookup for", Name)
	if d.files != nil {
		for _, n := range *d.files {
			if n.Name == Name {
				return n, nil
			}
		}
	}
	if d.directories != nil {
		for _, n := range *d.directories {
			if n.Name == Name {
				return n, nil
			}
		}
	}
	return nil, fuse.ENOENT
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	
	log.Println("Reading all directory")
	var content []fuse.Dirent
	if d.files != nil {
		for _, f := range *d.files {
			content = append(content, fuse.Dirent{Inode: f.Attributes.Inode, Type: fuse.DT_File, Name: f.Name})
		}
	}
	if d.directories != nil {
		for _, dir := range *d.directories {
			content = append(content, fuse.Dirent{Inode: dir.Attributes.Inode, Type: fuse.DT_Dir, Name: dir.Name})
		}
	}
	return content, nil
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	
	log.Println("Create request for Name", req.Name)
	f := &File{Node: Node{Name: req.Name}}
	f.DataNodes = make(map[uint64][]*OneBlockInfo)
	f.Replicas = 2		// number of replicas under user's control
	f.InitNode()
	if d.files != nil {
		(*d.files) = append(*d.files, f)
	}
	return f, f, nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	
	log.Println("Remove request for ", req.Name)
	if req.Dir && *d.directories != nil {
		for index,value:= range *d.directories {
			if req.Name == value.Name {
				*d.directories = append((*d.directories)[:index], (*d.directories)[index+1:]...)
				return nil
			}
		}

	} else if !req.Dir && *d.files != nil {
		for index,value:= range *d.files {
			if req.Name == value.Name {
				*d.files = append((*d.files)[:index], (*d.files)[index+1:]...)
				for i := value.Attributes.Blocks-1; i >= 0; i-- {
					for j := 0; j < value.Replicas; j++{
						deleteBlock(value.DataNodes[i][j].PeerInfo, value.DataNodes[i][j].Name)
					}
					value.DataNodes[i] = value.DataNodes[i][:1]
					value.DataNodes[i][0].Used = false
					if i == 0 {
						break
					}
				}
				return nil
			}
		}
	}
	return fuse.ENOENT
}