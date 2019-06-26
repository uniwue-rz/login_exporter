# Login Exporter

Login Exporter is a simple Prometheus exporter that uses Chrome and 
Chrome-driver to open a website, to log in and then check for a given
text in the result. It simulate the first user interactions for the
given application. This can be used for the general availability test
of the application for the end users and for status management.

## Installation

This application it self does not need any installation, copy the binary
in the given directory should be enough. For running the logins and test
it uses chrome-headless (chromium-headless) and the chrome-driver, both
of them should be installed on the given machine. To Install them you
need to follow these guides:

### Install on MacOS

To install the needed dependencies, use brew:

```bash
brew install chromedriver
brew cask install google-chrome
```

### Install on Ubuntu

To install the needed dependencies on Ubuntu:

```bash
# Install Chrome.
sudo curl -sS -o - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add
sudo echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list
sudo apt-get -y update
sudo apt-get -y install google-chrome-stable

# Install ChromeDriver.
CHROME_DRIVER_VERSION=`curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE`
wget -N http://chromedriver.storage.googleapis.com/$CHROME_DRIVER_VERSION/chromedriver_linux64.zip -P ~/
unzip ~/chromedriver_linux64.zip -d ~/
rm ~/chromedriver_linux64.zip
sudo mv -f ~/chromedriver /usr/local/bin/chromedriver
sudo chown root:root /usr/local/bin/chromedriver
sudo chmod 0755 /usr/local/bin/chromedriver
```

## Configuration

The following parameters can be set for the application:

|  parameter  | default           | help|
|------------ |-------------------|-----|
| -config     | login.yml         | The path for the configuration file. This should be readable by running user|
| -listen_ip  | 127.0.0.1         | The IP address the application should listen to|
| -listen_port| 9980              | The port the application is listening to|
| -log_level  | INFO              | The log level for the application|
| -log_path   | login_exporter.log| The path for the log file. This should be writable by the running user|
| -timeout    | 60                | The timeout for the application to stop the check in seconds|


## Login.YML

There is an example file `login.yml.dist` available which can be used as
template for the `login.yml`. For every type of login the following 
parameters should be set:

### Simple Form

The simple is the most common type of login that is used for different
application. As long as no `CAPTCHAS` are used, this should work for every
kind of login forms. The following parameters should be set:

| parameter      | help           |
|----------------|----------------|
| login_type     | simple_form    |
| target         | The target that is searched for to find the config|
| url            | The url that login form is included|
| username       | The username that should be used for the login|
| password       | The password that should be used for the login|
| username_xpath | The xpath address of the username field (must be unique)|
| password_xpath | The xpath address of the password field (must be unique)|
| submit_xpath   | The xpath address of the submit button (must be unique)|
| expected_text  | The text that is expected to be there|

The following parameters are optional:

| parameter          | help           |
|--------------------|----------------|
| expected_text_xpath| The expected text xpath (must be unique), if not given the whole text is searched for the string|
| submit_type        | The type of submission that should be used it can be click or submit, default is submit|
| logout_xpath       | The xpath that should be used for the logout button|
| wait_time          | The time that should the form submitter wait for the page to load in seconds|


### Shibboleth 

Shibboleth is the type used for the services that use the Shibboleth IDP
for the authentication. As all of them at the end land on the same form,
this method can be used for all of them. The following parameters must
be set for the Shibboleth:

| parameter      | help           |
|----------------|----------------|
| login_type     | shibboleth     |
| target         | The target that is searched for to find the config|
| url            | The url that login form is included|
| username       | The username that should be used for the login|
| password       | The password that should be used for the login|
| expected_text  | The text that is expected to be there|

The following parameters are optional:

| parameter          | help           |
|--------------------|----------------|
| username_xpath | The xpath address of the username field (must be unique) for Shibboleth default value is the `//input[@id='username']`|
| password_xpath | The xpath address of the password field (must be unique) for Shibboleth default value is the `//input[@id='password']`|
| submit_xpath   | The xpath address of the submit button (must be unique) for Shibboleth default value is the `//button[@class='aai_login_button']`|
| expected_text_xpath| The expected text xpath (must be unique), if not given the whole text is searched for the string|
| submit_type        | The type of submission that should be used it can be click or submit, default is submit|

