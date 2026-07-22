package segment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// WriteJSONL writes segments one-per-line, creating parent directories.
func WriteJSONL(path string, segs []Segment) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	for _, s := range segs {
		if err := enc.Encode(s); err != nil {
			return err
		}
	}
	return w.Flush()
}

// LoadJSONL reads a segments file written by WriteJSONL.
func LoadJSONL(path string) ([]Segment, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Segment
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)
	line := 0
	for sc.Scan() {
		line++
		var s Segment
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, line, err)
		}
		out = append(out, s)
	}
	return out, sc.Err()
}
