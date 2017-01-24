# Introduction

This tool has been developed to help network engineer to deal with undesirable traffic that is passing through their Internet network. This tool has been design to propose a graphical user interface to manage network features like:

* Filter traffic with BGP flowspec,
* Drop malicious traffic with BGP blackhole,
* Design and configuration of RTBH (Remote Triggered Blackhole),
* Interface analytics system,
* ... 

This first version is currently Alpha and needs to go through a set of test to make it an usable version. For now, it is only supporting BGP flowspec (RFC5575). __This tool is not suppose to be installed in production network but rather be used for lab / test purposes.__

# Building blocks and dependencies

This tool is relying on two open source software and related API.

* The UI is using a Qt binding for Golang: https://github.com/therecipe/qt

* BGP protocol stack: https://github.com/osrg/gobgp

## Using this tool in your network

This tool provides a BGP route reflector (RR using GoBGP as BGP stack) with an UI to inject BGP flowspec updates. The BGP RR propagates those updates to all its peers. In the current version, the UI can only connect to a local Go BGP daemon and doesn't support BGP clustering. The UI is using the gRPC API to interface GoBGP.

GoBGP needs to be installed manually (not yet via the UI). All neighbors have to be configured in the GoBGP configuration file. Once the GoBGP daemon is running, you launch the application. The UI will then connect to the BGP daemon and will be ready to push BGP flowspec updates to all involved neighbors.

## Typical configuration

You will have to create a VM or use a bare metal server, install the OS of your choice (software has been developed on ubuntu server), install all dependencies (mainly Qt / GoBGP / Unity) and you are ready to go. Make sure that your server is using is reachable to all BGP neighbors.

# Tutorial and features
* [Basic main window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/main_win.md)
* [Flowspec window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/flowspec_win.md)

![flowspec-win](/docs/flowspec-win.png)

* [Console window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/console_win.md)

Next step will be about testing this tool against vendor routers.

# Install and configure your development machine

* Follow the installation process of GoBGP: https://github.com/osrg/gobgp/blob/master/docs/sources/getting-started.md
* Follow the installation process of Qt Golang binding: https://github.com/therecipe/qt
  * Make sure that you alocated 8 GB of RAM to your VM
* ENV variables: The following example needs to be updated acccording to your machine but here is a snippet of my .bashrc file as an example

```
# go-lang variable
export GOPATH=/home/matthieu/go-work
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig
export LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu
export GOBIN=/home/matthieu/go-work/bin
PATH=/home/matthieu/go-work/bin:$PATH
PATH=/usr/local/go/bin:$PATH
# qt variable
export QT_DIR=/home/matthieu/Qt
export QT_VERSION=5.7.0
PATH=/home/matthieu/Qt/5.7/gcc_64/bin:$PATH
```


# Tool name and special dedicace

This tool is named "Gabu". This name is just the nickname of a brilliant young french engineer that passed away way too early and miss to anybody who crossed his way.

To enjoy a nice and still very interesting BGP flowspec presentation done by Frederic Gabut-Deloraine about the use of this protocol for DDoS mitigation, please follow this link: [Frederic FRnOG presentation](http://www.dailymotion.com/video/xtngjg_frnog-18-flowspec-frederic-gabut-deloraine-neo-telecoms_tech)

# Licensing

This tool is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/LICENSE) for the full license text.


