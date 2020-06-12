package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xiaoenai/xmodel/cmd/create"
	"github.com/xiaoenai/xmodel/cmd/info"
)

func main() {
	app := cli.NewApp()
	app.Name = "XModel"
	app.Version = "v1.0.0"
	app.Author = "xiaoenai"
	app.Usage = "a deployment tools of xmodel frameware"

	// new a project
	newCom := cli.Command{
		Name:  "gen",
		Usage: "Generate a xmodel code",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
			cli.BoolFlag{
				Name:  "force, f",
				Usage: "Forced to rebuild the whole project",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.CreateProject(c.Bool("force"))
			return nil
		},
	}

	// add mysql model struct code to project template
	tplCom := cli.Command{
		Name:  "tpl",
		Usage: "Add mysql model struct code to project template",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "app_path, p",
				Usage: "The path(relative/absolute) of the project",
			},
			cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "mysql host ip",
			},
			cli.StringFlag{
				Name:  "port",
				Value: "3306",
				Usage: "mysql host port",
			},
			cli.StringFlag{
				Name:  "username, user",
				Value: "root",
				Usage: "mysql username",
			},
			cli.StringFlag{
				Name:  "password, pwd",
				Value: "",
				Usage: "mysql password",
			},
			cli.StringFlag{
				Name:  "db",
				Value: "test",
				Usage: "mysql database",
			},
			cli.StringSliceFlag{
				Name:  "table",
				Usage: "mysql table",
			},
			cli.StringFlag{
				Name:  "ssh_user",
				Value: "",
				Usage: "ssh user",
			},
			cli.StringFlag{
				Name:  "ssh_host",
				Value: "",
				Usage: "ssh host ip",
			},
			cli.StringFlag{
				Name:  "ssh_port",
				Value: "",
				Usage: "ssh host port",
			},
		},
		Before: initProject,
		Action: func(c *cli.Context) error {
			create.AddTableStructToTpl(create.ConnConfig{
				MysqlConfig: create.MysqlConfig{
					Host:     c.String("host"),
					Port:     c.String("port"),
					User:     c.String("user"),
					Password: c.String("password"),
					Db:       c.String("db"),
				},
				Tables:  c.StringSlice("table"),
				SshHost: c.String("ssh_host"),
				SshPort: c.String("ssh_port"),
				SshUser: c.String("ssh_user"),
			})
			return nil
		},
	}

	app.Commands = []cli.Command{newCom, tplCom}
	app.Run(os.Args)
}

func initProject(c *cli.Context) error {
	appPath := c.String("app_path")
	if len(appPath) == 0 {
		appPath = c.Args().First()
	}
	if len(appPath) == 0 {
		appPath = "./"
	}
	return info.Init(appPath)
}
