Roadmap discussion here: https://github.com/cometkim/mattermost-typegen/issues/1

# mattermost-typegen
Typing generator for Mattermost

## Why

This first step of a plan to provide strong typings to Mattermost world.

![](road-to-type-mattermost.png)

Mattermost core is already typed by Go, but not in JS clients and integrations.

I had started [mattermost-typed](https://github.com/cometkim/mattermost-typed) to support types also in clients and integrations, but there are problems I met:

- mattermost-redux mostly targets on browsers, but it's not fit in server-like integrations.
- It's hard to keep a huge amount of types **up-to-date**.
- Flow is effective for functional programming like redux codebase. However, community support is much better with TypeScript.
- Have to **do all this again**, If someone need types in different ways like [TypeScript](https://github.com/cometkim/mattermost-typescript) and [GraphQL](https://github.com/cometkim/matterql)

## How this works

1. Parses Go AST from [mattermost-server/model](https://github.com/mattermost/mattermost-server/tree/master/model)
2. Imports only struct type declarations and json-labeled properties.
3. Converts to valid types for a target. for example in TypeScript, `int64` in Go should be `number` in TypeScript
4. Generates type definition files using Go templates
