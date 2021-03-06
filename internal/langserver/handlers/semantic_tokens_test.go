package handlers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-ls/internal/langserver"
	"github.com/hashicorp/terraform-ls/internal/terraform/exec"
	"github.com/stretchr/testify/mock"
)

func TestSemanticTokensFull(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	var testSchema tfjson.ProviderSchemas
	err := json.Unmarshal([]byte(testModuleSchemaOutput), &testSchema)
	if err != nil {
		t.Fatal(err)
	}

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls: &exec.TerraformMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Dir(): {
					{
						Method:        "Version",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							version.Must(version.NewVersion("0.12.0")),
							nil,
							nil,
						},
					},
					{
						Method:        "GetExecPath",
						Repeatability: 1,
						ReturnArguments: []interface{}{
							"",
						},
					},
					{
						Method:        "ProviderSchemas",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							&testSchema,
							nil,
						},
					},
				},
			},
		}}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
		"capabilities": {
			"textDocument": {
				"semanticTokens": {
					"tokenTypes": [
						"type",
						"property",
						"string"
					],
					"tokenModifiers": [
						"deprecated",
						"modification"
					],
					"requests": {
						"full": true
					}
				}
			}
		},
		"rootUri": %q,
		"processId": 12345
	}`, TempDir(t).URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"test\" {\n\n}\n",
			"uri": "%s/main.tf"
		}
	}`, TempDir(t).URI())})

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "textDocument/semanticTokens/full",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"uri": "%s/main.tf"
			}
		}`, TempDir(t).URI())}, `{
			"jsonrpc": "2.0",
			"id": 3,
			"result": {
				"data": [
					0,0,8,0,0,
					0,9,6,1,2
				]
			}
		}`)
}

func TestSemanticTokensFull_clientSupportsDelta(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	var testSchema tfjson.ProviderSchemas
	err := json.Unmarshal([]byte(testModuleSchemaOutput), &testSchema)
	if err != nil {
		t.Fatal(err)
	}

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls: &exec.TerraformMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Dir(): {
					{
						Method:        "Version",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							version.Must(version.NewVersion("0.12.0")),
							nil,
							nil,
						},
					},
					{
						Method:        "GetExecPath",
						Repeatability: 1,
						ReturnArguments: []interface{}{
							"",
						},
					},
					{
						Method:        "ProviderSchemas",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							&testSchema,
							nil,
						},
					},
				},
			},
		}}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
		"capabilities": {
			"textDocument": {
				"semanticTokens": {
					"tokenTypes": [
						"type",
						"property",
						"string"
					],
					"tokenModifiers": [
						"deprecated",
						"modification"
					],
					"requests": {
						"full": {
							"delta": true
						}
					}
				}
			}
		},
		"rootUri": %q,
		"processId": 12345
	}`, TempDir(t).URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"test\" {\n\n}\n",
			"uri": "%s/main.tf"
		}
	}`, TempDir(t).URI())})

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "textDocument/semanticTokens/full",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"uri": "%s/main.tf"
			}
		}`, TempDir(t).URI())}, `{
			"jsonrpc": "2.0",
			"id": 3,
			"result": {
				"data": [
					0,0,8,0,0,
					0,9,6,1,2
				]
			}
		}`)
}

func TestVarsSemanticTokensFull(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	var testSchema tfjson.ProviderSchemas
	err := json.Unmarshal([]byte(testModuleSchemaOutput), &testSchema)
	if err != nil {
		t.Fatal(err)
	}

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls: &exec.TerraformMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Dir(): {
					{
						Method:        "Version",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							version.Must(version.NewVersion("0.12.0")),
							nil,
							nil,
						},
					},
					{
						Method:        "GetExecPath",
						Repeatability: 1,
						ReturnArguments: []interface{}{
							"",
						},
					},
					{
						Method:        "ProviderSchemas",
						Repeatability: 1,
						Arguments: []interface{}{
							mock.AnythingOfType(""),
						},
						ReturnArguments: []interface{}{
							&testSchema,
							nil,
						},
					},
				},
			},
		}}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
		"capabilities": {
			"textDocument": {
				"semanticTokens": {
					"tokenTypes": [
						"type",
						"property",
						"string"
					],
					"tokenModifiers": [
						"deprecated",
						"modification"
					],
					"requests": {
						"full": true
					}
				}
			}
		},
		"rootUri": %q,
		"processId": 12345
	}`, TempDir(t).URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "variable \"test\" {\n type=string\n}\n",
			"uri": "%s/variables.tf"
		}
	}`, tmpDir.URI())})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"version": 0,
				"languageId": "terraform-vars",
				"text": "test = \"dev\"\n",
				"uri": "%s/terraform.tfvars"
			}
	}`, tmpDir.URI())})

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "textDocument/semanticTokens/full",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"uri": "%s/terraform.tfvars"
			}
		}`, TempDir(t).URI())}, `{
			"jsonrpc": "2.0",
			"id": 4,
			"result": {
				"data": [
					0,0,4,2,0,
					0,7,5,1,0
				]
			}
		}`)
}
