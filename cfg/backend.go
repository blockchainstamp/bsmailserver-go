package cfg

import "path/filepath"

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
}

func (bc *BackConfig) prepare(homeDir, dbHome string) error {
	if !filepath.IsAbs(dbHome) {
		dbHome = filepath.Join(homeDir, string(filepath.Separator), dbHome)
	}
	bc.CurDBHome = dbHome
	return nil
}
