
# .info:
# 	ll /usr/local/opt/gcc@12/bin /usr/local/opt/llvm/bin

printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

## list: list all available targets in this Makefile
list:
	@printf "%-20s %s\n" "Target" "Description"
	@printf "%-20s %s\n" "------" "-----------"
	@$(MAKE) -pqR : 2>/dev/null \
	    | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' \
	    | sort \
	    | grep -vE -e '^[^[:alnum:]]' -e '^$@$$' \
	    | xargs -I _ sh -c 'printf "%-20s " _; make _ -nB | (grep -i "^# Help:" || echo "") | tail -1 | sed "s/^# Help: //g"'

#: list all available targets in this Makefile, and current build env
## help: list all available targets in this Makefile, and current build env
i info help: .help1 # for '#: xxx\nxxx:'
	@echo "> The environment detected:"
	@echo
	@echo "          Current OS (COS) = $(COS), OS = $(OS), GCC_PRIOR = $(GCC_PRIOR)"
	@echo "                      Arch = $(ARCH) ($(CARCH))"
	@echo "        uname -p | -s | -m = $(UNAME_P) | $(UNAME_S) | $(UNAME_M)"
	@echo "                 uid / gid = $(CURRENT_UID) / $(CURRENT_GID)"
	@echo "                    CCTYPE = $(CCTYPE)"
	@[ "$(GCC_PREFIX)" != "" ]  && echo "                       gcc = $(GCC_PREFIX)/, $(shell $(GCC_PREFIX)/bin/gcc-12 -v 2>&1 | tail -1)"  || echo "                       gcc = $(GCC), $(GCC_VER)"
	@[ "$(LLVM_PREFIX)" != "" ] && echo "                llvm clang = $(LLVM_PREFIX)/, $(shell $(LLVM_PREFIX)/bin/clang -v 2>&1 | head -1)" || \
	 { [ "$(CLANG_VER)" != "" ] && echo "                llvm clang = $(CLANG), $(CLANG_VER)" || echo "                llvm clang = "; }
	@echo "              CC/GCC/CLANG = $(CC) | $(GCC) | $(CLANG)"
	@echo "           CXX/GXX/CLANGXX = $(CXX) | $(GXX) | $(CLANGXX)"
	@echo "                    CFLAGS = $(CFLAGS)"
	@echo "                  CXXFLAGS = $(CXXFLAGS)"
	@echo " CPPFLAGS (for both c/c++) = $(CPPFLAGS)"
	@echo "                FLEX/BISON = $(LEX) | $(YACC)"
	@echo "                    LFLAGS = $(LFLAGS)"
	@echo "                    YFLAGS = $(YFLAGS)"
	@echo "                       LLD = $(LLD)"
	@echo "                   LDFLAGS = $(LDFLAGS)"
	@echo "                   OBJDUMP = $(OBJDUMP)"
	@echo "                   READELF = $(READELF)"
	@echo "                CLANG-TIDY = $(CLANG_TIDY)"
	@echo "              CLANG-FORMAT = $(CLANG_FORMAT)"
	@echo "               NASM Format = $(NASM_FMT) (suffix: $(NASM_FMT_SUFFIX))"
	@-$(MAKE) help-extras
	@echo "END."

.help1: Makefile # for '## xx: xx'
	@echo
	@echo "> Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.help-help: # help of help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	  | column -t  -s ' '

.help2: Makefile
	@grep -B2 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" $(MAKEFILE_LIST) \
	  | grep -v -- -- \
	  | sed 'N;s/\n/###/' \
	  | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	  | column -t  -s '###'


# .DEFAULT_GOAL := help
.PHONY: help list info i
