package main

import (
	"context"
	"github.com/aserto-dev/go-aserto/client"
	"github.com/aserto-dev/go-aserto/middleware"
	mh "github.com/aserto-dev/go-aserto/middleware/http/macaron"
	authz "github.com/aserto-dev/go-authorizer/aserto/authorizer/v2"
	dir_apis "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dir_reader "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dir_writer "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gopkg.in/macaron.v1"

	"log"
	"net/http"
	"os"
)

func main() {
	options := loadOptions()
	ctx := context.Background()

	// directory reader
	var directoryReader dir_reader.ReaderClient
	var err error
	err = retry.Do(func() error {
		// Create a directory reader client
		directoryReader, err = NewDirectoryReader(ctx, &options.directory)
		if err != nil {
			log.Println("Retry: Failed to create directory reader client:", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to create directory reader client:", err)
	}

	// directory writer
	var directoryWriter dir_writer.WriterClient
	err = retry.Do(func() error {
		directoryWriter, err = NewDirectoryWriter(ctx, &options.directory)
		if err != nil {
			log.Println("Retry: Failed to create directory writer client", err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal("Failed to create directory writer client:", err)
	}

	// authorizer client
	var authorizerClient authz.AuthorizerClient
	err = retry.Do(func() error {
		// Create an authorizer client
		authorizerClient, err = NewAuthorizerClient(ctx, &options.authorizer)
		if err != nil {
			log.Println("Retry: Failed to create authorizer client:", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to create authorizer client:", err)
	}

	if err := createDummyObjectsAndRelations(authorizerClient, directoryReader, directoryWriter); err != nil {
		panic(err)
	}

	mw := mh.New(
		authorizerClient,
		middleware.Policy{
			Decision: "enable",
			Path:     "policies.bytebuilders.GET.users",
		})

	m := macaron.Classic()

	m.Use(mw.Handler())
	m.Get("/", myHandler)

	log.Println("Server is running...")
	log.Println(http.ListenAndServe("0.0.0.0:4000", m))
}

func myHandler(ctx *macaron.Context) string {
	return "the request path is: " + ctx.Req.RequestURI
}

func createDummyObjectsAndRelations(client authz.AuthorizerClient, reader dir_reader.ReaderClient, writer dir_writer.WriterClient) error {
	// User object
	newUserObject := &dir_apis.Object{
		Key:         "abc@gmail.com",
		Type:        "user",
		DisplayName: "abc",
		Properties: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"email": {
					Kind: &structpb.Value_StringValue{
						StringValue: "abc@gmail.com",
					},
				},
				"isAdmin": {
					Kind: &structpb.Value_BoolValue{
						BoolValue: false,
					},
				},
			},
		},
	}

	// identity object
	newIdentityObject := &dir_apis.Object{
		Key:  "abc@gmail.com",
		Type: "identity",
		Properties: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"kind": {
					Kind: &structpb.Value_StringValue{
						StringValue: "IDENTITY_TYPE_EMAIL",
					},
				},
				"provider": {
					Kind: &structpb.Value_StringValue{
						StringValue: "local",
					},
				},
				"verified": {
					Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					},
				},
			},
		},
	}

	// relation object

	newRelationObject := &dir_apis.Relation{
		Subject: &dir_apis.ObjectIdentifier{
			Type: newTypeStringAddr("user"),
			Key:  newTypeStringAddr("abc@gmail.com"),
		},
		Relation: "identifier",
		Object: &dir_apis.ObjectIdentifier{
			Type: newTypeStringAddr("identity"),
			Key:  newTypeStringAddr("abc@gmail.com"),
		},
	}

	_, err := writer.SetObject(context.Background(), &dir_writer.SetObjectRequest{
		Object: newUserObject,
	})

	if err != nil {
		return err
	}

	_, err = writer.SetObject(context.Background(), &dir_writer.SetObjectRequest{
		Object: newIdentityObject,
	})

	if err != nil {
		return err
	}

	_, err = writer.SetRelation(context.Background(), &dir_writer.SetRelationRequest{
		Relation: newRelationObject,
	})

	if err != nil {
		return err
	}

	return nil
}

func newTypeStringAddr(x string) *string {
	return &x
}

type options struct {
	authorizer client.Config
	directory  client.Config

	policyInstanceName  string
	policyInstanceLabel string
	policyRoot          string

	jwksKeysURL string
}

func loadOptions() *options {
	if envFileError := godotenv.Load(); envFileError != nil {
		log.Fatal("Error loading .env file")
	}

	authorizerAddr := os.Getenv("ASERTO_AUTHORIZER_SERVICE_URL")
	if authorizerAddr == "" {
		authorizerAddr = "authorizer.prod.aserto.com:8443"
	}

	directoryAddr := os.Getenv("ASERTO_DIRECTORY_SERVICE_URL")
	if directoryAddr == "" {
		directoryAddr = "directory.prod.aserto.com:8443"
	}

	log.Printf("Authorizer: %s\n", authorizerAddr)
	log.Printf("Directory:  %s\n", directoryAddr)

	return &options{
		authorizer: client.Config{
			Address:    authorizerAddr,
			APIKey:     os.Getenv("ASERTO_AUTHORIZER_API_KEY"),
			CACertPath: os.ExpandEnv(os.Getenv("ASERTO_AUTHORIZER_CERT_PATH")),
			TenantID:   os.Getenv("ASERTO_TENANT_ID"),
		},
		directory: client.Config{
			Address:    directoryAddr,
			APIKey:     os.Getenv("ASERTO_DIRECTORY_API_KEY"),
			CACertPath: os.ExpandEnv(os.Getenv("ASERTO_DIRECTORY_GRPC_CERT_PATH")),
			TenantID:   os.Getenv("ASERTO_TENANT_ID"),
		},
		jwksKeysURL:         os.Getenv("JWKS_URI"),
		policyInstanceName:  os.Getenv("ASERTO_POLICY_INSTANCE_NAME"),
		policyInstanceLabel: os.Getenv("ASERTO_POLICY_INSTANCE_LABEL"),
		policyRoot:          os.Getenv("ASERTO_POLICY_ROOT"),
	}
}

func NewDirectoryReader(ctx context.Context, cfg *client.Config) (dir_reader.ReaderClient, error) {
	conn, err := newConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return dir_reader.NewReaderClient(conn), nil
}

func NewDirectoryWriter(ctx context.Context, cfg *client.Config) (dir_writer.WriterClient, error) {
	conn, err := newConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return dir_writer.NewWriterClient(conn), nil
}

func NewAuthorizerClient(ctx context.Context, cfg *client.Config) (authz.AuthorizerClient, error) {
	conn, err := newConnection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return authz.NewAuthorizerClient(conn), nil
}

func newConnection(ctx context.Context, cfg *client.Config) (grpc.ClientConnInterface, error) {
	connectionOpts, err := cfg.ToConnectionOptions(client.NewDialOptionsProvider())
	if err != nil {
		return nil, err
	}

	conn, err := client.NewConnection(ctx, connectionOpts...)
	if err != nil {
		return nil, err
	}

	return conn.Conn, nil
}

//type UserObject struct {
//	Key         string                 `json:"key"`
//	Type        string                 `json:"type"`
//	DisplayName string                 `json:"displayName"`
//	Properties  map[string]interface{} `json:"properties"`
//}
//
//type IdentityObject struct {
//	Key        string                 `json:"key"`
//	Type       string                 `json:"type"`
//	Properties map[string]interface{} `json:"properties"`
//}
//
//type Relation struct {
//	Subject  Indicator `json:"subject"`
//	Relation string    `json:"relation"`
//	Object   Indicator `json:"object"`
//}
//
//type Indicator struct {
//	Type string `json:"type"`
//	Key  string `json:"key"`
//}
//
//type ObjectList struct {
//	Objects []UserObject `json:"objects"`
//}
//
//type RelationList struct {
//	Relations []Relation `json:"relations"`
//}
