package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"

	"github.com/golang/glog"

	"github.com/tebeka/selenium"
)

const (
	// These paths will be different on your system.
	seleniumPath            = "vendor/selenium-server.jar"
	defaultGeckoDriverPath  = "C:/Users/lenovo/Documents/geckodriver.exe" // firefox
	defaultChromeDriverPath = "C:/Users/lenovo/Documents/chromedriver.exe"
	port                    = 19999
)

var (
	help               = flag.Bool("help", false, "this help")
	debug              = flag.Bool("debug", false, "enable debug")
	seleniumServerPath = flag.String("selenium_server_path", "../vendor/selenium-server.jar", "The path to the Selenium 3 server JAR.")
	firefoxBinary      = flag.String("firefox_binary", "../vendor/firefox/firefox", "The path of the Firefox binary.")
	geckoDriverPath    = flag.String("geckodriver_path", "", "The path to the geckodriver binary. required.")
	chromeDriverPath   = flag.String("chromedriver_path", "C:/Users/lenovo/Documents/chromedriver.exe", "The path to the ChromeDriver binary.")
	chromeBinary       = flag.String("chrome_binary", "C:/Program Files (x86)/Google/Chrome/Application/chrome.exe", "The name of the Chrome binary")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
	}
	if *debug {
		selenium.SetDebug(true)
	}

	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(defaultGeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                   // Output debug information to STDERR.
		selenium.ChromeDriver(defaultChromeDriverPath),
		selenium.Display("1", ""),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()
	println(`selenium server start success listening on ` + strconv.Itoa(port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT)

	<-stop
	println("selenium server stop.")
}

func setDriverPaths() error {
	if *seleniumServerPath == "" {
		*seleniumServerPath = findBestPath("vendor/selenium-server*" /*binary=*/, false)
	}

	if *geckoDriverPath == "" {
		*geckoDriverPath = findBestPath("vendor/geckodriver*" /*binary=*/, true)
	}

	if *chromeDriverPath == "" {
		*chromeDriverPath = findBestPath("vendor/chromedriver*" /*binary=*/, true)
	}

	return nil
}

func findBestPath(glob string, binary bool) string {
	matches, err := filepath.Glob(glob)
	if err != nil {
		glog.Warningf("Error globbing %q: %s", glob, err)
		return ""
	}
	if len(matches) == 0 {
		return ""
	}
	// Iterate backwards: newer versions should be sorted to the end.
	sort.Strings(matches)
	for i := len(matches) - 1; i >= 0; i-- {
		path := matches[i]
		fi, err := os.Stat(path)
		if err != nil {
			glog.Warningf("Error statting %q: %s", path, err)
			continue
		}
		if !fi.Mode().IsRegular() {
			continue
		}
		if binary && fi.Mode().Perm()&0111 == 0 {
			continue
		}
		return path
	}
	return ""
}
