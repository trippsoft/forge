// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscoveryServer struct {
	UnimplementedDiscoveryPluginServer
}

func (s *DiscoveryServer) DiscoverInfo(
	ctx context.Context,
	request *DiscoverInfoRequest,
) (*DiscoverInfoResponse, error) {

	osInfo, err := discoverOSInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fipsInfo, err := discoverFIPSInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	appArmorInfo, err := discoverAppArmorInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	seLinuxInfo, err := discoverSELinuxInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	packageManagerInfo, err := discoverPackageManagerInfo(osInfo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	serviceManagerInfo, err := discoverServiceManagerInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &DiscoverInfoResponse{
		Os:             osInfo,
		Fips:           fipsInfo,
		AppArmor:       appArmorInfo,
		Selinux:        seLinuxInfo,
		PackageManager: packageManagerInfo,
		ServiceManager: serviceManagerInfo,
	}, nil
}
