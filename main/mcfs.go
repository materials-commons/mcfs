package main

import (
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/materials-commons/mcfs"
	"os"
	"path/filepath"
)

func main() {
	home := os.Getenv("HOME")
	mountPoint := filepath.Join(home, "fuse", "materialscommons")

	fuse.Unmount(mountPoint)
	conn, err := fuse.Mount(
		mountPoint,
		fuse.FSName("mcfs"),
		fuse.Subtype("mcfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("mcfs"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	err = fs.Serve(conn, mcfs.FS{})

	<-conn.Ready
	if conn.MountError != nil {
		log.Fatal(conn.MountError)
	}
}
