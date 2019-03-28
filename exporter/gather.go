package exporter

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	is_service = regexp.MustCompile(`.service`)
)

func (e *Exporter) GetPIDState(pid string) string {
	cmd := exec.Command("ps", "-p", pid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "dead"
	}
	if len(output) > 1 {
		return "alive"
	}
	return "dead"
}

func (e *Exporter) IsService(name string) bool {
	m := is_service.FindAllStringSubmatch(name, 1)
	if len(m) > 0 {
		return true
	}
	return false
}

func (e *Exporter) IsWhitelistedService(name string) bool {
    for _, v := range e.ServiceWhitelist {
    	m := v.FindAllStringSubmatch(name, 1)
    	if len(m) > 0 {
			return true
		}
	}
    return false
}

func (e *Exporter) IsWhitelistedPID(name string) bool {
	for _, v := range e.PIDWhitelist {
		m := v.FindAllStringSubmatch(name, 1)
		if len(m) > 0 {
			return true
		}
	}
	return false
}

func (e *Exporter) systemd() ([]*Service, error){

	services := []*Service{}

	cmd := exec.Command("systemctl", "list-units", "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return services, err
	}
	lines := strings.Split(string(output), "\n")
	for _, v := range lines[1:len(lines) - 8] {
		parts := strings.Fields(v)
		if e.IsWhitelistedService(parts[0]) && e.IsService(parts[0]) {
			service := &Service{Name: parts[0], State: parts[2], Substate: parts[3]}
			services = append(services, service)
		}
	}
	return services, nil
}

func (e *Exporter) pid() ([]*PidFile, error) {

	pids := []*PidFile{}

	files := []string{}

	err := filepath.Walk("/var/run/", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pid" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return pids, err
	}

	for _, file := range files {
		cmd := exec.Command("cat", file)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return pids, err
		}
		if e.IsWhitelistedPID(file) {
			pid := &PidFile{Name: file, PID: string(output), State: e.GetPIDState(string(output))}
			pids = append(pids, pid)
		}
	}

	return pids, nil
}
