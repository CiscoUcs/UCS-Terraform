# Terraform UCS Provider

## Overview of solution

The major goal of this project is to be able to **C**reate, **R**ead, **U**pdate & **D**estroy
different kind of resources in UCS, from Service Profiles to Boot Policies, Network Policies,
Service Profile Templates, etc. For such purpose we have made a Terraform provider.
Visit Terraform's [official website](https://terraform.io) for more information on what
exactly is it and how it works.

On a **feature level**, this will allow the user to:

  * Use Terraform to seamlessly deploy an arbitrary amount of Cisco machines into UCS at scale.
  * Specify predefined service profiles from Terraform.

Given a configuration file whose purpose is to create a new Service Profile in UCS, the way it all clicks
together looks more or less like this:

  1. `terraform-provider-ucs` talks to UCS, requesting to create a new Service Profile.
  2. Once the Service Profile has been created the provider then fetches its MAC address.
  3. A `cobbler system` is then created (also via Terraform as there is a `terraform-provider-cobbler`),
     passing along the Service Profile as its hostname, an IP address (read on to learn how the IP is generated)
     and its MAC address.

As you can see in the steps aboved described the process is pretty straight forward. The only problem here is that
as soon as the Service Profile is created the machine that gets associated immediately boots up so there could be
a moment where the machine is up looking for some PXE server to provide an installation kernel and the machine can
not find it. But as of today this really is not an issue as the machine must first load Cisco's propietary system
first, which takes about 15 minutes in average. There have been discussions with the Hardware team at Cisco to see
if this initial time can be reduced but nothing concise has come out of it.  When that time decreases to milliseconds
then we will need to think of a better solution.

It is also important to point out that steps 2 and 3 are completely optional. The user is free to chose a different way
to bootstrap their servers. We chose Cobbler because its relative ease of use and also because it allows integration
with the `terraform-provider-cobbler`. The machine where Cobbler is running should have the same network access to UCS
as the machine where the `terraform-provider-ucs`

### How the IP is generated?
It's a combination between the `inventory` file and the `CIDR` configuration field.
If the UCS setup is pristine, that is, if there is no inventory file whatsoever then an IP will be generated from
the given `CIDR`, the inventory file will be created and the IP will be appended to it.
If there is an existing inventory file with at least an IP in it then that IP will be interpreted as the last IP
known to the system and the new Service Profile that is being created will have the IP that follows the previous one.

**Important note: The `inventory` file should not be manually modified as it may lead to odd behaviour of
the provider.**

## Roadmap

### Features

  - [x] Add cross-platform compilation task in Makefile.
  - [ ] Add support for configuring multiple network interfaces.
  - [ ] Automatically push dist files to Github.
  - [ ] Improve error message logging. Currently errors are only dumped into the log file,
        which makes things difficult to debug. Better if we could report them on the same
        screen from which Terraform was ran on.
  - [ ] Extract UCS client into its own library and repository.
  - [ ] Remove unused fields from the different `struct` types that we have in the codebase.
        Right now there are a bunch of fields being returned by UCS that we don't really need
        but we still map in memory but perhaps we can just get rid of that to avoid dust
        piling up in the code.

### Dev Tools
  - [x] Vagrant.

### Bugs

## System Dependencies
  * [Terraform]

---

Currently the supported arguments are:

* UCS provider parameters:

  * **ip_addres** the IP where the UCS manager service is running on.
  * **username** username used for authentication.
  * **password** password used for authentication.
  * **log_level** default: 1.  
    Valid values:  
      * 0 (TRACE)
      * 1 (DEBUG)
      * 2 (INFO)
      * 3 (WARN)
      * 4 (ERROR)
      * 5 (FATAL)
  * **log_filename** default: stderr.

    Example:

    ```
    provider "ucs" {
      ip_address   = "1.2.3.4"
      username     = "john"
      password     = "supersecret"
      log_level    = 6
      log_filename = "terraform.log"
    }
    ```

* Service Profile:

  * **name** the name of the Service Profile.
  * **target_org** the target organization of the Service Profile.
  * **service_profile_template** the Service Profile Template of the Service Profile.

    Example:

    ```
    resource "ucs_service_profile" "master-server" {
      name                     = "master-server"
      target_org               = "some-target-org"
      service_profile_template = "some-template"
      metadata { # This field is pretty much free style. Values must always be strings.
        role             = "master" # This is useful when creating a Mantl cluster
        ansible_ssh_user = "root"
        foo              = "bar"
      }
    }
    ```

## How to compile and install

The only development dependency are:

  * [Git] (obviously)
  * [Go]
  * [Terraform]

If you are not familiar with the folder structure for Go projects, check out
[this](http://golang.org/doc/code.html#Organization) document.

There is a Makefile with several tasks to ease the development process.
When you are starting out with this project, use the task `make bootstrap` to
get the required dependencies as well as the required Terraform's configuration file.

The default make task will `test` the project, `clean` any distribution file and
`build` the project.

## Known Issues
  * Mounting the ISO image inside the Docker's HTTP Server image
    will fail if Docker is using anything different than `devicemapper`
    as the storage driver.

[Terraform]: https://www.terraform.io/
[Docker]: https://www.docker.com/
[Vagrant]: https://www.vagrantup.com/
[Go]: https://www.golang.org/
[Git]: https://git-scm.com/
