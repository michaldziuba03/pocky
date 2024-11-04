> ⚠️ Experimental project, idk if I will develop it seriously.

# pocky

Pocky - lightweight containers ⚓

![image](https://github.com/user-attachments/assets/8ec1e6ab-2e64-4a81-9a85-7603a3288dfd)


> Basic container with Alpine Linux (I use Ubuntu WSL as host btw).

## Build

You can use `make` utility to build project

First run may require configuring `resolve.conf`. Currently it shares all networking with host (works like `host` parameter in Docker)

```shell
make
cd build
# currently it always installs alpine 
sudo ./pocky download alpine
# run any binary on installed distro
sudo ./pocky run /bin/sh
```

## License

Distributed under the MIT License. See `LICENSE` for more information.
