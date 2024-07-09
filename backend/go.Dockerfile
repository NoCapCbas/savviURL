FROM golang:1.22
# declare working dir
WORKDIR /app 
# copy all go files
COPY go.* ./
# download dependencies
RUN go mod download
# copy remaining files
COPY . .


# Live Hot Reloading using Air (https://github.com/air-verse/air)
# binary will be $(go env GOPATH)/bin/air
# RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# # Install Neovim
# RUN apt-get update && \
#     apt-get install -y neovim && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

CMD ["go", "run", "main.go"]
