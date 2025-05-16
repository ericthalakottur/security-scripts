package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PANIC = zerolog.PanicLevel
	FATAL = zerolog.FatalLevel
	ERROR = zerolog.ErrorLevel
	WARN  = zerolog.WarnLevel
	INFO  = zerolog.InfoLevel
	DEBUG = zerolog.DebugLevel
	TRACE = zerolog.TraceLevel
)

func checkIfFileExists(filepath string) error {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func checkPathForUrl(targetUrl, wordlistPath string) error {
	if err := checkIfFileExists(wordlistPath); err != nil {
		log.Error().
			Msg("File does not exist")
		return err
	}

	file, err := os.OpenFile(wordlistPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Msg("Failed to read text file")
		return fmt.Errorf("failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.Debug().
			Msgf("Checking %s", scanner.Text())
	}

	return nil
}

func getLogLevel(level string) zerolog.Level {
	logLevel := INFO
	switch level {
	case "PANIC":
		logLevel = PANIC
	case "FATAL":
		logLevel = FATAL
	case "ERROR":
		logLevel = ERROR
	case "WARN":
		logLevel = WARN
	case "INFO":
		logLevel = INFO
	case "DEBUG":
		logLevel = DEBUG
	case "TRACE":
		logLevel = TRACE
	default:
		log.Fatal().
			Msgf("%s is not a valid log level", level)
		os.Exit(1)
	}
	return logLevel
}

func main() {
	logLevel := flag.String("log", "INFO", "Log levels: PANIC, FATAL, ERROR, WARN, INFO, DEBUG, TRACE")
	targetUrl := flag.String("url", "", "Target Url")
	pathWordlist := flag.String("wl-path", "", "Wordlist to be used to query paths")
	vhostWordlist := flag.String("wl-vhost", "", "Wordlist to be used to query VHOST")

	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(getLogLevel(strings.ToUpper(*logLevel)))

	if *targetUrl == "" {
		log.Fatal().
			Msg("-url is a required argument")
		os.Exit(1)
	} else if (*pathWordlist == "") && (*vhostWordlist == "") {
		log.Fatal().
			Msg("-wl-path or -wl-vhost is a required argument")
		os.Exit(1)
	}

	if *pathWordlist != "" {
		log.Info().
			Str("target_url", *targetUrl).
			Msg("Testing the path for the target url")
		checkPathForUrl(*targetUrl, *pathWordlist)
	}

	fmt.Println("Path Wordlist: ", *pathWordlist)
	fmt.Println("VHOST Wordlist: ", *vhostWordlist)
}
