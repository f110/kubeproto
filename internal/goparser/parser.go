package goparser

import (
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

var packageMap = map[string]string{
	// Go package to protobuf package
	"k8s.io/apimachinery/pkg/runtime":      "k8s.io.apimachinery.pkg.runtime",
	"k8s.io/apimachinery/pkg/types":        "k8s.io.apimachinery.pkg.types",
	"k8s.io/apimachinery/pkg/apis/meta/v1": "k8s.io.apimachinery.pkg.apis.meta.v1",
	"k8s.io/apimachinery/pkg/api/resource": "k8s.io.apimachinery.pkg.api.resource",
	"k8s.io/apimachinery/pkg/util/intstr":  "k8s.io.apimachinery.pkg.util.intstr",
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
	protobufFile *ProtobufFile

	packageMap          map[string]string
	enumCandidates      map[string]struct{}
	enumValueCandidates map[string][]*enumValue

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

			for _, d := range f.Decls {
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
					gen := d.(*ast.GenDecl)
					msg := g.structToProtobufMessage(gen)
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
	enumMap := make(map[string]string)
	for _, v := range g.typeDeclaration {
		if _, ok := g.enumValueCandidates[v.Name]; ok {
			// This type is the enum
			enumMap[v.Name] = v.ProtobufKind
		}

		typeMap[v.Name] = v
	}
	for _, v := range typeMap {
		if v.ProtobufMapKeyKind != "" {
			switch v.ProtobufMapKeyKind {
			case "string", "int32":
			default:
				if t, ok := typeMap[v.ProtobufMapKeyKind]; ok {
					v.ProtobufMapKeyKind = t.ProtobufKind
				}
			}
		}
	}
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if v, ok := enumMap[f.Kind]; ok {
				f.Kind = v
				continue
			}
			if v, ok := typeMap[f.Kind]; ok {
				f.Kind = v.ProtobufTypeDeclaration()
			}
		}
	}

	g.protobufFile = p
	return nil
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

func (g *Generator) WriteFile(path string) error {
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
		protoFile := k
		if strings.HasPrefix(protoFile, "k8s.io.") {
			protoFile = "k8s.io/" + strings.Replace(strings.TrimPrefix(protoFile, "k8s.io."), ".", "/", -1) + "/generated.proto"
		}
		externalProtos = append(externalProtos, protoFile)
	}

	w := codegeneration.NewWriter()
	w.F("// Generated by: gen-go-to-protobuf")
	w.F("syntax = \"proto3\";")
	w.F("package %s;", g.protoPackage)
	w.F("option go_package = %q;", g.goPackage)
	w.F("")
	w.F("import \"kube.proto\";")
	for _, v := range externalProtos {
		w.F("import %q;", v)
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
				for _, v := range []string{"-", ".", "/"} {
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
			if f.Optional {
				w.Fn("optional ")
			}
			if f.Repeated {
				w.Fn("repeated ")
			}
			w.F(
				"%s %s = %d [(dev.f110.kubeproto.field) = {go_name: %q, api_field_name: %q, inline: %v}];",
				f.Kind, f.Name, f.Index,
				f.GoName,
				f.APIFieldName,
				f.Inline,
			)
		}
		w.F("}")
		w.F("")
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
		cmd := exec.CommandContext(context.Background(), "clang-format", "-i", path)
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

func (g *Generator) structToProtobufMessage(v *ast.GenDecl) *ProtobufMessage {
	typeSpec := v.Specs[0].(*ast.TypeSpec)
	if unicode.IsLower(rune(typeSpec.Name.String()[0])) {
		// Private struct
		return nil
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

		var inline bool
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
			if len(s) == 2 && strings.Contains(s[1], "inline") {
				inline = true
			}
			apiFieldName = s[0]
		}

		var kind, externalPackage string
		var optional, repeated bool
		switch v := f.Type.(type) {
		case *ast.Ident:
			kind = g.goTypeToProtobufKind(v)
		case *ast.StarExpr:
			optional = true
			kind = g.goTypeToProtobufKind(v.X)
		case *ast.ArrayType:
			repeated = true
			kind = g.goTypeToProtobufKind(v.Elt)
		case *ast.SelectorExpr:
			kind = g.goTypeToProtobufKind(v)
			externalPackage = g.resolveProtobufPackage(v.X.(*ast.Ident))
		case *ast.MapType:
			keyKind := g.goTypeToProtobufKind(v.Key)
			valueKind := g.goTypeToProtobufKind(v.Value)
			kind = fmt.Sprintf("map<%s, %s>", keyKind, valueKind)
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

		m.Fields = append(m.Fields, &ProtobufField{
			Name:            stringsutil.ToLowerSnakeCase(name),
			GoName:          name,
			APIFieldName:    apiFieldName,
			Kind:            kind,
			ExternalPackage: externalPackage,
			Index:           i,
			Optional:        optional,
			Repeated:        repeated,
			Inline:          inline,
		})
		i++
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
		default:
			return v.Name
		}
	case *ast.SelectorExpr:
		nameIdent := v.X.(*ast.Ident)
		// Special type resolve for time
		if nameIdent.Name == "time" && v.Sel.Name == "Duration" {
			return "int64"
		}

		protobufPackage := g.resolveProtobufPackage(nameIdent)
		// Special type resolve for types.UID.
		if protobufPackage == "k8s.io.apimachinery.pkg.types" && v.Sel.Name == "UID" {
			return "string"
		}
		if protobufPackage != "" {
			return protobufPackage + "." + v.Sel.Name
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
	default:
		return ""
	}
}

func (g *Generator) resolveProtobufPackage(in *ast.Ident) string {
	if packageFullPath, ok := g.packageMap[in.Name]; ok {
		if protobufPackage, ok := packageMap[packageFullPath]; ok {
			return protobufPackage
		} else {
			log.Printf("Not found protobuf package corresponding to %s", packageFullPath)
		}
	} else {
		log.Printf("Package full path not found: %s", in.Name)
	}

	return ""
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
	return false
}
