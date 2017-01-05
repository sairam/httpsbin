Making a replica of requestb.in as httpsb.in in Golang

## Running the service
* `go run main.go config.toml`
* `config.toml` is optional. You can load from env variables as well
## Alternative
* Set environment variables `set -a ; . .env ; set +a`
* `go run main.go`

## HomeBrew install

### Installation

  brew install sairam/formulae/httpsbin

### Starting

* `httpsbin` command to start the server

### Customisation

If you'd like to customise the port / hostname, you can modify at `/usr/local/etc/httpsbin`.

You need to restart the server to reflect any changes.

Data Directory exists at `/usr/local/var/httpsbin/data/` . Changing the `data_dir` attribute will modify the directory

### Upgrading
NOTE: Config and `tmpl/` files are overridden when you upgrade. You will not lose any data
