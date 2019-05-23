#!/usr/bin/env bash



build_core () {
	build_one $*
}

build_win () {
	build_one amd64 windows win $*
}

build_linux () {
	build_one amd64 linux bin $*
}

build_darwin () {
	build_one amd64 darwin dar $*
}

build_all () {
	for ARCH in amd64; do
		for OS in darwin linux windows; do
			build_one $ARCH $OS $OS-$ARCH $*
		done
	done
}

build_one () {
	local ARCH=${1:-amd64}
	local OS=${2:-darwin}
	local suffix=${3:-dar}
	local S=''
	case $suffix in
		dar) 	S="" ;;
		bin|linux-*|darwin-*) 	S="-$OS-$ARCH" ;;
		win|windows-*)	S="-$OS-$ARCH.exe" ;;
		*) 		S="-$suffix" ;;
	esac
	shift;shift;shift;
	echo "PWD=$(pwd)"
	headline "---- Building the binary, for $PROJ_DIR | S='$S', OS='$OS' | suffix=$suffix"

	cat <<-EOF
	GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build -tags="gocql_debug" -ldflags "$LDFLAGS" -o "$PROJ_DIR/bin/${APPNAME}$S" $* $PKG_SRC
	EOF
	GOARCH=$ARCH GOOS=$OS CGO_ENABLED=1 go build -tags="gocql_debug" -ldflags "$LDFLAGS" -o "$PROJ_DIR/bin/${APPNAME}$S" $* $PKG_SRC && \
	chmod +x "$PROJ_DIR/bin/${APPNAME}$S" && \
	echo && ls -l "$PROJ_DIR/bin/${APPNAME}$S" # && echo $?
}

app_name () {
	local temp=$(grep -E "APP_NAME[ \t]+=[ \t]+" doc.go|grep -Eo "\\\".+\\\"")
	temp="${temp%\"}"
	temp="${temp#\"}"
	echo $temp
}

app_version () {
	grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+"
}

app_version_int () {
	grep -E "VersionInt[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9x]+"
}


build_fmt () {
	gofmt -l -w -s .
}



build_mod() { commander 'build_mod' "$@"; }
build_mod_usage () {
	cat <<-EOF
	'mod' Usages:
	
	SUB-COMMANDS:
	  run		run 
	  test,dump	dump the information
	  
	Help:
	  # go mod helpers:
	  export GO111MODULE=on
	  export GOPROXY=https://athens.azurefd.net
	  export GOPATH=$(dirname $(dirname $(dirname $(dirname $CD))))
	EOF
}
build_mod_run () {
	export GO111MODULE=on
	export GOPROXY=https://athens.azurefd.net
	export GOPATH=$(dirname $(dirname $CD))
}
build_mod_test () {
	cat <<-EOF
	GO111MODULE=$GO111MODULE
	GOPROXY=$GOPROXY
	GOPATH=$GOPATH
	EOF
}
build_mod_dump () {
	build_mod_test
}


build_help(){
	cat <<-EOF
	Usages: $0 [command]
	Commands:
	  (nothing)	go build "consul-tags" to executable
	  mod		(( go mod sub-commands set ))
	  help		Show this screen
	EOF
}
build_usage(){ build_help; }
build_usages(){ build_help; }
build_h(){ build_help; }
build_-h(){ build_help; }
build_--help(){ build_help; }




#### write your functions here, and invoke them by: `./bash.sh <your-func-name>`
cool(){ echo cool; }
sleeping(){ echo sleeping; }



