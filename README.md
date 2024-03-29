This repo uses aws-cdk-python and golang for lambdas


## Python tools to install:
- [uv for virtualenv](https://github.com/astral-sh/uv)
- cdk for aws IaC (npm install -g aws-cdk)
- ruff and mypy for static code checking


### uv:

```bash
uv venv
source .venv/bin/activate
pip install -r requirements.txt -r requirements-dev.txt
```

### ruff:
```bash
ruff .
```

### mypy:

```bash
mypy .
```


## For lambdas golang code:

To check complilation for a lambda, go to the lambda subdirectory and run:

```bash
go build -o main main.go
```

### Go format:
```bash
gofmt -s -w .
```



DOMAINS:

- keycharm.io
- keyify.io
- keylite.io