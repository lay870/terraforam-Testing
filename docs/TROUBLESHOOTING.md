# Troubleshooting

## Logging

The language server produces detailed logs which are send to stderr by default.
Most IDEs provide a way of inspecting these logs when server is launched in the standard
stdin/stdout mode.

Logs can also be redirected into file using flags of the `serve` command, e.g.

```sh
$ terraform-ls serve \
	-log-file=/tmp/terraform-ls-{{pid}}.log \
	-tf-log-file=/tmp/tf-exec-{{lsPid}}-{{args}}.log
```

It may be helpful to share these logs when reporting bugs.

### How To Share Logs

It is recommended to avoid pasting logs into the body of an issue,
unless you are trying to draw attention to a selected line or two.

It's always better to upload the log as [GitHub Gist](https://gist.github.com/)
and attach the link to your issue/comment, or [attach the file to your issue/comment](https://docs.github.com/en/github/managing-your-work-on-github/file-attachments-on-issues-and-pull-requests).

### Sensitive Data

Logs may contain sensitive data (such as content of the files being edited in the editor).
If you consider the content sensitive you may PGP encrypt it using [HashiCorp's key](https://www.hashicorp.com/security#secure-communications)
to reduce the exposure of the sensitive data to HashiCorp.

### Log Rotation

Keep in mind that the language server itself does not have any log rotation facility,
but the destination path will be truncated before being logged into.

Static paths may produce large files over the lifetime of the server and
templated paths (as described below) may produce many log files over time.

### Log Path Templating

Log paths support template syntax. This allows for separation of logs while accounting for:

 - multiple server instances
 - multiple clients
 - multiple Terraform executions which may happen in parallel

**`-log-file`** supports the following functions:

 - `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
 - `pid` - process ID of the language server
 - `ppid` - parent process ID (typically editor's or editor plugin's PID)

 **`-tf-log-file`** supports the following functions:

  - `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
  - `lsPid` - process ID of the language server
  - `lsPpid` - parent process ID of the language server (typically editor's or editor plugin's PID)
  - `method` - [`terraform-exec`](https://pkg.go.dev/github.com/hashicorp/terraform-exec) method (e.g. `Format` or `Version`)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/terraform-ls-{{timestamp}}.log`.

## CPU Profiling

If the bug you are reporting is related to high CPU usage it may be helpful
to collect and share CPU profile which can be done via `cpuprofile` flag.

For example you can modify the launch arguments in your editor to:

```sh
$ terraform-ls serve \
	-cpuprofile=/tmp/terraform-ls-cpu.prof
```

The target file will be truncated before being written into.

### Path Templating

Path supports template syntax. This allows for separation of logs while accounting for multiple server or client instances.

**`-cpuprofile`** supports the following functions:

 - `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
 - `pid` - process ID of the language server
 - `ppid` - parent process ID (typically editor's or editor plugin's PID)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/terraform-ls-cpuprofile-{{timestamp}}.log`.

## Memory Profiling

If the bug you are reporting is related to high memory usage it may be helpful
to collect and share memory profile which can be done via `memprofile` flag.

For example you can modify the launch arguments in your editor to:

```sh
$ terraform-ls serve \
	-memprofile=/tmp/terraform-ls-mem.prof
```

The target file will be truncated before being written into.

### Path Templating

Path supports template syntax. This allows for separation of logs while accounting for multiple server or client instances.

**`-memprofile`** supports the following functions:

 - `timestamp` - current timestamp (formatted as [`Time.Unix()`](https://golang.org/pkg/time/#Time.Unix), i.e. the number of seconds elapsed since January 1, 1970 UTC)
 - `pid` - process ID of the language server
 - `ppid` - parent process ID (typically editor's or editor plugin's PID)

The path is interpreted as [Go template](https://golang.org/pkg/text/template/), e.g. `/tmp/terraform-ls-memprofile-{{timestamp}}.log`.

## "No root module found for ... functionality may be limited"

Most of the language server features depend on initialized root modules
(i.e. folder with `*.tf` files where you ran `terraform init` successfully).
and server's ability to discover them within the hierarchy and match them
with files being open in the editor.

This functionality should cover many hierarchies, but it may not cover yours.
If it appears that root modules aren't being discovered or matched the way
they should be, it can be useful to use `inspect-module` to obtain
the discovery results and provide them to maintainers in a bug report.

Point it to the same directory that you tried to open in your IDE/editor
and wait for the output - it may take some seconds or low minutes
depending on the complexity of your hierarchy and number of root modules in it.

```
$ terraform-ls inspect-module /path/to/dir
```

## "Unable to retrieve schemas for ..."

The process of obtaining the schema currently requires access to the state,
which in turn means that if the code itself doesn't have enough context
to obtain the state and/or there isn't context available from config file(s)
in standard locations you may need to provide that extra context.

See https://github.com/hashicorp/terraform-ls/issues/128 for more.
