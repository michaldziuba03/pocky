> ⚠️ Experimental project, idk if I will develop it seriously.

# pocky

Pocky - lightweight, pocket containers ⚓

![image](https://github.com/user-attachments/assets/8ec1e6ab-2e64-4a81-9a85-7603a3288dfd)


> Basic container with Alpine Linux (I use Ubuntu WSL as host btw).

## TODO:

- [x] Ability to run container (process with specific **namespace** flags)
- [x] Download Alpine Linux rootfs and chroot to it (currently acts like **"image"**).
- [ ] Configurable limitations (via **cgroups**)
- [ ] Experiment with `pivot_root` over `chroot` as more secure alternative (or make it configurable)
- [ ] Implement something what acts as actual, configurable images
- [ ] More networking seperation options, bridge etc
- [ ] Real CLI interface and maybe REST API

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
