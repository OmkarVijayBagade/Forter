// Package organizer handles file organization and movement operations
package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OmkarVijayBagade/forter/internal/config"
	"github.com/OmkarVijayBagade/forter/internal/scanner"
)

// Operation represents a file organization operation
type Operation struct {
	Source      string
	Destination string
	Category    string
	Subfolder   string
	Status      OpStatus
	Error       error
	Timestamp   time.Time
}

// OpStatus represents operation status
type OpStatus int

const (
	OpPending OpStatus = iota
	OpSuccess
	OpFailed
	OpSkipped
)

func (s OpStatus) String() string {
	switch s {
	case OpPending:
		return "pending"
	case OpSuccess:
		return "success"
	case OpFailed:
		return "failed"
	case OpSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// Organizer handles file organization
type Organizer struct {
	config      *config.Config
	dryRun      bool
	rootPath    string
	operations  []Operation
	mu          sync.RWMutex
	onProgress  func(Operation)
}

// NewOrganizer creates a new organizer instance
func NewOrganizer(cfg *config.Config, dryRun bool) *Organizer {
	return &Organizer{
		config:     cfg,
		dryRun:     dryRun,
		operations: make([]Operation, 0),
	}
}

// SetProgressCallback sets a callback for operation progress updates
func (o *Organizer) SetProgressCallback(callback func(Operation)) {
	o.onProgress = callback
}

// PreviewDestination generates the destination path for a file
func (o *Organizer) PreviewDestination(file *scanner.FileInfo, rootPath string) string {
	ext := strings.TrimPrefix(strings.ToLower(file.Ext), ".")
	category := file.Category

	// Create subfolder based on extension
	var subfolder string
	if ext != "" {
		subfolder = strings.ToUpper(ext)
	} else {
		subfolder = "misc"
	}

	destDir := filepath.Join(rootPath, category, subfolder)
	destPath := filepath.Join(destDir, file.Name)

	// Handle duplicates
	destPath = o.getUniquePath(destPath)

	return destPath
}

// getUniquePath generates a unique path by appending (1), (2), etc.
func (o *Organizer) getUniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	counter := 1
	for {
		newName := fmt.Sprintf("%s (%d)%s", name, counter, ext)
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
		// Safety limit
		if counter > 9999 {
			return path + "." + fmt.Sprintf("%d", time.Now().Unix())
		}
	}
}

// PreviewOperations generates a list of operations without executing them
func (o *Organizer) PreviewOperations(files []*scanner.FileInfo, rootPath string) []Operation {
	ops := make([]Operation, 0, len(files))

	for _, file := range files {
		if !file.IsSelected || file.Status == scanner.StatusMoved {
			continue
		}

		dest := o.PreviewDestination(file, rootPath)
	ext := strings.TrimPrefix(strings.ToLower(file.Ext), ".")

		op := Operation{
			Source:      file.Path,
			Destination: dest,
			Category:    file.Category,
			Subfolder:   strings.ToUpper(ext),
			Status:      OpPending,
			Timestamp:   time.Now(),
		}
		ops = append(ops, op)
	}

	return ops
}

// Execute performs the file organization
func (o *Organizer) Execute(files []*scanner.FileInfo, rootPath string) error {
	o.mu.Lock()
	o.operations = make([]Operation, 0)
	o.rootPath = rootPath
	o.mu.Unlock()

	// Create channels for parallel processing
	fileChan := make(chan *scanner.FileInfo, 100)
	resultChan := make(chan Operation, 100)

	var wg sync.WaitGroup
	maxWorkers := 5 // Limit concurrent file operations

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go o.worker(fileChan, resultChan, rootPath, &wg)
	}

	// Collector
	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		defer collectWg.Done()
		for op := range resultChan {
			o.mu.Lock()
			o.operations = append(o.operations, op)
			o.mu.Unlock()

			if o.onProgress != nil {
				o.onProgress(op)
			}
		}
	}()

	// Send files to workers
	go func() {
		defer close(fileChan)
		for _, file := range files {
			if file.IsSelected && file.Status != scanner.StatusMoved {
				fileChan <- file
			}
		}
	}()

	// Wait for completion
	wg.Wait()
	close(resultChan)
	collectWg.Wait()

	return nil
}

