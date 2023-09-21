package yhcyasdb

import (
	"fmt"

	"yhc/utils/yasqlutil"
)

const (
	QUERY_YASDB_PARAMETER_BY_NAME = "select name,value from v$parameter where name='%s'"
)

const (
	LISTEN_ADDR = "LISTEN_ADDR"
)

type VParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func QueryParameter(tx *yasqlutil.Yasql, item string) (string, error) {
	tmp := &yasqlutil.SelectRaw{
		RawSql: fmt.Sprintf(QUERY_YASDB_PARAMETER_BY_NAME, item),
	}
	pv := make([]*VParameter, 0)
	err := tx.SelectRaw(tmp).Find(&pv).Error()
	if err != nil {
		return "", err
	}
	if len(pv) == 0 {
		return "", yasqlutil.ErrRecordNotFound
	}
	return pv[0].Value, nil
}
