package interceptor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"os"
)

type Interceptor struct {
	Prefix        string
	HasStdOutput  bool
	Base64Encoded bool
	Output        io.Writer
}

// Write implements io.Writer.
func (i *Interceptor) Write(p []byte) (n int, err error) {
	if i.HasStdOutput {
		go func() {
			var data []byte
			if i.Base64Encoded {
				data = encodeB64(p)
			} else {
				data = p
			}
			out := []byte(fmt.Sprintf("%s\n%s\n=====================\n", i.Prefix, data))
			os.Stdout.Write(out)
		}()
	}
	return i.Output.Write(p)
}

func encodeB64(p []byte) []byte {
	base4 := float64(len(p)) / 4
	size := int(math.Ceil(base4) * 4)
	data := make([]byte, 0, size)
	buf := bytes.NewBuffer(data)
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	defer enc.Close()
	enc.Write(p)
	return buf.Bytes()
}

// func encodeB64(p []byte) []byte {
// 	ss := float64(len(p) / 4)
// 	sizef := math.Ceil(ss)
// 	size := int(sizef * 4)
// 	data := make([]byte, 0, size)
// 	buf := bytes.NewBuffer(data)
// 	enc := base64.NewEncoder(base64.StdEncoding, buf)
// 	defer enc.Close()
// 	enc.Write(p)
// 	fmt.Println(data)
// 	return data
// }
