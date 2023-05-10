package configuration

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	ErrInvalidParse = errors.New("couldn't parse")
	ErrInvalidHour  = errors.New("value is invalid for hours")
)

// Provider is type which is responsible to deliver configuration to the
// various components within the project.
type Provider struct {
	// MaxNodes sets what is maximal allowed size of generated graphs.
	// it is configurable by the `GENERATOR_MAX_NODES` environment variables
	// and `maxNodes` command line options.
	MaxNodes int `env:"GENERATOR_MAX_NODES" flag:"maxNodes"`

	// MaxBatchSize sets what is maximal number of graphs requested within one batch.
	// this option can be set by `GENERATOR_MAX_BATCH_SIZE` environment variable or
	// by `maxBatchSize` command line flag.
	MaxBatchSize int `env:"GENERATOR_MAX_BATCH_SIZE" flag:"maxBatchSize"`

	// Workers sets how many worker goroutines should be started for generation.
	// it is reasonable to use some fraction of available threads on target system
	// in order to ensure enough resources for the http interface to be responsive.
	Workers int `env:"GENERATOR_WORKERS" flag:"workers"`

	// DbRoot sets where the data for database should be stored, or load from.
	DbRoot string `env:"GENERATOR_DB_ROOT" flag:"dbRoot"`

	// MaintenanceInterval duration set to periodically start the cleaning goroutine of database.
	// currently is used to start checking if the maintenance mode should be started.
	MaintenanceInterval time.Duration `env:"GENERATOR_MAINTENANCE_INTERVAL" flag:"maintenanceInterval"`

	// MaintenanceHour sets the hour during day, where the clean is expected to happen.
	// If the MaintenanceInterval is lower than 24 hours the first maintenance is calculated
	// by adding MaintenanceInterval to the MaintenanceHour in current date, until it is in the future.
	MaintenanceHour BaseHourValue `env:"GENERATOR_MAINTENANCE_HOUR" flag:"maintenanceHour"`

	// RequestTTL duration set to each request telling the database for how long should be the
	// request and all dependent resources retained in persistent storage.
	RequestTTL time.Duration `env:"GENERATOR_REQ_TTL" flag:"ttl"`

	// LogLevel set the log level of the program - expected in format understood by the logrus library.
	LogLevel log.Level `env:"GENERATOR_LOG_LEVEL" flag:"logLevel"`

	// BindAddr sets IP address/address of the interface to bind to.
	BindAddr string `env:"GENERATOR_BIND_HOST" flag:"bindAddr"`

	// Host sets the expected host used by the application.
	Host string `env:"GENERATOR_HOST" flag:"host"`

	// Port which will be opened by the HTTP server
	Port string `env:"GENERATOR_PORT" flag:"port"`

	// SecureMode sets the cookie provider to strict and secure if passed as true,
	// if false the cookies are sent as not secured with Lax StrictMode.
	SecureMode bool `env:"GENERATOR_SECURE" flag:"secure"`

	// MailServer address to the server capable of accepting insecure connection
	// and sending emails. This is optional parameter and by default set to nil.
	// If not set server doesn't send any emails.
	MailServer *string `env:"GENERATOR_MAIL_SERVER" flag:"mailServer"`
	// AdminMail email address of administrator of system which will be notified about
	// any damaging failures e.g. process failures, this is optional parameter.
	// If not set email sending capabilities are not enabled.
	AdminMail *string `env:"GENERATOR_ADMIN_MAIL" flag:"adminMail"`
	// ServerMail email address used as FROM in the email this is optional parameter.
	// if not set email sending capabilities are not enabled.
	ServerMail *string `env:"GENERATOR_SENDER_MAIL" flag:"serverMail"`

	// UiLocation sets the location of generated web assets, if this variable is set, the server
	// delivers the UI in the requests to the / endpoint, if this variable is not set,
	// web assets are not delivered and requests to / will result in 404 response.
	UiLocation *string `env:"GENERATOR_UI_DIR"`
}

