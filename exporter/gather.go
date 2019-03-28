package exporter

import (
	"fmt"
	"github.com/nfrush/prometheus_service_exporter/config"
	"os/exec"
	"regexp"
	"strings"
)

var (
	isSysvinit = regexp.MustCompile(`(\S+).*\sis\s(.*[^.])`)
)

func (e *Exporter) GetPIDState(pid string) string {
	cmd := exec.Command("/bin/sh", "-c", "ps -p " + pid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return "dead"
	}
	if len(output) > 1 {
		return "alive"
	}
	return "dead"
}

func (e *Exporter) DeriveState(state string) string {
	if state == "running" {
		return "active"
	}
	return "inactive"
}

func (e *Exporter) IsWhitelistedService(name string, target_group config.TargetGroup) bool {
    for _, v := range target_group.Daemon_Whitelist {
    	if name == v || name == v + ".service" {
    		return true
		}
	}
    return false
}

func (e *Exporter) sysvinit() ([]*Service, error){

	services := []*Service{}

	cmd := exec.Command("service", "--status-all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return services, err
	}
	lines := strings.Split(string(output), "\n")
	for _, v := range lines {
		m := isSysvinit.FindStringSubmatch(v)
		if len(m) == 3 {
			for _, target_group := range e.Targets {
				if e.IsWhitelistedService(m[1], target_group) {
					service := &Service{Name: m[1], State: e.DeriveState(m[2]), Substate: m[2], Group: target_group.Name}
					services = append(services, service)
				}
			}
		}
	}
	return services, nil
}

func (e *Exporter) systemd() (services []*Service, err error){
	services = []*Service{}

	cmd := exec.Command("systemctl", "list-units", "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return services, err
	}
	lines := strings.Split(string(output), "\n")
	for _, v := range lines[1:len(lines) - 8] {
		parts := strings.Fields(v)
		for _, target_group := range e.Targets {
			if e.IsWhitelistedService(parts[0], target_group) {
				service := &Service{Name: parts[0], State: parts[2], Substate: parts[3], Group: target_group.Name}
				services = append(services, service)
			}
		}
	}
	return services, nil
}

func (e *Exporter) pid() (pids []*PidFile, err error) {
	pids = []*PidFile{}

	for _, target_group := range e.Targets {
		for _, filename := range target_group.Pid_Whitelist {
			cmd := exec.Command("cat", fmt.Sprintf("/var/run/%s.pid", filename.Name))
			output, err := cmd.CombinedOutput()
			if err != nil {
				continue
			}
			pid := &PidFile{
				Name: fmt.Sprintf("%s.pid", filename.Name),
				PID: string(output),
				State: e.GetPIDState(string(output)),
				Service: filename.Service,
				Group: target_group.Name,
			}
			pids = append(pids, pid)
		}
	}

	return pids, err
}
