package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func NewFileUploader(config *Config) *FileUploader {
	return &FileUploader{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (fu *FileUploader) validateFile(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file '%s' not found", filename)
		}
		return fmt.Errorf("cannot access file '%s': %w", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filename)
	}

	return nil
}

func (fu *FileUploader) UploadFile(filename string, private bool) (*Response, error) {
	if err := fu.validateFile(filename); err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	bar := pb.Full.Start64(fileInfo.Size())
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	defer bar.Finish()

	baseName := filepath.Base(filename)
	url := fmt.Sprintf("%s/%s", fu.config.BaseURL, baseName)
	if private {
		url += "?q=private"
	}

	req, err := http.NewRequest("POST", url, bar.NewProxyReader(file))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = fileInfo.Size()

	resp, err := fu.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

func PrintUsage() {
	fmt.Printf(`%sUsage:%s
  %s <filename> [options]

%sOptions:%s
  -p, --private    Upload as private file
  -h, --help       Show this help message

%sExamples:%s
  %s myfile.txt                    # Upload public file
  %s myfile.txt -p                 # Upload private file
  %s myfile.txt --private          # Upload private file (alternative)

`,
		ColorBold, ColorReset,
		os.Args[0],
		ColorBold, ColorReset,
		ColorBold, ColorReset,
		ColorYellow+os.Args[0]+ColorReset,
		ColorYellow+os.Args[0]+ColorReset,
		ColorYellow+os.Args[0]+ColorReset,
	)
}

func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, ColorRed+"Error: "+format+ColorReset+"\n", args...)
}

func PrintSuccess(response *Response, isPrivate bool) {
	fmt.Printf(`
%sUpload Successful!%s

%sFile Information:%s
  Name:     %s%s%s
  Size:     %s%s%s
  ID:       %s%s%s
  Created:  %s%s%s

%sAccess Details:%s
  Private:  %s%v%s
  URL:      %s%s%s

`,
		ColorGreen+ColorBold, ColorReset,
		ColorCyan+ColorBold, ColorReset,
		ColorWhite, response.FileName, ColorReset,
		ColorGreen, response.Metadata.Size.Formatted, ColorReset,
		ColorYellow, response.ID, ColorReset,
		ColorBlue, response.CreateAt, ColorReset,
		ColorCyan+ColorBold, ColorReset,
		ColorRed, isPrivate, ColorReset,
		ColorPurple, response.Metadata.URL, ColorReset,
	)
}

func ParseArgs(args []string) (filename string, private bool, showHelp bool, err error) {
	if len(args) < 2 {
		return "", false, false, fmt.Errorf("missing required parameter 'filename'")
	}

	filename = args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "-p", "--private":
			private = true
		case "-h", "--help":
			showHelp = true
			return
		default:
			return "", false, false, fmt.Errorf("unknown option: %s", args[i])
		}
	}

	return filename, private, false, nil
}
