package mcfs

import (
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"github.com/levigross/grequests"
	"fmt"
)

type Dir struct{}

type File struct{}

func (_ Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	//fmt.Println("Calling Attr for dir")
	//attr.Inode = 1
	attr.Mode = os.ModeDir
	attr.Size = uint64(4096)
	return nil
}

func (_ Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	//fmt.Println("Lookup ", name)
	return File{}, nil
}

type Project struct {
	Name string `json:"name"`
}

var opts = grequests.RequestOptions{InsecureSkipVerify: true}

func (_ Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var projects []Project
	var entries []fuse.Dirent
	var apikey = os.Getenv("MCAPIKEY")
	var mcurl = os.Getenv("MCURL")
	var urlstr = fmt.Sprintf("%s/api/v2/projects?apikey=%s", mcurl, apikey)
	fmt.Println("urlstr", urlstr)
	resp, err := grequests.Get(urlstr, &opts)
	if (err != nil) {
		fmt.Println("Error", err)
		return entries, nil
	}
	resp.JSON(&projects)

	for _, proj := range projects {
		entries = append(entries, fuse.Dirent{Name: proj.Name, Type: fuse.DT_Dir,})
	}
	return entries, nil
}

func (f File) Attr(ctx context.Context, attr *fuse.Attr) error {
	//attr.Inode = 2
	attr.Mode = os.ModeDir // 0444
	attr.Size = 4096 // uint64(len("helloworld"))
	return nil
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
		return []byte("helloworld\n"), nil
}
