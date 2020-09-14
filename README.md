# Mergo

### Presentation

This tool comes from my need not to leave my terminal and that, although there is the hub client, you can not create a pull request for any git host you'd like from your terminal  
So this binary creates pull request from terminal, it has been thought so that it doesn't have your specific host, you can add it easily  
Its behavior is trying to be as close as possible to the git commands behavior and its default values try to follow the github pull request defaults  

### Install
```bash
go get -u gitlab.com/jfaucherre/mergo
```
### Usage
```bash
mergo --help
Usage:
  mergo [OPTIONS]

Application Options:
  -v, --verbose        Add logs, you can have more logs by calling it more times
  -d, --head=          The head branch you want to merge into the base (default: the actual checked out branch)
  -b, --base=          The base branch you want to merge into (default: master)
  -m, --message=       The pull request message (default: If you have only one commit, it takes this commit's message)
  -f, --force          Force the pull request, doesn't ask you if you have unstaged changes or things like that
  -c, --clipboard      Copy the URLs of your merge requests to your clipboard
      --remote=        The remote to use (default: origin)
  -r, --remote-url=    The remote URLs to use. Note that this overwrite the "remote" option
  -e, --force-edition  Force the edition of the message event it already have a value
      --delete-creds=  Use this option when you want to delete the credentials of an host

Help Options:
  -h, --help           Show this help message
```

```bash
# placed in a repository on a branch you want to pull into master
mergo -m "My first pull request with mergo !"
```

### Configuration

You can configure the default behavior of the binary through the git config by writing to the `mergo` section  
For example suppose the default branch you want to merge into is not master but staging you can write
```bash
git config add mergo.base staging
```
After that all the pull requests you are going to create with mergo are going to be merged into staging and not master

All the configurable values are:
  * head
  * base
  * force
  * clipboard
  * remote
  * remote-urls
  * force-edition
  * verbose

Note that if you want to configure verbose for a certain level (from 1 to 5) you must give an array of boolean and not an int, because of how [the argument parsing lib](https://github.com/jessevdk/go-flags#example) works:
```bash
# set log level to 3
git config add mergo.verbose true,true,true
```

### Personal config
In order to use it as any other git command I made an alias of it in git like that
```bash
git config --global alias.pr '!mergo'
```
So that I can then run
```bash
git pr
```

### Support

Mergo actually support:
- github
- gitlab

Do not hesitate to propose pull request to support more git hosts
