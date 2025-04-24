package toolprovider

import version "github.com/hashicorp/go-version"

func ParseVersionString(v string) (*version.Version, error) {
	return version.NewVersion(v)
}
