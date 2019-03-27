package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/yarikbratashchuk/cht/internal/cht"
)

func readMessage(r io.Reader) cht.Message {
	reader := bufio.NewReader(r)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	return cht.Message{
		Text: text,
	}
}

func printMessage(w io.Writer, m cht.Message) {
	buf := new(bytes.Buffer)

	buf.WriteString(color.New(color.FgCyan, color.Bold).Sprintf(m.Author))
	buf.WriteString(": ")
	buf.WriteString(m.Text)
	buf.WriteRune('\n')

	fmt.Fprint(w, buf.String())
}
