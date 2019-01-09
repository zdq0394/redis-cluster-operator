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
	}
}

// Action of sub command `cluster`
func Action(ctx *cli.Context) {
	develop := ctx.Bool("develop")
	kubeconfig := ctx.String("kubeconfig")
	bootImage := ctx.String("bootimg")
	clusterDomain := ctx.String("clusterdomain")
	rediscluster.Start(develop, kubeconfig, bootImage, clusterDomain)
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
