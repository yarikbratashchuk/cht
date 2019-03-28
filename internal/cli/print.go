package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/yarikbratashchuk/cht/internal/cht"
)

func readMessage(r io.Reader) cht.Message {
	reader := bufio.NewReader(r)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	fmt.Print("\r\033[F")

	return cht.Message{
		Text: text,
	}
}

func printMessage(w io.Writer, m cht.Message) {
	buf := new(bytes.Buffer)

	buf.WriteRune('\r')
	buf.WriteString(m.String())
	buf.WriteRune('\n')

	fmt.Fprint(w, buf.String())
}
