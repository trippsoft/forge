package info

import (
	"errors"
	"maps"

	"github.com/trippsoft/forge/internal/log"
	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type HostInfo struct {
	osInfo             *osInfo
	selinuxInfo        *selinuxInfo
	appArmorInfo       *appArmorInfo
	fipsInfo           *fipsInfo
	packageManagerInfo *packageManagerInfo
	serviceManagerInfo *serviceManagerInfo
}

func NewHostInfo() *HostInfo {
	return &HostInfo{
		osInfo:             newOSInfo(),
		selinuxInfo:        newSELinuxInfo(),
		appArmorInfo:       newAppArmorInfo(),
		fipsInfo:           newFipsInfo(),
		packageManagerInfo: newPackageManagerInfo(),
		serviceManagerInfo: newServiceManagerInfo(),
	}
}

func (i *HostInfo) Populate(transport transport.Transport) error {

	if transport == nil {
		return errors.New("transport cannot be nil")
	}

	fileSystem := transport.FileSystem()

	if fileSystem == nil || fileSystem.IsNull() {
		return errors.New("file system is null or not supported")
	}

	defer fileSystem.Close()

	osInfoFailed := false
	err := i.osInfo.populateOSInfo(transport, fileSystem)
	if err != nil {
		osInfoFailed = true
		log.Errorf("failed to populate OS info: %v", err)
	}

	if !osInfoFailed {
		err = i.selinuxInfo.populateSelinuxInfo(i.osInfo, fileSystem)
		if err != nil {
			log.Errorf("failed to populate SELinux info: %v", err)
		}
	} else {
		log.Warn("SELinux info population skipped due to OS info failure")
	}

	if !osInfoFailed {
		err = i.appArmorInfo.populateAppArmorInfo(i.osInfo, fileSystem)
		if err != nil {
			log.Errorf("failed to populate AppArmor info: %v", err)
		}
	} else {
		log.Warn("AppArmor info population skipped due to OS info failure")
	}

	if !osInfoFailed {
		err = i.fipsInfo.populateFipsInfo(i.osInfo, transport)
		if err != nil {
			log.Errorf("failed to populate FIPS info: %v", err)
		}
	} else {
		log.Warn("FIPS info population skipped due to OS info failure")
	}

	if !osInfoFailed {
		err = i.packageManagerInfo.populatePackageManagerInfo(i.osInfo, transport, fileSystem)
		if err != nil {
			log.Errorf("failed to populate Package Manager info: %v", err)
		}
	} else {
		log.Warn("Package Manager info population skipped due to OS info failure")
	}

	if !osInfoFailed {
		err = i.serviceManagerInfo.populateServiceManagerInfo(i.osInfo, transport, fileSystem)
		if err != nil {
			log.Errorf("failed to populate Service Manager info: %v", err)
		}
	} else {
		log.Warn("Service Manager info population skipped due to OS info failure")
	}

	return nil
}

func (i *HostInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	maps.Copy(values, i.osInfo.toMapOfCtyValues())
	maps.Copy(values, i.selinuxInfo.toMapOfCtyValues())
	maps.Copy(values, i.appArmorInfo.toMapOfCtyValues())
	maps.Copy(values, i.fipsInfo.toMapOfCtyValues())
	maps.Copy(values, i.packageManagerInfo.toMapOfCtyValues())
	maps.Copy(values, i.serviceManagerInfo.toMapOfCtyValues())

	return values
}
