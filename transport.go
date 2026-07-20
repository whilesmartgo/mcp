package agentsmcp

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/whilesmartgo/agents"
)

// StreamableHTTPHandler serves a tool set over MCP's streamable-HTTP transport.
// name and version identify the server to clients.
//
// getReg supplies the registry for each request, so a host can authenticate the
// request and hand back a caller-scoped tool set; returning nil serves 400. The
// handler is stateless: every call is served with the registry for that
// request, which is what request-scoped authorization needs. Wrapping the SDK
// here keeps the official MCP SDK out of consumer imports.
func StreamableHTTPHandler(name, version string, getReg func(*http.Request) *agents.Registry) http.Handler {
	return mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		reg := getReg(r)
		if reg == nil {
			return nil
		}
		return NewServer(name, version, reg)
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true})
}
