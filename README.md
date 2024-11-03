# pocky

Containers experiment.

## Build

You can use `make` utility to build project

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
