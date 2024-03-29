package k8s

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v2"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type CRDGenerator struct {
	files  protoreflect.FileDescriptor
	lister *definition.Lister
}

func NewCRDGenerator(fileToGenerate []string, files *protoregistry.Files) (*CRDGenerator, error) {
	desc, err := files.FindFileByPath(fileToGenerate[0])
	if err != nil {
		return nil, err
	}

	nsm := definition.NewPackageNamespaceManager()
	return &CRDGenerator{
		files:  desc,
		lister: definition.NewLister(fileToGenerate, files, nsm),
	}, nil
}

func (g *CRDGenerator) Generate(out io.Writer) error {
	messages := g.lister.GetMessages()
	var keys []string
	kinds := make(map[string][]*definition.Message)
	for _, m := range messages.FilterKind() {
		// If the message is virtual, it is a List object.
		// We don't need to make the manifest for it.
		if m.Virtual {
			continue
		}
		if _, ok := kinds[m.ShortName]; !ok {
			keys = append(keys, m.ShortName)
		}
		kinds[m.ShortName] = append(kinds[m.ShortName], m)
	}
	sort.Strings(keys)

	for key, msgs := range kinds {
		served := false
		storage := false
		sort.Slice(kinds[key], func(i, j int) bool {
			return kinds[key][i].Version < kinds[key][j].Version
		})
		for _, m := range msgs {
			k8sExt, err := m.Kubernetes()
			if err != nil {
				return err
			}
			if k8sExt.Served {
				served = true
			}
			if k8sExt.Storage {
				storage = true
			}
		}
		if !served {
			m := msgs[len(msgs)-1]
			k8sExt, _ := m.Kubernetes()
			k8sExt.Served = true
		}
		if !storage {
			m := msgs[len(msgs)-1]
			k8sExt, _ := m.Kubernetes()
			k8sExt.Storage = true
		}
	}

	i := 0
	for _, name := range keys {
		msgs := kinds[name]
		ext, err := msgs[0].Kubernetes()
		if err != nil {
			return err
		}

		crd := customResourceDefinition{
			APIVersion: "apiextensions.k8s.io/v1",
			Kind:       "CustomResourceDefinition",
			Metadata: metadata{
				Name: fmt.Sprintf("%s.%s.%s", strings.ToLower(stringsutil.Plural(name)), ext.SubGroup, ext.Domain)},
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Group: fmt.Sprintf("%s.%s", ext.SubGroup, ext.Domain),
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind:     name,
					ListKind: fmt.Sprintf("%sList", name),
					Plural:   strings.ToLower(stringsutil.Plural(name)),
					Singular: strings.ToLower(stringsutil.Singular(name)),
				},
				Scope: apiextensionsv1.NamespaceScoped,
			},
		}
		if msgs[0].Scope == definition.ScopeTypeCluster {
			crd.Spec.Scope = apiextensionsv1.ClusterScoped
		}
		for _, m := range msgs {
			k8sExt, err := m.Kubernetes()
			if err != nil {
				return err
			}

			var printerColumns []apiextensionsv1.CustomResourceColumnDefinition
			for _, p := range m.AdditionalPrinterColumns {
				printerColumns = append(printerColumns, apiextensionsv1.CustomResourceColumnDefinition{
					Name:        p.GetName(),
					Description: p.GetDescription(),
					JSONPath:    p.JsonPath,
					Priority:    p.Priority,
					Type:        p.Type,
					Format:      p.Format,
				})
			}

			var subResources *apiextensionsv1.CustomResourceSubresources
			for _, f := range m.Fields {
				if f.SubResource {
					if subResources == nil {
						subResources = &apiextensionsv1.CustomResourceSubresources{}
					}
					switch f.Name {
					case "Status":
						subResources.Status = &apiextensionsv1.CustomResourceSubresourceStatus{}
					case "Scale":
						subResources.Scale = &apiextensionsv1.CustomResourceSubresourceScale{}
					}
				}
			}

			schema := g.ToOpenAPISchema(m)
			ver := apiextensionsv1.CustomResourceDefinitionVersion{
				Name:                     k8sExt.Version,
				Served:                   k8sExt.Served,
				Storage:                  k8sExt.Storage,
				AdditionalPrinterColumns: printerColumns,
				Subresources:             subResources,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: schema,
				},
			}
			crd.Spec.Versions = append(crd.Spec.Versions, ver)
		}

		sort.Slice(crd.Spec.Versions, func(i, j int) bool {
			return crd.Spec.Versions[i].Name < crd.Spec.Versions[j].Name
		})

		tmp, err := json.Marshal(crd)
		if err != nil {
			return err
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(tmp, &raw); err != nil {
			return err
		}
		if i != 0 {
			io.WriteString(out, "---\n")
		}
		if err := yaml.NewEncoder(out).Encode(raw); err != nil {
			return err
		}
		i++
	}

	return nil
}

