package cluster

import (
	"github.com/urfave/cli"
	"github.com/zdq0394/redis-cluster-operator/operator/rediscluster"
)

// Flags of sub command `cluster`
var Flags []cli.Flag

func init() {
	Flags = []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEVELOP",
			Name:   "develop",
			Usage:  "start the operator in develop mode",
		},
		cli.StringFlag{
			EnvVar: "KUBECONFIG",
			Name:   "kubeconfig",
			Usage:  "kubeconfig of the kubernetes cluster",
		},
		cli.StringFlag{
			EnvVar: "BOOTIMG",
			Name:   "bootimg",
			Value:  "zdq0394/redis-cluster-boot:1.2",
			Usage:  "Redis Cluster Boot Image",
		},
		cli.StringFlag{
			EnvVar: "ClusterDomain",
			Name:   "clusterdomain",
			Value:  "cluster.local",
			Usage:  "Kubernetes cluster domain: e.g. cluster.local",
		},
		cli.IntFlag{
			EnvVar: "ConcurrentWorkers",
			Name:   "concurrentworkers",
			Value:  3,
			Usage:  "Concurrent goroutines to process crd management",
		},
	}
}

// Action of sub command `cluster`
func Action(ctx *cli.Context) {
	conf := rediscluster.Config{}
	conf.Development = ctx.Bool("develop")
	conf.Kubeconfig = ctx.String("kubeconfig")
	conf.BootImg = ctx.String("bootimg")
	conf.ClusterDomain = ctx.String("clusterdomain")
	conf.ConcurrentWorkers = ctx.Int("concurrentworkers")

	rediscluster.Start(&conf)
}

// Command Cluster Sub Command
func Command() cli.Command {
	return cli.Command{
		Name:    "cluster",
		Aliases: []string{"c"},
		Usage:   "start redis cluster operator",
		Flags:   Flags,
		Action:  Action,
	}
}
