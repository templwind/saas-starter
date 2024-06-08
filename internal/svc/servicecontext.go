package svc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/templwind/sass-starter/internal/config"
	"github.com/templwind/sass-starter/internal/db"
	"github.com/templwind/sass-starter/internal/menu"
	"github.com/templwind/sass-starter/internal/types"

	"github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	ctx    context.Context
	Config *config.Config
	SqlxDB *sqlx.DB
	Menus  menu.Menus
}

func NewServiceContext(ctx context.Context, c *config.Config) *ServiceContext {
	// Connect to the database
	dirPath := filepath.Join(c.DefaultDataDir)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err != nil {
			panic(fmt.Errorf("failed to check if database directory exists: %w", err))
		}
		// Attempt to create the directory if it doesn't exist
		err := os.MkdirAll(dirPath, 0755) // Adjust permissions as necessary
		if err != nil {
			panic(fmt.Errorf("failed to create database directory: %w", err))
		}
	}

	dbFilePath := filepath.Join(c.DefaultDataDir, c.DatabaseFileName)
	dsn := fmt.Sprintf("sqlite:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbFilePath)

	conn, err := db.NewPersistentSQLx(dsn, c)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	return &ServiceContext{
		ctx:    ctx,
		Config: c,
		SqlxDB: conn.GetDB(),
		Menus:  c.Menus,
	}
}

// AdjustMenus adjusts the menus based on the user context
// this will be performed on every request
func (svcCtx *ServiceContext) AdjustMenus(userContext types.ACLContext) {
	for name := range svcCtx.Menus {
		svcCtx.Menus[name] = svcCtx.Menus.GetContextualMenu(name, userContext)
	}
}
