FROM scratch

COPY build/scraparr /scraparr
ENTRYPOINT [ "/scraparr" ]