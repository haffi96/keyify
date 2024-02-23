FROM golang:1.20 as build
# Copy dependencies list
COPY go.mod go.sum /home/handler/
# Build with optional lambda.norpc tag
COPY main.go /home/handler/
WORKDIR /home/handler
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /home/handler/main ./main
ENTRYPOINT [ "./main" ]