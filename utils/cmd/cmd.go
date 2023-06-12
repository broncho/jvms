package cmd

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/ystyle/jvms/utils/file"
	"github.com/ystyle/jvms/utils/jdk"
	"log"
	"path/filepath"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "jvms"
	app.Usage = `JDK Version Manager (JVMS) for Windows`
	app.Version = Version

	app.CommandNotFound = func(c *cli.Context, command string) {
		log.Fatal("Command Not Found")
	}
	app.Commands = commands()
	app.Before = startup
	app.After = shutdown
	return app
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "init",
			Usage:       "Initialize config file",
			Description: `before init you should clear JAVA_HOME, PATH Environment variable.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "java_home",
					Usage: "the JAVA_HOME location",
					Value: filepath.Join(file.GetCurrentPath(), "jdk"),
				},
				cli.StringFlag{
					Name:  "jvms_home",
					Usage: "the jvms work home.",
					Value: file.GetCurrentPath(),
				},
				cli.StringFlag{
					Name:  "original_path",
					Usage: "the jdk download index file url.",
					Value: jdk.DefaultOriginalPath,
				},
			},
			Action: initAction,
		},
		{
			Name:   "config",
			Usage:  "Show config file",
			Action: configAction,
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "List current JDK installations.",
			Action:    listAction,
		},
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "Install available remote jdk",
			Action:    installAction,
		},
		{
			Name:      "add",
			ShortName: "s",
			Usage:     "Install local JDK Symlink to store location.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version,vn",
					Usage: "the local JDK version name",
				},
			},
			Action: addAction,
		},
		{
			Name:      "switch",
			ShortName: "s",
			Usage:     "Switch to use the specified version.",
			Action:    switchAction,
		},
		{
			Name:      "remove",
			ShortName: "rm",
			Usage:     "Remove a specific version.",
			Action:    removeAction,
		},
		{
			Name:      "origins",
			ShortName: "lso",
			Usage:     "Show a list of origins available for download url. ",
			Action:    originsAction,
		},
		{
			Name:      "versions",
			ShortName: "lsv",
			Usage:     "Show a list of versions available for download. ",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all,a",
					Usage: "list all the version",
				},
				cli.BoolFlag{
					Name:  "cache,c",
					Usage: "cache all the version in local.",
				},
				cli.BoolFlag{
					Name:  "force,f",
					Usage: "force fetch remote versions, not use local",
				},
			},
			Action: versionsAction,
		},
		{
			Name:  "proxy",
			Usage: "Set a proxy to use for downloads.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "show",
					Usage: "show proxy.",
				},
				cli.StringFlag{
					Name:  "set",
					Usage: "set proxy.",
				},
			},
			Action: proxyAction,
		},
	}
}

func startup(c *cli.Context) error {
	_ = InitConfig()
	return nil
}

func shutdown(c *cli.Context) error {
	err := StoreConfig()
	if err != nil {
		return errors.New("failed to save the config:" + err.Error())
	}
	return nil
}
