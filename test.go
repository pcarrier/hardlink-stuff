package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var src = flag.String("src", "/usr", "source path")
var dst = flag.String("dst", "/usr.copy", "destination path")
var maxDirs = flag.Uint("maxdirs", 256, "maximum number of open directories")

type sem struct{}

func copyEntries(src, dst string, sc chan sem) error {
	sc <- sem{}
	f, err := os.Open(src)

	if err != nil {
		<-sc
		return err
	}

	infos, err := f.Readdir(-1)
	f.Close()
	<-sc

	if err != nil {
		return err
	}

	errs := make(chan error, len(infos))
	dircount := 0

	for _, info := range infos {
		csrc := path.Join(src, info.Name())
		cdst := path.Join(dst, info.Name())

		// paralelize dirs
		if info.IsDir() {
			dircount++
			go func(info os.FileInfo, csrc, cdst string) {
				errs <- copy(csrc, cdst, info, sc)
			}(info, csrc, cdst)

		} else {
			copy(csrc, cdst, info, sc)
		}
	}

	for i := 0; i < dircount; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}

	return nil
}

func copy(src string, dst string, fi os.FileInfo, sc chan sem) error {
	mod := fi.Mode()

	switch {
	case mod&os.ModeSymlink != 0:
		if tgt, err := os.Readlink(src); err != nil {
			return err
		} else {
			if err := os.Symlink(tgt, dst); err != nil {
				return err
			}
		}
	case mod&os.ModeDir != 0:
		if err := os.Mkdir(dst, mod); err != nil {
			return err
		}
		if err := copyEntries(src, dst, sc); err != nil {
			return err
		}
	default:
		if err := os.Link(src, dst); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if err := os.Mkdir(*dst, 0755); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	sc := make(chan sem, int(*maxDirs))

	if err := copyEntries(*src, *dst, sc); err != nil {
		log.Fatal(err)
	}
}
