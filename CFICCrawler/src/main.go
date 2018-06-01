package main

import (
	"crawler"
	"routingpool"
	"runtime"
	"utility"
	"flag"
	"os"
	"runtime/pprof"

	"log"
	"github.com/spf13/viper"
	"time"
)

/************************************************
1. Download pages
2. Download history data
3. Analyse
************************************************/

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

// proxy //http://203.17.66.133:8000   http://203.17.66.134:8000
func main() {
	flag.Parse()
	if *cpuprofile != "" {
		log.Println("Received cpuprofile=", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Create trace file error", err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		log.Println("Received memprofile=", *memprofile)
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	utility.NewConfig("E:\\Programing\\golang\\CFICCrawler\\resource\\configuration\\")

	// Log setting
	utility.Init_Logger()
	logger := utility.GetLogger()

	routingpool.SetPoolSize(viper.GetInt("routinepool.number"), viper.GetInt("routinepool.capacity"))
	routingpool.Start()
	defer routingpool.Wait() 	// Waiting for all threads finish and exit
	//var proxy *httpcontroller.Proxy = nil
	//proxy := &httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"}
	//folder := "D:/Work/MyDemo/go/golang/CFICCrawler/resource/download/"

	crawler.StartCrawl([]string{"600089", "600096", "600036", "123456"})




	elapsed := time.Since(start)
	logger.Debug("Exit...........................%d", elapsed)
}
