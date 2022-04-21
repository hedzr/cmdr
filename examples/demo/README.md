# demo

- with external configuration files
  checkout `<project>/ci/etc/demo/demo.yml` and `<project>/ci/etc/demo/conf.d/*.yml`.
- normalize app structure
-

```bash

[ -d ci/certs ] || mkdir -p ci/certs
openssl req -newkey rsa:2048 -nodes -keyout ci/certs/server.key -x509 -days 3650 -out ci/certs/server.crt

```




