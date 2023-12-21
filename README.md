# GitHub Proxy

Start a GitHub proxy

## Requirements

1.  Github proxy requires writable permission to `/etc/hosts`

    ```sh
    sudo chmod a+w /etc/hosts
    ```

2.  Make self signed certifications

    ```sh
    mkcert github.com gist.github.com gist.githubusercontent.com raw.githubusercontent.com 
    ```

3.  Rename certifications to `github.com.pem` and `github.com-key.pem` if you don't want to add them in the command line

## Quick start

1.  Start github proxy

    ```sh
    Usage:
    github-proxy [flags]

    Flags:
        --cert-file string   Certificate for the proxy (default "github.com.pem")
        --key-file string    Private key for the proxy (default "github.com-key.pem")
        --proxy string       Proxy target (default "https://mirror.ghproxy.com")
    ```

    Github proxy will use port `443`, try `sudo github-proxy` if you start failed.

2.  Perform a request

    ```sh
    curl https://raw.githubusercontent.com/fatindeed/github-proxy/master/README.md
    ```
