# GoTorWeb

This project as has the objective of providing an easy and automatic way of eposing your containers to the Tor network.

## Running GoTorWeb

```bash
docker run --rm --name gtw  -v /var/run/docker.sock:/var/run/docker.sock  -e "GTW-USERNAME=<auth user>" -e "GTW-PASSWORD=<auth password>" -v <config location>:/config knoker/go_tor_web
```

### Runtime Parameters

* `GTW-USERNAME`: web ui username
* `GTW-PASSWORD`: web ui password
* config location: location of the tor generated certificates and hostnames

### Finding the web ui

After tor has initialized there will be a file `config/local/hostname` containning the web ui address, and the following log will appear in the container stdout

``` text
------------------------------------
Local: something.onion 80
------------------------------------
```

## Other containers

``` bash
docker run --rm -it -l "GTW=1" -l "GWT-PORT=8080" -l "GWT-PORT-FW=80"  hashicorp/http-echo -listen=:8080 -text="hello world"
```

### Other containers parameters

* `-l "GTW=1"`: GoTorWeb enabled

* `-l "GWT-PORT=8080"`: Internal container port

* `-l "GWT-PORT-FW=80"`: Tor exposed port

## TODO

* Improve the web ui's current "design"

* Improve the name discovery routine to run only when a new container is launched and until it's name is discovered

* Allow disabling the web ui

***

Feedback and pull requests are apreciated
