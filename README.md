# mcp

Expose a [whilesmartgo/agents](https://github.com/whilesmartgo/agents) tool set
over the [Model Context Protocol](https://modelcontextprotocol.io), built on the
official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk).

The same tools your assistant calls become an MCP server, so any MCP client (an
IDE, a desktop app, another agent) can drive them too. One tool set, two front
doors.

## Install

```sh
go get github.com/whilesmartgo/mcp
```

## Quickstart

```go
package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/whilesmartgo/agents"
	agentsmcp "github.com/whilesmartgo/mcp"
)

func main() {
	reg := agents.NewRegistry(/* your agents.Tool values */)
	server := agentsmcp.NewServer("my-app", "v1.0.0", reg)

	// Serve over any MCP transport; stdio is the common one for local clients.
	_ = server.Run(context.Background(), &mcp.StdioTransport{})
}
```

`AddRegistry(server, reg)` adds the same tools to a server you already built,
for hosts that also carry tools from other sources.

## Authorization

The bridge adds no auth of its own. A tool that returns an error yields an MCP
error result (`IsError`) rather than a protocol error, so the caller sees the
failure and can adapt. Gate access inside each `agents.Tool` handler, or with
the `agents` package's `Runner.Authorize` hook before a tool reaches the
registry.

## Status

v0, one-directional: it serves an `agents.Registry` as an MCP server. Consuming
external MCP servers as agent tools is not here yet.

## License

MIT.
