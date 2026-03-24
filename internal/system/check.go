package system

import (
	"errors"
	"fmt"
	"os"
)

func CheckSys(syst string) error {
	switch syst {
	case "windows":
		return fmt.Errorf("Windows is not supported. Please use WSL instead.")
	case "macos":
		return fmt.Errorf("macOS is not supported. Please use Linux instead.")
	case "linux":
		return nil
	default:
		return fmt.Errorf("system %s is not supported", syst)
	}
}
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("The path %s is not a directory", path)
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	// 其他错误（如权限不足、路径无效等）
	return false, err
}
