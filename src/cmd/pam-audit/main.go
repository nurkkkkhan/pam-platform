// Command pam-audit verifies an exported audit log (JSON Lines) offline.
//
//	pam-audit verify -f audit.jsonl [-anchor <head-hash>]
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nurkkkkhan/pam-platform/internal/audit"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "verify" {
		fmt.Fprintln(os.Stderr, "usage: pam-audit verify -f audit.jsonl [-anchor HASH]")
		os.Exit(2)
	}
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	file := fs.String("f", "", "audit export in JSON Lines format")
	anchor := fs.String("anchor", "", "externally anchored head hash (optional)")
	_ = fs.Parse(os.Args[2:])

	f, err := os.Open(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer f.Close()

	var events []audit.Event
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	for sc.Scan() {
		var e audit.Event
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			fmt.Fprintf(os.Stderr, "line %d: %v\n", len(events)+1, err)
			os.Exit(1)
		}
		events = append(events, e)
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read audit export:", err)
		os.Exit(1)
	}
	if err := audit.Verify(events, *anchor); err != nil {
		fmt.Println("FAIL:", err)
		os.Exit(1)
	}
	head := audit.GenesisHash
	if len(events) > 0 {
		head = events[len(events)-1].Hash
	}
	if *anchor == "" {
		fmt.Printf("OK: %d events, chain internally consistent; no external anchor supplied. Head hash %s\n", len(events), head)
		return
	}
	fmt.Printf("OK: %d events, chain matches external anchor. Head hash %s\n", len(events), head)
}
