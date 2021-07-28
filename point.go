package simpletsdb

import (
	"bytes"
	"fmt"
	"strings"
)

func (p *InsertPointRequest) toLineProtocol() []byte {
	tagsBuilder := &strings.Builder{}
	buf := &bytes.Buffer{}
	var i int
	for k, v := range p.Tags {
		tagsBuilder.WriteString(fmt.Sprintf("%s=%s", k, v))
		if i+1 < len(p.Tags) {
			tagsBuilder.WriteString(" ")
		}
	}
	buf.WriteString(p.Metric)
	buf.WriteRune(',')
	buf.WriteString(tagsBuilder.String())
	buf.WriteRune(',')
	buf.WriteString(fmt.Sprint(p.Point.Value))
	buf.WriteRune(' ')
	buf.WriteString(fmt.Sprint(p.Point.Timestamp))
	return buf.Bytes()
}
