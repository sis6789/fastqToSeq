package main

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/text/message"
)

var _ = flag.BoolP("help", "h", false, "show parameters")
var _ = flag.StringP("pattern", "p", "*.fastq", "enter path and pattern")
var _ = flag.StringP("out", "o", "read.txt", "fastq sequence file")

func main() {
	CommaP := message.NewPrinter(message.MatchLanguage("kr"))
	var err error
	wd, err := os.Getwd()
	log.Printf("start %v", wd)

	{
		// 명령줄 flag를 설정한다.
		flag.Parse()
		err = viper.BindPFlags(flag.CommandLine)
		if err != nil {
			log.Printf("%v", err)
		}
		if viper.GetBool("help") {
			flag.PrintDefaults()
			return
		}
	}

	files, err := filepath.Glob(viper.GetString("pattern"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(files) == 0 {
		log.Printf("no files for %v", viper.GetString("pattern"))
		return
	}

	outHandle, err := os.Create(viper.GetString("out"))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer func() {
		_ = outHandle.Close()
	}()

	outCount := 0
	for _, fn := range files {
		fnLines := 0
		var scanner *bufio.Scanner
		inHandle, err := os.Open(fn)
		if err != nil {
			log.Printf("%v %v", fn, err)
			continue
		}
		if strings.HasSuffix(fn, "gz") {
			gzReader, err := gzip.NewReader(inHandle)
			if err != nil {
				log.Printf("%v %v", fn, err)
				continue
			}
			scanner = bufio.NewScanner(gzReader)
		} else {
			scanner = bufio.NewScanner(inHandle)
		}
		recOrder := 2
		for scanner.Scan() {
			recOrder++
			if recOrder%4 == 0 {
				_, _ = outHandle.WriteString(scanner.Text())
				_, _ = outHandle.WriteString("\n")
				outCount++
				fnLines++
			}
		}
		_ = inHandle.Close()
		log.Printf("%v %15v", fn, CommaP.Sprint(fnLines))
	}

	log.Printf("%v %15v", "EOJ:"+strings.Repeat(" ", len(files[0])-4),
		CommaP.Sprint(outCount))
}
