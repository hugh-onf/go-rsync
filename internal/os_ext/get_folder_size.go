package os_ext

import (
	folderSize "github.com/markthree/go-get-folder-size/src"
)

func GetFolderSize(path string) (int64, error) {
	size := folderSize.LooseParallel(path)
	return size, nil
}
