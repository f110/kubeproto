package goparser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type packageInformation struct {
	GoPackage       string
	ProtobufPackage string
	ImportPath      string
}

var preDefinedPackages = []packageInformation{
	{
		GoPackage:       "k8s.io/apimachinery/pkg/runtime",
		ProtobufPackage: "k8s.io.apimachinery.pkg.runtime",
		ImportPath:      "k8s.io/apimachinery/pkg/runtime",
	},
	{
		GoPackage:       "k8s.io/apimachinery/pkg/types",
		ProtobufPackage: "k8s.io.apimachinery.pkg.types",
		ImportPath:      "k8s.io/apimachinery/pkg/types",
	},
	{
		GoPackage:       "k8s.io/apimachinery/pkg/apis/meta/v1",
		ProtobufPackage: "k8s.io.apimachinery.pkg.apis.meta.v1",
		ImportPath:      "k8s.io/apimachinery/pkg/apis/meta/v1",
	},
	{
		GoPackage:       "k8s.io/apimachinery/pkg/api/resource",
		ProtobufPackage: "k8s.io.apimachinery.pkg.api.resource",
		ImportPath:      "k8s.io/apimachinery/pkg/api/resource",
	},
	{
		GoPackage:       "k8s.io/apimachinery/pkg/util/intstr",
		ProtobufPackage: "k8s.io.apimachinery.pkg.util.intstr",
		ImportPath:      "k8s.io/apimachinery/pkg/util/intstr",
	},
	{
		GoPackage:       "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
		ProtobufPackage: "k8s.io.apiextensions_apiserver.pkg.apis.apiextensions.v1",
		ImportPath:      "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
	},
	{
		GoPackage:       "k8s.io/api/core/v1",
		ProtobufPackage: "k8s.io.api.core.v1",
		ImportPath:      "k8s.io/api/core/v1",
	},
	{
		GoPackage:       "k8s.io/api/apps/v1",
		ProtobufPackage: "k8s.io.api.apps.v1",
		ImportPath:      "k8s.io/api/apps/v1",
	},
	{
		GoPackage:       "k8s.io/api/batch/v1",
		ProtobufPackage: "k8s.io.api.batch.v1",
		ImportPath:      "k8s.io/api/batch/v1",
	},
	{
		GoPackage:       "k8s.io/api/authentication/v1",
		ProtobufPackage: "k8s.io.api.authentication.v1",
		ImportPath:      "k8s.io/api/authentication/v1",
	},
	{
		GoPackage:       "k8s.io/api/admission/v1",
		ProtobufPackage: "k8s.io.api.admission.v1",
		ImportPath:      "k8s.io/api/admission/v1",
	},
	{
		GoPackage:       "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ProtobufPackage: "sigs.k8s.io.gateway_api.apis.v1alpha2",
		ImportPath:      "sigs.k8s.io/gateway-api/apis/v1alpha2",
	},
}

// Go package to protobuf package
var packageMap = map[string]string{}

func init() {
	for _, v := range preDefinedPackages {
		packageMap[v.GoPackage] = v.ProtobufPackage
	}
}

type typeDeclaration struct {
	Name                 string
	ProtobufKind         string
	ProtobufMapKeyKind   string
	ProtobufMapValueKind string
}

func (t *typeDeclaration) ProtobufTypeDeclaration() string {
	if t.ProtobufMapKeyKind != "" && t.ProtobufMapValueKind != "" {
		return fmt.Sprintf("map<%s, %s>", t.ProtobufMapKeyKind, t.ProtobufMapValueKind)
	}
	return t.ProtobufKind
}

type Generator struct {
	allStructs bool

	protoPackage string
	goPackage    string
	apiDomain    string
	apiSubGroup  string
	apiVersion   string
	protobufFile *ProtobufFile

	// packageMap is a map for type resolving.
	// The key of map is an import alias.
	// The value of map is an import path.
	packageMap          map[string]string
	enumCandidates      map[string]struct{}
	enumValueCandidates map[string][]*enumValue

	// importedPackages is a list of imported packages
	importedPackages []*packageInformation
	importPrefix     string

	// kubeprotoGoPackage is a package name that used to used instead of go_package.
	kubeprotoGoPackage string

	typeDeclaration []*typeDeclaration
}

func New() *Generator {
	return &Generator{
		packageMap:          make(map[string]string),
		enumCandidates:      make(map[string]struct{}),
		enumValueCandidates: make(map[string][]*enumValue),
	}
}

