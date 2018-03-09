require 'rbconfig'
class Hclq < Formula
  desc "CLI tool for querying and modifying HashiCorp HCL files, similar to jq"
  homepage "https://github.com/mattolenik/hclq"
  version ""

  case RbConfig::CONFIG['host_os']
  when /darwin|mac os/
    url "https://github.com/mattolenik/hclq/releases/download//"
    sha256 ""
  when /linux/
    url "https://github.com/mattolenik/hclq/releases/download//"
    sha256 ""
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
