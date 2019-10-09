#!/usr/bin/env bash

# aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

# go tool dist list

build-all-full() {
  for xx in $(go tool dist list); do
    IFS='/' read -r -a array <<<"$xx"

    case ${array[0]} in
    aix | android | js) ;;

    darwin | linux)
      case ${array[1]} in
      arm*) ;;
      *)
        local dist="${array[0]}-${array[1]}"
        local GOOS=${array[0]}
        local GOARCH=${array[1]}
        echo $dist
        one
        ;;
      esac
      ;;

    *)
      local dist="${array[0]}-${array[1]}"
      local GOOS=${array[0]}
      local GOARCH=${array[1]}
      echo $dist
      one
      ;;
    esac
  done
}

one() {
  echo "-- build for $GOOS - $GOARCH ..."
  GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/nosig-$GOOS-$GOARCH plugin/daemon/nosig/main.go
  ls -la bin/nosig-$GOOS-$GOARCH
  echo -n " > file $GOOS - $GOARCH ..."
  file bin/nosig-$GOOS-$GOARCH
}

build-all() {
  local A=(darwin
    dragonfly
    freebsd
    linux
    nacl
    netbsd
    openbsd
    plan9
    solaris
    windows)

  for GOOS in "${A[@]}"; do
    if [ "$GOOS" == "nacl" ]; then
      local B=(386 amd64p32 arm)
    else
      local B=(amd64)
    fi
    for GOARCH in "${B[@]}"; do
      echo "-- build for $GOOS ..."
      GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/nosig-$GOOS-$GOARCH plugin/daemon/nosig/main.go
    done
  done

  ls -la bin/nosig*

  for GOOS in "${A[@]}"; do
    if [ "$GOOS" == "nacl" ]; then
      local B=(386 amd64p32 arm)
    else
      local B=(amd64)
    fi
    echo -n " > file $GOOS - $GOARCH ..."
    file bin/nosig-$GOOS-$GOARCH
  done
}

build-one() {
  local GOOS=${1:-}

  if [ "$GOOS" == "nacl" ]; then
    local B=(386 amd64p32 arm)
  else
    local B=(amd64)
  fi

  for GOARCH in "${B[@]}"; do
    echo "-- build for $GOOS - $GOARCH ..."
    GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/nosig-$GOOS-$GOARCH plugin/daemon/nosig/main.go
    ls -la bin/nosig-$GOOS-$GOARCH
    echo -n " > file $GOOS - $GOARCH ..."
    file bin/nosig-$GOOS-$GOARCH
  done
}

cmd=$1 && shift || {
  build-all
  exit
}
case $cmd in
full) build-all-full ;;
*) build-one $cmd ;;
esac
