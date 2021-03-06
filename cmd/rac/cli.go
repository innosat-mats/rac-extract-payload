package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/innosat-mats/rac-extract-payload/internal/awstools"
	"github.com/innosat-mats/rac-extract-payload/internal/common"
	"github.com/innosat-mats/rac-extract-payload/internal/exports"
	"github.com/innosat-mats/rac-extract-payload/internal/extractors"
)

// Version is the version of the source code
var Version string

// Head is the short commit id of head
var Head string

// Buildtime is the time of the build
var Buildtime string

var skipImages *bool
var skipTimeseries *bool
var project *string
var stdout *bool
var aws *bool
var awsDescription *string
var version *bool

//myUsage replaces default usage since it doesn't include information on non-flags
func myUsage() {
	fmt.Println("Extracts information from Innosat-MATS rac-files")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTIONS] rac-file ...\n", os.Args[0])
	if len(os.Args) > 2 {
		switch helpSection := strings.ToUpper(os.Args[2]); helpSection {
		case "OUTPUT":
			infoGeneral()
		case "CCD":
			infoCCD()
		case "CPRU":
			infoCPRU()
		case "HTR":
			infoHTR()
		case "PWR":
			infoPWR()
		case "STAT":
			infoSTAT()
		case "TCV":
			infoTCV()
		case "PM":
			infoPM()
		case "MATS", "SPACE", "M.A.T.S.", "SATELLITE":
			infoSpace()
		default:
			fmt.Printf("\nUnrecognized help section %s\n", helpSection)
		}
		return
	}
	flag.PrintDefaults()
	fmt.Printf(
		"\nFor extra information about the output CSV:s type \"%s -help output\"\n",
		os.Args[0],
	)

	fmt.Println(`
The tool can be used to scan rac files for contents. Use the -stdout flag and
use command line tools to scan for interesting information, e.g.:
	rac -stdout my.rac  | grep STAT

Tip for finding parsing errors:
	rac -stdout my.rac | grep -E -e".*Error:[^<:]+" -o

or if you want the Buffer contents which can be rather large if you are unlucky:
	rac -stdout my.rac | grep -E -e".*Error:[^<]+" -o
	`)
}

func getCallback(
	toStdout bool,
	toAws bool,
	project string,
	skipImages bool,
	skipTimeseries bool,
	awsDescription string,
	wg *sync.WaitGroup,
) (common.Callback, common.CallbackTeardown, error) {
	if project == "" && !toStdout {
		flag.Usage()
		fmt.Println("\nExpected a project")
		return nil, nil, errors.New("Invalid arguments")
	}
	if skipTimeseries && (skipImages || toStdout) {
		fmt.Println("Nothing will be extracted, only validating integrity of rac-file(s)")
	}

	if toStdout {
		callback, teardown := exports.StdoutCallbackFactory(os.Stdout, !skipTimeseries)
		return callback, teardown, nil
	} else if toAws {
		callback, teardown := exports.AWSS3CallbackFactory(
			awstools.AWSUpload,
			project,
			awsDescription,
			!skipImages,
			!skipTimeseries,
			wg,
		)
		return callback, teardown, nil
	}
	callback, teardown := exports.DiskCallbackFactory(
		project,
		!skipImages,
		!skipTimeseries,
		wg,
	)
	return callback, teardown, nil
}

func processFiles(
	extractor extractors.ExtractFunction,
	inputFiles []string,
	callback common.Callback,
) error {
	batch := make([]extractors.StreamBatch, len(inputFiles))
	for n, filename := range inputFiles {
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}
		batch[n] = extractors.StreamBatch{
			Buf: f,
			Origin: &common.OriginDescription{
				Name:           filename,
				ProcessingDate: time.Now(),
			},
		}

	}
	extractor(callback, batch...)
	return nil
}

func init() {
	common.Version = Version
	common.Head = Head
	common.Buildtime = Buildtime

	skipImages = flag.Bool("skip-images", false, "Extract images from rac-files.\n(Default: false)")
	skipTimeseries = flag.Bool(
		"skip-timeseries",
		false,
		"Extract timeseries from rac-files.\n(Default: false)",
	)
	project = flag.String(
		"project",
		"",
		"Name for experiments, when outputting to disk a directory will be created with this name, when sending to AWS files will have this as a prefix",
	)
	stdout = flag.Bool(
		"stdout",
		false,
		"Output to standard out instead of to disk (only timeseries)\n(Default: false)",
	)
	aws = flag.Bool(
		"aws",
		false,
		"Output to aws instead of disk (requires credentials and permissions)",
	)
	awsDescription = flag.String(
		"description",
		"",
		"Path to a file containing a project description to be uploaded to AWS",
	)
	version = flag.Bool(
		"version",
		false,
		"Only display current version of the program",
	)

	flag.Usage = myUsage
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	if *version {
		fmt.Println("Version", Version, "Commit", Head, "@", Buildtime)
		return
	}

	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		flag.Usage()
		log.Fatal("No rac-files supplied")
	}
	callback, teardown, err := getCallback(
		*stdout,
		*aws,
		*project,
		*skipImages,
		*skipTimeseries,
		*awsDescription,
		&wg,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = processFiles(extractors.ExtractData, inputFiles, callback)
	if err != nil {
		log.Fatal(err)
	}
	teardown()
}
