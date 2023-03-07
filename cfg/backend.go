package cfg

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"path/filepath"
)

const (
	DBTypMem = iota
	DBTypJson
	DBTypLevelDB
	DBTypSqlite
)

type BackConfig struct {
	DBType    int8   `json:"db_type"` //0:memory db; 1:json file db; 2:leveldb; 3:sqlite db
	CurDBHome string `json:"-"`
	DBParam   string `json:"db_param"`
	DBHome    string `json:"db_home"`
}

func (bc *BackConfig) prepare(homeDir, dbCfgFilePath string) error {
	if !filepath.IsAbs(dbCfgFilePath) {
		dbCfgFilePath = filepath.Join(homeDir, string(filepath.Separator), dbCfgFilePath)
	}
	if err := util.ReadJsonFile(dbCfgFilePath, bc); err != nil {
		fmt.Println("parse backend store config failed:=>", err)
		return err
	}
	if !filepath.IsAbs(bc.DBHome) {
		bc.CurDBHome = filepath.Join(homeDir, string(filepath.Separator), bc.DBHome)
	} else {
		bc.CurDBHome = bc.DBParam
	}
	return nil
}
