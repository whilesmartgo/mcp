package agentsmcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/whilesmartgo/agents"
)

// NewServer builds an MCP server that exposes every tool in reg. name and
// version identify the server to MCP clients.
func NewServer(name, version string, reg *agents.Registry) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil)
	AddRegistry(server, reg)
	return server
}

// AddRegistry adds every tool in reg to an existing MCP server. Use it when the
// server also carries tools from other sources.
func AddRegistry(server *mcp.Server, reg *agents.Registry) {
	for _, schema := range reg.Schemas() {
		tool, ok := reg.Get(schema.Name)
		if !ok {
			continue
		}
		server.AddTool(mcpTool(schema), handlerFor(tool))
	}
}

func mcpTool(schema agents.ToolSchema) *mcp.Tool {
	return &mcp.Tool{
		Name:        schema.Name,
		Description: schema.Description,
		// InputSchema accepts any value that marshals to a JSON Schema object;
		// agents.ToolSchema.Parameters already is one.
		InputSchema: schema.Parameters,
	}
}

// handlerFor runs an agents.Tool for an MCP call. A tool error becomes an error
// result the caller can see, not a protocol error, matching MCP's convention so
// a model can read the failure and adapt.
func handlerFor(tool agents.Tool) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		out, err := tool.Handler(ctx, req.Params.Arguments)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			}, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: out}},
		}, nil
	}
}
