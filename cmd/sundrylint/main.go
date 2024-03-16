package main

import (
	"log"

	"github.com/alingse/sundrylint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	setting := sundrylint.LinterSetting{}
	analyzer, err := sundrylint.NewAnalyzer(setting)
	if err != nil {
		log.Fatal(err)
	}
	singlechecker.Main(analyzer)
}
