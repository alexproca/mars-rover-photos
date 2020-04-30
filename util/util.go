package util

import (
	"fmt"
	"html/template"
	"io/ioutil"
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

func GetBody(url string) []byte {

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
		fmt.Println(safeString(err.Error()))
		return body, false
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		fmt.Println(safeString(getErr.Error()))
		return body, false
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(safeString(readErr.Error()))
		return body, false
	}

	return body, true
}

func SystemSignalsHandler(signals chan os.Signal) {
	go func() {
		sig := <-signals
		switch sig {
		case syscall.SIGINT:
			fmt.Printf("\nCtrl-C signalled\n")
			os.Exit(0)
		}
	}()
}


func TemplateHandler(templateName string, items interface{}, writer http.ResponseWriter) error {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", templateName))
	if err != nil {
		fmt.Println(err)
	}
	return t.Execute(writer, items)
}

func safeString(in string) (out string) {
	out = strings.Replace(in, config.Config.EndpointData.ApiKey, "**********", -1)
	return
}