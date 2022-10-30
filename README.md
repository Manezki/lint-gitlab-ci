# Gitlab CI script linter

Lint the script sections of Gitlab CI configurations (.gitlab-ci.yml) 

## Docker
Docker runtime can be build with:

```
docker build -t lint-gitlab-ci .
```

Using the docker runtime requires mounting the `.gitlab-ci.yml` into to container. E.g:

```
docker run -v PATH_TO_GITLAB_CI_YML:/mnt/ lint-gitlab-ci /mnt/.gitlab-ci.yml
```

## Dependencies

* [**shellcheck**](https://github.com/koalaman/shellcheck#installing) needs to be installed in the system

## Running

```bash
go run main.go PATH_TO_GITLAB_CI
```
