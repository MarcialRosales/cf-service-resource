package out

import (
	"os"
	"os/exec"
)

type PAAS interface {
	Login(api string, username string, password string, insecure bool) error
	Target(organization string, space string) error
	CreateService(service string, plan string, instanceName string) error
	CreateUserProvidedService(instanceName string, credentials string) error
	BindService(currentAppName string, instanceName string) error
	RestageApp(currentAppName string) error
}

type CloudFoundry struct{}

func NewCloudFoundry() *CloudFoundry {
	return &CloudFoundry{}
}

func (cf *CloudFoundry) Login(api string, username string, password string, insecure bool) error {
	args := []string{"api", api}
	if insecure {
		args = append(args, "--skip-ssl-validation")
	}

	err := cf.cf(args...).Run()
	if err != nil {
		return err
	}

	return cf.cf("auth", username, password).Run()
}

func (cf *CloudFoundry) Target(organization string, space string) error {
	return cf.cf("target", "-o", organization, "-s", space).Run()
}

func (cf *CloudFoundry) CreateService(service string, plan string, instanceName string) error {
	args := []string{}
	args = append(args, "create-service", service, plan, instanceName)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) CreateUserProvidedService(instanceName string, credentials string) error {
	args := []string{}
	args = append(args, "create-user-provided-service", instanceName, "-p", credentials)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) BindService(currentAppName string, instanceName string) error {
	args := []string{}
	args = append(args, "bind-service", currentAppName, instanceName)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) RestageApp(currentAppName string) error {
	args := []string{}
	args = append(args, "restage", currentAppName)

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) cf(args ...string) *exec.Cmd {
	cmd := exec.Command("cf", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CF_COLOR=true")

	return cmd
}
