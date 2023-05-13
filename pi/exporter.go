package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Sensor struct {
	Name     string  // Name of sensor
	Path     string  // Path to sensor
	Temp     float64 // Temp in celsus
	lastTick int     // Data loss tick counter
}

var (
	globSensors   = []Sensor{} // List of sensors
	verboseMode   = flag.Bool("v", false, "verbose mode")
	resolution    = flag.Int("r", 12, "data resolution")
	saveFile      = flag.String("f", "data.log", "file to save data to")
	skipProbe     = flag.Bool("s", false, "skip modprobe")
	logFileHandle *os.File // Logfile handle
)

// Check for errors and log them
func check(e error) {
	if e != nil {
		log.Println("[!]", e)
	}
}

// Log http requests
func httpLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		w.Header().Add("Content-Type", "text/plain")
		handler.ServeHTTP(w, r)
	})
}

// Preperations for execution
func _init() {
	if *verboseMode { // Enable logging
		var err error
		if _, chkFile := os.Open(*saveFile); chkFile != nil {
			os.Create(*saveFile)
		}
		// Open the logfile
		if logFileHandle, err = os.OpenFile(*saveFile, os.O_APPEND|os.O_WRONLY, os.ModeDevice); err != nil {
			panic(err)
		}
		defer logFileHandle.Close()
		log.SetOutput(logFileHandle)
	}

	// Search for sensors
	if !*skipProbe {
		err := exec.Command("/usr/sbin/modprobe", "ws-gpio").Run()
		check(err)
		err = exec.Command("/usr/sbin/modprobe", "w1-therm").Run()
		check(err)
	}

	_globSensors, err := filepath.Glob("/sys/bus/w1/devices/28*")
	check(err)
	for i, s := range _globSensors {
		slaveHandle, err := os.OpenFile(s+"/w1_slave", os.O_WRONLY, os.ModeDevice)
		check(err)
		_, err = slaveHandle.Write([]byte(strconv.Itoa(*resolution)))
		check(err)
		slaveHandle.Close()

		globSensors = append(globSensors, Sensor{Name: "Sensor" + strconv.Itoa(i), Temp: 0, Path: s + "/temperature"})
	}
}

// Read the temperature from a sensor
func readTemp(s Sensor) Sensor {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		log.Println("[!]", err)
	} else if string(data) != "0\n" {
		_dataS := strings.Split(string(data), "\n")
		fTemp, _ := strconv.ParseFloat(_dataS[0], 64)
		s.Temp = fTemp / 1000 // Convert to celsus
		return s
	}
	s.lastTick += 1
	return s
}

// Handle http requests
func httpMetricsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusForbidden)
		return
	}
	// Help menu (is it really needed?)
	fmt.Fprint(w, "# HELP sauna_temperature Check temperature from the sensors\n# TYPE sauna_temperature gauge\n")
	for _, s := range globSensors {
		fmt.Fprintln(w, "sauna_temperature{sensor=\""+s.Name+"\"} "+strconv.FormatFloat(s.Temp, 'f', 2, 64))
	}
}

func main() {
	wsPort := flag.String("p", "8000", "port to run on")
	flag.Parse()
	_init()

	// Start reading the sensors as a separete thread
	go func() {
		for range time.Tick(time.Second * 5) { // read the sensors every 5 seconds
			for i := range globSensors {
				globSensors[i] = readTemp(globSensors[i])
			}
		}
	}()

	// Start the webserver
	http.HandleFunc("/metrics", httpMetricsHandle)
	log.Println("Starting webserver at port", *wsPort, "(http://localhost:"+*wsPort+")")
	if err := http.ListenAndServe(":"+*wsPort, httpLogHandler(http.DefaultServeMux)); err != nil {
		check(err)
	}
}
