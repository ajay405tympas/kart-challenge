package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbURL     = "postgres://user:12345@localhost:5432/couponsdb"
	batchSize = 50000
)

var files = []string{
	"couponbase1.txt",
	"couponbase2.txt",
	"couponbase3.txt",
}

func main() {
	start := time.Now()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			loadFile(ctx, pool, f)
		}(file)
	}

	wg.Wait()

	fmt.Println("✅ All files loaded")
	fmt.Println("⏱ Time:", time.Since(start))
}

func loadFile(ctx context.Context, pool *pgxpool.Pool, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	batch := make([]string, 0, batchSize)
	total := 0

	for scanner.Scan() {
		batch = append(batch, scanner.Text())

		if len(batch) >= batchSize {
			copyBatch(ctx, pool, batch)
			total += len(batch)
			fmt.Printf("[%s] Inserted: %d\n", filePath, total)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		copyBatch(ctx, pool, batch)
		total += len(batch)
	}

	fmt.Printf("[%s] Done. Total: %d\n", filePath, total)
}

func copyBatch(ctx context.Context, pool *pgxpool.Pool, data []string) {
	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"coupons"},
		[]string{"code"},
		pgx.CopyFromSlice(len(data), func(i int) ([]interface{}, error) {
			return []interface{}{data[i]}, nil
		}),
	)

	if err != nil {
		log.Fatal("Copy failed:", err)
	}
}