func (g *Generator) AddDir(dir string, allStructs bool) error {
	g.allStructs = allStructs

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, dir, nil, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return err
	}

	p := &ProtobufFile{}
	for _, v := range pkgs {
		for _, f := range v.Files {
			for _, v := range f.Imports {
				g.addPackageMap(v)
			}

			for i, d := range f.Decls {
				switch {
				case isConstDeclaration(d):
					gen := d.(*ast.GenDecl)
					defs := g.constToEnumValue(gen)
					for name, values := range defs {
						g.enumValueCandidates[name] = append(g.enumValueCandidates[name], values...)
					}
				case isStringTypeDeclaration(d):
					gen := d.(*ast.GenDecl)
					enumName := g.typeToEnumName(gen)
					if enumName != "" {
						g.enumCandidates[enumName] = struct{}{}
					}
					// string type is not only used enum but also just type declaration
					g.typeDeclaration = append(g.typeDeclaration, g.parseTypeDeclaration(gen))
				case isStructDeclaration(d):
					commentGroup := &ast.CommentGroup{}
					if i != 0 {
						var prevDecl ast.Decl
						prevDecl = f.Decls[i-1]
						for _, v := range f.Comments {
							if v.Pos() > prevDecl.Pos() {
								commentGroup.List = append(commentGroup.List, v.List...)
							}
							if d.Pos() < v.Pos() {
								break
							}
						}

					}
					gen := d.(*ast.GenDecl)
					msg := g.structToProtobufMessage(gen, commentGroup, g.allStructs)
					if msg == nil {
						continue
					}

					p.Messages = append(p.Messages, msg)
				case isProtobufDefinableDeclaration(d):
					gen := d.(*ast.GenDecl)
					msg := g.protobufDefinableToMessage(gen)
					if msg == nil {
						continue
					}

					p.Messages = append(p.Messages, msg)
				case isTypeDeclaration(d):
					gen := d.(*ast.GenDecl)
					g.typeDeclaration = append(g.typeDeclaration, g.parseTypeDeclaration(gen))
				}
			}
		}
	}

	for k := range g.enumValueCandidates {
		if _, ok := g.enumCandidates[k]; !ok {
			delete(g.enumValueCandidates, k)
		}
	}

	// Type resolving
	typeMap := make(map[string]*typeDeclaration)
	for _, v := range g.typeDeclaration {
		typeMap[v.Name] = v
	}
	for _, v := range typeMap {
		if v.ProtobufMapKeyKind != "" {
			switch v.ProtobufMapKeyKind {
			case "string", "int32", "int64":
			default:
				if t, ok := typeMap[v.ProtobufMapKeyKind]; ok {
					v.ProtobufMapKeyKind = t.ProtobufKind
				}
			}
		}
	}
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if _, ok := g.enumValueCandidates[f.Kind]; ok {
				continue
			}
			if f.IsMap {
				if v, ok := typeMap[f.MapKeyKind]; ok {
					f.MapKeyKind = v.ProtobufTypeDeclaration()
				}
				if v, ok := typeMap[f.MapValueKind]; ok {
					f.MapValueKind = v.ProtobufTypeDeclaration()
				}
			}
			if v, ok := typeMap[f.Kind]; ok {
				f.Kind = v.ProtobufTypeDeclaration()
				// optional map is an invalid type
				if f.Optional {
					f.Optional = false
				}
			}
		}
	}

	g.protobufFile = p
	return nil
}

func (g *Generator) AddImport(imports ...string) {
	for _, v := range imports {
		s := strings.Split(v, ":")
		goPackage, protobufPkg, importPath := s[0], s[1], s[2]
		g.importedPackages = append(g.importedPackages, &packageInformation{
			GoPackage:       goPackage,
			ProtobufPackage: protobufPkg,
			ImportPath:      importPath,
		})
	}
}

func (g *Generator) addPackageMap(in *ast.ImportSpec) {
	importPath := in.Path.Value[1 : len(in.Path.Value)-1]
	var alias string
	if in.Name != nil {
		alias = in.Name.String()
	} else {
		alias = path.Base(importPath)
	}
	g.packageMap[alias] = importPath
}

func (g *Generator) SetProtoPackage(p string) {
	g.protoPackage = p
}

func (g *Generator) SetGoPackage(p string) {
	g.goPackage = p
}

func (g *Generator) SetAPIDomain(d, s string) {
	g.apiDomain = d
	g.apiSubGroup = s
}

func (g *Generator) SetAPIVersion(v string) {
	g.apiVersion = v
}

