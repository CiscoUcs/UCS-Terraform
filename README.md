# UCS Terraform Provider

## Overview of solution

The major goal of this project is to be able to Create, Read, Update & Destroy different kind of resources in UCS, from Service Profiles to Boot Policies, Network Policies, Service Profile Templates, etc. For such purpose we have made a Terraform provider. Visit [Terraform's official website](https://terraform.io) for more information on what exactly terraform it and how it works.

On a feature level, this will allow the user to:

- Use Terraform to seamlessly deploy an arbitrary amount of Cisco machines into UCS at scale.
- Specify predefined service profiles from Terraform.

Given a configuration file whose purpose is to create a new Service Profile in UCS, the way it all works is the ```terraform-provider-ucs``` talks to UCS, requesting to create a new Service Profile

This terraform provider will be enhance and iterated upon when we understand further requirement from the DevNet community. One of the first suggested additions is using [Cobbler](http://cobbler.github.io/) to add the Operating System into the new created UCS Server so please provider feedback if you feel this is interesting to you.

## How the Provider setup process works

- Make bootstrap pulls down the terraform binaries, config .tf,
- The terraform binary file (which is ```$GOPATH/bin``` directory) looks to the ```terraform-provider-ucs``` file (which you will compile in the setup process documented below) in order to interact with the specified resources UCSM API. 
- During this process the default make bootstrap task will download some of the dependencies including the actually terraform binaries, will test, clean any distribution file and build the required project files.
- It is the make build process that actually build out the terraform-provider-ucs file and this file needs to be in the same folder as the terraform binary when executing the actual terraform commands.


## How to compile and install

There are dependencies on the following to environments being setup prior to running the terraform provider creation process:

- Git
- Go

If you are not familiar with the folder structure for Go projects, check out this document.
For Git information https://github.com/



### Git

Make sure GIT is setup and working locally as this will be used during the setup process to clone files in different points in this setup process.



### Running the Go setup process (Mac OS)

You will need to download Go(lang) onto the machine from where you will run the terraform setup process.

You can download it from [here](https://golang.org/dl/)

Use the default setting during the install process for GOlang

#### Create the required directory structure

```
mkdir -p ~/Users/yourlocalusernamehere/go/src/github.com/micdoher
```

#### Setting up the Environment Variables


```
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/Users/yourlocalusernamehere/go
export PATH=$PATH:/Users/yourlocalusernamehere/go/bin
export  PATH=$PATH:/Users/yourlocalusernamehere/go
```
(putting this in the ```.bash_profile``` will probably be help!)


### Clone the Git binaries

Clone the terraform-provider-ucs into ```/Users/yourlocalusernamehere/go/src/github.com/micdoher/``` via the following command: -


```
cd ~/Users/yourlocalusernamehere/go/src/github.com/micdoher/
git clone https://github.com/micdoher/terraform-provider-ucs.git
```

Clone the “go-utils” into ```/Users/yourlocalusernamehere/go/src/github.com/micdoher``` with the following command: -

```
git clone https://github.com/micdoher/GoUtils.git
```


### Compiling and dependency setup


After the terraform provider has been cloned, the resulting directory structure should look like the following: -



Now navigate to:
```
cd /Users/yourlocalusernamehere/go/src/github.com/CiscoCloud/terraform-provider-ucs
```


and run: 

```
make bootstrap
```

An additional file should have been added called ```config.tf``` in the same folder.

Now you need to run the following command from the same directory: 

```
make build
```

The output should return no errors.


Make a copy over a copy of the ```terraform-provider-ucs``` file into the same directory where the terraform binary lives  (normally this is ```/Users/yourlocalusernamehere/go/bin```)


## Running Terraform

Make a copy of the ```config.tf``` file into the same directory where the terraform binary lives  (normally this is ```/Users/yourlocalusernamehere/go/bin```) and customise this to reflect your UCSM environment.


Customize the config.tf file to reflect your desired UCS state.

### UCS provider parameters

* ```ip_address``` the IP where the UCS manager service is running on.
* ```username``` username used for authentication.
* ```password``` password used for authentication.
* ```log_filename``` default: stderr.
* ```log_level``` default: 1.

| log_level | Level amount |
| --- | --- |  
| 0 | TRACE |
| 1 | DEBUG |
| 2 | INFO |
| 3 | WARN | 
| 4 | ERROR |
| 5 | FATAL |
   

#### Example
```
provider "ucs" {
  ip_address  = "1.2.3.4"
  username    = "john"
  password    = "supersecret"
  log_level    = 6
  log_filename = "terraform.log"
}
```

### Service Profile

* ```name``` the name of the Service Profile.
* ```target_org``` the target organization of the Service Profile.
* ```service_profile_template``` the Service Profile Template of the Service Profile.

#### Example

```
resource "ucs_service_profile" "master-server" {
  name                     = “terraserver1"
  target_org               = “root-org"
  service_profile_template = “terraformprofiletemplate"
  metadata { # This field is pretty much free style. Values must always be strings.
    role             = "master" # This is useful when creating a Mantl cluster
    ansible_ssh_user = "root"
    foo              = "bar"
  }
}
```

Make sure you have a pre-defined Service Profile Template available in UCSM for which you align the config.tf with

Once customised, run the following commands in the order given below: 

```
terraform get

terraform plan

terraform apply
```

BTW, Terraform will default to the files in the directory from which it is run with .ft extensions

After the terraform plan you should see the something like following output: 


After the terraform apply you should see the something like following output: 

After terraform apply is run you will find 2 additional files in the ```/Users/yourlocalusernamehere/go/bin``` directory. You should you the xxxx.log file as the primary point of troubleshooting.


If you wish to delete the resources in UCSM you have just created you can run: -

```
terraform destroy
```

## Additional Info for troubleshooting

### Error Messages in Setup

The following error is when the “ucs-terraform-provider, provider file has not been compiled or made available to the terraform binary in ```/Users/yourlocalusernamehere/go/bin``` .

```
Error configuring: 1 error(s) occurred:

* unknown provider "ucs"
```

The following error means the directory structure is not setup correclty: 

```
terraform-provider-ucs-master micdoher$ go build -o terraform-provider-ucs
resource_ucs_service_profile.go:8:2: cannot find package "github.com/CiscoCloud/terraform-provider-ucs/ipman" in any of:
    /usr/local/go/src/github.com/CiscoCloud/terraform-provider-ucs/ipman (from $GOROOT)
    /Users/yourname/terraform/terraform-provider-ucs-master/src/github.com/CiscoCloud/terraform-provider-ucs/ipman (from $GOPATH)
provider.go:4:2: cannot find package "github.com/CiscoCloud/terraform-provider-ucs/ucsclient" in any of:
    /usr/local/go/src/github.com/CiscoCloud/terraform-provider-ucs/ucsclient (from $GOROOT)
    /Users/yourname/terraform/terraform-provider-ucs-master/src/github.com/CiscoCloud/terraform-provider-ucs/ucsclient (from $GOPATH)
```

If you get the following error when running this from a Mac ```xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun```

Try Opening a Terminal, and run the following:

```
xcode-select --install
```

This will download and install xcode developer tools and fix the problem. The problem is that one needs to explicitly agree to the license agreement.