### Basic Auth

The application can be used to login in Basic Auth systems. The 
following parameters should be set for the application:

| parameter      | help           |
|----------------|----------------|
| login_type     | basic_auth     |
| url            | The url that login form is included|
| target         | The target that is searched for to find the config|
| username       | The username that should be used for the login|
| password       | The password that should be used for the login|
| expected_text  | The text that is expected to be there|

The following parameters are optional:

| parameter          | help           |
|--------------------|----------------|
| expected_text_xpath| The expected text xpath (must be unique), if not given the whole text is searched for the string|

### Password Only

Some login fields like mailman only have a single password or token
field, this kind of authentications are also supported. The following parameters
should be set:

| parameter      | help           |
|----------------|----------------|
| login_type     | password_only     |
| url            | The url that login form is included|
| target         | The target that is searched for to find the config|
| password       | The password that should be used for the login|
| password_xpath | The xpath address of the password field (must be unique)|
| expected_text  | The text that is expected to be there|

The following parameters are optional

| parameter          | help           |
|--------------------|----------------|
| expected_text_xpath| The expected text xpath (must be unique), if not given the whole text is searched for the string|
| submit_type        | The type of submission that should be used it can be click or submit, default is submit|
| logout_xpath       | The xpath that should be used for the logout button|

### API

This login exporter has limited support of API logins and approaches. 
The API part of the application does not use the chrome driver and uses
the `http.client` of go directly for the procedures. The following 
parameters should be set:
 
| parameter      | help           |
|----------------|----------------|
| login_type     | api     |
| url            | The url that login form is included|
| target         | The target that is searched for to find the config|
| password       | This could be key, password or token that should be used|
| password_xpath | This is the parameter that should be used for the token or password|
| expected_text  | The text that is expected to be there|

The following parameters are optional:

| parameter      | help           |
|----------------|----------------|
| username       | The username that should be used for the login|
| username_xpath | The parameter that is used for the username|
| method         | The method that should be used for the call default is POST|

### No Auth

When the page is only protected by IP address firewall or does not
have any authentication, but still should be checked for a given text
this method can be used.

| parameter      | help           |
|----------------|----------------|
| login_type     | no_auth     |
| url            | The url that login form is included|
| target         | The target that is searched for to find the config|
| expected_text  | The text that is expected to be there|

These parameters are optional:

| parameter          | help           |
|--------------------|----------------|
| expected_text_xpath| The expected text xpath (must be unique), if not given the whole text is searched for the string|

## Configuring Prometheus

This exporter works the same way as
[blackbox exporter](https://github.com/prometheus/blackbox_exporter). As
it uses a full browser to run the queries, The queries should
be done in bigger intervals and timeout should be set to a higher number.

In `prometheus.yaml the following settings should be enough:

```yaml
  - job_name: 'login_exporter'
    scrape_interval: 5m
    metrics_path: /probe
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9980
    file_sd_configs:
    - refresh_interval: 2m
      files:
         - '/etc/prometheus/login_targets.json'
```

In `login_targets.json` the following settings would be enough:

```json
[
  {
    "labels": {
      "group": "apps",
      "host": "hostname",
      "ip": "ip_address",
      "job": "login_exporter"
    },
    "targets": [
      "target_which_is_defined_in_login.yml_before"
    ]
  }
]
```

## Development

This application is open-source and can be extended. This repository
is a mirror of our home owned repository, as a result the pull request
here can not be directly merged. But pull requests are still welcome. 
I will extract the patch and add it manually to our internal repo, when
it is acceptable patch. You still can fork this repository and add your 
changes too.

### Build

To build this application you need to install the needed requirements
first with go get. Dont forget to install chrome and chrome driver
as it was discussed above.

```bash
go get
go build -o ./login_exporter
```

Or if you want to create a binary for several platforms at once, you can
use the `go_build.sh` script.

## Change Log

See the CHANGELOG file.

## License

See the LICENSE file.