#!/usr/bin/env bash

# aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

build-all(){
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

cmd=$1 && shift || { build-all; exit; }
build-one $cmd