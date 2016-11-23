package mcfs

import "bazil.org/fuse/fs"

type FS struct{}

func (_ FS) Root() (fs.Node, error) {
	return Dir{}, nil
}
