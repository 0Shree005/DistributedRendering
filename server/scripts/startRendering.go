package scripts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"

	"github.com/0Shree005/DistributedRendering/server/types"
)

func StartRendering(fileName string, renderJobs *sync.Map) error {
	fmt.Println("Rendering file: ", fileName)

	cmd := exec.Command("./startRender.sh", fileName)
	cmd.Dir = "/home/biggiecheese/code/DistributedRendering/server/scripts"
	cmd.Env = os.Environ() // inherit env so blender is found

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	// This regex captures the progress from Blender's "Rendered X/Y Tiles" output.
	re := regexp.MustCompile(`Rendered ([0-9]+)\/([0-9]+) Tiles`)

	// Read from stdout in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Log the output for debugging
			match := re.FindStringSubmatch(line)
			if len(match) == 3 {
				renderedTiles, err1 := strconv.ParseFloat(match[1], 64)
				totalTiles, err2 := strconv.ParseFloat(match[2], 64)
				if err1 == nil && err2 == nil && totalTiles > 0 {
					progress := (renderedTiles / totalTiles) * 100
					renderJobs.Store(fileName, &types.JobStatus{Status: "in-progress", Progress: int(progress)})
				}
			}
		}
	}()

	go io.Copy(os.Stderr, stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	fmt.Println("Rendering process completed.")
	return nil
}
