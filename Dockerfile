FROM golang:1.19.0 as build-env
LABEL stage="builder"

WORKDIR /app
COPY . /app/

RUN go build -o lint-gitlab-ci .

FROM koalaman/shellcheck:v0.8.0
LABEL maintainer="Janne Holopainen <manezki@gmail.com>"
COPY --from=build-env /app/lint-gitlab-ci /bin/
ENTRYPOINT [ "/bin/lint-gitlab-ci" ]
