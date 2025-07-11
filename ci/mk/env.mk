
ifeq ($(OS),Windows_NT)
    LS_OPT=
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    LS_OPT=
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        OS = Linux
        CCFLAGS += -D LINUX
        LS_OPT=--color
    endif
    ifeq ($(UNAME_S),Darwin)
        OS = macOS
        CCFLAGS += -D OSX
        LS_OPT=-G
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif


INSTALL            ?= 0
UNINSTALL          ?= 0
ifeq (INSTALL,1)
  INSTALL_CMD      :=
  UNINSTALL_CMD    :=
  INSTALL_HELP     := 
  UNINSTALL_HELP   := 
else
  INSTALL_CMD      := echo
  UNINSTALL_CMD    := echo
  INSTALL_HELP     := @printf "\n\e[0;38;2;133;133;133m>>> %s\e[0m\n" "Use 'make INSTALL=1 install' to commit installing action to your local system."
  UNINSTALL_HELP   := @printf "\n\e[0;38;2;133;133;133m>>> %s\e[0m\n" "Use 'make UNINSTALL=1 install' to commit installing action to your local system."
endif
INSTALL_PREFIX     := "/usr/local"
LS_OPT             := "--color"
TIMESTAMP          := $(shell date -u '+%Y-%m-%dT%H:%M:%S.%N')
# TIMESTAMP        := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
# TIMESTAMP        := $(shell date -u '+%Y-%mm-%ddT%HH:%MM:%SS')

ifeq ($(OS),macOS)
    INSTALL_PREFIX := $(shell brew --prefix)
    LS_OPT         := "-G"
	TIMESTAMP      := $(shell date -Iseconds)
endif

LS                 := ls $(LS_OPT)
LL                 := ls -la $(LS_OPT)
LA                 := ls -la $(LS_OPT)
M                  := $(shell printf "\033[34;1m▶\033[0m")
TIP                := printf "\e[0;38;2;133;133;133m>>> %s\e[0m\n"
ERR                := printf "\e[0;33;1;133;133;133m>>> %s\e[0m\n"
DBG                := printf ">>> \e[0;38;2;133;133;133m%s\e[0m\n"
HEADLINE           := printf "\e[0;1m%s\e[0m:\n"
MM                 := printf "\033[34;1m▶\033[0m %s\e[0m\n"