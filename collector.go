package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/sclevine/agouti"

	"gopkg.in/yaml.v2"
)

/// LoginConfigs Is the list of configuration
/// That should be read
type LoginConfigs struct {
	Configs []SingleLoginConfig `yaml:"targets"`
}

/// SingleLoginConfig Is the login configuration settings
/// that is used to read the yaml files.
type SingleLoginConfig struct {
	Url               string `yaml:"url"`
	Target            string `yaml:"target"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Certificate       string `yaml:"certificate"`
	UsernameXpath     string `yaml:"username_xpath"`
	PasswordXpath     string `yaml:"password_xpath"`
	CertificateXpath  string `yaml:"certificate_xpath"`
	SubmitXpath       string `yaml:"submit_xpath"`
	LoginType         string `yaml:"login_type"`
	ExpectedText      string `yaml:"expected_text"`
	ExpectedTextXpath string `yaml:"expected_text_xpath"`
	SSLCheck          bool   `yaml:"ssl_check"`
	Debug             bool   `yaml:"debug"`
	Method            string `yaml:"method"`
	SubmitType        string `yaml:"submit_type"`
	LogoutXpath       string `yaml:"logout_xpath"`
}

/// getChromeOptions Returns the options for the chrome driver that is used
/// to fetch the data from the server
func getChromeOptions() agouti.Option {
	return agouti.ChromeOptions("args", []string{
		"--headless",
		"--disable-gpu",
		"--no-first-run",
		"--no-default-browser-check",
		"--allow-insecure-localhost",
		"--no-sandbox",
	})
}

///startDriver Starts the given driver on the machine
func startDriver() *agouti.WebDriver {
	driver := agouti.ChromeDriver(getChromeOptions())
	if err := driver.Start(); err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "creation",
			}).Panicln(err.Error())
	}
	return driver
}

/// stopDriver Stops the given driver running. This
func stopDriver(driver *agouti.WebDriver) {
	err := driver.Stop()
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "destruction",
			}).Panicln(err.Error())
	}
}

/// readConfig Reads the yaml configuration from the given server
func readConfig(path string) LoginConfigs {
	var loginConfigs LoginConfigs
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "config_loader",
				"part":      "read_file",
			}).Panicln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &loginConfigs)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "config_loader",
				"part":      "parse_file",
			}).Panicln(err.Error())
	}
	return loginConfigs
}

// Configuration Params
var configFilePath string
var listenIp string
var listenPort int

// Logging Params
var logPath string
var logLevel string

// Timeout settings
var timeout int

var logger = log.New()

/// getCommandLineOptions Returns the command options from the terminal
func getCommandLineOptions() {
	flag.StringVar(&configFilePath, "config", "/etc/prometheus/login.yml", "Configuration file path")
	flag.StringVar(&listenIp, "listen_ip", "127.0.0.1", "Listen IP Address")
	flag.IntVar(&listenPort, "listen_port", 9980, "Listen Port")

	flag.StringVar(&logPath, "log_file", "login_exporter.log", "Log file path")
	flag.StringVar(&logLevel, "log_level", "INFO", "Log level")

	flag.IntVar(&timeout, "timeout", 60, "Timeout in seconds")

	flag.Parse()
}

/// getLogger Returns the logger that is used to log the data
func getLogger() *log.Logger {
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(f)
	parsedLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(parsedLevel)
	return logger
}

/// loginSimpleForm Logs in the simple for using the username, password and the submit button
func loginSimpleForm(page *agouti.Page, urlText string, usernameXpath string, passwordXpath string, submitXpath string,
	username string, password string, submitType string) {
	err := page.Navigate(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "navigation_error",
			}).Warningln(err.Error())
		panic(err)
	}
	usernameField := page.FindByXPath(usernameXpath)
	err = usernameField.SendKeys(username)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_simple_form",
				"part":      "username_field",
			}).Warningln(err.Error())
	}
	passwordField := page.FindByXPath(passwordXpath)
	err = passwordField.SendKeys(password)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_simple_form",
				"part":      "password_field",
			}).Warningln(err.Error())
	}
	submitField := page.FindByXPath(submitXpath)
	if submitType == "click" {
		err = submitField.Click()
	} else {
		err = submitField.Submit()
	}
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_simple_form",
				"part":      "submit_field",
			}).Warningln(err.Error())
	}
}

/// loginShibboleth Logs in the shibboleth system using the given username and password
func loginShibboleth(page *agouti.Page, urlText string, username string, password string, usernameXpath string,
	passwordXpath string, submitXpath string, submitType string) {
	err := page.Navigate(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "navigation_error",
			}).Warningln(err.Error())
	}
	if usernameXpath == "" && passwordXpath == "" && submitXpath == "" {
		usernameXpath = "//input[@id='username']"
		passwordXpath = "//input[@id='password']"
		submitXpath = "//button[@class='aai_login_button']"
	}
	usernameField := page.FindByXPath(usernameXpath)
	err = usernameField.SendKeys(username)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_shibboleth",
				"part":      "username_field",
			}).Warningln(err.Error())
	}
	passwordField := page.FindByXPath(passwordXpath)
	err = passwordField.SendKeys(password)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_shibboleth",
				"part":      "password_field",
			}).Warningln(err.Error())
	}
	if submitType == "" {
		submitType = "click"
	}
	submitField := page.FindByXPath(submitXpath)
	if submitType == "click" {
		err = submitField.Click()
	} else {
		err = submitField.Submit()
	}
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_shibboleth",
				"part":      "submit_field",
			}).Warningln(err.Error())
	}
}

/// loginApi Logs in the given API system with the given parameter for username/password and the method for the call
func loginApi(urlText string, usernameXpath string, passwordXpath string, username string,
	password string, method string) *http.Response {
	data := map[string]string{}
	if usernameXpath != "" {
		data[usernameXpath] = username
	}
	if passwordXpath != "" {
		data[passwordXpath] = password
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_api",
				"part":      "json_marshal",
			}).Warningln(err.Error())
	}
	req, err := http.NewRequest(method, urlText, bytes.NewBuffer(jsonValue))
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_api",
				"part":      "http_request_new",
			}).Warningln(err.Error())
	} else {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: time.Duration(time.Duration(timeout) * time.Second),
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.WithFields(
				log.Fields{
					"subsystem": "login_api",
					"part":      "http_request_do",
				}).Warningln(err.Error())
		} else {
			defer resp.Body.Close()
		}
		return resp
	}
	return nil
}

/// loginBasicAuth Logs in the given HTTPAuth page with the username and password given
func loginBasicAuth(page *agouti.Page, urlText string, username string, password string) {
	u, err := url.Parse(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_basic_auth",
				"part":      "url_parser",
			}).Warningln(err.Error())
	} else {
		urlText = u.Scheme + "://" + username + ":" + password + "@" + u.Hostname() + u.Path + "?" + u.RawQuery
	}
	err = page.Navigate(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "navigation_error",
			}).Warningln(err.Error())
	}
}

// logOut Logs out of the given page using the xpath that is given for the logout.
func logOut(page *agouti.Page, logoutXpath string, submitType string) {
	logoutField := page.FindByXPath(logoutXpath)
	if submitType == "click" {
		err := logoutField.Click()
		if err != nil {
			logger.WithFields(
				log.Fields{
					"subsystem": "logout",
					"part":      "click",
				}).Warningln(err.Error())
		}
	} else {
		err := logoutField.Submit()
		if err != nil {
			logger.WithFields(
				log.Fields{
					"subsystem": "logout",
					"part":      "submit",
				}).Warningln(err.Error())
		}
	}
}

/// loginPasswordOnly Logs in the given website with only password and the xpath to find the field
func loginPasswordOnly(page *agouti.Page, urlText string, passwordXPath string, submitXpath string, password string,
	submitType string) {
	err := page.Navigate(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "navigation_error",
			}).Warningln(err.Error())
	}
	passwordField := page.FindByXPath(passwordXPath)
	err = passwordField.SendKeys(password)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_password_only",
				"part":      "password_field",
			}).Warningln(err.Error())
	}
	submitField := page.FindByXPath(submitXpath)
	if submitType == "click" {
		err = submitField.Click()
	} else {
		err = submitField.Submit()
	}
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "login_password_only",
				"part":      "submit_field",
			}).Warningln(err.Error())
	}
}

/// checkExpected Checks if the expected text exists in the given path or page
func checkExpected(page *agouti.Page, expectedXPath string, expectedText string) bool {
	if expectedXPath != "" {
		expectedElement := page.FindByXPath(expectedXPath)
		expectedFieldText, err := expectedElement.Text()
		if err != nil {
			logger.WithFields(
				log.Fields{
					"subsystem": "check_expected",
					"part":      "xpath_data",
				}).Warningln(err.Error())
			return false
		}
		return expectedFieldText == expectedText
	}
	content, err := page.HTML()
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "check_expected",
				"part":      "match_all_data",
			}).Warningln(err.Error())
		return false
	}
	return strings.Contains(content, expectedText)
}

/// getNoLogin The no login function that only returns the page without any submissions
func getNoLogin(page *agouti.Page, urlText string) {
	err := page.Navigate(urlText)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "navigation_error",
			}).Warningln(err.Error())
	}
}

/// checkExpectedResponse Checks if the expected string exists in the http.Response
func checkExpectedResponse(response *http.Response, expectedText string) bool {
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "check_expected",
				"part":      "get_content",
			}).Warningln(err.Error())
	}
	return strings.Contains(string(content), expectedText)
}

///getStatus Returns the data from the server
func getStatus(config SingleLoginConfig) (status bool, elapsed float64) {
	status = false
	driver := startDriver()
	page, err := driver.NewPage()
	if err != nil {
		logger.WithFields(
			log.Fields{
				"subsystem": "driver",
				"part":      "new_page",
			}).Warningln(err.Error())
	} else {
		err = page.SetPageLoad(timeout * 1000)
	}

	start := time.Now()
	switch config.LoginType {
	case "simple_form":
		loginSimpleForm(page, config.Url, config.UsernameXpath, config.PasswordXpath, config.SubmitXpath,
			config.Username, config.Password, config.SubmitType)
		status = checkExpected(page, config.ExpectedTextXpath, config.ExpectedText)
		break
	case "shibboleth":
		loginShibboleth(page, config.Url, config.Username, config.Password, config.UsernameXpath,
			config.PasswordXpath, config.SubmitXpath, config.SubmitType)
		status = checkExpected(page, config.ExpectedTextXpath, config.ExpectedText)
		break
	case "basic_auth":
		loginBasicAuth(page, config.Url, config.Username, config.Password)
		status = checkExpected(page, config.ExpectedTextXpath, config.ExpectedText)
		break
	case "password_only":
		loginPasswordOnly(page, config.Url, config.PasswordXpath, config.SubmitXpath, config.Password, config.SubmitType)
		status = checkExpected(page, config.ExpectedTextXpath, config.ExpectedText)
		break
	case "api":
		response := loginApi(config.Url, config.UsernameXpath, config.PasswordXpath, config.Username, config.Password,
			config.Method)
		status = checkExpectedResponse(response, config.ExpectedText)
		break
	case "no_auth":
		getNoLogin(page, config.Url)
		status = checkExpected(page, config.ExpectedTextXpath, config.ExpectedText)
		break
	}
	end := time.Now()
	elapsed = end.Sub(start).Seconds()

	// logout if the value is set
	if config.LogoutXpath != "" {
		logOut(page, config.LogoutXpath, config.SubmitType)
	}

	if err == nil {
		err = page.CloseWindow()
		if err != nil {
			logger.WithFields(
				log.Fields{
					"subsystem": "driver",
					"part":      "new_page_close_error",
				}).Warningln(err.Error())
		}
	}
	stopDriver(driver)
	return status, elapsed
}

func init() {
	getCommandLineOptions()
	getLogger()
}
