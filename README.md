Roadmap discussion is open!: https://github.com/cometkim/mattermost-typegen/issues/1

# mattermost-typegen
Automatically generates type definitions from Mattermost core

## Why
Mattermost core is already typed by Go, but not in JS clients and integrations.

Before Mattermost client started using flow, I had started 3rd-party flow type library as [mattermost-typed](https://github.com/cometkim/mattermost-typed) to support types also in clients and integrations, but there are some problems I met:

- mattermost-redux mostly targets on browsers, but it's not fit in server-like integrations (at least I felt).
- It's hard to keep a huge amount of type definitions **up-to-date**.
- Flow is effective for functional programming like redux codebase. However, community support is much better with TypeScript.
- Have to **do all this again**, If someone need types in different ways like TypeScript, GraphQL, etc.

## How this works

1. Parses [Go AST](https://golang.org/pkg/go/ast/) from [mattermost-server/model](https://github.com/mattermost/mattermost-server/tree/master/model)
2. Imports only struct type declarations and json-tagged fields.
3. Converts to valid types for a target. for example in TypeScript, `int` in Go should be `number` in TypeScript
4. Generates type definition files using [Go template](https://golang.org/pkg/text/template/)
