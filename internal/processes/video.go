package processes

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	BG_VIDEO_WITH_AUDIO = "background-video-with-audio.mp4"
	VIDEO_LENGTH_CROP   = "cropped-length-video.mp4"
	VIDEO_FOLDER        = "./resources/videos"
)

type VideoGenerationProcess struct {
	next Process
}

func NewVideoGenerationProcess() *VideoGenerationProcess {
	return &VideoGenerationProcess{}
}

func (p *VideoGenerationProcess) Execute(request interface{}) (interface{}, error) {
	context, ok := request.(*GenerationContext)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	finalVideo, err := p.generateVideo(context)
	if err != nil {
		return nil, err
	}

	context.GeneratedVideo = finalVideo
	if p.next != nil {
		return p.next.Execute(context)
	}

	return context.GeneratedVideo, nil
}

func (p *VideoGenerationProcess) SetNext(handler Process) {
	p.next = handler
}

func (p *VideoGenerationProcess) generateVideo(context *GenerationContext) (string, error) {
	backgroundVideo, err := getRandomVideo(VIDEO_FOLDER)
	if err != nil {
		return "", err
	}

	bgVideoWithAudio, err := overwriteVideoAudio(context.TempDir, backgroundVideo, context.SpeechFile)
	if err != nil {
		return "", fmt.Errorf("error executing overwriteVideoAudio: %w", err)
	}

	croppedVideo, err := cropVideoLength(context.TempDir, bgVideoWithAudio, context.SpeechFile)
	if err != nil {
		return "", fmt.Errorf("error executing cropVideoLength: %w", err)
	}

	finalPath, err := runAutocap(croppedVideo)
	if err != nil {
		return "", fmt.Errorf("error executing runAutocap: %w", err)
	}

	return finalPath, nil
}

func getRandomVideo(folder string) (string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	var videos []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mp4") {
			videos = append(videos, filepath.Join(folder, file.Name()))
		}
	}

	if len(videos) == 0 {
		return "", fmt.Errorf("no video files found in folder %s", folder)
	}

	rand.Seed(time.Now().Unix())
	return videos[rand.Intn(len(videos))], nil
}

func overwriteVideoAudio(tempDir, videoPath, audioPath string) (string, error) {
	videoWithAudio := filepath.Join(tempDir, BG_VIDEO_WITH_AUDIO)
	args := []string{
		"-i", videoPath,
		"-i", audioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-y",
		videoWithAudio,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", fmt.Errorf("error overwriting video audio: %w", err)
	}

	return videoWithAudio, nil
}

func cropVideoLength(tempDir, videoPath, audioPath string) (string, error) {
	audioDuration, err := getAudioDuration(audioPath)
	if err != nil {
		return "", fmt.Errorf("error getting audio duration: %w", err)
	}

	croppedVideo := filepath.Join(tempDir, VIDEO_LENGTH_CROP)
	args := []string{
		"-i", videoPath,
		"-t", audioDuration,
		"-c", "copy",
		"-y",
		croppedVideo,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", fmt.Errorf("error cropping video length: %w", err)
	}

	return croppedVideo, nil
}

func getAudioDuration(audioPath string) (string, error) {
	args := []string{
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		audioPath,
	}

	output, err := executeCommandOutput("ffprobe", args)
	if err != nil {
		return "", fmt.Errorf("error getting audio duration: %w", err)
	}

	return strings.TrimSpace(output), nil
}

func runAutocap(videoPath string) (string, error) {
	// Create generated directory if it doesn't exist
	generatedDir := "./generated"
	if err := os.MkdirAll(generatedDir, 0755); err != nil {
		return "", fmt.Errorf("error creating generated directory: %w", err)
	}

	finalVideoName := fmt.Sprintf("%s.mp4", generateRandomWord(7))
	finalVideoPath := filepath.Join(generatedDir, finalVideoName)

	// Use absolute path for Python script execution
	cmd := exec.Command("python3", "./scripts/captions.py", "attach", videoPath, finalVideoName)
	cmd.Dir = "/app" // Set working directory explicitly
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Python script error: %s", string(output))
		log.Printf("Command: %v", cmd.Args)
		return "", fmt.Errorf("error running captions script: %w", err)
	}

	log.Printf("Captions script output: %s", string(output))
	return finalVideoPath, nil
}

func executeCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command %s: %v, output: %s", name, err, string(output))
		return err
	}
	return nil
}

func executeCommandOutput(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command %s: %v", name, err)
		return "", err
	}
	return string(output), nil
}

func generateRandomWord(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	word := make([]byte, length)
	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}

	return string(word)
}
