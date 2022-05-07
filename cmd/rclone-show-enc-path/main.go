package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/rclone/rclone/backend/crypt"
	_ "github.com/rclone/rclone/backend/drive"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
)

func main() {
	usr, _ := user.Current()
	config.SetConfigPath(filepath.Join(usr.HomeDir, ".rclone.conf"))
	configfile.Install()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if len(os.Args) != 2 {
		log.Fatal("path not specified")
	}

	remotePath := os.Args[1]
	pathComponents := strings.Split(remotePath, ":")
	if len(pathComponents) != 2 {
		log.Fatal("path must be of the following format: remote:/path/to/file/or/directory")
	}
	remote := pathComponents[0] + ":"
	path := pathComponents[1]

	f, err := fs.NewFs(ctx, remote)
	if err != nil {
		log.Fatal(err)
	}

	cryptFS, ok := f.(*crypt.Fs)
	if !ok {
		log.Fatal("not a crypt remote")
	}

	encryptedPath := cryptFS.EncryptFileName(path)
	fmt.Println(encryptedPath)
}
