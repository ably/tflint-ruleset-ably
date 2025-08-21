# Ably Custom TFLint Ruleset

## Requirements

- TFLint v0.42+
- Go v1.24

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "ably" {
  enabled = true
  version = "0.1.0"
  source = "github.com/ably/tflint-ruleset-ably"

  signing_key = <<-KEY
  -----BEGIN PGP PUBLIC KEY BLOCK-----

  mQINBGinOwkBEADAER8xJUqQfwgH9ViJ9g38sJEFbaJsUMCp4BJ+kKSx4BIU5Qkh
  +k+VoFBxIinvqzsCwJTgaN2enW/Jnag50QJhrDdheUKaw7EvQsb7wEgzSvKf/uhM
  zFoFZ8IysXL/scyz+57EsSs+PiBS8hdzQJz/swFYl4r+CQZQzsE5HSgA0FvgBa5f
  002ilyWSJja/7vO5tFMrwAGgcNnjPFyA5QtdWrm0Miqo/cmqoVolwJbvYh6DHj0p
  ONQpyGI55Ht5ZwYj45L8gq1ziSLSnMTykYXIz8O3vsOYfTwGQGytNcTNzDcMJnYp
  EojjB9w3pKdFhtGuvQkOl9tBwUZYHkaVLR1NxWaCLVRNAl9VdVM+75PV6mInEnWa
  P5QY4nyiTEaUdmwoL3yfokbFUth8L9AMnt59krks2u9EkZgXBMFeH0cgkQYrIoIc
  LaOxX6ZaQslrwmsczZfM0EAoKnzk0Gfsu40YHyNbLJGmTXONAK+y5qMsii3QQ1EH
  7Gcytj8/b3geVpTbF/7fQr6A1AK1+nJfCf45bzytUQfCa/KblTXugA9/gvOGIpiT
  Hfi+a61BamPsDEwXe8zBAKYVWOcYeAf9GSEJ80En7yqUs4SudGFi7eWKIqJO0W1h
  yW0i70h9RBwf3NHq0DXFTxvAzvj1bvffs19UwH4KV/VkizQl8WIRMZqu+wARAQAB
  tEdBYmx5IFJlYWx0aW1lIEx0ZCAoQWJseSBSZWFsdGltZSBTaWduaW5nIEtleSkg
  PGFjY291bnQub3duZXJzQGFibHkuY29tPokCTgQTAQoAOBYhBBbaxqgy7cIGoiJd
  pOu+ZcbeBDPTBQJopzsJAhsDBQsJCAcCBhUKCQgLAgQWAgMBAh4BAheAAAoJEOu+
  ZcbeBDPTJDAP/2MU6FSvcOAEnCMpoSxc5LDBq/T2sMIt76eM86qO6k02cSABXFj8
  hwOLp8M7tyED3srQO2SR79llfjGUWaGkTF50PuuLJIPZX0cXA04CAL9x8q48DmB/
  2K2YWGVDatR6yWOtBNzvotQ3z3QRwftsoMpRzlVBtoxlaXQ8JTZRMgxqf+47yOp4
  b4L8OGgBVttMxvUob4cppzKqlpF7dU1jM3qA0GlMWlLQfb7cW+7A69izJMlmmpP6
  a/seVde7WIMYG663Mk69j0gK27pj2RIENj8it9hP6g837ozRlltFPg+wnVqY9R2i
  PxVXgh+qT/blQRcUrHoXM+hW3FJ28Er3xu7pktbfUYhaa7+80JXRheqGMgD5wE0I
  BCsZ38jbdZaqf94i6RHxRY9OcAKFoZ8WC67IKoFRNnjkWHYezO1Y5zOxajnhSKrF
  yjyFYZQkcXGH9jRDZXFj7iPI6KNLqF8ilodT214fVp5+WWivXFOrm+YrxvlIeh6m
  oA7H1IaLtpp8b1O7LzL5ZJjMp318uY6lF6AeH58FfCBjS/OKl8VpcK6Pc5QhIi8k
  dJUTDHIfXbhQdzexQyA65lPoGp9I6JqWNE7mMtCuoL7KPqQ8DPYN21CgGnSD2OVk
  MSoC8gjuDj+GeHAMoR/tbxvSMCWMjn3RE5KjRLWezP0yVznr5Ai/TGuA
  =a4cd
  -----END PGP PUBLIC KEY BLOCK-----
  KEY
}
```

## Rules

|Name|Description|Severity|
| --- | --- | --- |
|rightmost_operator_rule|Rule for enforcing required_provider version format: `~> x.y`|WARNING|

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

You can run the built plugin like the following:

```
$ cat << EOS > .tflint.hcl
plugin "ably" {
  enabled = true
}
EOS
$ tflint
```
