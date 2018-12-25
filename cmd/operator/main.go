package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/zdq0394/redis-cluster-operator/cmd/operator/cluster"
)

func main() {
	app := cli.NewApp()
	app.Name = "redisops"
	app.Description = "Redis Cluster Operator manages the creation/update/deletion of Redis Cluster."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "DqZhang",
			Email: "zdq123.hn@163.com",
		},
	}
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		cluster.Command(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
