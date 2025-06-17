package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)

	return
}

func (d Duration) MarshalJSON() (b []byte, err error) {
	return fmt.Appendf(nil, `"%s"`, d.String()), nil
}

func (d Duration) Value() (driver.Value, error) {
	return d.Duration.String(), nil
}

func (d *Duration) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		d.Duration, err = time.ParseDuration(v)
	default:
		return fmt.Errorf("cannot sql.Scan() Duration from: %#v", v)
	}
	return nil
}
