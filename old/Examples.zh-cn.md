# cmdr-examples

see also: https://github.com/hedzr/cmdr-examples

#### Concepts

- Options: flags and app settings

	- Tail Args

	  首先我们知道 cmdr 采用如下的命令行格式：

	  ```bash
	  app <sub-commands> <all options or ascestor of it> <tail args>
	  ```

	  所以 Tail Args 也就是惯常所说的 Positional Arguments。为了与 POSIX 一致的原因，我们采用 TailArgs 来表示除开 Subcommands 序列以及 Flags
	  序列之外的剩余命令行参数。

	  一般来说，我们混用 Tail Arguments，Remained Arguments 等术语，无需要严格地区分他们。


- Action: [`actions`](https://github.com/hedzr/cmdr-examples/tree/master/examples/actions)

  cmdr 以 Action 来响应一条命令。

  我们可以这样定义一个 Subcommand：

  ```go
	  root := cmdr.Root(appName, "1.0.1").
		  Copyright(copyright, "hedzr").
		  Description(desc, longDesc).
		  Examples(examples)
	  // rootCmd = root.RootCommand()
	  soundex(root)
	// ...
  
  func soundex(root cmdr.OptCmd) {
	  // soundex
  
	// To test for this subcommand, type command in shell:
	// $ go run ./examples/simple soundex quick fox
	  root.NewSubCommand("soundex", "snd", "sndx", "sound").
		  Description("soundex test").
		  Group("Test").
		  TailPlaceholder("[text1, text2, ...]").
		  Action(func(cmd *cmdr.Command, remainArgs []string) (err error) {
			  for ix, s := range remainArgs {
				  fmt.Printf("%5d. %s => %s\n", ix, s, cmdr.Soundex(s))
			  }
			  return
		  })
  }
  ```

  所以，我们可以使用如下命令行对其进行测试：

  ```bash
  $ go run ./examples/simple soundex quick fox
  ```

  对于 `soundex` 这个子命令来说，其 Action 函数回调将会处理其具体逻辑。

  我们可以注意到回调函数的入口参数中包含已经被解析到的当前命令对象，以及该命令的剩余 TailArgs，你可以将其理解为 Remained Arguments。


- 

