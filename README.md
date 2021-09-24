A simple wallet manager of lotus.

Refer to https://github.com/filecoin-project/lotus

# Build
Reference to https://docs.filecoin.io/get-started/lotus/installation/
Mac  
```
brew install go bzr jq pkg-config rustup hwloc
```

Ubuntu/Debian:
```
sudo apt install mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev wget -y && sudo apt upgrade -y
```

Bin
```
make  
./lotus-wallet --help
```

# TODO
* Add a MIT lisence 
* Build encrypt wallet
* Send filecoin message to the network.


