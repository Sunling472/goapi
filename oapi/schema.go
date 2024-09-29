package oapi

import "encoding/json"

type (
	Contact struct {
		Name  string `json:"name" default:"Sunling"`
		Url   string `json:"url"`
		Email string `json:"email"`
	}
	License struct {
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
		Url        string `json:"url"`
	}
	Info struct {
		Title          string  `json:"title" default:"Goapi App"`
		Summary        string  `json:"summary"`
		Description    string  `json:"description"`
		TermsOfService string  `json:"termsOfService"`
		Contact        Contact `json:"contact"`
		License        License `json:"license"`
		Version        string  `json:"version"`
	}

	ServerVariable struct {
		Enum        []string `json:"enum"`
		Default     string   `json:"default"`
		Description string   `json:"description"`
	}
	Server struct {
		Url         string                    `json:"url"`
		Description string                    `json:"description"`
		Variables   map[string]ServerVariable `json:"variables"`
	}

	PathItemOperation struct{}
	PathItemParam     struct{}
	PathItem          struct {
		Ref         string            `json:"$ref"`
		Summary     string            `json:"summary"`
		Description string            `json:"description"`
		Get         PathItemOperation `json:"get"`
		Put         PathItemOperation `json:"put"`
		Post        PathItemOperation `json:"post"`
		Delete      PathItemOperation `json:"delete"`
		Options     PathItemOperation `json:"options"`
		Head        PathItemOperation `json:"head"`
		Patch       PathItemOperation `json:"patch"`
		Trace       PathItemOperation `json:"trace"`
		Servers     []Server          `json:"servers"`
		Parameters  []PathItemParam   `json:"parameters"`
	}
	Path map[string]PathItem

	ContentMediaType struct {
		Schema   SchemaComponent              `json:"schema"`
		Example  ExampleComponent             `json:"example"`
		Examples map[string]ExampleComponent  `json:"examples"`
		Encoding map[string]EncodingComponent `json:"encoding"`
	}
	SchemaComponent   struct{}
	EncodingComponent struct{}
	ResponseComponent struct {
		Description string                      `json:"description"`
		Headers     map[string]HeaderComponent  `json:"headers"`
		Content     map[string]ContentMediaType `json:"content"`
		Links       map[string]LinkComponent    `json:"links"`
	}
	ExampleComponent     struct{}
	RequestBodyComponent struct{}
	HeaderComponent      struct{}
	SecurityComponent    struct{}
	LinkComponent        struct{}
	CallbackComponent    struct{}
	Components           struct {
		Schemas         map[string]SchemaComponent      `json:"schemas"`
		Responses       map[string]ResponseComponent    `json:"responses"`
		Parameters      map[string]PathItemParam        `json:"parameters"`
		Examples        map[string]ExampleComponent     `json:"examples"`
		RequestBodies   map[string]RequestBodyComponent `json:"requestBodies"`
		Headers         map[string]HeaderComponent      `json:"headers"`
		SecuritySchemes map[string]SecurityComponent    `json:"securitySchemes"`
		Links           map[string]LinkComponent        `json:"links"`
		Callbacks       map[string]CallbackComponent    `json:"callbacks"`
		PathItems       map[string]PathItem             `json:"pathItems"`
	}
	Security     struct{}
	Tag          struct{}
	ExternalDocs struct{}

	OpenAPI struct {
		OpenApi           string              `json:"openapi"` //Version
		Info              Info                `json:"info"`
		JsonSchemaDialect string              `json:"jsonSchemaDialect"`
		Servers           []Server            `json:"servers"`
		Paths             Path                `json:"paths"`
		Webhooks          map[string]PathItem `json:"webhooks"`
		Components        Components          `json:"components"`
		Security          []Security          `json:"security"`
		Tags              []Tag               `json:"tags"`
		ExternalDocs      ExternalDocs        `json:"externalDocs"`
	}
)

func (o OpenAPI) ToJson() ([]byte, error) {
	return json.Marshal(o)
}
