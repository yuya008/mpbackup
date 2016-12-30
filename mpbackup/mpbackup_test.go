package mpbackup

import (
	"testing"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
)

func TestMpbackup(t *testing.T) {
	db, err := xorm.NewEngine("mysql", "root:12345678@tcp(127.0.0.1:3306)/mepai_v2?charset=utf8mb4")
	if err != nil {
		t.Fatal(err.Error())
	}
	db.Ping()
	db.SetColumnMapper(core.GonicMapper{})
	users := make([]MepaiUser, 0)
	
	err = db.Limit(1000, (1000000 - 1) * 1000).Asc("id").Find(&users)
	t.Log(err)
	t.Log(users)
}
