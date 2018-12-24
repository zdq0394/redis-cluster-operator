package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/zdq0394/redis-cluster-operator/operator"
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
		clusterCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func clusterCommand() cli.Command {
	return cli.Command{
		Name:    "cluster",
		Aliases: []string{"c"},
		Usage:   "start redis cluster operator",
		Flags: []cli.Flag{
			cli.BoolFlag{
				EnvVar: "DEVELOP",
				Name:   "develop",
				Usage:  "start the operator in develop mode",
			},
		},
		Action: clusterAction,
	}
}

func clusterAction(ctx *cli.Context) {
	operator.Start(false)
}
