package info

import (
	"maps"
	"strings"

	"github.com/trippsoft/forge/pkg/diag"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

type HostInfo struct {
	osInfo             *OSInfo
	selinuxInfo        *SELinuxInfo
	appArmorInfo       *AppArmorInfo
	fipsInfo           *FIPSInfo
	packageManagerInfo *PackageManagerInfo
	serviceManagerInfo *ServiceManagerInfo
	userInfo           *UserInfo
}

func NewHostInfo() *HostInfo {
	return &HostInfo{
		osInfo:             newOSInfo(),
		selinuxInfo:        newSELinuxInfo(),
		appArmorInfo:       newAppArmorInfo(),
		fipsInfo:           newFipsInfo(),
		packageManagerInfo: newPackageManagerInfo(),
		serviceManagerInfo: newServiceManagerInfo(),
		userInfo:           newUserInfo(),
	}
}

func (i *HostInfo) OSInfo() *OSInfo {
	return i.osInfo
}

func (i *HostInfo) SELinuxInfo() *SELinuxInfo {
	return i.selinuxInfo
}

func (i *HostInfo) AppArmorInfo() *AppArmorInfo {
	return i.appArmorInfo
}

func (i *HostInfo) FipsInfo() *FIPSInfo {
	return i.fipsInfo
}

func (i *HostInfo) PackageManagerInfo() *PackageManagerInfo {
	return i.packageManagerInfo
}

func (i *HostInfo) ServiceManagerInfo() *ServiceManagerInfo {
	return i.serviceManagerInfo
}

func (i *HostInfo) UserInfo() *UserInfo {
	return i.userInfo
}

func (i *HostInfo) Populate(transport transport.Transport) diag.Diags {

	if transport == nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Invalid transport",
			Detail:   "Transport cannot be nil",
		}}
	}

	diags := diag.Diags{}

	moreDiags := i.osInfo.populateOSInfo(transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.selinuxInfo.populateSelinuxInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.appArmorInfo.populateAppArmorInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.fipsInfo.populateFipsInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.packageManagerInfo.populatePackageManagerInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.serviceManagerInfo.populateServiceManagerInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	moreDiags = i.userInfo.populateUserInfo(i.osInfo, transport)
	diags = diags.AppendAll(moreDiags)

	return diags
}

func (i *HostInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	maps.Copy(values, i.osInfo.toMapOfCtyValues())
	maps.Copy(values, i.selinuxInfo.toMapOfCtyValues())
	maps.Copy(values, i.appArmorInfo.toMapOfCtyValues())
	maps.Copy(values, i.fipsInfo.toMapOfCtyValues())
	maps.Copy(values, i.packageManagerInfo.toMapOfCtyValues())
	maps.Copy(values, i.serviceManagerInfo.toMapOfCtyValues())
	maps.Copy(values, i.userInfo.toMapOfCtyValues())

	return values
}

func (i *HostInfo) String() string {
	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString(i.OSInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.SELinuxInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.AppArmorInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.FipsInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.PackageManagerInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.ServiceManagerInfo().String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.UserInfo().String())

	return stringBuilder.String()
}
