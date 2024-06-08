# Go Randomify

<p align="center">
    <img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/fc0aa87e909cad008209a08bbf8c19d436730d05.png" height="15%" width="15%" >
</p>
Go Randomify is a lightweight command-line tool written in Go that allows you to generate randomized data within JSON templates. With Go Randomify, you can easily create dynamic JSON files for testing, prototyping,automations and more.

## Features

- **Template Filling**: Easily populate JSON templates with randomized data.
- **Customization**: Easily customize the type and format of randomized data.
- **Output format**: Add output to Json files or pipe it to another tool
- **Current Support**: int, char,uuid.

## Installation

To install Go Randomify, you can use `go get`:

```bash
go get github.com/rohanchavan1918/gorandomify
```

Alternatively, you can clone the repository and build from source:

```bash
git clone https://github.com/rohanchavan1918/gorandomify.git
cd gorandomify
go build
```

Easiest way is to simply download the binary from the releases page.

## Usage

Go Randomify provides a command-line interface for generating randomized JSON files. Here's how to use it:

```
gorandomify -t <template-file> -o <output-file: optional>
```

`-t, --template`: Specify the path to the JSON template file.  
`-o, --output`: Specify the path to the output JSON file (optional), if not passed prints JSON to stdout.

For more options and examples, please refer to the documentation.

## Configuration / Template

### Placeholders

| Type            | Description                                                                              |
| --------------- | ---------------------------------------------------------------------------------------- |
| `$UUID`         | Inserts a UUID                                                                           |
| `$INT`          | Insert a random int (max 10000)                                                          |
| `$INT(MIN:MAX)` | Insert a random int between a lower and upper limit (limit of 10000 is not applied here) |
| `$CHAR(LIMIT)`  | Add random characters upto passed LIMIT                                                  |

### Examples

You can customize the type and format of randomized data using a configuration file. Here's an example configuration file:

```
{
  "uuid": "$UUID",
  "name": "rohan",
  "request_id": "$UUID",
  "some_int": "$INT(100:200)",
  "rand_int": "$INT",
  "rand_char": "$CHAR(100)",
  "data": {
    "lol": "asd",
    "foo": "asad",
    "aas": "asd"
  }
}

```

Above template generates below json

```
{
  "data": {
    "aas": "asd",
    "foo": "asad",
    "lol": "asd"
  },
  "name": "rohan",
  "rand_char": "jXEiPwYZ8rYYN8VEvlsW6f2E42HfYJma9EEEzgRCNO0V7IFxg6f1jg5arrEbPeop0xLKWjuGhnI8bcxfDJhWozl0IIDqcKwfrvZw",
  "rand_int": 540,
  "request_id": "4afab242-7d54-48cd-8f6a-501fe936ce66",
  "some_int": 181,
  "uuid": "3edea700-0cf1-46d6-8381-e137b4b3e65f"
}
```

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
