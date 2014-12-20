package main

import (
	"errors"
	"github.com/ZurgInq/go-paste/fpaste"
	"github.com/ZurgInq/go-paste/pastebin"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
)

var errUnknownService = errors.New("unknown paste service")

func checkError(err error) {
	if err != nil {
		println("ERROR:", err.Error())
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Usage = "get and put pastes from pastebin and other paste sites."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "service, s", Value: "pastebin", Usage: "the pastebin service to use"},
	}
	app.Commands = []cli.Command{
		{
			Name:      "put",
			ShortName: "p",
			Usage:     "put a paste",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "id", Usage: "return the paste id not the url"},
				cli.StringFlag{Name: "title, t", Value: "", Usage: "the title for the paste"},
			},
			Action: put,
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a paste from its url",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "id", Usage: "get a paste from its ID instead of its URL"},
			},
			Action: get,
		},
	}
	app.Run(os.Args)
}

func put(c *cli.Context) {
	srv, err := convertService(c.GlobalString("service"))
	checkError(err)

	title := c.String("title")
	ext := c.String("ext")

	var text []byte
	if c.Args().First() == "-" || c.Args().First() == "" {
		text, err = ioutil.ReadAll(os.Stdin)
	} else {
		fileName := c.Args().First()
		if title == "" {
			title = fileName
		}
		ext = filepath.Ext(fileName)
		text, err = ioutil.ReadFile(fileName)
	}
	checkError(err)

	code, err := srv.Put(string(text), title, ext)
	checkError(err)

	if c.Bool("id") {
		println(code)
	} else {
		println(srv.WrapID(code))
	}
}

func get(c *cli.Context) {
	srv, err := convertService(c.GlobalString("service"))
	if err != nil {
		println("ERROR:", err.Error())
		os.Exit(1)
	}
	var id string
	if c.Bool("id") {
		id = c.Args().First()
	} else {
		id = srv.StripURL(c.Args().First())
	}
	text, err := srv.Get(id)
	if err != nil {
		println("ERROR:", err.Error())
		os.Exit(1)
	}
	println(text)
}

func convertService(srv string) (service, error) {
	switch {
	case srv == "pastebin" || srv == "pastebin.com" || srv == "http://pastebin.com":
		return pastebin.Pastebin{}, nil
	case srv == "fpaste" || srv == "fpaste.org" || srv == "http://fpaste.org":
		return fpaste.Fpaste{}, nil
	}
	return nil, errUnknownService
}
