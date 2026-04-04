#!/bin/bash

if ! command -v go >/dev/null 2>&1; then
  echo -e "\033[31mError \033[0mgo binary not found"
  exit 1
fi

if ! command -v git >/dev/null 2>&1; then
  echo -e "\033[31mError \033[0mgit binary not found"
  exit 1
fi

echo -e "\033[32mInfo \033[0mclone miver from github"
rm -rf /tmp/miver
git clone https://github.com/hunter-hongg/miver.git /tmp/miver --depth=1
cd /tmp/miver 

echo -e "\033[32mInfo \033[0mbuild miver"
go build -o miver

echo -ne "\033[34mChoice \033[0minstall miver to /usr/local/bin? (y/n) " 
read choice
if [ "$choice" == "y" -o "$choice" == "Y" ]; then
  echo -e "\033[32mInfo \033[0minstall miver to /usr/local/bin"
  sudo cp miver /usr/local/bin
else 
  echo -e "\033[32mInfo \033[0minstall miver to ~/.local/bin"
  cp miver ~/.local/bin
fi

mkdir -p ~/.miver/
echo 'export PATH=$PATH:$HOME/.local/bin:$HOME/.miver/bin' > ~/.miver/env
echo -e "\033[32mInfo \033[0mplease source ~/.miver/env to add miver to your PATH"
echo -e "\033[32mInfo \033[0mthen you can use miver after running miver init"