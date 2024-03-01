description = "standalone executable for tailwind"
binaries = [ "tailwindcss" ]
test = "tailwindcss --version"
dont-extract = true

// https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
// https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-macos-arm64

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
  dont-extract = true
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

version "3.4.1" {
}
