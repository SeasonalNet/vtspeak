FROM ubuntu:24.04

ARG DEBIAN_FRONTEND=noninteractive
ARG GHIDRA_VERSION=12.1.4
ARG GHIDRA_BUILD_DATE=20260921
ARG GHIDRA_SHA256=ddac49f903da9d5bac833e5cc79395098b9c33cfd3279be5f31bd00387d2d4db
ARG CAPA_RULES_VERSION=9.4.0
ARG CAPA_RULES_SHA256=79e37bb648dd7a912ff49e8d89bff022ea156d2567f8fe5080f9d6d03067ee02
ARG CAPA_SIGS_VERSION=9.4.0
ARG CAPA_SIGS_SHA256=98867ceeb2f0c530970759fafaa706aae280bf9e8d290fbfe02f139ce8f214d1

ENV GHIDRA_HOME=/opt/ghidra \
    CAPA_RULES_DIR=/opt/capa-rules \
    CAPA_SIGS_DIR=/opt/capa-sigs \
    JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64 \
    PATH=/opt/capa-venv/bin:/opt/floss-venv/bin:/opt/ghidra/support:$PATH \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        binutils \
        build-essential \
        ca-certificates \
        curl \
        file \
        openjdk-21-jre-headless \
        python3 \
        python3-dev \
        python3-pip \
        python3-venv \
        unzip \
    && rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/floss-venv \
    && /opt/floss-venv/bin/pip install --no-cache-dir flare-floss==3.1.1 \
    && python3 -m venv /opt/capa-venv \
    && /opt/capa-venv/bin/pip install --no-cache-dir \
        flare-capa==9.4.0 \
        pefile==2024.8.26

RUN set -eu; \
    archive="capa-rules-${CAPA_RULES_VERSION}.zip"; \
    curl -fsSL "https://github.com/mandiant/capa-rules/archive/refs/tags/v${CAPA_RULES_VERSION}.zip" -o "/tmp/${archive}"; \
    echo "${CAPA_RULES_SHA256}  /tmp/${archive}" | sha256sum -c -; \
    unzip -q "/tmp/${archive}" -d /opt; \
    mv "/opt/capa-rules-${CAPA_RULES_VERSION}" "${CAPA_RULES_DIR}"; \
    rm -f "/tmp/${archive}"

RUN set -eu; \
    archive="ghidra_${GHIDRA_VERSION}_PUBLIC_${GHIDRA_BUILD_DATE}.zip"; \
    curl -fsSL "https://github.com/NationalSecurityAgency/ghidra/releases/download/Ghidra_${GHIDRA_VERSION}_build/${archive}" -o "/tmp/${archive}"; \
    echo "${GHIDRA_SHA256}  /tmp/${archive}" | sha256sum -c -; \
    unzip -q "/tmp/${archive}" -d /opt; \
    mv "/opt/ghidra_${GHIDRA_VERSION}_PUBLIC" "${GHIDRA_HOME}"; \
    rm -f "/tmp/${archive}"

RUN apt-get update \
    && apt-get install -y --no-install-recommends openjdk-21-jdk-headless \
    && rm -rf /var/lib/apt/lists/*

RUN set -eu; \
    archive="capa-${CAPA_SIGS_VERSION}.zip"; \
    curl -fsSL "https://github.com/mandiant/capa/archive/refs/tags/v${CAPA_SIGS_VERSION}.zip" -o "/tmp/${archive}"; \
    echo "${CAPA_SIGS_SHA256}  /tmp/${archive}" | sha256sum -c -; \
    mkdir -p "${CAPA_SIGS_DIR}"; \
    unzip -q -j "/tmp/${archive}" "capa-${CAPA_SIGS_VERSION}/sigs/*.sig" -d "${CAPA_SIGS_DIR}"; \
    rm -f "/tmp/${archive}"

WORKDIR /work
CMD ["/bin/bash"]
