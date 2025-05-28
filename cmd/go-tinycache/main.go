package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gleicon/go-tinycache/internal/metrics"
	"github.com/gleicon/go-tinycache/internal/networking"
	"github.com/pkg/profile"
	log "github.com/sirupsen/logrus"
)

func main() {
	address := flag.String("s", "127.0.0.1", "Bind Address")
	port := flag.String("p", "11211", "Bind Port")
	filename := flag.String("f", "./memcached.db", "path and file for database")
	backendType := flag.String("b", "boltdb", "backend: boltdb, inmem")
	pf := flag.Bool("q", false, "Enable profiling")
	dumpLogs := flag.Bool("m", false, "Enable metric dump each 60 seconds")

	flag.Usage = func() {
		fmt.Println("Usage: go-tinycache [-s ip] [-p port] [-f /path/to/db/file -q -b boltdb|inmem]")
		fmt.Println("default ip: 127.0.0.1")
		fmt.Println("default port: 11211")
		fmt.Println("default backend: boltdb")
		fmt.Println("default file: ./memcached.db")
		fmt.Println("-q enables profiling to /tmp/*.prof")
		os.Exit(1)
	}

	flag.Parse()

	if *pf == true {
		c := profile.Start(profile.CPUProfile, profile.ProfilePath("/tmp"), profile.NoShutdownHook)
		defer c.Stop()
	}

	metricsInstance := metrics.NewInternalMetrics(*filename, *dumpLogs)

	server := networking.NewServer(metricsInstance, *backendType, log.StandardLogger(), *filename)
	log.Infof("Starting server at %s port: %s | Storage backend: %s", *address, *port, *backendType)

	server.Serve(*address, *port, *filename)
}
