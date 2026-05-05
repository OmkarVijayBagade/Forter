// Package scanner handles directory scanning and file discovery
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/OmkarVijayBagade/forter/internal/config"
)

// FileInfo represents a file with its metadata
type FileInfo struct {
	Path         string
	Name         string
	Ext          string
	Category     string
	Size         int64
	ModTime      int64
	IsSelected   bool
	Destination  string
	Status       FileStatus
}

// FileStatus represents the current status of a file
type FileStatus int

const (
	StatusPending FileStatus = iota
	StatusSelected
	StatusMoved
	StatusSkipped
	StatusError
)

func (s FileStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusSelected:
		return "selected"
	case StatusMoved:
		return "moved"
	case StatusSkipped:
		return "skipped"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
}

// Result holds the scan results
type Result struct {
	Files     []*FileInfo
	TotalSize int64
	Errors    []error
}

// Scanner handles file scanning operations
type Scanner struct {
	config     *config.Config
	recursive  bool
	maxWorkers int
}

// NewScanner creates a new scanner instance
func NewScanner(cfg *config.Config, recursive bool) *Scanner {
	return &Scanner{
		config:     cfg,
		recursive:  recursive,
		maxWorkers: 10,
	}
}

// Scan scans a directory and returns all files
func (s *Scanner) Scan(rootPath string) (*Result, error) {
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", rootPath)
	}

	result := &Result{
		Files:  make([]*FileInfo, 0),
		Errors: make([]error, 0),
	}

	// Use worker pool for concurrent processing
	filesChan := make(chan string, 100)
	resultChan := make(chan *scanWorkerResult, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.maxWorkers; i++ {
		wg.Add(1)
		go s.scanWorker(filesChan, resultChan, &wg)
	}

	// Collector goroutine
	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		defer collectWg.Done()
		for res := range resultChan {
			if res.err != nil {
				result.Errors = append(result.Errors, res.err)
			} else if res.file != nil {
				result.Files = append(result.Files, res.file)
				result.TotalSize += res.file.Size
			}
		}
	}()

	// Walk directory and send files to workers
	go func() {
		defer close(filesChan)
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("error accessing %s: %w", path, err))
				return nil // Continue walking
			}

			// Skip directories if configured
			if info.IsDir() {
				if s.config.SkipHidden && config.IsHiddenFile(path) {
					return filepath.SkipDir
				}
				// Skip non-recursive subdirectories
				if !s.recursive && path != rootPath {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip hidden files if configured
			if s.config.SkipHidden && config.IsHiddenFile(info.Name()) {
				return nil
			}

			filesChan <- path
			return nil
		})
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)
	collectWg.Wait()

	return result, nil
}

type scanWorkerResult struct {
	file *FileInfo
	err  error
}

func (s *Scanner) scanWorker(files <-chan string, results chan<- *scanWorkerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range files {
		info, err := os.Stat(path)
		if err != nil {
			results <- &scanWorkerResult{err: fmt.Errorf("failed to stat %s: %w", path, err)}
			continue
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		category := s.config.GetCategoryForExtension(ext)

		fileInfo := &FileInfo{
			Path:     path,
			Name:     info.Name(),
			Ext:      ext,
			Category: category,
			Size:     info.Size(),
			ModTime:  info.ModTime().Unix(),
			Status:   StatusPending,
		}

		results <- &scanWorkerResult{file: fileInfo}
	}
}

// FilterByCategory filters files by category
func FilterByCategory(files []*FileInfo, category string) []*FileInfo {
	if category == "" || category == "All" {
		return files
	}

	filtered := make([]*FileInfo, 0)
	for _, f := range files {
		if f.Category == category {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// FilterByExtension filters files by extension
func FilterByExtension(files []*FileInfo, ext string) []*FileInfo {
	if ext == "" {
		return files
	}

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	filtered := make([]*FileInfo, 0)
	for _, f := range files {
		if strings.TrimPrefix(strings.ToLower(f.Ext), ".") == ext {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
