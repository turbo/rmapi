package shell

import (
	"path"

	"github.com/abiosoft/ishell"
)

func mvCmd(ctx *ShellCtxt) *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mv",
		Help: "mv file or directory",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				return
			}

			src := c.Args[0]

			srcNode, err := ctx.api.Filetree.NodeByPath(src, ctx.node)

			if err != nil {
				c.Println("source entry doesn't exist")
				return
			}

			dst := c.Args[1]

			dstNode, err := ctx.api.Filetree.NodeByPath(dst, ctx.node)

			if dstNode != nil && dstNode.IsFile() {
				c.Println("destination entry already exists")
				return
			}

			// We are moving the node to antoher directory
			if dstNode != nil && dstNode.IsDirectory() {
				n, err := ctx.api.MoveEntry(srcNode, dstNode, srcNode.Name())

				if err != nil {
					c.Println("failed to move entry", err)
					return
				}

				ctx.api.Filetree.MoveNode(srcNode, n)
				return
			}

			// We are renaming the node
			parentDir := path.Dir(dst)
			newEntry := path.Base(dst)

			parentNode, err := ctx.api.Filetree.NodeByPath(parentDir, ctx.node)

			if err != nil || parentNode.IsFile() {
				c.Println("directory doesn't exist")
				return
			}

			n, err := ctx.api.MoveEntry(srcNode, parentNode, newEntry)

			if err != nil {
				c.Println("failed to move entry", err)
				return
			}

			ctx.api.Filetree.MoveNode(srcNode, n)
		},
	}
}