func (g *CRDGenerator) ToOpenAPISchema(m *definition.Message) *apiextensionsv1.JSONSchemaProps {
	required := make([]string, 0)
	properties := make(map[string]apiextensionsv1.JSONSchemaProps)
	for _, f := range m.Fields {
		switch f.Kind {
		case protoreflect.BoolKind, protoreflect.StringKind, protoreflect.Int64Kind, protoreflect.Int32Kind:
			properties[f.FieldName] = g.fieldToJSONSchemaProps(f)
			if !f.Optional {
				required = append(required, f.FieldName)
			}
		case protoreflect.MessageKind:
			props := g.messageToJSONSchemaProps(f)
			if f.Inline {
				for k, v := range props.Properties {
					properties[k] = v
				}
				if len(props.Required) > 0 {
					required = append(required, props.Required...)
				}
			} else {
				properties[f.FieldName] = *props
				if !f.Optional && !m.Kind && !f.IsMap() && !f.Repeated {
					required = append(required, f.FieldName)
				}
			}
		case protoreflect.EnumKind:
			enum := g.lister.GetEnums().Find(f.MessageName)
			if enum != nil {
				var values []apiextensionsv1.JSON
				for _, v := range enum.Values {
					values = append(values, apiextensionsv1.JSON{Raw: []byte(fmt.Sprintf("%q", v))})
				}
				properties[f.FieldName] = apiextensionsv1.JSONSchemaProps{
					Description: f.Description,
					Type:        "string",
					Enum:        values,
				}
				if !f.Optional {
					required = append(required, f.FieldName)
				}
			}
		}
	}
	props := &apiextensionsv1.JSONSchemaProps{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}

	return props
}

func (g *CRDGenerator) fieldToJSONSchemaProps(f *definition.Field) apiextensionsv1.JSONSchemaProps {
	props := apiextensionsv1.JSONSchemaProps{
		Description: f.Description,
	}
	props.Type = definition.ProtoreflectKindToJSONSchemaType[f.Kind]
	switch f.Kind {
	case protoreflect.Int64Kind:
		props.Format = "int64"
	}

	if f.Repeated {
		props.Description = ""
		return apiextensionsv1.JSONSchemaProps{
			Type:        "array",
			Description: f.Description,
			Items: &apiextensionsv1.JSONSchemaPropsOrArray{
				Schema: &props,
			},
		}
	}

	return props
}

func (g *CRDGenerator) messageToJSONSchemaProps(f *definition.Field) *apiextensionsv1.JSONSchemaProps {
	props := &apiextensionsv1.JSONSchemaProps{
		Description: f.Description,
	}

	if f.IsMap() {
		props.Type = "object"
		props.AdditionalProperties = &apiextensionsv1.JSONSchemaPropsOrBool{Schema: &apiextensionsv1.JSONSchemaProps{Type: "string"}}
		return props
	}

	switch f.MessageName {
	case "k8s.io.apimachinery.pkg.apis.meta.v1.Time":
		props.Type = "string"
		props.Format = "date-time"
	default:
		child := g.lister.GetMessages().Find(f.MessageName)
		props = g.ToOpenAPISchema(child)
	}

	if f.Repeated {
		props.Description = ""
		return &apiextensionsv1.JSONSchemaProps{
			Type:        "array",
			Description: f.Description,
			Items: &apiextensionsv1.JSONSchemaPropsOrArray{
				Schema: props,
			},
		}
	}

	return props
}

type customResourceDefinition struct {
	APIVersion string                                       `json:"apiVersion"`
	Kind       string                                       `json:"kind"`
	Metadata   metadata                                     `json:"metadata"`
	Spec       apiextensionsv1.CustomResourceDefinitionSpec `json:"spec"`
}

type metadata struct {
	Name string `json:"name"`
}
