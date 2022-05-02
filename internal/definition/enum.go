package definition

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type Enum struct {
	// Name is a fully qualified enum name.
	Name string
	// ShortName is a name of enum
	ShortName string
	Values    []string
	Package   ImportPackage
}

func NewEnum(f *descriptorpb.FileDescriptorProto, enum *descriptorpb.EnumDescriptorProto) *Enum {
	var values []string
	prefix := stringsutil.ToUpperSnakeCase(enum.GetName()) + "_"
	for _, v := range enum.Value {
		e := proto.GetExtension(v.GetOptions(), kubeproto.E_Value)
		if ext := e.(*kubeproto.EnumValue); ext != nil {
			values = append(values, ext.Value)
		} else {
			values = append(values, stringsutil.ToUpperCamelCase(strings.TrimPrefix(v.GetName(), prefix)))
		}
	}

	return &Enum{
		Name:      fmt.Sprintf(".%s.%s", f.GetPackage(), enum.GetName()),
		ShortName: enum.GetName(),
		Values:    values,
		Package: ImportPackage{
			Name: path.Base(f.GetOptions().GetGoPackage()),
			Path: f.GetOptions().GetGoPackage(),
		},
	}
}

type Enums []*Enum

func (e Enums) Find(name string) *Enum {
	for _, v := range e {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (e *Enums) Own() Enums {
	m := make(map[string]*Enum)
	for _, v := range *e {
		m[v.ShortName] = v
	}

	var own []*Enum
	for _, v := range m {
		own = append(own, v)
	}

	sort.Slice(own, func(i, j int) bool {
		return own[i].ShortName < own[j].ShortName
	})
	return own
}