func (g *Generator) SetImportPrefix(v string) {
	g.importPrefix = v
}

func (g *Generator) SetKubeprotoPackage(v string) {
	g.kubeprotoGoPackage = v
}

func (g *Generator) WriteFile(outputFilePath string) error {
	if g.protobufFile == nil {
		return errors.New("not loaded any files. please call AddDir first")
	}
	g.protobufFile.Package = g.protoPackage

	sort.Slice(g.protobufFile.Messages, func(i, j int) bool {
		return g.protobufFile.Messages[i].Name < g.protobufFile.Messages[j].Name
	})

	imports := make(map[string]struct{})
	for _, m := range g.protobufFile.Messages {
		for _, f := range m.Fields {
			if f.ExternalPackage != "" {
				imports[f.ExternalPackage] = struct{}{}
			}
		}
	}
	var externalProtos []string
	for k := range imports {
		if k == "k8s.io.apimachinery.pkg.types" {
			continue
		}
		protoFile := g.resolveImportPathFromProtobufPackage(k) + "/generated.proto"
		externalProtos = append(externalProtos, protoFile)
	}
	sort.Strings(externalProtos)

	w := codegeneration.NewWriter()
	w.F("// Generated by: gen-go-to-protobuf")
	w.F("syntax = \"proto3\";")
	w.F("package %s;", g.protoPackage)
	w.F("option go_package = %q;", g.goPackage)
	if g.apiDomain != "" || g.apiSubGroup != "" || g.apiVersion != "" {
		w.F("option (dev.f110.kubeproto.k8s) = {")
		if g.apiDomain != "" {
			w.F("domain: %q,", g.apiDomain)
		}
		if g.apiSubGroup != "" {
			w.F("sub_group: %q,", g.apiSubGroup)
		}
		if g.apiVersion != "" {
			w.F("version: %q,", g.apiVersion)
		}
		w.F("};")
	}
	if g.kubeprotoGoPackage != "" {
		w.F("option (dev.f110.kubeproto.kubeproto_go_package) = %q;", g.kubeprotoGoPackage)
	}
	w.F("")
	w.F("import \"kube.proto\";")
	for _, v := range externalProtos {
		if g.isPredefinedPackage(v) {
			w.F("import %q;", v)
		} else {
			w.F("import %q;", path.Join(g.importPrefix, v))
		}
	}
	w.F("")

	var enums []string
	for name := range g.enumValueCandidates {
		enums = append(enums, name)
	}
	sort.Strings(enums)
	for _, name := range enums {
		values := g.enumValueCandidates[name]

		w.F("enum %s {", name)
		for i := 0; i < len(values); i++ {
			enumName := values[i].Value
			if enumName == "" {
				enumName = strings.TrimPrefix(values[i].Name, name)
			}
			enumName = stringsutil.ToUpperSnakeCase(enumName)
			if stringsutil.ToUpperCamelCase(enumName) == values[i].Value {
				w.F("%s_%s = %d;", stringsutil.ToUpperSnakeCase(name), enumName, i)
			} else {
				for _, v := range []string{"-", ".", "/", " "} {
					enumName = strings.Replace(enumName, v, "_", -1)
				}
				w.F(
					"%s_%s = %d [(dev.f110.kubeproto.value) = {value: %q}];",
					stringsutil.ToUpperSnakeCase(name),
					enumName,
					i,
					values[i].Value,
				)
			}
		}
		w.F("}")
		w.F("")
	}

	for _, m := range g.protobufFile.Messages {
		w.F("message %s {", m.Name)

		fields := m.Fields
		if m.UseFieldsOf != "" {
			for _, v := range g.protobufFile.Messages {
				if v.Name == m.UseFieldsOf {
					fields = v.Fields
					break
				}
			}
		}
		for _, f := range fields {
			if f.Doc != "" {
				s := bufio.NewScanner(strings.NewReader(f.Doc))
				for s.Scan() {
					t := s.Text()
					if len(t) == 0 {
						continue
					}
					if t[0] == '+' || strings.HasPrefix(t, "TODO:") {
						continue
					}
					w.F("// %s", t)
				}
			}

			if f.Optional && !f.IsMap {
				w.Fn("optional ")
			}
			if f.Repeated {
				w.Fn("repeated ")
			}
			if f.IsMap {
				if f.InvalidProtobuf {
					w.F("// This field can not be represented by protobuf.")
					w.Fn("// ")
				}
				w.Fn("map<%s, %s> %s = %d ", f.MapKeyKind, f.MapValueKind, f.Name, f.Index)
			} else {
				w.Fn("%s %s = %d ", f.Kind, f.Name, f.Index)
			}
			w.Fn("[(dev.f110.kubeproto.field) = {go_name: %q, ", f.GoName)
			if f.APIFieldName != "" {
				w.Fn("api_field_name: %q, ", f.APIFieldName)
			}
			w.Fn("inline: %v}];", f.Inline)
			w.F("")
		}

		if m.Option != nil {
			w.F("")
			w.F("option (dev.f110.kubeproto.kind) = {")
			if m.Option.ClusterScope {
				w.F("scope: SCOPE_CLUSTER")
			}
			w.F("};")
		}
		w.F("}")
		w.F("")
	}

	f, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := w.WriteTo(f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if _, err := exec.LookPath("clang-format"); err == nil {
		cmd := exec.CommandContext(context.Background(), "clang-format", "-i", outputFilePath)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) typeToEnumName(v *ast.GenDecl) string {
	typeSpec := v.Specs[0].(*ast.TypeSpec)
	return typeSpec.Name.String()
}

type enumValue struct {
	Name  string
	Value string
}

func (g *Generator) constToEnumValue(v *ast.GenDecl) map[string][]*enumValue {
	constDefinitions := make(map[string][]*enumValue)
	for _, s := range v.Specs {
		valueSpec := s.(*ast.ValueSpec)

		var name, value string
		for _, v := range valueSpec.Names {
			value += v.String()
		}
		if ident, ok := valueSpec.Type.(*ast.Ident); ok {
			name = ident.Name
		}
		if name != "" && value != "" {
			if len(valueSpec.Values) == 0 {
				continue
			}
			lit, ok := valueSpec.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			if strings.HasSuffix(value, "Prefix") {
				continue
			}
			strValue := lit.Value[1 : len(lit.Value)-1]
			constDefinitions[name] = append(constDefinitions[name], &enumValue{Name: value, Value: strValue})
		}
	}

	return constDefinitions
}

func (g *Generator) structToProtobufMessage(v *ast.GenDecl, comment *ast.CommentGroup, allStruct bool) *ProtobufMessage {
	typeSpec := v.Specs[0].(*ast.TypeSpec)
	if unicode.IsLower(rune(typeSpec.Name.String()[0])) {
		// Private struct
		return nil
	}

	if !allStruct {
		var toGenerate bool
		if v.Doc != nil {
			for _, line := range v.Doc.List {
				if strings.Contains(line.Text, "+protobuf=true") {
					toGenerate = true
				}
			}
		}
		if !toGenerate {
			return nil
		}
	}

	var useFieldsOf string
	if v.Doc != nil {
		for _, line := range v.Doc.List {
			if strings.Contains(line.Text, "+protobuf=false") {
				return nil
			}
			if strings.Contains(line.Text, "+protobuf.as=") {
				s := strings.Split(line.Text, "=")
				useFieldsOf = s[1]
			}
		}
	}

	m := &ProtobufMessage{
		Name:        typeSpec.Name.String(),
		UseFieldsOf: useFieldsOf,
	}
	i := 1
	for _, f := range typeSpec.Type.(*ast.StructType).Fields.List {
		var name string
		for _, v := range f.Names {
			name += v.String()
		}
		if name == "" {
			switch t := f.Type.(type) {
			case *ast.Ident:
				name = t.Name
			case *ast.SelectorExpr:
				name = t.Sel.Name
				// TODO: handle package
			}
		}
		if unicode.IsLower(rune(name[0])) {
			// Private field is ignored
			continue
		}

		var inline, optional bool
		var apiFieldName string
		if f.Tag != nil && f.Tag.Value != "" {
			tag := reflect.StructTag(f.Tag.Value[1 : len(f.Tag.Value)-1])

			if jsonTag, protobufTag := tag.Get("json"), tag.Get("protobuf"); jsonTag == "-" && protobufTag == "" {
				continue
			}
			if v := tag.Get("protobuf"); v == "-" {
				continue
			}

			s := strings.Split(tag.Get("json"), ",")
			if len(s) == 2 {
				if strings.Contains(s[1], "inline") {
					inline = true
				}
				if strings.Contains(s[1], "omitempty") {
					optional = true
				}
			}
			apiFieldName = s[0]
		}

		var kind, mapKeyKind, mapValueKind, externalPackage string
		var repeated, isMap, invalidProtobuf bool
		switch v := f.Type.(type) {
		case *ast.Ident:
			kind = g.goTypeToProtobufKind(v)
		case *ast.StarExpr:
			optional = true
			kind = g.goTypeToProtobufKind(v.X)
			if s, ok := v.X.(*ast.SelectorExpr); ok {
				externalPackage = g.resolveProtobufPackageFromGoPackage(s.X.(*ast.Ident))
			}
		case *ast.ArrayType:
			repeated = true
			kind = g.goTypeToProtobufKind(v.Elt)
			if s, ok := v.Elt.(*ast.SelectorExpr); ok {
				externalPackage = g.resolveProtobufPackageFromGoPackage(s.X.(*ast.Ident))
			}
		case *ast.SelectorExpr:
			kind = g.goTypeToProtobufKind(v)
			externalPackage = g.resolveProtobufPackageFromGoPackage(v.X.(*ast.Ident))
		case *ast.MapType:
			mapKeyKind, mapValueKind = g.goTypeToProtobufKind(v.Key), g.goTypeToProtobufKind(v.Value)
			isMap = true
			if mapKeyKind == "" || mapValueKind == "" {
				invalidProtobuf = true
			}
		default:
			log.Printf("%T", v)
			kind = "string"
		}

		// []byte in GO is "optional bytes" in protobuf
		if repeated && kind == "byte" {
			kind = "bytes"
			repeated = false
			optional = true
		}

		// optional repeated is an invalid type in protobuf
		if repeated && optional {
			optional = false
		}
		// optional map is an invalid type in protobuf
		if len(kind) > 4 && kind[:4] == "map<" && optional {
			optional = false
		}

		m.Fields = append(m.Fields, &ProtobufField{
			Name:            stringsutil.ToLowerSnakeCase(name),
			GoName:          name,
			APIFieldName:    apiFieldName,
			Kind:            kind,
			IsMap:           isMap,
			InvalidProtobuf: invalidProtobuf,
			MapKeyKind:      mapKeyKind,
			MapValueKind:    mapValueKind,
			ExternalPackage: externalPackage,
			Index:           i,
			Optional:        optional,
			Repeated:        repeated,
			Inline:          inline,
			Doc:             f.Doc.Text(),
		})
		i++
	}

	if m.IsRuntimeObject() {
		m.Option = &ProtobufMessageOption{}
		fields := make([]*ProtobufField, 0)
		for _, f := range m.Fields {
			switch f.Name {
			case "type_meta", "object_meta":
			default:
				fields = append(fields, f)
			}
		}
		m.Fields = fields

		for _, v := range comment.List {
			if strings.HasPrefix(v.Text, "// +genclient") {
				if strings.Contains(v.Text, "nonNamespaced") {
					m.Option.ClusterScope = true
				}
			}
		}
	}

	return m
}

func (g *Generator) goTypeToProtobufKind(in ast.Expr) string {
	switch v := in.(type) {
	case *ast.Ident:
		switch v.Name {
		case "string", "int64", "int32", "bool":
			return v.Name
		case "int":
			return "int32"
		case "float64":
			return "float"
		default:
			return v.Name
		}
	case *ast.SelectorExpr:
		nameIdent := v.X.(*ast.Ident)
		// Special type resolve for time
		if nameIdent.Name == "time" && v.Sel.Name == "Duration" {
			return "int64"
		}

		protobufPackage := g.resolveProtobufPackageFromGoPackage(nameIdent)
		// Special type resolve for types.UID.
		if protobufPackage == "k8s.io.apimachinery.pkg.types" && v.Sel.Name == "UID" {
			return "string"
		}
		if protobufPackage != "" {
			return "." + protobufPackage + "." + v.Sel.Name
		}
		return ""
	case *ast.ArrayType:
		ident, ok := v.Elt.(*ast.Ident)
		if ok {
			if ident.Name == "byte" {
				return "bytes"
			}
		}
		return ""
	case *ast.StarExpr:
		return g.goTypeToProtobufKind(v.X)
	default:
		return ""
	}
}

func (g *Generator) resolveProtobufPackageFromGoPackage(in *ast.Ident) string {
	if packageFullPath, ok := g.packageMap[in.Name]; ok {
		if protobufPackage, ok := packageMap[packageFullPath]; ok {
			return protobufPackage
		} else {
			for _, v := range g.importedPackages {
				if v.GoPackage == packageFullPath {
					return v.ProtobufPackage
				}
			}
			log.Printf("Not found protobuf package corresponding to %s", packageFullPath)
		}
	} else {
		log.Printf("Package full path not found: %s", in.Name)
	}

	return ""
}

func (g *Generator) resolveImportPathFromProtobufPackage(in string) string {
	for _, v := range preDefinedPackages {
		if v.ProtobufPackage == in {
			return v.ImportPath
		}
	}
	if importPath, ok := g.packageMap[in]; ok {
		return importPath
	}
	for _, v := range g.importedPackages {
		if v.ProtobufPackage == in {
			return v.ImportPath
		}
	}

	return ""
}

func (g *Generator) isPredefinedPackage(importPath string) bool {
	for _, v := range preDefinedPackages {
		if strings.HasPrefix(importPath, v.ImportPath) {
			return true
		}
	}

	return false
}

func (g *Generator) protobufDefinableToMessage(v *ast.GenDecl) *ProtobufMessage {
	if v.Doc != nil {
		for _, line := range v.Doc.List {
			if strings.Contains(line.Text, "+protobuf=false") {
				return nil
			}
		}
	}

	if !g.allStructs {
		var toGenerate bool
		if v.Doc != nil {
			for _, line := range v.Doc.List {
				if strings.Contains(line.Text, "+protobuf=true") {
					toGenerate = true
				}
			}
		}
		if !toGenerate {
			return nil
		}
	}

	typeSpec := v.Specs[0].(*ast.TypeSpec)
	m := &ProtobufMessage{
		Name: typeSpec.Name.String(),
	}

	switch v := typeSpec.Type.(type) {
	case *ast.ArrayType:
		kind := g.goTypeToProtobufKind(v.Elt)
		m.Fields = append(m.Fields, &ProtobufField{
			Name:     "items",
			Kind:     kind,
			Repeated: true,
			Index:    1,
		})
	default:
		return nil
	}

	return m
}

func (g *Generator) parseTypeDeclaration(v *ast.GenDecl) *typeDeclaration {
	typeSpec := v.Specs[0].(*ast.TypeSpec)

	var protobufKind, keyKind, valueKind string
	switch v := typeSpec.Type.(type) {
	case *ast.MapType:
		keyKind = g.goTypeToProtobufKind(v.Key)
		valueKind = g.goTypeToProtobufKind(v.Value)
	case *ast.Ident:
		protobufKind = g.goTypeToProtobufKind(v)
	}

	return &typeDeclaration{
		Name:                 typeSpec.Name.String(),
		ProtobufKind:         protobufKind,
		ProtobufMapKeyKind:   keyKind,
		ProtobufMapValueKind: valueKind,
	}
}

func isStructDeclaration(v ast.Decl) bool {
	gen, ok := v.(*ast.GenDecl)
	if !ok {
		return false
	}
	if len(gen.Specs) == 0 {
		return false
	}
	typeSpec, ok := gen.Specs[0].(*ast.TypeSpec)
	if !ok {
		return false
	}
	_, ok = typeSpec.Type.(*ast.StructType)
	if !ok {
		return false
	}

	return true
}

func isStringTypeDeclaration(v ast.Decl) bool {
	gen, ok := v.(*ast.GenDecl)
	if !ok {
		return false
	}
	if len(gen.Specs) == 0 {
		return false
	}
	typeSpec, ok := gen.Specs[0].(*ast.TypeSpec)
	if !ok {
		return false
	}
	ident, ok := typeSpec.Type.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name == "string" {
		return true
	}

	return false
}

func isConstDeclaration(v ast.Decl) bool {
	gen, ok := v.(*ast.GenDecl)
	if !ok {
		return false
	}
	if gen.Tok != token.CONST {
		return false
	}
	return true
}

func isProtobufDefinableDeclaration(v ast.Decl) bool {
	gen, ok := v.(*ast.GenDecl)
	if !ok {
		return false
	}
	if gen.Tok != token.TYPE {
		return false
	}
	typeSpec := gen.Specs[0].(*ast.TypeSpec)
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return false
	}
	if _, ok := typeSpec.Type.(*ast.ArrayType); ok {
		return true
	}
	return false
}

func isTypeDeclaration(v ast.Decl) bool {
	gen, ok := v.(*ast.GenDecl)
	if !ok {
		return false
	}
	if gen.Tok != token.TYPE {
		return false
	}
	typeSpec := gen.Specs[0].(*ast.TypeSpec)
	if _, ok := typeSpec.Type.(*ast.MapType); ok {
		return true
	}
	if _, ok := typeSpec.Type.(*ast.Ident); ok {
		return true
	}
	return false
}
