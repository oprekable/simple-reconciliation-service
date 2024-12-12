package cconfig

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"simple-reconciliation-service/internal/app/config"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/creasty/defaults"
	"github.com/denisbrodbeck/machineid"
	godotenvFS "github.com/driftprogramming/godotenv"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	// TZ ...
	TZ string = "TZ"
)

type ConfigPaths []string
type WorkDirPath string
type AppName string
type TimeZone string
type TimeOffset int

type Config struct {
	*config.Data
	timeLocation     *time.Location
	appName          AppName
	defaultQueueName string
	workDirPath      WorkDirPath
	machineID        string
	timeZone         TimeZone
	machineIPs       []string
	timeOffset       TimeOffset
}

func initTimeZone(tzArgs TimeZone) (tz string, loc *time.Location, offset int, err error) {
	tz = os.Getenv(TZ)
	if tz == "" {
		err = os.Setenv(TZ, string(tzArgs))
		if err != nil {
			return
		}

		tz = string(tzArgs)
	}

	tzString, offset1 := time.Now().Zone()
	loc, err = time.LoadLocation(os.Getenv(TZ))
	if err != nil {
		return tzString, time.Local, offset1, nil
	}

	_, offset2 := time.Now().In(loc).Zone()

	offset = offset1

	if offset1 != offset2 {
		time.Local = loc
		offset = offset2
	}

	return tz, loc, offset, err
}

func initWorkDirPath() WorkDirPath {
	w, _ := os.UserHomeDir()
	if ex, er := os.Executable(); er == nil {
		w = filepath.Dir(ex)
	}

	return WorkDirPath(w)
}

func initMachineID() (machineID string) {
	var er error
	machineID = "localhost"
	if machineID, er = machineid.ID(); er != nil {
		if machineID, er = os.Hostname(); er != nil {
			machineID = "localhost"
		}
	}

	return machineID
}

func initMachineIP() (machineIPs []string) {
	var er error
	var netInterfaces []net.Interface
	if netInterfaces, er = net.Interfaces(); er != nil {
		return
	}

	for _, i := range netInterfaces {
		var adders []net.Addr
		var e error

		if adders, e = i.Addrs(); e != nil {
			fmt.Printf("failed to load machine IP information %v\n", e)
			continue
		}

		for _, addr := range adders {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil {
				ip4 := ip.To4()
				if ip4 != nil && ip4.String() != "127.0.0.1" {
					machineIPs = append(machineIPs, ip4.String())
				}
			}
		}
	}

	return machineIPs
}

func fromFS(embedFS *embed.FS, patterns []string, conf interface{}) (err error) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	for i := range patterns {
		matches, er := fs.Glob(embedFS, patterns[i])
		if er != nil {
			err = er
			return
		}

		for i2 := range matches {
			fileData, er := embedFS.ReadFile(matches[i2])
			if er != nil {
				err = er
				return
			}

			if err = viper.MergeConfig(bytes.NewReader(fileData)); err != nil {
				return
			}
		}
	}

	return viper.Unmarshal(conf)
}

func fromFiles(patterns []string, conf interface{}) (err error) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	for i := range patterns {
		matches, er := filepath.Glob(patterns[i])
		if er != nil {
			return er
		}

		for i2 := range matches {
			if _, err := os.Stat(matches[i2]); err == nil {
				viper.SetConfigFile(matches[i2])
				if err := viper.MergeInConfig(); err != nil {
					return err
				}
			}
		}
	}

	return viper.Unmarshal(conf)
}

func NewConfig(ctx context.Context, embedFS *embed.FS, configPaths ConfigPaths, appName AppName, tzArgs TimeZone) (rd *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %s", r)
			rd = nil
			return
		}
	}()

	rd = &Config{}
	rd.workDirPath = initWorkDirPath()
	tz, loc, offset, err := initTimeZone(tzArgs)
	if err != nil {
		return nil, err
	}

	rd.timeZone = TimeZone(tz)
	rd.timeLocation = loc
	rd.timeOffset = TimeOffset(offset)
	rd.machineID = initMachineID()
	rd.machineIPs = initMachineIP()

	var cfg config.Data
	_, er := hunch.Waterfall(
		ctx,
		// Set env from embedFS file
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			fileEnvPath := "embeds/envs/.env"
			e = godotenvFS.Load(*embedFS, fileEnvPath)
			return
		},
		// Set env from regular file
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if _, er := os.Stat(string(rd.workDirPath) + "/params/.env"); er == nil {
				e = godotenv.Overload(string(rd.workDirPath) + "/params/.env")
			}

			return
		},
		// Load config from embed files
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			cPFS := append(configPaths, "embeds/params/*.toml")
			e = fromFS(embedFS, cPFS, &cfg)
			return
		},
		// Load config from regular files
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			cP := append(configPaths, fmt.Sprintf("%s/params/*.toml", rd.workDirPath))
			e = fromFiles(cP, &cfg)
			return
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			e = defaults.Set(&cfg)
			return
		},
	)

	if er != nil {
		panic(fmt.Errorf("failed to init config %v", er))
	}

	rd.Data = &cfg
	rd.appName = appName
	rd.defaultQueueName = fmt.Sprintf("%s_%s", appName, "queue")

	return rd, nil
}

func (c *Config) GetDefaultQueueName() string {
	return c.defaultQueueName
}

func (c *Config) GetWorkDirPath() WorkDirPath {
	return c.workDirPath
}

func (c *Config) GetMachineID() string {
	return c.machineID
}

func (c *Config) GetMachineIPs() []string {
	return c.machineIPs
}

func (c *Config) GetTimeLocation() *time.Location {
	return c.timeLocation
}

func (c *Config) GetTimeOffset() TimeOffset {
	return c.timeOffset
}

func (c *Config) GetTimeZone() TimeZone {
	return c.timeZone
}

func (c *Config) GetAppName() AppName {
	return c.appName
}
