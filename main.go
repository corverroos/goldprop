package main

import (
	"log"
	"os"

	"github.com/corverroos/goldprop/models"
	"github.com/corverroos/goldprop/p24"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "goldprop"
	app.Usage = "cli for scraping, searching and other goldprop magic"

	app.Commands = []cli.Command{
		{
			Name:     "property24",
			Aliases:  []string{"p24"},
			Category: "property24.com actions",
			Subcommands: []cli.Command{
				{
					Name:   "scrape",
					Usage:  "scrape properties for area(s)",
					Action: p24Scrape,
				},
				{
					Name:   "genareas",
					Usage:  "generate area code mapping",
					Action: p24GenAreas,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func p24Scrape(ctx *cli.Context) error {
	db, err := setupDB()
	if err != nil {
		return err
	}

	return p24.Scrape(db, ctx.Args()...)
}

func p24GenAreas(ctx *cli.Context) error {
	return p24.GenAreas()
}

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root@tcp(127.0.0.1:3306)/goldprop?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&models.Listing{})
	db.AutoMigrate(&models.Features{})

	return db, nil
}
