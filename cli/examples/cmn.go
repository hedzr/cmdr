package examples

import (
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
)

func AttachModifyFlags(bdr cli.CommandBuilder) {
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("=").
			Titles("delim", "d", "delimiter").
			Description("delimiter char in `non-plain` mode.").
			PlaceHolder("").
			CompJustOnce(true).
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("clear", "c", "clr").
			Description("clear all tags.").
			PlaceHolder("").
			Group("Operate").
			Hidden(false, false).
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("string", "g", "string-mode").
			Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.").
			PlaceHolder("").
			ToggleGroup("Mode").
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("meta", "m", "meta-mode").
			Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).").
			PlaceHolder("").
			ToggleGroup("Mode").
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("both", "2", "both-mode").
			Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.").
			PlaceHolder("").
			ToggleGroup("Mode").
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("plain", "p", "plain-mode").
			Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.").
			ToggleGroup("Mode").
			Build()
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("tag", "t", "tag-mode").
			Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.").
			ToggleGroup("Mode").
			Build()
	})
}

func AttachConsulConnectFlags(bdr cli.CommandBuilder) {
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("localhost").
			Titles("addr", "a").
			Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
			PlaceHolder("HOST[:PORT]").
			Group("Consul")
	})

	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(8500).
			Titles("port", "p").
			Description("Consul port", ``).
			PlaceHolder("PORT").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default(false).
			Titles("insecure", "K").
			Description("Skip TLS host verification", ``).
			PlaceHolder("").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("/").
			Titles("prefix", "px").
			Description("Root key prefix", ``).
			PlaceHolder("ROOT").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").
			Titles("cacert", "").
			Description("Consul Client CA cert)", ``).
			PlaceHolder("FILE").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").
			Titles("cert", "").
			Description("Consul Client cert", ``).
			PlaceHolder("FILE").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("http").
			Titles("scheme", "").
			Description("Consul connection protocol", ``).
			PlaceHolder("SCHEME").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").
			Titles("username", "u", "user", "usr", "uid").
			Description("HTTP Basic auth user", ``).
			PlaceHolder("USERNAME").
			Group("Consul")
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").
			Titles("password", "pw", "passwd", "pass", "pwd").
			Description("HTTP Basic auth password", ``).
			PlaceHolder("PASSWORD").
			Group("Consul").
			ExternalEditor(cli.ExternalToolPasswordInput)
	})
	bdr.AddFlg(func(b cli.FlagBuilder) {
		b.Default("").
			Titles("message", "m", "msg", "mesg").
			Description("The commit message", ``).
			PlaceHolder("MSG").
			Group("Git").
			ExternalEditor(cli.ExternalToolEditor)
	})
}

func serverStartup(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverStop(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverShutdown(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverRestart(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverLiveReload(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverInstall(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverUninstall(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverStatus(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverPause(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverResume(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func kvBackup(cmd *cli.Command, args []string) (err error) {
	// err = consul.Backup()
	fmt.Printf(`
cmd: %v
args: %v

`, cmd, args)
	return
}

func kvRestore(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	// err = consul.Restore()
	return
}

func msList(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.ServiceList()
	return
}

func msTagsList(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.TagsList()
	return
}

func msTagsAdd(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsRemove(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsModify(cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsToggle(cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.TagsToggle()
	return
}