func (o *Organizer) worker(files <-chan *scanner.FileInfo, results chan<- Operation, rootPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		op := o.processFile(file, rootPath)
		results <- op
	}
}

func (o *Organizer) processFile(file *scanner.FileInfo, rootPath string) Operation {
	ext := strings.TrimPrefix(strings.ToLower(file.Ext), ".")
	dest := o.PreviewDestination(file, rootPath)
	destDir := filepath.Dir(dest)

	op := Operation{
		Source:      file.Path,
		Destination: dest,
		Category:    file.Category,
		Subfolder:   strings.ToUpper(ext),
		Status:      OpPending,
		Timestamp:   time.Now(),
	}

	if o.dryRun {
		op.Status = OpSkipped
		return op
	}

	// Check if source file exists and is accessible
	srcInfo, err := os.Stat(file.Path)
	if err != nil {
		op.Status = OpFailed
		op.Error = fmt.Errorf("cannot access source file: %w", err)
		return op
	}

	// Check if file is locked/busy (basic check)
	if !srcInfo.Mode().IsRegular() {
		op.Status = OpSkipped
		op.Error = fmt.Errorf("not a regular file")
		return op
	}

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		op.Status = OpFailed
		op.Error = fmt.Errorf("failed to create directory: %w", err)
		return op
	}

	// Move the file
	if err := os.Rename(file.Path, dest); err != nil {
		// Try copy + delete as fallback (for cross-device moves)
		if err := o.copyAndDelete(file.Path, dest); err != nil {
			op.Status = OpFailed
			op.Error = fmt.Errorf("failed to move file: %w", err)
			return op
		}
	}

	op.Status = OpSuccess
	file.Status = scanner.StatusMoved
	return op
}

// copyAndDelete copies a file then deletes the source (for cross-device moves)
func (o *Organizer) copyAndDelete(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := sourceFile.Read(buf)
		if n > 0 {
			if _, writeErr := destFile.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			break
		}
	}

	// Sync to ensure data is written
	if err := destFile.Sync(); err != nil {
		return err
	}

	// Close files before deleting
	sourceFile.Close()
	destFile.Close()

	// Delete source
	return os.Remove(src)
}

// GetOperations returns the list of operations
func (o *Organizer) GetOperations() []Operation {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]Operation, len(o.operations))
	copy(result, o.operations)
	return result
}

// GetStats returns statistics about operations
func (o *Organizer) GetStats() (total, success, failed, skipped int) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, op := range o.operations {
		total++
		switch op.Status {
		case OpSuccess:
			success++
		case OpFailed:
			failed++
		case OpSkipped:
			skipped++
		}
	}
	return
}

// Undo attempts to undo the last set of operations
func (o *Organizer) Undo() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Process in reverse order
	for i := len(o.operations) - 1; i >= 0; i-- {
		op := o.operations[i]
		if op.Status != OpSuccess {
			continue
		}

		// Check if destination still exists
		if _, err := os.Stat(op.Destination); os.IsNotExist(err) {
			continue // Already gone, skip
		}

		// Ensure source directory exists
		srcDir := filepath.Dir(op.Source)
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			return fmt.Errorf("failed to recreate source directory: %w", err)
		}

		// Move back
		if err := os.Rename(op.Destination, op.Source); err != nil {
			// Try copy + delete fallback
			if err := o.copyAndDelete(op.Destination, op.Source); err != nil {
				return fmt.Errorf("failed to undo move for %s: %w", op.Destination, err)
			}
		}
	}

	// Clear operations after undo
	o.operations = make([]Operation, 0)
	return nil
}

// Atomic counter helpers for statistics
type Stats struct {
	Total   int32
	Success int32
	Failed  int32
	Skipped int32
}

func (s *Stats) IncrementTotal() {
	atomic.AddInt32(&s.Total, 1)
}

func (s *Stats) IncrementSuccess() {
	atomic.AddInt32(&s.Success, 1)
}

func (s *Stats) IncrementFailed() {
	atomic.AddInt32(&s.Failed, 1)
}

func (s *Stats) IncrementSkipped() {
	atomic.AddInt32(&s.Skipped, 1)
}
