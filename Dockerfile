FROM scratch
#LABEL authors="developer"
ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
CMD ["/main"]

#ENTRYPOINT ["top", "-b"]