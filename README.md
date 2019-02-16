# Mergo

### Presentation

Have you ever been in the situation where you code in a group project, pushed your work on your remote branch and then comes the moment where you need to create a pull request on your git client in order for your code to be be present on the main branch.  

You have to leave your terminal to go on a web interface and then click several button just to create a "github like pull request". This is very annoying and I have always thought that it would be great if there was a way to create a pull request without having to leave my terminal.  

This is the purpose of mergo.
### Install
```bash
go get -u gitlab.com/jfaucherre/mergo
```
### Usage
```bash
Usage:
  mergo [OPTIONS]

Application Options:
  -d, --head=       The head branch you want to merge into the base
  -b, --base=       The base branch you want to merge into (default: master)
      --host=       The git host you use, ie github, gitlab, etc.
      --remote=     The remote to use (default: origin)
      --repository= The name of the repository on which you want to make the pull request
      --owner=      The owner of the repository

Help Options:
  -h, --help        Show this help message
```
### Personal config
As I like to keep all things in one place I have run the following command
```bash
git config --global alias.pr '!mergo'
```
So that I can then run
```bash
git pr
```
just after I have pushed code on my repository