_my_main_do_sth(){
	local cmd=${1:-core} && { [ $# -ge 1 ] && shift; } || :
	# for linux only:
	# local cmd=${1:-sleeping} && shift || :

	is_vagrant && CD=/vagrant && SCRIPT=$CD/bootstrap.sh
	headline "CD = $CD, SCRIPT = $SCRIPT, PWD=$(pwd)"
	# [ -d $CD/.boot.d ] || mkdir -p $CD/.boot.d
	# [ -z "$(ls -A $CD/.boot.d/)" ] && for f in $CD/.boot.d/*; do source $f; done

	local GOSRC=$GOPATH/src/

	local APPNAME=${APPNAME:-$(app_name)}
	local PWD=$(pwd)
	local PKG=${PWD/$GOSRC/}
	local PROJ=$(cat doc.go|grep -o 'package .*'|awk '{print $2}')
	local PROJ_DIR=$PWD
	local VERSION=$(app_version)
	local versionInt=$(app_version_int)
	echo "appName=$APPNAME | pkg=$PKG, $PROJ | version = $VERSION, $versionInt"

	local W_PKG="github.com/hedzr/cmdr"
	local LDFLAGS="-s -w -X ${W_PKG}.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X ${W_PKG}.Githash=`git rev-parse HEAD` -X ${W_PKG}.Version=$VERSION -X ${W_PKG}.AppName=$APPNAME"
	local PKG_SRC=${PKG_SRC:-cli/main.go}
	#echo "LDFLAGS=$LDFLAGS"

	debug "build_$cmd - $@"
	eval "build_$cmd $@" || { echo 'not true'; :; }
}

is_vagrant() { [[ -d /vagrant ]]; }













#### HZ Tail BEGIN ####
in_debug()       { [[ $DEBUG -eq 1 ]]; }
is_root()        { [ "$(id -u)" = "0" ]; }
is_bash()        { is_bash_t1 && is_bush_t2; }
is_bash_t1()     { [ -n "$BASH_VERSION" ]; }
is_bash_t2()     { [ ! -n "$BASH" ]; }
is_zsh()         { [[ $SHELL == */zsh ]]; }
is_zsh_t2()      { [ -n "$ZSH_NAME" ]; }
is_darwin()      { [[ $OSTYPE == *darwin* ]]; }
is_linux()       { [[ $OSTYPE == *linux* ]]; }
in_sourcing()    { is_zsh && [[ "$ZSH_EVAL_CONTEXT" == toplevel* ]] || [[ $(basename -- "$0") != $(basename -- "${BASH_SOURCE[0]}") ]]; }
is_interactive_shell () { [[ $- == *i* ]]; }
is_not_interactive_shell () { [[ $- != *i* ]]; }
is_ps1 () { [ -z "$PS1" ]; }
is_not_ps1 () { [ ! -z "$PS1" ]; }
is_stdin () { [ -t 0 ]; }
is_not_stdin () { [ ! -t 0 ]; }
headline()       { printf "\e[0;1m$@\e[0m:\n"; }
headline_begin() { printf "\e[0;1m"; }  # for more color, see: shttps://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
headline_end()   { printf "\e[0m:\n"; } # https://misc.flogisoft.com/bash/tip_colors_and_formatting
printf_black()   { printf "\e[0;30m$@\e[0m:\n"; }
printf_red()     { printf "\e[0;31m$@\e[0m:\n"; }
printf_green()   { printf "\e[0;32m$@\e[0m:\n"; }
printf_yellow()  { printf "\e[0;33m$@\e[0m:\n"; }
printf_blue()    { printf "\e[0;34m$@\e[0m:\n"; }
printf_purple()  { printf "\e[0;35m$@\e[0m:\n"; }
printf_cyan()    { printf "\e[0;36m$@\e[0m:\n"; }
printf_white()   { printf "\e[0;37m$@\e[0m:\n"; }
debug()          { in_debug && printf "\e[0;38;2;133;133;133m$@\e[0m\n" || :; }
debug_begin()    { printf "\e[0;38;2;133;133;133m"; }
debug_end()      { printf "\e[0m\n"; }
dbg()            { ((DEBUG)) && printf ">>> \e[0;38;2;133;133;133m$@\e[0m\n" || :; }
debug_info()     {
	debug_begin
	cat <<-EOF
	             in_debug: $(in_debug && echo Y || echo '-')
	              is_root: $(is_root && echo Y || echo '-')
	              is_bash: $(is_bash && echo Y || echo '-')
	               is_zsh: $(is_zsh && echo Y || echo '-')
	          in_sourcing: $(in_sourcing && echo Y || echo '-')   # ZSH_EVAL_CONTEXT = $ZSH_EVAL_CONTEXT
	 is_interactive_shell: $(is_interactive_shell && echo Y || echo '-')
	EOF
	debug_end
}
commander ()    {
	local self=$1; shift;
	local cmd=${1:-usage}; [ $# -eq 0 ] || shift;
	#local self=${FUNCNAME[0]}
	case $cmd in
	help|usage|--help|-h|-H) "${self}_usage" "$@"; ;;
	funcs|--funcs|--functions|--fn|-fn)  script_functions "^$self"; ;;
	*)
		if [ "$(type -t ${self}_${cmd}_entry)" == "function" ]; then
		"${self}_${cmd}_entry" "$@"
		else
		"${self}_${cmd}" "$@"
		fi
		;;
	esac
}
script_functions () {
	# shellcheck disable=SC2155
	local fncs=$(declare -F -p | cut -d " " -f 3|grep -vP "^[_-]"|grep -vP "\\."|grep -vP "^[A-Z]"); # Get function list
	if [ $# -eq 0 ]; then
	echo "$fncs"; # not quoted here to create shell "argument list" of funcs.
	else
	echo "$fncs"|grep -P "$@"
	fi
	#declare MyFuncs=($(script.functions));
}
main_do_sth()    {
	set -e
	trap 'previous_command=$this_command; this_command=$BASH_COMMAND' DEBUG
	trap '[ $? -ne 0 ] && echo FAILED COMMAND: $previous_command with exit code $?' EXIT
	MAIN_DEV=${MAIN_DEV:-eth0}
	MAIN_ENTRY=${MAIN_ENTRY:-_my_main_do_sth}
	# echo $MAIN_ENTRY - "$@"
	in_debug && { debug_info; echo "$SHELL : $ZSH_NAME - $ZSH_VERSION | BASH_VERSION = $BASH_VERSION"; [ -n "$ZSH_NAME" ] && echo "x!"; }
	$MAIN_ENTRY "$@"
	trap - EXIT
	${HAS_END:-$(false)} && { debug_begin;echo -n 'Success!';debug_end; } || :
}
DEBUG=${DEBUG:-0}
is_darwin && realpathx(){ [[ $1 == /* ]] && echo "$1" || echo "$PWD/${1#./}"; } || realpathx () { readlink -f $*; }
in_sourcing && { CD=${CD}; debug ">> IN SOURCING, \$0=$0, \$_=$_"; } || { SCRIPT=$(realpathx $0) && CD=$(dirname $SCRIPT) && debug ">> '$SCRIPT' in '$CD', \$0='$0','$1'."; }
main_do_sth "$@"
#### HZ Tail END ####
