package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
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

func getWordlist(wordlistPath string) ([]string, error) {
	if err := checkIfFileExists(wordlistPath); err != nil {
		log.Error().
			Msg("File does not exist")
		return nil, err
	}

	file, err := os.OpenFile(wordlistPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Msg("Failed to read text file")
		return nil, fmt.Errorf("failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordlist := []string{}

	for scanner.Scan() {
		wordlist = append(wordlist, scanner.Text())
	}

	return wordlist, nil
}

func checkPathForUrl(targetUrl, wordlistPath string) error {
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		log.Error().
			Msg("Not able to Parse url")
		return err
	}

	urlPaths, _ := getWordlist(wordlistPath)
	for _, urlPath := range urlPaths {
		currentUrl, err := url.JoinPath(parsedUrl.String(), urlPath)
		if err != nil {
			log.Error().
				Msgf("Failed to join %q with %q", parsedUrl.String(), urlPath)
			continue
		}

		req, err := http.NewRequest(http.MethodGet, currentUrl, nil)
		if err != nil {
			log.Error().
				Msg("Failed to create request")
			continue
		}

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().
				Msgf("Failed to send request to %q with error: %s", req.URL.String(), err)
			continue
		}

		if response.StatusCode == http.StatusOK {
			log.Info().
				Msgf("Found: %q", urlPath)
		} else {
			log.Debug().
				Msgf("Not Found: %q", urlPath)
		}
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

	if *vhostWordlist != "" {
		log.Info().
			Str("target_url", *targetUrl).
			Msg("Testing the host for any vhost")
		// checkForVHost(*targetUrl, *vhostWordlist)
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
