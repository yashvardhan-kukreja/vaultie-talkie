#!/bin/bash

function main() {
     curl -L "$(download_url)" | tar xvz vaultie-talkie
}

function download_url() {
     tag=$(latest_release)
     kernel=$(system_kernel)
     arch=$(system_arch)

     echo $tag
     echo $kernel
     echo $arch

     echo "https://github.com/yashvardhan-kukreja/vaultie-talkie/releases/download/v${tag}/vaultie-talkie_${tag}_${kernel}_${arch}.tar.gz"
}

function latest_release() {
     curl -s "https://api.github.com/repos/yashvardhan-kukreja/vaultie-talkie/releases/latest" \
          | jq .tag_name \
          | grep -o '[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+'
}

function system_kernel() {
     case "$(uname -s)" in
          "Linux")
          echo "linux"
          ;;
          "Darwin")
          echo "darwin"
          ;;
          *)
          echo ""
          ;;
     esac
}

function system_arch() {
     case "$(uname -m)" in
          "x86_64")
          echo "amd64"
          ;;
          .*386.*)
          echo "386"
          ;;
          "armv8")
          echo "arm64"
          ;;
          *)
          echo ""
          ;;
     esac
}

main