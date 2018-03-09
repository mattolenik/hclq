require 'rbconfig'
class Hclq < Formula
  desc "CLI tool for querying and modifying HashiCorp HCL files, similar to jq"
  homepage "https://github.com/mattolenik/hclq"
  version "{{VERSION}}"

  case RbConfig::CONFIG['host_os']
  when /darwin|mac os/
    url "https://github.com/mattolenik/hclq/releases/download/{{VERSION}}/{{DARWIN_FILENAME}}"
    sha256 "{{DARWIN_HASH}}"
  when /linux/
    url "https://github.com/mattolenik/hclq/releases/download/{{VERSION}}/{{LINUX_FILENAME}}"
    sha256 "{{LINUX_HASH}}"
  else
    :unknown
  end

  def install
    bin.install "hclq"
  end

  test do
    system "hclq --version"
  end

end
