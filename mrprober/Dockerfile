FROM scratch

COPY conf/samples/mrprober.yaml /etc/mrprober/mrprober.yaml
COPY mrprober /bin/mrprober

ENTRYPOINT ["/bin/mrprober"]
CMD ["serve"]
