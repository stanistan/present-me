description = "Trunk is a CI runner."
binaries = [ "trunk" ]
test = "trunk --version"

linux {
  source = "https://trunk.io/releases/${version}/trunk-${version}-linux-x86_64.tar.gz"
}

darwin {
  source = "https://trunk.io/releases/${version}/trunk-${version}-linux-x86_64.tar.gz"
}

version "1.9.1" {
  auto-version {
    github-release = "getzola/zola"
  }
}
