/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	mycmdr "github.com/hedzr/cmdr/examples/fluent/cmdr"
)

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// defer func() {
	// 	fmt.Println("defer caller")
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("recover success. error: %v", err)
	// 	}
	// }()

	mycmdr.Entry()
}
