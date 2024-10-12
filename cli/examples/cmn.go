package examples

import (
	"context"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
)

func AttachModifyFlags(bdr cli.CommandBuilder) {
	bdr.Flg("delim", "d", "delimiter").
		Default("=").
		Description("delimiter char in `non-plain` mode.").
		PlaceHolder("").
		CompJustOnce(true).
		Build()

	bdr.Flg("clear", "c", "clr").
		Default(false).
		Description("clear all tags.").
		PlaceHolder("").
		Group("Operate").
		Hidden(false, false).
		Build()

	bdr.Flg("string", "g", "string-mode").
		Default(false).
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.").
		PlaceHolder("").
		ToggleGroup("Mode").
		Build()

	bdr.Flg("meta", "m", "meta-mode").
		Default(false).
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).").
		PlaceHolder("").
		ToggleGroup("Mode").
		Build()

	bdr.Flg("both", "2", "both-mode").
		Default(false).
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.").
		PlaceHolder("").
		ToggleGroup("Mode").
		Build()

	bdr.Flg("plain", "p", "plain-mode").
		Default(false).
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.").
		ToggleGroup("Mode").
		Build()

	bdr.Flg("tag", "t", "tag-mode").
		Default(false).
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.").
		ToggleGroup("Mode").
		Build()
}

func AttachConsulConnectFlags(bdr cli.CommandBuilder) {
	bdr.Flg("addr", "a").
		Default("localhost").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		PlaceHolder("HOST[:PORT]").
		Group("Consul").
		Build()
	bdr.Flg("port", "p").
		Default(8500).
		Description("Consul port", ``).
		PlaceHolder("PORT").
		Group("Consul").
		Build()
	bdr.Flg("insecure", "K").
		Default(false).
		Description("Skip TLS host verification", ``).
		PlaceHolder("").
		Group("Consul").
		Build()
	bdr.Flg("prefix", "px").
		Default("/").
		Description("Root key prefix", ``).
		PlaceHolder("ROOT").
		Group("Consul").
		Build()
	bdr.Flg("cacert", "").
		Default("").
		Description("Consul Client CA cert)", ``).
		PlaceHolder("FILE").
		Group("Consul").
		Build()
	bdr.Flg("cert", "").
		Default("").
		Description("Consul Client cert", ``).
		PlaceHolder("FILE").
		Group("Consul").
		Build()
	bdr.Flg("scheme", "").
		Default("http").
		Description("Consul connection protocol", ``).
		PlaceHolder("SCHEME").
		Group("Consul").
		Build()
	bdr.Flg("username", "u", "user", "usr", "uid").
		Default("").
		Description("HTTP Basic auth user", ``).
		PlaceHolder("USERNAME").
		Group("Consul").
		Build()
	bdr.Flg("password", "pw", "passwd", "pass", "pwd").
		Default("").
		Description("HTTP Basic auth password", ``).
		PlaceHolder("PASSWORD").
		Group("Consul").
		ExternalEditor(cli.ExternalToolPasswordInput).
		Build()
	bdr.Flg("message", "m", "msg", "mesg").
		Default("").
		Description("The commit message", ``).
		PlaceHolder("MSG").
		Group("Git").
		ExternalEditor(cli.ExternalToolEditor).
		Build()
}

func serverStartup(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverStop(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverShutdown(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverRestart(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverLiveReload(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverInstall(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverUninstall(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverStatus(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	return
}

func serverPause(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func serverResume(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	return
}

func kvBackup(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	// err = consul.Backup()
	fmt.Printf(`
cmd: %v
args: %v

`, cmd, args)
	return
}

func kvRestore(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	// err = consul.Restore()
	return
}

func msList(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.ServiceList()
	return
}

func msTagsList(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.TagsList()
	return
}

func msTagsAdd(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsRemove(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsModify(ctx context.Context, cmd *cli.Command, args []string) (err error) {
	_, _ = cmd, args
	// err = consul.Tags()
	return
}

func msTagsToggle(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:unused
	_, _ = cmd, args
	// err = consul.TagsToggle()
	return
}
