# regenv

A command-line tool for updating Windows environment variables.

## Description

This tool allows you to easily append or remove values from system or user environment variables. It is useful for scripting and automating updates to variables like `Path`.

## Usage

### Flags

- `-scope`: The scope to read variables from. Can be `machine` or `user`. Defaults to `machine`.
- `-append`: The variable to append to.
- `-remove`: The variable to remove from.
- `-value`: The value to append or remove.
- `-separator`: The separator to use. Defaults to `;`.

### Examples

#### Append to a variable

To append a value to the system `Path` variable:

```
go run main.go -scope machine -append Path -value "C:\new\path"
```

#### Remove from a variable

To remove a value from the user `Path` variable:

```
go run main.go -scope user -remove Path -value "C:\old\path"
```

## Build

To build the executable, run:

```
go build
```
