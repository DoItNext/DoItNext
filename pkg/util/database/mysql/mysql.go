package mysql

import (
	"fmt"
	"github.com/DoItNext/DoItNext/pkg/util/config"
	"github.com/go-mods/wait-host"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func New() (*gorm.DB, error) {
	// database configuration
	cfg := config.Configuration.Database

	// database uri
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Charset)

	// Waiting for database to be online
	_ = waithost.Wait(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	// open database
	if db, err := gorm.Open(cfg.Type, uri); err != nil {
		return nil, err
	} else {
		return db, err
	}
}
