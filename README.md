# websocket-reuse-connection-repro

Reproduction code for Websocket reused connection issue

## Steps to reproduce

1. Clone this repo
2. Change directory to `assets`
3. Run `pnpm install` to install dependencies for the frontend
4. Run `pnpm dev` to start the frontend dev server
5. Change directory to `..`
6. Run `go run cmd/ws/main.go ./cmd/ws/gin.go ./cmd/ws/server.go ./cmd/ws/controller.go` to start the backend server
7. Open a new tab to access `localhost:3333` in your browser
8. Type in a random number of room id and click `Go`
9. Repeat the steps 5 and 6 for 2 to 4 times with different room id
10. You will notice the `Chat Room: <some number>` mismatched the `Room <some number>` down below

## Steps to fix

1. Open `cmd/ws/controller.go`
2. Uncomment out the line 27 `// roomIDStr := ctx.Param("roomId")` and comment out the line 42 `roomIDStr := ctx.Param("roomId")`
3. Run `go run cmd/ws/main.go ./cmd/ws/gin.go ./cmd/ws/server.go ./cmd/ws/controller.go` to start the backend server
4. Open a new tab to access `localhost:3333` in your browser
5. Type in a random number of room id and click `Go`
6. Repeat the steps 5 and 6 for 2 to 4 times with different room id
7. You will notice the `Chat Room: <some number>` will never mismatched the `Room <some number>` down below

## Issue demo

![Figure1](https://github.com/nekomeowww/websocket-reuse-connection-repro/blob/main/images/screenshot-01.png)
