// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"os"

	"github.com/hashicorp/go-cty-funcs/cidr"
	"github.com/hashicorp/go-cty-funcs/collection"
	"github.com/hashicorp/go-cty-funcs/crypto"
	"github.com/hashicorp/go-cty-funcs/encoding"
	"github.com/hashicorp/go-cty-funcs/filesystem"
	"github.com/hashicorp/go-cty-funcs/uuid"
	"github.com/hashicorp/hcl/v2/ext/tryfunc"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	yaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func HCLFunctions() map[string]function.Function {
	workingDir, _ := os.Getwd()

	return map[string]function.Function{
		"abs":              stdlib.AbsoluteFunc,
		"abspath":          filesystem.AbsPathFunc,
		"alltrue":          AllTrueFunc,
		"anytrue":          AnyTrueFunc,
		"basename":         filesystem.BasenameFunc,
		"base64decode":     encoding.Base64DecodeFunc,
		"base64encode":     encoding.Base64EncodeFunc,
		"bcrypt":           crypto.BcryptFunc,
		"can":              tryfunc.CanFunc,
		"ceil":             stdlib.CeilFunc,
		"chomp":            stdlib.ChompFunc,
		"chunklist":        stdlib.ChunklistFunc,
		"cidrhost":         cidr.HostFunc,
		"cidrnetmask":      cidr.NetmaskFunc,
		"cidrsubnet":       cidr.SubnetFunc,
		"cidrsubnets":      cidr.SubnetsFunc,
		"coalesce":         collection.CoalesceFunc,
		"coalescelist":     stdlib.CoalesceListFunc,
		"compact":          stdlib.CompactFunc,
		"concat":           stdlib.ConcatFunc,
		"contains":         stdlib.ContainsFunc,
		"convert":          typeexpr.ConvertFunc,
		"csvdecode":        stdlib.CSVDecodeFunc,
		"dirname":          filesystem.DirnameFunc,
		"distinct":         stdlib.DistinctFunc,
		"element":          stdlib.ElementFunc,
		"endswith":         EndsWithFunc,
		"env":              EnvFunc,
		"file":             filesystem.MakeFileFunc(workingDir, false),
		"filebase64":       FileBase64Func,
		"fileexists":       filesystem.MakeFileExistsFunc(workingDir),
		"fileset":          filesystem.MakeFileSetFunc(workingDir),
		"flatten":          stdlib.FlattenFunc,
		"floor":            stdlib.FloorFunc,
		"format":           stdlib.FormatFunc,
		"formatdate":       stdlib.FormatDateFunc,
		"formatlist":       stdlib.FormatListFunc,
		"indent":           stdlib.IndentFunc,
		"index":            IndexFunc,
		"join":             stdlib.JoinFunc,
		"jsondecode":       stdlib.JSONDecodeFunc,
		"jsonencode":       stdlib.JSONEncodeFunc,
		"keys":             stdlib.KeysFunc,
		"length":           LengthFunc,
		"log":              stdlib.LogFunc,
		"lookup":           stdlib.LookupFunc,
		"lower":            stdlib.LowerFunc,
		"max":              stdlib.MaxFunc,
		"md5":              crypto.Md5Func,
		"merge":            stdlib.MergeFunc,
		"min":              stdlib.MinFunc,
		"parseint":         stdlib.ParseIntFunc,
		"pathexpand":       filesystem.PathExpandFunc,
		"pow":              stdlib.PowFunc,
		"range":            stdlib.RangeFunc,
		"reverse":          stdlib.ReverseListFunc,
		"replace":          stdlib.ReplaceFunc,
		"regex":            stdlib.RegexFunc,
		"regexall":         stdlib.RegexAllFunc,
		"regex_replace":    stdlib.RegexReplaceFunc,
		"rsadecrypt":       crypto.RsaDecryptFunc,
		"sensitive":        SensitiveFunc,
		"setintersection":  stdlib.SetIntersectionFunc,
		"setproduct":       stdlib.SetProductFunc,
		"setunion":         stdlib.SetUnionFunc,
		"sha1":             crypto.Sha1Func,
		"sha256":           crypto.Sha256Func,
		"sha512":           crypto.Sha512Func,
		"signum":           stdlib.SignumFunc,
		"slice":            stdlib.SliceFunc,
		"sort":             stdlib.SortFunc,
		"split":            stdlib.SplitFunc,
		"startswith":       StartsWithFunc,
		"strcontains":      StrContainsFunc,
		"strrev":           stdlib.ReverseFunc,
		"substr":           stdlib.SubstrFunc,
		"sum":              SumFunc,
		"textdecodebase64": TextDecodeBase64Func,
		"textencodebase64": TextEncodeBase64Func,
		"timestamp":        TimestampFunc,
		"timeadd":          stdlib.TimeAddFunc,
		"title":            stdlib.TitleFunc,
		"trim":             stdlib.TrimFunc,
		"trimprefix":       stdlib.TrimPrefixFunc,
		"trimspace":        stdlib.TrimSpaceFunc,
		"trimsuffix":       stdlib.TrimSuffixFunc,
		"try":              tryfunc.TryFunc,
		"upper":            stdlib.UpperFunc,
		"urlencode":        encoding.URLEncodeFunc,
		"uuidv4":           uuid.V4Func,
		"uuidv5":           uuid.V5Func,
		"values":           stdlib.ValuesFunc,
		"yamldecode":       yaml.YAMLDecodeFunc,
		"yamlencode":       yaml.YAMLEncodeFunc,
		"zipmap":           stdlib.ZipmapFunc,
	}
}
