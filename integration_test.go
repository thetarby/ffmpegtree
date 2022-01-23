package ffmpegtree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

func randFile() string {
	str := uuid.New().String()
	return fmt.Sprintf("./test_assets/tmp/%v.mp4", str)
}

func runFFmpeg(cmd []string) (exitCode int, err error) {
	cmd = append([]string{"/ffmpeg_binary/ffmpeg"}, cmd...)
	buf := bytes.Buffer{}
	return ffmpeg.Exec(cmd, dockertest.ExecOptions{
		StdIn:  nil,
		StdOut: &buf,
		StdErr: &buf,
		TTY:    false,
	})
}

type resp struct {
	AvgSim float32 `json:"avg_sim"`
	Error  string  `json:"error"`
}

func cmp(p1, p2 string) (float32, error) {
	buf := bytes.Buffer{}
	_, err := python.Exec([]string{
		"python",
		"./test_assets/compare.py", "-v", p1, p2,
	}, dockertest.ExecOptions{
		StdIn:  nil,
		StdOut: &buf,
		StdErr: &buf,
		TTY:    false,
	})
	if err != nil {
		return 0, err
	}

	m := resp{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		return 0, errors.New(buf.String())
	}
	if m.Error != ""{
		return 0, errors.New(m.Error)
	}

	return m.AvgSim, err
}

var python *dockertest.Resource
var ffmpeg *dockertest.Resource

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	wd, _ := os.Getwd()

	// pulls an image, creates a container based on it and runs it
	ffmpegContainer, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		Dockerfile: "./test_assets/Dockerfile",
		ContextDir: ".",
	}, &dockertest.RunOptions{
		Name:       "ffmpeg",
		Mounts:     []string{fmt.Sprintf("%v:/outs", wd)},
		WorkingDir: "/outs",
		Tty:        true,
	})
	if err != nil {
		fmt.Printf("Could not start resource: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	pythonContainer, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		Dockerfile: "./test_assets/Dockerfile.python",
		ContextDir: ".",
	}, &dockertest.RunOptions{
		Name:       "python",
		Mounts:     []string{fmt.Sprintf("%v:/wd", wd)},
		WorkingDir: "/wd",
		Tty:        true,
	})
	if err != nil {
		fmt.Printf("Could not start resource: %s", err)
	}

	python, ffmpeg = pythonContainer, ffmpegContainer
	c := m.Run()

	//You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(ffmpegContainer); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	if err := pool.Purge(pythonContainer); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(c)
}

func TestCurvesFilter(t *testing.T) {
	videoFile := randFile()

	i := NewSelectStreamNode(NewInputNode("./test_assets/test-vid.mp4", nil, nil), 0)
	c := NewCurvesFilter(i, "vintage")
	res := NewScaleFilterNode(c, 100, 100, true)

	args := Select(res, videoFile, nil)
	println(strings.Join(args, " "))

	_, err := runFFmpeg(args)
	require.NoError(t, err)

	sim, err := cmp(videoFile, "test_assets/results/curves-applied.mp4")
	require.NoError(t, err)

	require.GreaterOrEqual(t, sim, float32(0.99))
}

func TestResultsAreNotSimilar(t *testing.T) {
	videoFile := randFile()

	i := NewSelectStreamNode(NewInputNode("./test_assets/test-vid.mp4", nil, nil), 0)
	c := NewCurvesFilter(i, "vintage")
	d := NewCurvesFilter(c, "lighter")
	res := NewScaleFilterNode(d, 100, 100, true)

	args := Select(res, videoFile, nil)
	
	runFFmpeg(args)
	sim, err := cmp(videoFile, "./test_assets/results/curves-applied.mp4")
	require.NoError(t, err)
	require.Less(t, sim, float32(0.99))
}

func TestDrawbox(t *testing.T) {
	videoFile := randFile()

	i := NewSelectStreamNode(NewInputNode("./test_assets/test-vid.mp4", nil, nil), 0)
	c := NewCurvesFilter(i, "vintage")
	res := NewScaleFilterNode(c, 100, 100, true)
	draw := NewDrawBoxFilter(res, 0, 0, 30, 30, "#ff0000", "fill")

	args := Select(draw, videoFile, nil)
	println(strings.Join(args, " "))

	_, err := runFFmpeg(args)
	require.NoError(t, err)

	sim, err := cmp(videoFile, "test_assets/results/draw-box.mp4")
	require.NoError(t, err)

	require.GreaterOrEqual(t, sim, float32(0.99))
}

func TestTwoInputsComplexTree(t *testing.T) {
	videoFile := randFile()

	i := NewSelectStreamNode(NewInputNode("./test_assets/test-vid.mp4", nil, nil), 0)
	s := NewScaleFilterNode(i, 200, -2, true)
	b := NewBoxBlurFilter(s, "min(w\\,h)/5", "min(cw\\,ch)/5", 1)
	sb := NewScaleFilterNode(b, 200, 250, true)

	ov := NewOverlayIntoMiddleFilterNode(sb, s)
	i2 := NewSelectStreamNode(NewInputNode("./test_assets/test-vid-2.mp4", nil, nil), 0)
	s2 := NewScaleFilterNode(i2, 50, 50, true)
	NewRotateFilter(s2, "PI")
	var res INode = NewOverlayFilterNode(ov, s2, "0", "0")                                             // top left
	res = NewOverlayFilterNode(res, s2, "main_w-overlay_w", "0")                                       // top right
	res = NewOverlayFilterNode(res, NewRotateFilter(s2, "PI"), "0", "main_h-overlay_h")                // bottom left
	res = NewOverlayFilterNode(res, NewRotateFilter(s2, "PI"), "main_w-overlay_w", "main_h-overlay_h") // bottom right

	args := Select(res, videoFile, []string{"-shortest"})
	println(strings.Join(args, " "))

	_, err := runFFmpeg(args)
	require.NoError(t, err)

	sim, err := cmp(videoFile, "test_assets/results/two-inputs-complex-tree.mp4")
	require.NoError(t, err)

	require.GreaterOrEqual(t, sim, float32(0.99))
}
