package definition

import (
	"fmt"
	"path"
)

type PackageNamespaceManager struct {
	// packages is a map for imported packages
	// The key is the import path.
	// The value is an alias.
	packages map[string]string
}

func NewPackageNamespaceManager() *PackageNamespaceManager {
	return &PackageNamespaceManager{packages: make(map[string]string)}
}

// Add will manages new package namespace and returns the alias for importPath.
// packageName argument is an optional.
func (m *PackageNamespaceManager) Add(importPath, packageName string) string {
	if importPath == "" {
		return ""
	}
	if v, ok := m.packages[importPath]; ok {
		return v
	}
	if packageName == "" {
		_, packageName = path.Split(importPath)
	}

	return m.add(importPath, packageName)
}

func (m *PackageNamespaceManager) All() map[string]string {
	n := make(map[string]string)
	for k, v := range m.packages {
		n[k] = v
	}
	return n
}

func (m *PackageNamespaceManager) add(importPath, packageName string) string {
	if m.isNotUsedPackageName(packageName) {
		m.packages[importPath] = packageName
		return packageName
	} else {
		i := 1
		for {
			if m.isNotUsedPackageName(fmt.Sprintf("%s_%d", packageName, i)) {
				m.packages[importPath] = fmt.Sprintf("%s_%d", packageName, i)
				return fmt.Sprintf("%s_%d", packageName, i)
			}
			i++
		}
	}
}

func (m *PackageNamespaceManager) isNotUsedPackageName(packageName string) bool {
	for _, v := range m.packages {
		if v == packageName {
			return false
		}
	}
	return true
}
