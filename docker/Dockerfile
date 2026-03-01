# Vox Builder Image
# Cross-compiles Vox for all 6 target platforms
# Inherits Zig cross-compilation and Go from powertech-center/alpine-go

FROM ghcr.io/powertech-center/alpine-go:latest

# Project-specific: Linux x64 dev libraries for CGo (ALSA, GTK, systray)
RUN apk add --no-cache \
    xvfb \
    alsa-lib-dev \
    libx11-dev \
    gtk+3.0-dev \
    libayatana-appindicator-dev \
    musl-fts-dev

# Project-specific: ARM64 sysroot with dev libraries for Linux ARM64 cross-compilation
RUN mkdir -p /opt/aarch64-sysroot/etc/apk && \
    cp /etc/apk/repositories /opt/aarch64-sysroot/etc/apk/ && \
    cp -r /etc/apk/keys /opt/aarch64-sysroot/etc/apk/ && \
    apk --arch aarch64 --root /opt/aarch64-sysroot \
    --allow-untrusted --initdb --no-cache \
    add alsa-lib-dev libx11-dev gtk+3.0-dev libayatana-appindicator-dev zlib-dev && \
    rm -f /opt/aarch64-sysroot/usr/include/stdc-predef.h

WORKDIR /workspace

CMD ["/bin/bash"]

