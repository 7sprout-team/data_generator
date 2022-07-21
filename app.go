package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// App struct
type App struct {
	ctx      context.Context
	progress float64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) Progress() float64 {
	return a.progress
}

// Greet returns a greeting for the given name
func (a *App) Generate(INTER_SIZE int, PARTY_POOL []int) {
	if len(PARTY_POOL)%2 != 0 || len(PARTY_POOL) == 0 {
		log.Fatal("Party 必须是偶数个，第一个数代表行数，第二个表示特征数，2个一组")
	}
	type G struct {
		Name string
		Row  int
		Col  int
	}
	var parties []G = []G{}

	var REAL_SIZE int64 = 0
	for i := 0; i < len(PARTY_POOL); i += 2 {
		row := PARTY_POOL[i]
		if row < INTER_SIZE {
			log.Fatalf("%d小于交集%d", row, INTER_SIZE)
		}
		col := PARTY_POOL[i+1]
		REAL_SIZE += int64(row)

		parties = append(parties, G{Name: fmt.Sprintf("参与方%d 行%d_特征%d.csv", i/2+1, row, col), Row: row, Col: col})
	}

	inter, _ := os.Create(fmt.Sprintf("交集 %d.csv", INTER_SIZE))
	defer inter.Close()
	interWriter := bufio.NewWriterSize(inter, 1024*1024)
	defer interWriter.Flush()
	total := float64(REAL_SIZE) + float64(INTER_SIZE)
	current := float64(0)
	end := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-end:
				return
			default:
				a.progress = current / total
			}
		}
	}()

	var i int64 = 0
	//交集就是0~INTER_SIZE
	for ; i < int64(INTER_SIZE); i++ {
		interWriter.WriteString(fmt.Sprintf("%d\n", i))
		current++
	}
	for pid, party := range parties {
		start := i
		f, _ := os.Create(party.Name)
		writer := bufio.NewWriterSize(f, 1024*1024)
		header := []string{"id"}
		hasY := pid == 0
		if hasY {
			header = append(header, "y")
		}

		for headerCount := 1; headerCount <= party.Col; headerCount++ {
			header = append(header, fmt.Sprintf("x_%d_%d", pid, headerCount))
		}

		writer.WriteString(fmt.Sprintf("%s\n", strings.Join(header, ",")))
		cache := []string{}
		for l := 0; l < 1000; l++ {
			inner := []string{}
			if hasY {
				inner = append(inner, strconv.Itoa(rand.Intn(2)))
			}
			for k := 1; k <= party.Col; k++ {
				inner = append(inner, strconv.Itoa(rand.Intn(9)))
			}
			cache = append(cache, strings.Join(inner, ","))
		}
		for ; i < start+int64(party.Row-INTER_SIZE); i++ {
			writer.WriteString(fmt.Sprintf("%d,%s\n", i, cache[i%1000]))
			current++
		}
		for j := 0; j < INTER_SIZE; j++ {
			writer.WriteString(fmt.Sprintf("%d,%s\n", j, cache[j%1000]))
			current++
		}

		writer.Flush()
		f.Close()

	}

	a.progress = 1
	end <- struct{}{}

}
