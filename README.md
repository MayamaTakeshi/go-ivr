# go-ivr


## overview
An IVR platform written in golang using freeswitch

This is a work in progress. There is nothing to see yet.


## host preparation
We will need a server running freeswitch.

In debian 10 you can install freeswitch this way:
```
apt-get update && apt-get install -y gnupg2 wget lsb-release

wget -O - https://files.freeswitch.org/repo/deb/debian-release/fsstretch-archive-keyring.asc | apt-key add -

apt-get update && apt-get install -y freeswitch-meta-all
```

Then set freeswitch to use our configuration files
```
mv /etc/freeswitch /etc/freeswitch_old

mkdir /etc/freeswitch
cp test/artifacts/freeswitch_configuration/* /etc/freeswitch
```

And then you can start freeswitch with
```
service freeswitch start
```

## installing golang
In the same server we will need golang installed. Do:
```
# install asdf (a tool that permits to install/manage different versions of apps/tools like node, golang, rust, java, scala etc via plugins):
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.9.0

# add line in .bashrc to enable asdf on login:
cat >> ~/.bashrc <<EOF
. $HOME/.asdf/asdf.sh
EOF

# enable asdf in the current session via .bashrc
. ~/.bashrc

# add asdf golang plugin:
asdf plugin-add golang

# install golang
asdf install golang 1.17.6

# set golang version to use:
asdf global golang 1.17.6

# test golang installation:
go version
```
## preparing testing framework
In the same server do:
```

apt install build-essential automake autoconf libtool libspeex-dev libopus-dev libsdl2-dev libavdevice-dev libswscale-dev libv4l-dev libopencore-amrnb-dev libopencore-amrwb-dev libvo-amrwbenc-dev libopus-dev libsdl2-dev libopencore-amrnb-dev libopencore-amrwb-dev libvo-amrwbenc-dev libboost-dev libspandsp-dev libpcap-dev libssl-dev uuid-dev

cd tests/functional
apt install jq
nvm install `jq -r .engines.node package.json`
nvm use `jq -r .engines.node package.json`

git clone https://github.com/MayamaTakeshi/bcg729
cd bcg729
git checkout faaa895862165acde6df8add722ba4f85a25007d
cmake . 
make
make install
ldconfig

npm install

```
## testing
We don't have tests for the ivr engine yet.
But to test if the test infra itself is OK do:
```
node first.js
```

