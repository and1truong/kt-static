package activities

import (
	"context"
	"fmt"
	"os"
	"path"
)

var baseDir = ""

// SetBaseDir set base directory
func SetBaseDir(value string) {
	baseDir = value
}

func GetBaseDir() string {
	return baseDir
}

func FileWritingActivity(ctx context.Context, filePath string, content []byte, perm os.FileMode) error {
	writePath := path.Join(baseDir, filePath)
	
	fmt.Println("os.WriteFile", writePath, string(content), perm)
	
	return os.WriteFile(writePath, content, perm)
}
