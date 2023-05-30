# Build weshnet

These are instructions to build weshnet.

<!-- markdownlint-disable MD034 -->

## Prerequisites

* Required: asdf
* Required on macOS: Command Line Developer Tools

Following are the steps to install each prerequisite (if it's needed for your
build target).

### macOS 11, macOS 12 and macOS 13

To install the Command Line Developer Tools, in a terminal enter:

    xcode-select --install

After the Developer Tools are installed, we need to make sure it is updated. In
System Preferences, click Software Update and update it if needed.

To install asdf using brew, follow instructions at https://asdf-vm.com . In short,
first install brew following the instructions at https://brew.sh . Then, in
a terminal enter:

    brew install asdf

If your terminal is zsh, enter:

    echo -e "\n. $(brew --prefix asdf)/libexec/asdf.sh" >> ${ZDOTDIR:-~}/.zshrc

If your terminal is bash, enter:

    echo -e "\n. \"$(brew --prefix asdf)/libexec/asdf.sh\"" >> ~/.bash_profile

Start a new terminal to get the changes to .zshrc .

### Ubuntu 18.04, 20.04 and 22.04

To install asdf, follow instructions at https://asdf-vm.com . In short, in
a terminal enter:

    sudo apt install curl git
    git clone https://github.com/asdf-vm/asdf.git ~/.asdf
    echo '. "$HOME/.asdf/asdf.sh"' >> ~/.bashrc

Start a new terminal to get the changes to .bashrc .

## Build

In a terminal, enter:

    git clone https://github.com/berty/weshnet
    cd weshnet

First time only (or after updating .tool-versions), enter:

    make asdf.install_tools

To run the tests, enter:

    make test

Or you can make other targets. See:

    make help
