description = "standalone executable for tailwindcss"
binaries = [ "tailwindcss" ]
test = "tailwindcss -h"
dont-extract = true

darwin {
  source = "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-macos-arm64"
  on "unpack" {
    rename {
      from = "${root}/tailwindcss-macos-arm64"
      to = "${root}/tailwindcss"
    }
    chmod {
      file = "${root}/tailwindcss"
      mode = 448
    }
  }
}

linux {
  source = "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-linux-arm64"
  on "unpack" {
    rename {
      from = "${root}/tailwindcss-linux-arm64"
      to = "${root}/tailwindcss"
    }
    chmod {
      file = "${root}/tailwindcss"
      mode = 448
    }
  }
}

version "3.4.1" { }
