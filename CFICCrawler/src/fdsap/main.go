package main

import (
	//"crawler"
	//"routingpool"
	"runtime"
	//"fdsap/utility"
	"flag"
	//"os"
	//"runtime/pprof"

	"log"
	//"github.com/spf13/viper"
	"time"
	"fdsap/utility"
	"fdsap/crawler"
)



/************************************************
Financial Data Statistics & Analysis Platform
1. Download pages
2. Download history data
3. Analyse
************************************************/

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {
/*
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
    */
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	utility.NewConfig("D:\\programing\\golang\\CFICCrawler\\resource\\configuration\\")

	// Log setting
	utility.Init_Logger()

	//routingpool.SetPoolSize(viper.GetInt("routinepool.number"), viper.GetInt("routinepool.capacity"))
	//routingpool.Start()

	//var proxy *httpcontroller.Proxy = nil
	//proxy := &httpcontroller.Proxy{"HTTP", "203.17.66.134", "8000"}
	//folder := "D:/Work/MyDemo/go/golang/CFICCrawler/resource/download/"

	//crawler.StartCrawl([]string{})
	//routingpool.Wait() 	// Waiting for all threads finish and exit

	crawler.StartCrawl([]string{})
	elapsed := time.Since(start)

    log.Println("Main end...", elapsed)
	//logger.Debugf("Exit...........................%d", elapsed)
}