type logLevelValue struct {
	containingLevel *log.Level
}

type BaseHourValue int

func (h *BaseHourValue) String() string {
	return fmt.Sprintf("%d", int(*h))
}

func (h *BaseHourValue) Set(s string) error {
	hour, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if hour < 0 || hour >= 24 {
		return ErrInvalidHour
	}

	*h = BaseHourValue(hour)
	return nil
}

func (h BaseHourValue) Value() int {
	return int(h)
}

func (i *logLevelValue) String() string {
	if i.containingLevel == nil {
		return log.InfoLevel.String()
	}
	return i.containingLevel.String()
}

func (i *logLevelValue) Set(s string) error {
	newL, err := log.ParseLevel(s)
	if err != nil {
		return err
	}
	*i.containingLevel = newL
	return nil
}

var (
	prodConf Provider = Provider{
		MaxNodes:            100,
		MaxBatchSize:        50,
		Workers:             4,
		MaintenanceInterval: 24 * time.Hour,
		MaintenanceHour:     2,
		RequestTTL:          1 * time.Hour,
		DbRoot:              "/var/badger/db",
		LogLevel:            log.InfoLevel,
		UiLocation:          nil,
		Port:                "8080",
		BindAddr:            "0.0.0.0",
		SecureMode:          true,
	}

	develUiLocation string = "./ui/dist/graphGenerator"

	develConf Provider = Provider{
		MaxNodes:            1000,
		MaxBatchSize:        1000,
		Workers:             20,
		DbRoot:              "./tmp/database",
		MaintenanceInterval: 10 * time.Minute,
		MaintenanceHour:     BaseHourValue(time.Now().Add(-1 * time.Hour).Hour()),
		RequestTTL:          10 * time.Minute,
		LogLevel:            log.TraceLevel,
		UiLocation:          &develUiLocation,
		Port:                "8080",
		BindAddr:            "localhost",
		Host:                "localhost",
		SecureMode:          false,
	}

	testingConf = Provider{
		DbRoot:              "./tmp/test",
		MaxNodes:            100,
		MaxBatchSize:        50,
		Workers:             4,
		MaintenanceInterval: 60 * time.Second,
		RequestTTL:          30 * time.Second,
		LogLevel:            log.TraceLevel,
		UiLocation:          &develUiLocation,
		Port:                "8080",
		BindAddr:            "localhost",
		Host:                "localhost",
		SecureMode:          true,
	}

	cmdLineConf Provider = Provider{}

	passedData map[string]bool = map[string]bool{}

	defaultProvider *Provider = nil
	initialser      sync.Once
	develFlag       *bool
	testFlag        *bool
	flagSet         *flag.FlagSet = nil
)

func ParseFlags(set *flag.FlagSet) {

	flagSet = set
	develFlag = set.Bool("devel", false, "Start server in devel mode")
	testFlag = set.Bool("test", false, "Start server in testing mode")
	levelContainer := &logLevelValue{containingLevel: &cmdLineConf.LogLevel}

	set.StringVar(&cmdLineConf.Host, "host", "", "Set host to be expected by the web server (in case of proxy)")
	set.StringVar(&cmdLineConf.Port, "port", "", "Set port for the host to bind")
	set.StringVar(&cmdLineConf.BindAddr, "bindAddr", "", "sets interface to bind to")
	set.StringVar(&cmdLineConf.DbRoot, "dbRoot", "", "set root for database files")
	set.IntVar(&cmdLineConf.Workers, "workers", 0, "number of workers allocated for generation")
	set.IntVar(&cmdLineConf.MaxNodes, "maxNodes", 0, "maximal number of nodes in generated graphs")
	set.IntVar(&cmdLineConf.MaxBatchSize, "maxBatchSize", 0, "maximal number of graphs in batch request")
	set.DurationVar(&cmdLineConf.MaintenanceInterval, "maintInterval", 15*time.Minute, "Set clean interval")
	set.DurationVar(&cmdLineConf.RequestTTL, "ttl", 15*time.Minute, "Set timespan of requests")
	set.BoolVar(&cmdLineConf.SecureMode, "secure", true, "Set if cookie mode is secure")
	set.Var(levelContainer, "logLevel", "Set log level for message")
	set.Var(&cmdLineConf.MaintenanceHour, "maintHour", "Set the base hour for the cleaner, must be a whole number between <0-23>")
}

