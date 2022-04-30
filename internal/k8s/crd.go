package k8s

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v2"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type CRDGenerator struct {
	files  []*descriptorpb.FileDescriptorProto
	lister *definition.Lister
}

func NewCRDGenerator(files []*descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *CRDGenerator {
	return &CRDGenerator{files: files, lister: definition.NewLister(files, allProtos)}
}

func (g *CRDGenerator) Generate(out io.Writer) error {
	messages := g.lister.GetMessages()
	kinds := make(map[string][]*definition.Message)
	for _, m := range messages.FilterKind() {
		// If the message is virtual, it is a List object.
		// We don't need to make the manifest for it.
		if m.Virtual {
			continue
		}
		kinds[m.ShortName] = append(kinds[m.ShortName], m)
	}

	i := 0
	for name, msgs := range kinds {
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
					case "status":
						subResources.Status = &apiextensionsv1.CustomResourceSubresourceStatus{}
					case "scale":
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
	props := &apiextensionsv1.JSONSchemaProps{
		Type: "object",
	}
	properties := make(map[string]apiextensionsv1.JSONSchemaProps)
	for _, f := range m.Fields {
		switch f.Type {
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL,
			descriptorpb.FieldDescriptorProto_TYPE_STRING,
			descriptorpb.FieldDescriptorProto_TYPE_INT64:
			properties[f.FieldName] = g.fieldToJSONSchemaProps(f)
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			child := g.lister.GetMessages().Find(f.MessageName)
			props := g.ToOpenAPISchema(child)
			if f.Inline {
				for k, v := range props.Properties {
					properties[k] = v
				}
			} else {
				properties[f.FieldName] = *props
			}
		}
	}
	props.Properties = properties

	return props
}

func (g *CRDGenerator) fieldToJSONSchemaProps(f *definition.Field) apiextensionsv1.JSONSchemaProps {
	props := apiextensionsv1.JSONSchemaProps{
		Description: f.Description,
	}

	switch f.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		props.Type = "boolean"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		props.Type = "string"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		props.Type = "integer"
	}

	if f.Repeated {
		return apiextensionsv1.JSONSchemaProps{
			Type: "array",
			Items: &apiextensionsv1.JSONSchemaPropsOrArray{
				Schema: &props,
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
