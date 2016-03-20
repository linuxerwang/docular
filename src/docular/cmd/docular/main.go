package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"docular/server"
)

var (
	port                = flag.Int("port", 3455, "The http port.")
	docDir              = flag.String("doc-dir", ".", "The doc directory path.")
	webstatic           = flag.String("webstatic", "", "The web static directory.")
	allowExternalAccess = flag.Bool("allow-external", false, "Allow external access")

	allowedHosts                       = []string{"127.0.0.1", "localhost"}
	disallowedRootPath map[string]bool = map[string]bool{}
)

func init() {
	disallowedRootPath["/bin"] = true
	disallowedRootPath["/boot"] = true
	disallowedRootPath["/etc"] = true
	disallowedRootPath["/lib"] = true
	disallowedRootPath["/lib32"] = true
	disallowedRootPath["/lib64"] = true
	disallowedRootPath["/media"] = true
	disallowedRootPath["/mnt"] = true
	disallowedRootPath["/proc"] = true
	disallowedRootPath["/sbin"] = true
	disallowedRootPath["/srv"] = true
	disallowedRootPath["/sys"] = true
	disallowedRootPath["/usr"] = true
	disallowedRootPath["/var"] = true
}

func usage() {
	fmt.Println("Docular server.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  docular-server [options]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(0)
}

func checkFlags() {
	if *docDir == "" {
		fmt.Println("Error: flag doc-dir is required.")
		fmt.Println()
		usage()
	}
	if *webstatic == "" {
		current, err := user.Current()
		if err != nil {
			fmt.Println("Error: failed to get current user, ", err)
			os.Exit(2)
		}

		*webstatic = filepath.Join(current.HomeDir, "docular/webstatic")
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	checkFlags()

	docroot, err := filepath.Abs(*docDir)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(2)
	}
	isRootDir := docroot == "/"
	if _, found := disallowedRootPath[docroot]; found || isRootDir {
		fmt.Println("Error: serving files from this directory is not allowed.")
		os.Exit(2)
	}
	for p, _ := range disallowedRootPath {
		if strings.HasPrefix(docroot, p+"/") {
			fmt.Println("Error: serving files from this directory is not allowed.")
			os.Exit(2)
		}
	}

	server.Webstatic = *webstatic

	var target string
	if *allowExternalAccess {
		target = fmt.Sprintf(":%d", *port)
	} else {
		target = fmt.Sprintf("localhost:%d", *port)
	}

	fmt.Println("Serving files from " + docroot)
	fmt.Printf("Server running at http://localhost:%d. CTRL+C to shutdown\n", *port)
	err = http.ListenAndServe(target, server.NewWrapperHandler(http.Dir(docroot)))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