func Default() Provider {
	if defaultProvider == nil {
		initialser.Do(func() {
			provider := prodConf
			if develFlag != nil && *develFlag {
				provider = develConf
			}
			if testFlag != nil && *testFlag {
				provider = testingConf
			}
			defaultProvider = &provider
			defaultProvider.fillFromEnv()
			defaultProvider.fillFromCmdLine()
		})
	}
	return *defaultProvider
}

func parseValue(in string, p reflect.Type) reflect.Value {
	switch p {
	case reflect.TypeOf(""):
		return reflect.ValueOf(in)
	case reflect.TypeOf(0):
		v, err := strconv.Atoi(in)
		if err != nil {
			return reflect.ValueOf(ErrInvalidParse)
		}
		return reflect.ValueOf(v)
	case reflect.TypeOf(time.Second):
		d, err := time.ParseDuration(in)
		if err != nil {
			return reflect.ValueOf(ErrInvalidParse)
		}
		return reflect.ValueOf(d)
	case reflect.TypeOf(prodConf.MaintenanceHour):
		val := BaseHourValue(0)
		err := val.Set(in)
		if err != nil {
			return reflect.ValueOf(ErrInvalidParse)
		}
		return reflect.ValueOf(val)
	case reflect.TypeOf(log.TraceLevel):
		l, err := log.ParseLevel(in)
		if err != nil {
			return reflect.ValueOf(ErrInvalidParse)
		}
		return reflect.ValueOf(l)
	case reflect.TypeOf(true):
		v, err := strconv.ParseBool(in)
		if err != nil {
			return reflect.ValueOf(ErrInvalidParse)
		}
		return reflect.ValueOf(v)
	case reflect.TypeOf(&in):
		var res string = in
		post := &res
		return reflect.ValueOf(post)
	}
	return reflect.ValueOf(ErrInvalidParse)
}

func (p *Provider) fillFrom(filler func(v reflect.StructField, s reflect.Value)) {
	obj := reflect.ValueOf(p)

	s := obj.Elem()

	t := s.Type()
	fields := reflect.VisibleFields(t)

	for _, v := range fields {
		filler(v, s)
	}
}

func fromEnv(v reflect.StructField, s reflect.Value) {
	envVar := v.Tag.Get("env")
	if e, ok := os.LookupEnv(envVar); ok {
		parsed := parseValue(e, v.Type)
		if parsed.Type() == reflect.TypeOf(ErrInvalidParse) {
			log.Panicf("Can't parse value for env var %s %s", envVar, e)
		}
		s.FieldByName(v.Name).Set(parsed)
	}
}

func fromCmdLine(v reflect.StructField, s reflect.Value) {
	parent := reflect.ValueOf(cmdLineConf)
	envVar := v.Tag.Get("flag")
	if e, ok := passedData[envVar]; ok && e {
		s.FieldByName(v.Name).Set(parent.FieldByName(v.Name))
	}
}

func (p *Provider) fillFromEnv() {
	p.fillFrom(fromEnv)
}

func (p *Provider) fillFromCmdLine() {
	if flagSet == nil {
		return
	}
	flagSet.Visit(func(f *flag.Flag) {
		passedData[f.Name] = true
	})
	p.fillFrom(fromCmdLine)
}

func SetupTestingEnv() {
	initialser.Do(func() {
		defaultProvider = &testingConf
	})
}

func SetTestingDBRoot(value string) {
	testingConf.DbRoot = value
}

func SetCookieInsecure() {
	testingConf.SecureMode = false
}
