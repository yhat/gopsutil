// +build windows

package host

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/StackExchange/wmi"

	"github.com/yhat/gopsutil/internal/common"
	process "github.com/yhat/gopsutil/process"
)

var (
	procGetSystemTimeAsFileTime = common.Modkernel32.NewProc("GetSystemTimeAsFileTime")
	osInfo                      *Win32_OperatingSystem
)

type Win32_OperatingSystem struct {
	Version        string
	Caption        string
	ProductType    uint32
	BuildNumber    string
	LastBootUpTime time.Time
}

func HostInfo() (*HostInfoStat, error) {
	ret := &HostInfoStat{
		OS: runtime.GOOS,
	}

	hostname, err := os.Hostname()
	if err == nil {
		ret.Hostname = hostname
	}

	platform, family, version, err := GetPlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = version
	} else {
		return ret, err
	}

	boot, err := BootTime()
	if err == nil {
		ret.BootTime = boot
		ret.Uptime = uptime(boot)
	}

	procs, err := process.Pids()
	if err != nil {
		return ret, err
	}

	ret.Procs = uint64(len(procs))

	return ret, nil
}

func GetOSInfo() (Win32_OperatingSystem, error) {
	var dst []Win32_OperatingSystem
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return Win32_OperatingSystem{}, err
	}

	osInfo = &dst[0]

	return dst[0], nil
}

func BootTime() (uint64, error) {
	if osInfo == nil {
		_, err := GetOSInfo()
		if err != nil {
			return 0, err
		}
	}
	now := time.Now()
	t := osInfo.LastBootUpTime.Local()
	return uint64(now.Sub(t).Seconds()), nil
}

func uptime(boot uint64) uint64 {
	return uint64(time.Now().Unix()) - boot
}

func Uptime() (uint64, error) {
	boot, err := BootTime()
	if err != nil {
		return 0, err
	}
	return uptime(boot), nil
}

func GetPlatformInformation() (platform string, family string, version string, err error) {
	if osInfo == nil {
		_, err = GetOSInfo()
		if err != nil {
			return
		}
	}

	// Platform
	platform = strings.Trim(osInfo.Caption, " ")

	// PlatformFamily
	switch osInfo.ProductType {
	case 1:
		family = "Standalone Workstation"
	case 2:
		family = "Server (Domain Controller)"
	case 3:
		family = "Server"
	}

	// Platform Version
	version = fmt.Sprintf("%s Build %s", osInfo.Version, osInfo.BuildNumber)

	return
}

func Users() ([]UserStat, error) {

	var ret []UserStat

	return ret, nil
}
