#!/bin/bash

# Craft Launcher Linux Setup Script

echo "Checking for root permissions..."
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (sudo ./install_linux.sh)"
  exit
fi

echo "Updating package lists..."
# Ensure the required source is present for Ubuntu Jammy compatibility if needed
# Note: This is specific to the user's finding, might want to be careful about forcefully editing sources.list generally, 
# but sticking to the user's request for now.
if grep -q "jammy main" /etc/apt/sources.list; then
    echo "Source list already contains jammy main"
else
    echo "Adding jammy main to sources.list..."
    echo "deb http://gb.archive.ubuntu.com/ubuntu jammy main" >> /etc/apt/sources.list
fi

apt update

echo "Installing libwebkit2gtk-4.0-dev..."
apt install -y libwebkit2gtk-4.0-dev

echo "Setting permissions..."
if [ -f "craft-launcher-linux-amd64" ]; then
    chmod 777 craft-launcher-linux-amd64
    echo "Launcher is now executable."
    echo "You can run it with: ./craft-launcher-linux-amd64"
else
    echo "Warning: craft-launcher-linux-amd64 not found in current directory."
    echo "Please ensure this script is next to the launcher binary."
fi



echo "Done!"

