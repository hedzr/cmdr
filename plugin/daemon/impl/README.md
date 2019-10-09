# daemon impl, pure go

I rewrote this one because most of others are dead or deading.



To find out all goos and goarch values:

- https://github.com/golang/go/blob/master/src/go/build/syslist.go#L7





## ACK

Thanks them (incompletely):

- <https://github.com/sevlyar/go-daemon>
- <https://github.com/VividCortex/godaemon>
- <https://github.com/takama/daemon>
- <https://github.com/fiorix/go-daemon>
- <https://github.com/leprosus/golang-daemon>
- <https://github.com/kardianos/service>





## References

- [Daemon Definition](http://www.linfo.org/daemon.html)
- <https://en.wikipedia.org/wiki/Daemon_(computing)>
- <https://zh.wikipedia.org/wiki/守护进程>
- <https://web.archive.org/web/20061118065514/http://www.linuxprofilm.com/articles/linux-daemon-howto.html>
  - <http://www.netzmafia.de/skripten/unix/linux-daemon-howto.html>
- <https://www.thegeekstuff.com/2012/02/c-daemon-process/>
- <https://socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example>
- <https://stackoverflow.com/questions/23736046/how-to-create-a-daemon-process-in-golang>
- <https://www.reddit.com/r/golang/comments/35v5bm/best_way_to_run_go_server_as_a_daemon/>
- <https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/>
- <https://blog.burntsushi.net/golang-daemonize-bsd/>
- Systemd
  - [Integration of a Go service with systemd: socket activation](https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation)
  - [Integration of a Go service with systemd: readiness *&* liveness](https://vincent.bernat.ch/en/blog/2017-systemd-golang)
  - [**go-systemd**](https://github.com/coreos/go-systemd)
  - [**WRITING SYSTEMD ENABLED APPLICATIONS IN GO**](https://lxtreme.nl/blog/writing-systemd-enabled-applications-in-go/)

