package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr"
)

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmdr.NewString("=").Titles("delim", "d").
		// cmd.NewFlagV("=", "delim", "d").
		Description("delimiter char in `non-plain` mode.").
		Placeholder("").
		CompletionJustOnce(true).
		AttachTo(cmd)

	cmdr.NewBool().Titles("clear", "c").
		// cmd.NewFlagV(false, "clear", "c").
		Description("clear all tags.").
		Placeholder("").
		Group("Operate").
		VendorHidden(false).
		AttachTo(cmd)

	cmdr.NewBool().Titles("string", "g", "string-mode").
		// cmd.NewFlagV(false, "string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.").
		Placeholder("").
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().Titles("meta", "m", "meta-mode").
		// cmd.NewFlagV(false, "meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).").
		Placeholder("").
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().Titles("both", "2", "both-mode").
		// cmd.NewFlagV(false, "both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.").
		Placeholder("").
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().Titles("plain", "p", "plain-mode").
		// cmd.NewFlagV(false, "plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.").
		Placeholder("").
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().Titles("tag", "t", "tag-mode").
		// cmd.NewFlagV(true, "tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.").
		Placeholder("").
		ToggleGroup("Mode").
		AttachTo(cmd)

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {

	cmdr.NewString("localhost").Titles("addr", "a").
		// cmd.NewFlagV("localhost", "addr", "a").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		Placeholder("HOST[:PORT]").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewInt(8500).Titles("port", "p").
		Description("Consul port", ``).
		Placeholder("PORT").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewBool(false).Titles("insecure", "K").
		Description("Skip TLS host verification", ``).
		Placeholder("").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString("/").Titles("prefix", "px").
		Description("Root key prefix", ``).
		Placeholder("ROOT").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString("").Titles("cacert", "").
		Description("Consul Client CA cert)", ``).
		Placeholder("FILE").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString("").Titles("cert", "").
		Description("Consul Client cert", ``).
		Placeholder("FILE").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString("http").Titles("scheme", "").
		Description("Consul connection protocol", ``).
		Placeholder("SCHEME").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString().Titles("username", "u", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		Placeholder("USERNAME").
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString().Titles("password", "pw", "passwd", "pass", "pwd").
		// cmd.NewFlagV("", "password", "pw", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		Placeholder("PASSWORD").
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput).
		AttachTo(cmd)

}

func kvBackup(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Backup()
	fmt.Printf(`
cmd: %v
args: %v

`, cmd, args)
	return
}

func kvRestore(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Restore()
	return
}

func msList(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.ServiceList()
	return
}

func msTagsList(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.TagsList()
	return
}

func msTagsAdd(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsRemove(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsModify(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.Tags()
	return
}

func msTagsToggle(cmd *cmdr.Command, args []string) (err error) {
	// err = consul.TagsToggle()
	return
}
