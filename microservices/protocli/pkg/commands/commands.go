package commands

import (
	"fmt"
	_ "os"

	"protocli/pkg/client"
)

type CommandHandler struct {
	Client *client.ProtoKeyClient
}

func NewCommandHandler(c *client.ProtoKeyClient) *CommandHandler {
	return &CommandHandler{Client: c}
}

func (h *CommandHandler) Handle(args []string) {
	if len(args) < 1 {
		h.Usage()
		return
	}
	cmd := args[0]

	switch cmd {
	case "set":
		h.handleSet(args[1:])
	case "get":
		h.handleGet(args[1:])
	case "keys":
		h.handleKeys(args[1:])
	default:
		fmt.Println("Unknown command:", cmd)
		h.Usage()
	}
}

func (h *CommandHandler) Usage() {
	fmt.Println("Usage:")
	fmt.Println("  set <key> <value>  - Set a key to a value (int)")
	fmt.Println("  get <key>          - Get value by key")
	fmt.Println("  keys <prefix>      - Get all keys with prefix")
}

func (h *CommandHandler) handleSet(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: set <key> <value>")
		return
	}
	if err := h.Client.Set(args[0], args[1]); err != nil {
		fmt.Println("Error:", err)
	}
}

func (h *CommandHandler) handleGet(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: get <key>")
		return
	}
	val, err := h.Client.Get(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(val)
}

func (h *CommandHandler) handleKeys(args []string) {
	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}
	keys, err := h.Client.Keys(prefix)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, k := range keys {
		fmt.Println(k)
	}
}
