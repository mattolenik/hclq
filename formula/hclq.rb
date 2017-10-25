class Hclq < Formula
  desc "Small tool for querying and modifying HCL files, similar to jq"
  homepage "https://github.com/mattolenik/hclq"
  url "https://github.com/mattolenik/hclq", :tag => "0.1.1"

  depends_on "go" => :build

  def install
    system "make", "build-brew"
    bin.install "dist/hclq"
  end

  test do
    assert_equal "true", pipe_output("#{bin}/hclq get test", "test = true")
  end
end
