package scripts

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func StartRendering(fileName string) error {
	fmt.Println("Rendering file: ", fileName)

	cmd := exec.Command("./startRender.sh", fileName)
	cmd.Dir = "/home/biggiecheese/code/DistributedRendering/server/scripts"
	cmd.Env = os.Environ() // inherit env so blender is found

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	fmt.Println("Rendering process completed.")
	return nil
}
