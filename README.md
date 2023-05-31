# Terraform provider for NeuVector

![build](https://github.com/theobori/terraform-provider-neuvector/actions/workflows/build.yml/badge.svg)
![tests](https://github.com/theobori/terraform-provider-neuvector/actions/workflows/tests.yml/badge.svg)


## 📖 Build and run

You only need the following requirements:

-  [Terraform](https://www.terraform.io/downloads.html) 1.4.6+

To build and install the terraform plugin, you should run the following command.
Override the `OS_ARCH` environment variable if needed (the default one is in `Makefile`.

```bash
OS_ARCH="linux_adm64" \
make install
```

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🧪 Tests

There are some Terraform acceptance tests, it take the configuration from the files in `examples/`.
To run the tests, feel free to use the Docker allinone instance by running `make neuvector`, once it is done run `make testacc`.

If you want to override the default variables:

```bash
NEUVECTOR_BASE_URL="url" \
NEUVECTOR_USERNAME="username" \
NEUVECTOR_PASSWORD="password" \
make testacc
```

## 🎉 Tasks

- [x] Acceptance tests
- [x] Documentation
- [x] Registry update 
- [x] Admission rule update 
- [x] Supports `terraform import`
- [ ] resource service
