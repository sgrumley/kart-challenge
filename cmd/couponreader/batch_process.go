package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (u *Uploader) insertBatch(ctx context.Context, batch []string, fileNum, batchNum int) error {
	if len(batch) == 0 {
		return nil
	}

	start := time.Now()

	rows := make([][]any, len(batch))
	for i, id := range batch {
		newID := fmt.Sprintf("%s-%d", id, fileNum)
		rows[i] = []any{newID}
	}

	conn, err := u.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Conn().CopyFrom(
		ctx,
		pgx.Identifier{"coupons"},
		[]string{"id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	duration := time.Since(start)
	rowsPerSec := float64(len(batch)) / duration.Seconds()

	fmt.Printf("Batch %d: %d rows in %v (%.0f rows/sec)\n",
		batchNum, len(batch), duration, rowsPerSec)

	return nil
}
