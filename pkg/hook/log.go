package hook

import (
	"bufio"
	"io"

	"kubevirt.io/client-go/log"
)

func logStderr(reader io.Reader, hookName string) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), 512*1024)
	for scanner.Scan() {
		log.Log.With("hook", hookName).Info(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Log.Reason(err).Error("failed to read hook logs")
	}
}
