package prometheus

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"text/template"
)

type Config struct {
	Addresses []string
}
func GenerateProConfig(configs []string) {
	file, err := os.OpenFile("/usr/local/prometheus-2.52.0.linux-amd64/prometheus.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
    //file, err := os.OpenFile("/home/swung/桌面/minik8s/prometheus/prometheus.yml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error("[GenerateConfig] error opening file: ", err)
	}
	defaultHeader := `
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
`
	bytes, err := file.WriteString(defaultHeader)
	if err != nil {
		log.Error("[GenerateProConfig] error writing to file: ", err)
	}
	log.Debug("[GenerateProConfig] wrote ", bytes, " bytes to file")

	tmpl := template.Must(template.ParseFiles("/usr/local/prometheus-2.52.0.linux-amd64/prometheus.tmpl"))
	//tmpl := template.Must(template.ParseFiles("/home/swung/桌面/minik8s/prometheus/prometheus.tmpl"))

	config := Config{
		Addresses: configs,
	}
	err = tmpl.Execute(file, config)
	if err != nil {
		log.Error("[GenerateProConfig] error executing template: ", err)
	}
	file.Close()
}

func ReloadPrometheus() error{
	cmd := exec.Command("systemctl", "reload", "prometheus")
	err := cmd.Run()
	if err != nil {
		log.Error("[ReloadPrometheus] error starting prometheus: ", err)
		return err
	}
	return nil
}
