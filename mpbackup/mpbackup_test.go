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
	user := &MepaiUser{}
	db.Where("id=?", 10001).Get(user)
	t.Log(user)
	
	worksAppendix := &MepaiWorksAppendix{}
	db.Where("id=?", 629).Get(worksAppendix)
	t.Log(worksAppendix)
	
	activity := &MepaiActivity{}
	db.Where("id=?", 154).Get(activity)
	t.Log(activity)
}
