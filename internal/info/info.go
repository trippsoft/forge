package info

import (
	"errors"
	"maps"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type HostInfo struct {
	osInfo             *osInfo
	selinuxInfo        *selinuxInfo
	appArmorInfo       *appArmorInfo
	fipsInfo           *fipsInfo
	packageManagerInfo *packageManagerInfo
}

func NewHostInfo() *HostInfo {
	return &HostInfo{
		osInfo:             newOSInfo(),
		selinuxInfo:        newSELinuxInfo(),
		appArmorInfo:       newAppArmorInfo(),
		fipsInfo:           newFipsInfo(),
		packageManagerInfo: newPackageManagerInfo(),
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

	err := i.osInfo.populateOSInfo(transport, fileSystem)
	if err != nil {
		return err
	}

	err = i.selinuxInfo.populateSelinuxInfo(i.osInfo, fileSystem)
	if err != nil {
		return err
	}

	err = i.appArmorInfo.populateAppArmorInfo(i.osInfo, fileSystem)
	if err != nil {
		return err
	}

	err = i.fipsInfo.populateFipsInfo(i.osInfo, transport)
	if err != nil {
		return err
	}

	err = i.packageManagerInfo.populatePackageManagerInfo(i.osInfo, transport, fileSystem)
	if err != nil {
		return err
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

	return values
}
