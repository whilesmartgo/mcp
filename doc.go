// Package agentsmcp bridges a whilesmartgo/agents tool set to the Model Context
// Protocol, using the official MCP Go SDK for the wire protocol.
//
// The bridge is one-directional in v0: it exposes an agents.Registry as an MCP
// server, so an MCP client can call the same tools an assistant would. Each MCP
// tool call runs the corresponding agents.Tool handler. A tool that returns an
// error yields an error result (IsError) rather than a protocol error, so a
// model on the other end can see the failure and adapt.
//
//	reg := agents.NewRegistry(tool1, tool2)
//	server := agentsmcp.NewServer("my-app", "v1.0.0", reg)
//	// serve over any MCP transport, e.g. stdio:
//	_ = server.Run(ctx, &mcp.StdioTransport{})
//
// The consumer supplies authorization inside each agents.Tool handler (or via
// the Runner in the agents package); this bridge adds no auth of its own.
package agentsmcp
