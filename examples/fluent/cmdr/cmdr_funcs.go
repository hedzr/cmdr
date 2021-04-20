package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr"
)

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmd.NewFlagV("=", "delim", "d").
		Description("delimitor char in `non-plain` mode.").
		Placeholder("")

	cmd.NewFlagV(false, "clear", "c").
		Description("clear all tags.").
		Placeholder("").
		Group("Operate")

	cmd.NewFlagV(false, "string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.").
		Placeholder("").
		Group("Mode")

	cmd.NewFlagV(false, "meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).").
		Placeholder("").
		Group("Mode")

	cmd.NewFlagV(false, "both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.").
		Placeholder("").
		Group("Mode")

	cmd.NewFlagV(false, "plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.").
		Placeholder("").
		Group("Mode")

	cmd.NewFlagV(true, "tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.").
		Placeholder("").
		Group("Mode")

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {

	cmd.NewFlagV("localhost", "addr", "a").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		Placeholder("HOST[:PORT]").
		Group("Consul")
	cmd.NewFlagV(8500, "port", "p").
		Description("Consul port", ``).
		Placeholder("PORT").
		Group("Consul")
	cmd.NewFlagV(true, "insecure", "K").
		Description("Skip TLS host verification", ``).
		Placeholder("").
		Group("Consul")
	cmd.NewFlagV("/", "prefix", "px").
		Description("Root key prefix", ``).
		Placeholder("ROOT").
		Group("Consul")
	cmd.NewFlagV("", "cacert").
		Description("Consul Client CA cert)", ``).
		Placeholder("FILE").
		Group("Consul")
	cmd.NewFlagV("", "cert").
		Description("Consul Client cert", ``).
		Placeholder("FILE").
		Group("Consul")
	cmd.NewFlagV("http", "scheme").
		Description("Consul connection protocol", ``).
		Placeholder("SCHEME").
		Group("Consul")
	cmd.NewFlagV("", "username", "u", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		Placeholder("USERNAME").
		Group("Consul")
	cmd.NewFlagV("", "password", "pw", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		Placeholder("PASSWORD").
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput)

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
