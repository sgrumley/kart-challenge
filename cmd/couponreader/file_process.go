package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	dataDir    = "./coupons"
	batchSize  = 100_000
	charLimit  = 10
	maxWorkers = 12
)

func (u *Uploader) processFile(filename string, fileNum int) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open %s: %v\n", filename, err)
		return
	}
	defer f.Close()

	fmt.Printf("Processing file: %s (file number: %d)\n", filename, fileNum)
	start := time.Now()

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var totalRows int64
	batchNum := 1

	scanner := bufio.NewScanner(f)
	// small buffer since we truncate anyway
	// any valid coupon will be 8-10
	buf := make([]byte, 0, 64)
	scanner.Buffer(buf, 1024)
	batch := make([]string, 0, batchSize)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if len(line) > charLimit {
			line = line[:charLimit]
		}

		batch = append(batch, line)
		if len(batch) >= batchSize {
			c := make([]string, len(batch))
			copy(c, batch)

			sem <- struct{}{}
			wg.Add(1)
			go func(b []string, n int) {
				defer wg.Done()
				defer func() { <-sem }()
				err = u.insertBatch(context.Background(), b, fileNum, n)
				if err != nil {
					fmt.Printf("error on file %d:%s\n", fileNum, err.Error())
				}
				atomic.AddInt64(&totalRows, int64(len(b)))
			}(c, batchNum)

			batch = batch[:0]
			batchNum++
		}
	}

	if len(batch) > 0 {
		c := make([]string, len(batch))
		copy(c, batch)

		sem <- struct{}{}
		wg.Add(1)
		go func(b []string, n int) {
			defer wg.Done()
			defer func() { <-sem }()
			err = u.insertBatch(context.Background(), b, fileNum, n)
			if err != nil {
				fmt.Printf("error on leftover batch on file %d:%s\n", fileNum, err.Error())
			}
			atomic.AddInt64(&totalRows, int64(len(b)))
		}(c, batchNum)
	}

	if serr := scanner.Err(); serr != nil {
		fmt.Printf("warning: error scanning %s: %v\n", filename, serr)
	}

	wg.Wait()

	duration := time.Since(start)
	total := atomic.LoadInt64(&totalRows)
	rowsPerSec := float64(total) / duration.Seconds()

	fmt.Printf("✓ Completed %s: %d rows in %v (%.0f rows/sec)\n",
		filepath.Base(filename), total, duration, rowsPerSec)
}

func (u *Uploader) processFiles(couponFiles []string) {
	totalStart := time.Now()

	for _, filename := range couponFiles {
		fullPath := filepath.Join(dataDir, filename)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Printf("Warning: File %s does not exist, skipping...\n", fullPath)
			continue
		}

		fileNum, err := parseFileNumber(filename)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", filename, err)
			continue
		}

		u.processFile(fullPath, fileNum)
		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Printf("\n✓ All files processed in %v\n", time.Since(totalStart))
}

// func (u *Uploader) processFile(filename string, fileNum int) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to open %s: %w", filename, err)
// 	}
// 	defer file.Close()
//
// 	fmt.Printf("Processing file: %s (file number: %d)\n", filename, fileNum)
// 	start := time.Now()
//
// 	batchChan := make(chan []string, maxWorkers*2)
// 	errorChan := make(chan error, maxWorkers)
//
// 	var wg sync.WaitGroup
// 	var totalRows atomic.Int64
// 	var batchCount atomic.Int32
//
// 	// Start worker pool
// 	for i := 0; i < maxWorkers; i++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
// 			for batch := range batchChan {
// 				batchNum := batchCount.Add(1)
// 				if err := u.insertBatch(context.Background(), batch, fileNum, int(batchNum)); err != nil {
// 					errorChan <- fmt.Errorf("worker %d, batch %d: %w", workerID, batchNum, err)
// 					return
// 				}
// 				totalRows.Add(int64(len(batch)))
// 			}
// 		}(i)
// 	}
//
// 	// Read file and send batches
// 	go func() {
// 		scanner := bufio.NewScanner(file)
// 		buf := make([]byte, 0, 1024*1024)
// 		scanner.Buffer(buf, 10*1024*1024)
//
// 		batch := make([]string, 0, batchSize)
// 		for scanner.Scan() {
// 			line := strings.TrimSpace(scanner.Text())
// 			if line == "" {
// 				continue
// 			}
//
// 			batch = append(batch, line)
// 			if len(batch) >= batchSize {
// 				bcopy := make([]string, len(batch))
// 				copy(bcopy, batch)
// 				batchChan <- bcopy
// 				batch = batch[:0]
// 			}
// 		}
//
// 		if len(batch) > 0 {
// 			bcopy := make([]string, len(batch))
// 			copy(bcopy, batch)
// 			batchChan <- bcopy
// 		}
//
// 		if err := scanner.Err(); err != nil {
// 			errorChan <- fmt.Errorf("error reading file: %w", err)
// 		}
// 		close(batchChan)
// 	}()
//
// 	wg.Wait()
// 	close(errorChan)
//
// 	for err := range errorChan {
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	duration := time.Since(start)
// 	total := totalRows.Load()
// 	rowsPerSec := float64(total) / duration.Seconds()
//
// 	fmt.Printf("✓ Completed %s: %d rows in %v (%.0f rows/sec)\n",
// 		filepath.Base(filename), total, duration, rowsPerSec)
//
// 	return nil
// }

func parseFileNumber(filename string) (int, error) {
	re := regexp.MustCompile(`(\d+)$`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no integer found in filename: %s", filename)
	}
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return num, nil
}
