package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jguer/go-alpm/v2"
	"github.com/docker/go-units"
	"github.com/fatih/color"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	top := flag.Int("top", 0, "Show only X first packages")
	reverse := flag.Bool("r", false, "Reverse sorting order")
	sysroot := flag.String("sysroot", "/", "Operate on a mounted guest system")
	dbpath := flag.String("dbpath", "/var/lib/pacman", "Set an alternate database location")
	flag.Parse()

	h, e := alpm.Initialize(*sysroot, *dbpath)
	if e != nil {
		panic(e)
	}

	db, e := h.LocalDB()
	if e != nil {
		panic(e)
	}

	packagesBySize := db.PkgCache().SortBySize().Slice()

	if isFlagPassed("top") {
		if !*reverse {
			packagesBySize = packagesBySize[:*top]
		} else {
			packagesBySize = packagesBySize[len(packagesBySize)-*top:]
		}
	}

	printPkg := func(n int, pkg alpm.IPackage) {
		cNum := color.New(color.FgYellow).SprintFunc()
		cPkg := color.New(color.FgBlue).Add(color.Bold).SprintFunc()

		num := fmt.Sprintf("%5d", n)
		pkgName := fmt.Sprintf("%40s", pkg.Name())
		humanSize := units.HumanSize(float64(pkg.ISize()))

		fmt.Printf("%s  %s  %s\n", cNum(num), cPkg(pkgName), humanSize)
	}

	if *reverse {
		for i := len(packagesBySize) - 1; i >= 0; i-- {
			printPkg(len(packagesBySize)-i, packagesBySize[i])
		}
	} else {
		for i, pkg := range packagesBySize {
			printPkg(i+1, pkg)
		}
	}

	if h.Release() != nil {
		os.Exit(1)
	}
}
