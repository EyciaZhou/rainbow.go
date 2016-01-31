package PassThru

import (
	"io"
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"errors"
	"time"
)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
	io.Reader
	total    int64 // Total # of bytes transferred
	length   int64 // Expected length
	progress float64
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		fmt.Fprintf(os.Stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
		fmt.Fprintf(os.Stderr, "%s/%s", ByteSize(pt.total), ByteSize(pt.length))
	}
	if n <= 0 || pt.length == pt.total {
		fmt.Fprintf(os.Stderr, "\n")
	}

	return n, err
}

var (
	client = &http.Client{
		Timeout: time.Duration(45 * time.Second),
	}
)

func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	var e error

	for try := 0; try < 5; try++ {
		resp, err := client.Do(req)
		if err != nil {
			e = err
			continue
		}
		readerpt := &PassThru{Reader: resp.Body, length: resp.ContentLength}
		bs, err := ioutil.ReadAll(readerpt)
		resp.Body.Close()
		if err == nil {
			return bs, err
		}
		e = err
	}
	return nil, errors.New("Can't get even try 5 times reason : " + e.Error())
}