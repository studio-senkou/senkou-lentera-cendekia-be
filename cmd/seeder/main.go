package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/database/seeders"
)

type SeederClosure func(*sql.DB) error

type SeederContext struct {
	seeders map[string]SeederClosure
}

func NewSeederContext() *SeederContext {
	return &SeederContext{
		seeders: make(map[string]SeederClosure),
	}
}

func (ctx *SeederContext) Register(name string, fn SeederClosure) {
	ctx.seeders[name] = fn
}

func (ctx *SeederContext) Run(name string, db *sql.DB) error {
	fn, ok := ctx.seeders[name]
	if !ok {
		return fmt.Errorf("unknown seeder: %s", name)
	}

	return fn(db)
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Failed to load .env file: %v", err))
	}

	if err := database.InitializeDatabase(); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	defer database.CloseDatabase()

	db := database.GetDB()
	defer db.Close()

	ctx := NewSeederContext()
	register(ctx)

	seeder := "all"
	if len(os.Args) > 1 {
		seeder = strings.TrimSuffix(os.Args[1], ".go")
	}

	if err := ctx.Run(seeder, db); err != nil {
		panic(err)
	}
	fmt.Println("Database seeded successfully")
}

func register(c *SeederContext) {
	c.Register("user_seeder", seeders.SeedUsers)
	c.Register("all", func(db *sql.DB) error {
		if err := seeders.SeedUsers(db); err != nil {
			return err
		}
		return nil
	})
}
