package agentsmcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/whilesmartgo/agents"
)

func TestStreamableHTTPHandlerServesTools(t *testing.T) {
	// The registry is chosen per request, so a host can scope it to the caller.
	handler := StreamableHTTPHandler("test", "v0", func(*http.Request) *agents.Registry {
		return agents.NewRegistry(echoTool())
	})
	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	cs, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "echo", Arguments: map[string]any{"text": "hi"}})
	if err != nil {
		t.Fatal(err)
	}
	if text := res.Content[0].(*mcp.TextContent).Text; text != "echo:hi" {
		t.Errorf("tool output = %q", text)
	}
}

func TestStreamableHTTPHandlerNilRegistryRejected(t *testing.T) {
	handler := StreamableHTTPHandler("test", "v0", func(*http.Request) *agents.Registry { return nil })
	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	if _, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: httpServer.URL}, nil); err == nil {
		t.Fatal("connecting with no registry should fail")
	}
}
