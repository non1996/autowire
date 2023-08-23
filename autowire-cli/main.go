package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/non1996/go-autowire/autowire-cli/autowire"
)

var app = cli.App{
	Name:  "autowire-cli",
	Usage: "go 依赖注入工具",
	Commands: []*cli.Command{
		{
			Name:  "gen",
			Usage: "生成代码",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "module",
					Aliases:  []string{"m"},
					Required: true,
					Usage:    "指定包所属module",
				},
				&cli.StringFlag{
					Name:     "filename",
					Aliases:  []string{"f"},
					Value:    "autowire.go",
					Required: false,
					Usage:    "注解文件名称，默认为autowire.go",
				},
				&cli.StringFlag{
					Name:     "genfilename",
					Aliases:  []string{"g"},
					Value:    "autowire_gen.go",
					Required: false,
					Usage:    "生成文件名称，默认为autowire_gen.go",
				},
			},
			Action: func(context *cli.Context) error {
				autowire.Generate(
					autowire.Config{
						Module:           context.String("module"),
						Root:             context.Args().First(),
						AutowireFileName: context.String("filename"),
						GenFileName:      context.String("genfilename"),
					},
				)

				return nil
			},
		},
	},
}

func main() {
	_ = app.Run(os.Args)
}
