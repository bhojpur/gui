FROM moby/buildkit:v0.9.3
WORKDIR /gui
COPY gui README.md /gui/
ENV PATH=/gui:$PATH
ENTRYPOINT [ "/bhojpur/gui" ]