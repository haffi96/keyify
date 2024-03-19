# Use a base image with Python 3.11
FROM python:3.11

# Install dependencies for building Go and Neovim
RUN apt-get update && apt-get install -y \
  wget \
  curl \
  git \
  build-essential \
  cmake \
  unzip \
  zsh


# Install zim
RUN curl -fsSL https://raw.githubusercontent.com/zimfw/install/master/install.zsh | zsh

# Set default shell to Zsh
SHELL ["/bin/zsh", "-c"]

# Install Go
ENV GOLANG_VERSION 1.22.1
RUN wget -q https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
  rm go$GOLANG_VERSION.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"

# Install Node.js LTS
RUN curl -sL https://deb.nodesource.com/setup_lts.x | bash -
RUN apt-get install -y nodejs

# Install Neovim
RUN curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz && \
  rm -rf /opt/nvim && \
  tar -C /opt -xzf nvim-linux64.tar.gz && \
  rm nvim-linux64.tar.gz

# install uv
RUN curl -LsSf https://astral.sh/uv/install.sh | sh


# Add PATH to nvim binary to .zshrc
RUN echo 'export PATH="$PATH:/opt/nvim-linux64/bin"' >> ~/.zshrc

# Set entrypoint to zsh
ENTRYPOINT ["/bin/zsh"]
