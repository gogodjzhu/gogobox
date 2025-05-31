#!/bin/bash

# check if the script is run as root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# check if apt-get is installed
if ! [ -x "$(command -v apt-get)" ]; then
  echo 'Error: apt-get is not installed.' >&2
  exit 1
fi

echo "Step0. Create bin directory"
# create bin if not exist, exit if exist
if [ ! -d "/workspaces/bin" ]; then
  mkdir /workspaces/bin
else
  echo "bin directory already exist, skip init.sh"
  exit 0
fi

echo "Step1. Install basic tools"
echo "export HISTSIZE=10000000" >> /etc/profile
apt-get -y update
apt-get install -y vim curl iputils-ping net-tools zsh wget fzf
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
git config --global --add safe.directory '*'

echo "Step2. Install buzz tools"

echo "Step3. Build package"
go mod tidy
go mod vendor