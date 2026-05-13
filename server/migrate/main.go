package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"rabidaudio.com/coven-door/server/database"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db := must(database.Create())
	defer db.Close()

	m := must(database.NewMigrator(db))

	n := flag.Int("n", -1, "number of steps to migrate")
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "up":
		if *n < 0 {
			err := m.Up()
			if err != nil {
				panic(err)
			}
		} else {
			err := m.Steps(*n)
			if err != nil {
				panic(err)
			}
		}
	case "down":
		if *n < 0 {
			// assume 1
			err := m.Steps(-1)
			if err != nil {
				panic(err)
			}
		} else {
			err := m.Steps(-1 * *n)
			if err != nil {
				panic(err)
			}
		}

	case "new":
		name := flag.Arg(1)
		if name == "" {
			fmt.Println("name required")
			ShowHelp()
		}

		now := time.Now()
		version := now.Unix()
		date := now.Format("2006-01-02")
		for _, dir := range []string{"up", "down"} {
			path := fmt.Sprintf("server/migrations/%v_%v_%v.%v.sql", version, date, sanitize(name), dir)
			file := must(os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666))
			defer file.Close()

			file.WriteString("TODO UNIMPLEMENTED\n")

			fmt.Printf("Created file: %v\n", path)
		}
	case "to":
		version := flag.Arg(1)
		if version == "" {
			fmt.Println("version required")
			ShowHelp()
		}
		v := must(strconv.ParseUint(version, 10, 64))
		err := m.Migrate(uint(v))
		if err != nil {
			panic(err)
		}

	case "status":
		version, _, err := m.Version()
		if err != nil {
			panic(err)
		}
		fmt.Printf("current version: %v\n", version)

	case "":
		fmt.Println("command required")
		ShowHelp()
	default:
		fmt.Printf("unknown command:%v\n", cmd)
		ShowHelp()
	}
}

func ShowHelp() {
	fmt.Println("migrate <cmd> [flags]\n\n\tnew <name>\n\tup [-n 1] (default: all)\n\tdown [-n 3] (default 1)\n\tto <version>")
	os.Exit(1)
}

func sanitize(name string) string {
	return regexp.MustCompile("[^a-z0-9_-]+").ReplaceAllString(strings.ToLower(name), "_")
}

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
