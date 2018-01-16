# motki-cli Makefile

# Build the motki binary:
#   make
#
# Install the program:
#   make install
#
# Remove the program:
#   make uninstall
#
# Build for a different OS and arch:
#   make build GOOS=linux GOARCH=arm7
#
# Cross-compile the binaries for many platforms at once:
#   make matrix ARCHES="amd64 arm6 arm7 386" OSES="windows linux darwin"
#
# Clean up build files:
#   make clean

# Defines where build files are stored.
PREFIX ?= build/

# Define all the necessary binary dependencies.
deps := go git

# A template for defining a variable with the final form:
#   NAME ?= /path/to/bin/name
deps_tpl = $(shell echo $(1) | tr a-z A-Z) ?= $(shell which $(1))

# Initialize a variable for each dependency listed in deps.
# The variables are upper-cased and overridable from the command-lane.
# For example, the variable containing the path to the go binary is called GO,
# psql is called PSQL, etc.
# Note that this variable's value is meaningless, it serves only as a name for this procedure.
deps_initialized := $(foreach dep,$(deps),$(eval $(call deps_tpl,$(dep))))

# Error messages.
err_go_missing := unable to locate "go" binary. See https://golang.org/doc/ for more information.
err_default_missing = unable to locate "$1" binary. Ensure that it is installed and on your PATH. Specify a custom path to the binary with "make $(shell echo $1 | tr a-z A-Z)=/path/to/$1 $@"
err_binary_missing = $(or $(err_$1_missing),$(call err_default_missing,$1))

# This procedure throws a fatal error if the path to the given binary is empty or
# does not exist.
ensure_dep = $(if $(realpath $(value $(shell echo $(1) | tr a-z A-Z))),,$(error $(call err_binary_missing,$(1))),exit 0;)

# Make sure we have all our dependencies. If a dependency is missing, make will
# exit with an appropriate error message.
# Note that this variable's value is meaningless, it serves only as a name for this procedure.
deps_ensured := $(foreach dep,$(deps),$(call ensure_dep,$(dep)))

# By default, the system defined GOOS and GOARCH are used.
# These are overridable from the command line. For example:
#   make build GOOS=linux GOARCH=arm7
GOOS   ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

# When using the "matrix" target, these specify which OSes and arches to target.
# These are both overridable from the command line. For example:
#   make matrix ARCHES="amd64 arm6 arm7 386" OSES=linux
ARCHES ?= amd64
OSES   ?= linux darwin windows

# Components to build up a valid "go build" command.
build_version := $(if $(shell test -d .git && echo "1"),$(shell $(GIT) describe --always --tags),snapshot)
build_base    := $(GO) build -ldflags "-s -w -X main.Version=$(build_version)"
build_name     = $(PREFIX)$1$(if $(filter $(GOOS),windows),.exe,)
build_src      = ./cmd/$(word 1,$(subst _, ,$(subst ., ,$(subst $(PREFIX),,$1))))/*.go
build_cmd      = GOOS=$(GOOS) GOARCH=$(GOARCH) $(build_base) -o $1 $(call build_src,$1)
release_name   = $(call build_name,$1_$(GOOS)_$(GOARCH))
release_cmd    = $(subst build -ldflags, build -tags release -ldflags,$(call build_cmd,$1))

# These define the programs that get built. Adding more targets is
# automatic as long as the source code for the target exists in
# ./cmd/<target>/*.go.
binaries        := motki
binary_targets  := $(foreach bin,$(binaries),$(call build_name,$(bin)))
release_targets := $(foreach bin,$(binaries),$(call release_name,$(bin)))

# Print configuration information: paths, build options, and config params.
extra_params := GOOS GOARCH
define print_conf
	@$(foreach dep,$(deps),echo "$(shell echo $(dep) | tr a-z A-Z)=$(value $(shell echo $(dep) | tr a-z A-Z))";)
	@$(foreach val,$(extra_params),echo "$(val)=$($(val))";)
endef

# All of the files this generates.
files := $(PREFIX)motki_*_* $(binary_targets)

.PHONY: all
.PHONY: generate build release matrix
.PHONY: install uninstall
.PHONY: clean clean_files
.PHONY: debug


# Build all binaries.
build: $(binary_targets)

# Installs the program.
install:
	$(error "not implemented")


# This defines a target that matches any of the values listed in binary_targets.
$(binary_targets):
	$(call build_cmd,$@)
	@echo "Built $@"

# Make release builds for the specified OS and arch.
release: generate $(release_targets)

$(release_targets):
	$(call release_cmd,$@)
	@echo "Built $@"

# This target will build a binary for every combination of
# ARCHES and OSES specified.
matrix:
	@for arch in $(ARCHES); do                       \
		for os in $(OSES); do                        \
			echo "Building $$os $$arch...";          \
			$(MAKE) release GOOS=$$os GOARCH=$$arch; \
		done;                                        \
	done;                                            \
	echo "Done."

# Deletes build files.
clean: clean_files

# Deletes all build files.
clean_files:
	@for f in $(files); do (rm -r "$$f" 2> /dev/null && echo "Deleted $$f"; exit 0); done
	@echo "Cleaned files."

# Uninstalls the program.
uninstall:
	$(error "not implemented")


# Prints configuration information.
debug:
	$(print_conf)
