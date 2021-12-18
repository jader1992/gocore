package formatter

import (
	"bytes"
	"encoding/json"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/pkg/errors"
	"time"
)

// JsonFormatter json格式化
func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{}) // buffer写缓冲区
	fields["msg"] = msg
	fields["level"] = level
	fields["timestamp"] = t.Format(time.RFC3339)
	c, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.Wrap(err, "json format error")
	}

	bf.Write(c)
	return bf.Bytes(), nil
}
