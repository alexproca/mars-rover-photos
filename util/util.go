package util

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"nasa-api/config"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	TIMEOUT time.Duration = time.Second * 5
	RETRIES int           = 4
)

// HTTPGet - return http get content at url as []byte
func HTTPGet(url string) []byte {

	body := []byte{}

	success := false
	count := 0

	for success == false && count < RETRIES {
		count++
		body, success = getBody(url)
	}

	return body
}

func getBody(url string) ([]byte, bool) {

	body := []byte{}

	httpClient := http.Client{
		Timeout: TIMEOUT,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(SafeString(err.Error()))
		return body, false
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Println(SafeString(getErr.Error()))
		return body, false
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(SafeString(readErr.Error()))
		return body, false
	}

	return body, true
}

func SystemSignalsHandler(signals chan os.Signal) {
	go func() {
		sig := <-signals
		switch sig {
		case syscall.SIGINT:
			log.Println("\nCtrl-C signalled\n")
			os.Exit(0)
		}
	}()
}

//TemplateHandler - fill template and send result
func TemplateHandler(templateName string, items interface{}, writer http.ResponseWriter) error {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", templateName))
	if err != nil {
		log.Printf("Parse error for template '%s'\n %s \n", templateName, err)
	}
	return t.Execute(writer, items)
}

//SafeString - remove sensitive data from logged strings
func SafeString(in string) (out string) {
	out = strings.Replace(in, config.Config.EndpointData.ApiKey, "**********", -1)
	return
}
