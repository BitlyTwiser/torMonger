#!/bin/sh
#
# Installation script for torCrawl. It tries to move $bin in one of the
# directories stored in $binpaths.

INSTALL_DIR=$(dirname $0)

echo "Compiling go binary.."

go build -o torMonger

bin="$INSTALL_DIR/torMonger"
binpaths="/usr/local/bin /usr/bin"

echo "Creating data location for MongoDB: /var/tmp/tormonger_data "
mkdir /var/tmp/mongodb_data

# This variable contains a nonzero length string in case the script fails
# because of missing write permissions.
is_write_perm_missing=""

for binpath in $binpaths; do
  if mv "$bin" "$binpath/$bin" ; then
    echo "Moved $bin to $binpath"
    exit 0
  else
    if [ -d "$binpath" ] && [ ! -w "$binpath" ]; then
      is_write_perm_missing=1
    fi
  fi
done

echo "We cannot install $bin in one of the directories $binpaths"

if [ -n "$is_write_perm_missing" ]; then
  echo "It seems that we do not have the necessary write permissions."
  echo "Perhaps try running this script as a privileged user:"
  echo
  echo "    sudo $0"
  echo
fi

exit 1
