FROM gcr.io/distroless/static

COPY airmail /airmail

EXPOSE 9900

ENTRYPOINT [ "/airmail" ]
