FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD main /

ARG GIT_COMMIT=unknown
LABEL git-commit=$GIT_COMMIT
ARG GIT_BRANCH=unknown
LABEL git-branch=$GIT_BRANCH
ARG BUILD_TIME=unknown
LABEL build_time=$BUILD_TIME

CMD ["/main"]
