// Package env get env & app config, all the public field must after init()
// finished and flag.Parse().
package env

import (
	"flag"
	"os"
)

// deploy env.
const (
	DeployEnvDep  = "dev"
	deployEnvProd = "prod"
)

// env default value.
const (
	_deployEnv = "dev"
)

// env configuration.
var (
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string
	// AppID is global unique application id, register by service tree.
	// such as main.arch.demo-api.
	AppID string
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&AppID, "appid", os.Getenv("APP_ID"), "appid is global unique application id, register by service tree. or use APP_ID env variable.")
	fs.StringVar(&DeployEnv, "deploy.env", defaultString("DEPLOY_ENV", _deployEnv), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
