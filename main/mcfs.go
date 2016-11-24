package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"path/filepath"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/levigross/grequests"
	"github.com/materials-commons/mcfs"
)

type MCFS struct {
	pathfs.FileSystem
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var projects []Project = nil
var projectsMap map[string]Project = nil

var (
	opts   = grequests.RequestOptions{InsecureSkipVerify: true}
	apikey = os.Getenv("MCAPIKEY")
	mcurl  = os.Getenv("MCURL")
)

func initIt() {
	fmt.Println("Init called")
	var (
		urlstr = fmt.Sprintf("%s/api/v2/projects?apikey=%s", mcurl, apikey)
	)
	fmt.Println("urlstr", urlstr)
	resp, err := grequests.Get(urlstr, &opts)
	if err != nil {
		log.Fatal("Failed to call grequests")
	}
	resp.JSON(&projects)
	for _, proj := range projects {
		proj.Name = strings.Replace(proj.Name, "/", "_", -1)
	}
}

func fixName(name string) string {
	return strings.Replace(name, "/", "_", -1)
}

func (fs *MCFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	//fmt.Println("name", name)
	ext := filepath.Ext(name)
	switch {
	case name == "":
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755}, fuse.OK
	case ext == "":
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755}, fuse.OK
	case isFile(ext):
		fmt.Println("attr on file", name)
		return &fuse.Attr{Mode: fuse.S_IFREG | 0644, Size: uint64(4096)}, fuse.OK
	default:
		//fmt.Println("GetAttr for", name)
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755}, fuse.OK
	}
}

func isFile(ext string) bool {
	switch ext {
	case "json":
		return true
	case "html":
		return true
	default:
		return false
	}
}

func (fs *MCFS) GetXAttr(name string, attribute string, context *fuse.Context) ([]byte, fuse.Status) {
	return nil, fuse.OK
}

func (fs *MCFS) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	return nil, fuse.OK
}

func (fs *MCFS) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	if name == "" && projects == nil {
		// Top level dir show projects
		fmt.Println("OpenDir: Retrieving projects")
		var (
			entries []fuse.DirEntry
			urlstr  = fmt.Sprintf("%s/api/v2/projects?apikey=%s", mcurl, apikey)
		)
		fmt.Println("urlstr", urlstr)
		resp, err := grequests.Get(urlstr, &opts)
		if err != nil {
			return nil, fuse.ENOENT
		}
		resp.JSON(&projects)
		for _, proj := range projects {
			entries = append(entries, fuse.DirEntry{Name: fixName(proj.Name), Mode: fuse.S_IFDIR | 0755})
			mcfs.AddProject(filepath.Join(name, fixName(proj.Name)), proj.ID)
		}
		return entries, fuse.OK
	} else if name == "" {
		var entries []fuse.DirEntry
		for _, proj := range projects {
			entries = append(entries, fuse.DirEntry{Name: fixName(proj.Name), Mode: fuse.S_IFDIR | 0755})
		}
		return entries, fuse.OK
	} else if _, ptype := mcfs.PathEntry(name); ptype == mcfs.PMProject {
		fmt.Println("Returning top level project entries")
		return []fuse.DirEntry{
			{Name: "json", Mode: fuse.S_IFDIR | 0755},
			{Name: "html", Mode: fuse.S_IFDIR | 0755},
			{Name: "text", Mode: fuse.S_IFDIR | 0755},
			{Name: "files", Mode: fuse.S_IFDIR | 0755},
		}, fuse.OK
	} else if dir, last := filepath.Split(name); isPathType(last) {
		mcfs.PathEntry(filepath.Clean(dir))
		return []fuse.DirEntry{
			{Name: fmt.Sprintf("%s.json", filepath.Base(dir)), Mode: fuse.S_IFREG | 0644},
			{Name: "samples", Mode: fuse.S_IFDIR | 0755},
			{Name: "experiments", Mode: fuse.S_IFDIR | 0755},
		}, fuse.OK
	} else if name == "" {

	}

	fmt.Println("Returning fuse.ENOENT for", name)
	return nil, fuse.ENOENT
}

func isPathType(what string) bool {
	switch what {
	case "json":
		return true
	case "text":
		return true
	case "files":
		return true
	case "html":
		return true
	default:
		return false
	}
}

func (fs *MCFS) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	fmt.Println("Open", name)
	ID, _ := mcfs.PathEntry("WE43")
	var (
		urlstr  = fmt.Sprintf("%s/api/v2/projects/%s?apikey=%s", mcurl, ID, apikey)
	)
	fmt.Println("urlstr", urlstr)
	resp, err := grequests.Get(urlstr, &opts)
	if err != nil {
		return nil, fuse.ENOENT
	}
	val := resp.String()
	fmt.Println("Project Str", val)
	return nil, fuse.ENOENT
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage: \n mcfs MOUNTPOINT")
	}

	fs := pathfs.NewPathNodeFs(&MCFS{FileSystem: pathfs.NewLockingFileSystem(pathfs.NewDefaultFileSystem())}, &pathfs.PathNodeFsOptions{Debug: false})
	server, _, err := nodefs.MountRoot(flag.Arg(0), fs.Root(), &nodefs.Options{Debug: false})
	if err != nil {
		log.Fatal("Mount failed: %v\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			err := server.Unmount()
			if err != nil {
				fmt.Println("Error unmounting", err)
			}
			return
		}
	}()
	server.Serve()
}
