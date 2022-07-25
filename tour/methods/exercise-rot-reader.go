package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot13 *rot13Reader) Read(b []byte) (int, error) {
	n, err := rot13.r.Read(b)
	for i := 0; i < n; i++ {
		if (b[i] >= 65) && (b[i] < 65+13) {
			b[i] = b[i] + 13
		} else if (b[i] >= 97) && (b[i] < 97+13) {
			b[i] = b[i] + 13
		} else if (b[i] >= 77) && (b[i] < 77+13) {
			b[i] = b[i] - 13
		} else if (b[i] >= 110) && (b[i] < 110+13) {
			b[i] = b[i] - 13
		}
	}

	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
