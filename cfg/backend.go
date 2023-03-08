package cfg

import (
	"fmt"
	"github.com/blockchainstamp/bsmailserver-go/util"
	"path/filepath"
)

const (
	DBTypMem DBTyp = iota
	DBTypJson
	DBTypLevelDB
	DBTypSqlite
)

type DBTyp int8

func (dt DBTyp) String() string {
	switch dt {
	case DBTypMem:
		return "memory database"
	case DBTypJson:
		return "json file database"
	case DBTypLevelDB:
		return "leveldb"
	case DBTypSqlite:
		return "sqlite"
	default:
		return "unknown"
	}
}

type BackConfig struct {
	DBType    DBTyp  `json:"db_type"` //0:memory db; 1:json file db; 2:leveldb; 3:sqlite db
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
	fmt.Println(bc.String())
	return nil
}

func (bc *BackConfig) String() string {
	s := "\n======Backend Storage Config======"
	s += "\nDBType:   \t" + bc.DBType.String()
	s += "\nCurDBHome:\t" + bc.CurDBHome
	s += "\nDBParam:  \t" + bc.DBParam
	s += "\nDBHome:   \t" + bc.DBHome
	s += "\n==================================="
	return s

}
