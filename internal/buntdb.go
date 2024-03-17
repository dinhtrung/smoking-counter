package app

import (
	"github.com/knadh/koanf/providers/confmap"
	"github.com/tidwall/buntdb"
	"log/slog"
)

var (
	BuntDB         *buntdb.DB // DBConn hold the connection to database
	BuntDBInMemory *buntdb.DB // BuntDBInMemory provide an in-memory database
)

// BuntDBConfig configure application runtime
func BuntDBConfig() {
	// koanf defautl values
	Config.Load(confmap.Provider(map[string]interface{}{
		"buntdb.path": "assets/bunt.db",
	}, "."), nil)
}

// BuntDBInit initiate database
func BuntDBInit() {
	var err error
	slog.Info("Connecting to database", "path", Config.MustString("buntdb.path"))
	dbconn, err := buntdb.Open(Config.MustString("buntdb.path"))
	if err != nil {
		slog.Error("failed to connect database", "err", err)
	}
	BuntDB = dbconn
	BuntDBInMemory, _ = buntdb.Open(":memory:")
	slog.Info("connected to BuntDB", "path", Config.MustString("buntdb.path"))
}
