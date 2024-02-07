FROM scratch
#LABEL authors="developer"
ADD build/ca-certificates.crt /etc/ssl/certs/
ADD build/Moscow /etc/localtime
ADD build/main /
CMD ["/main"]

#ENTRYPOINT ["top", "-b"]