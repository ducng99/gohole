//go:build windows

package dns

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ducng99/gohole/internal/logger"
)

func RegisterAutostart() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)

	taskScheduleContent := `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Author>gohole</Author>
	<Description>Starts gohole DNS server</Description>
  </RegistrationInfo>
  <Principals>
    <Principal id="Author">
      <LogonType>S4U</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <UseUnifiedSchedulingEngine>true</UseUnifiedSchedulingEngine>
  </Settings>
  <Triggers>
    <BootTrigger />
  </Triggers>
  <Actions Context="Author">
    <Exec>
      <Command>"` + exePath + `"</Command>
      <Arguments>dns start</Arguments>
      <WorkingDirectory>` + exeDir + `</WorkingDirectory>
    </Exec>
  </Actions>
</Task>`

	taskFile, err := os.CreateTemp("", "gohole_*_task.xml")
	if err != nil {
		logger.Printf(logger.LogError, "Could not create task schedule file\n")
		return err
	}
	defer taskFile.Close()
	defer os.Remove(taskFile.Name())

	if _, err := taskFile.WriteString(taskScheduleContent); err != nil {
		logger.Printf(logger.LogError, "Could not write task schedule file\n")
		return err
	}

	taskFile.Close()

	command := exec.Command("schtasks", "/create", "/xml", taskFile.Name(), "/tn", "\\gohole DNS server")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		logger.Printf(logger.LogError, "Could not create a task in Task Scheduler\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Successfully register gohole DNS server autostart\n")

	if err := exec.Command("schtasks", "/run", "/tn", "\\gohole DNS server").Run(); err != nil {
		logger.Printf(logger.LogError, "Could not start gohole DNS server task\n")
		return err
	}

	logger.Printf(logger.LogNormal, "gohole DNS server started\n")

	return nil
}
