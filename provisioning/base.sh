#!/bin/bash
#
# Base provisioning script for the gogame box

# -------------------------------------
# Appends the given string to the specified file if it is not already found
# in that file
# Arguments:
#   filename
#   string
# -------------------------------------
function appendIfNotInFile {
  grep -qF "$2" "$1" || echo "$2" >> "$1"
}

# -------------------------------------
# Makes the progress output more visible
# -------------------------------------
function progressEcho {
  echo -e "\e[1m\e[48;5;27mProvisioning Progress: ${*}\e[0m"
}

# die if any command fails
set -e
export DEBIAN_FRONTEND=noninteractive

# make sure apt is up to date 
progressEcho "Updating APT package index"
apt-get update -qq --fix-missing

# base packages
progressEcho "Installing base packages"
apt-get install -y -qq git vim curl wget

# redis 
progressEcho "Installing redis"
apt-get install -y -qq redis-server

# Download go
goVersion="1.7.4"
goTmpDir="/tmp/go${goVersion}"
goSource="http://golang.org/dl/go${goVersion}.linux-amd64.tar.gz"

goDestination="${goTmpDir}/go${goVersion}.linux-amd64.tar.gz"
mkdir -p "$goTmpDir"

progressEcho "Downloading GO ${goVersion} - this can take a few minutes"
wget --no-verbose -O "$goDestination" "$goSource"

progressEcho "Installing GO"

# Extract go
cd "$goTmpDir"
tar xvf *.tar.gz
chown -R root:root ./go
# Remove any previously installed versions
rm -rf /usr/local/go
# move
mv go /usr/local

progressEcho "Configuring workspace"

# Setup path - if not already set
appendIfNotInFile \
  '/home/vagrant/.profile' \
  'export PATH=/vagrant/bin:/usr/local/go/bin:$PATH'

# Setup workspace - if not already set
appendIfNotInFile \
  '/home/vagrant/.profile' \
  'export GOPATH=/vagrant'

progressEcho "Done"