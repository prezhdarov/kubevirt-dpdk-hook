package hook

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/rand"
	"kubevirt.io/client-go/log"
	"kubevirt.io/kubevirt/pkg/hooks"
)

func getSocketPath() (string, error) {
	if _, err := os.Stat(hooks.HookSocketsSharedDirectory); err != nil {
		return "", fmt.Errorf("Failed dir %s due %s", hooks.HookSocketsSharedDirectory, err.Error())
	}

	// In case there are multiple shims being used, append random string and try a few times
	for i := 0; i < 10; i++ {
		socketName := fmt.Sprintf("shim-%s.sock", rand.String(4))
		socketPath := filepath.Join(hooks.HookSocketsSharedDirectory, socketName)
		if _, err := os.Stat(socketPath); !errors.Is(err, os.ErrNotExist) {
			log.Log.Infof("Failed socket %s due %s", socketName, err.Error())
			continue
		}
		return socketPath, nil
	}

	return "", fmt.Errorf("Failed generate socket path")
}
