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
	reProgress := regexp.MustCompile(`Rendered ([0-9]+)\/([0-9]+) Tiles`)
	// This regex captures the remaining time from Blender's output.
	reTime := regexp.MustCompile(`Remaining:([0-9:.]+)`)

	// Read from stdout in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(stdout)
		currentStatus := &types.JobStatus{Status: "in-progress", Progress: 0, RemainingTime: ""}
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Log the output for debugging

			// Check for progress
			matchProgress := reProgress.FindStringSubmatch(line)
			if len(matchProgress) == 3 {
				renderedTiles, err1 := strconv.ParseFloat(matchProgress[1], 64)
				totalTiles, err2 := strconv.ParseFloat(matchProgress[2], 64)
				if err1 == nil && err2 == nil && totalTiles > 0 {
					progress := (renderedTiles / totalTiles) * 100
					currentStatus.Progress = int(progress)
					renderJobs.Store(fileName, currentStatus)
				}
			}

			// Check for remaining time
			matchTime := reTime.FindStringSubmatch(line)
			if len(matchTime) == 2 {
				currentStatus.RemainingTime = matchTime[1]
				renderJobs.Store(fileName, currentStatus)
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
