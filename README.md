# Summary.go

Welcome to Summary.go, a GPT Chat Assistant CLI! This project is a command-line interface (CLI) that allows you to interact with OpenAI's GPT language models, which has impressive natural language understanding and generation capabilities. With this CLI, you can have engaging, dynamic, and thought-provoking conversations with the AI assistant. This tool is ideal for writers, developers, researchers, or anyone seeking an intelligent conversation partner.

## Features

- User-friendly command-line interface for seamless interaction
- Real-time chat experience with streaming API functionality
- Built with Cobra, a popular Go library for creating CLI applications
- Leverages the OpenAI GPT language model for high-quality AI-generated responses

## Installation

To install Summary, follow these steps:

1. Clone the repository:

```
git clone https://github.com/jharkins/summary.go.git
```

2. Change to the project directory:

```
cd summary.go
```

3. Install the required dependencies:

```
go get -u ...
```

4. Build the CLI application:

```
go build -o summary
```

## Usage

On first run, the app needs your OpenAI API key. You can get that here: [OpenAI Platform :: API Keys](https://platform.openai.com/account/api-keys)

To start a chat conversation, run the following command:

```
./summary chat
```

You'll be prompted to enter your message, and the AI assistant will respond accordingly.

## Credit

This project is a collaboration between Joe Harkins, the idea guy, and an AI assistant by OpenAI. The AI assistant has provided valuable input and implementation suggestions, while Joe has contributed creative vision and project direction.

## License

MIT License

## Disclaimer

This CLI uses OpenAI's GPT-4 language model, which can sometimes generate content that may be inappropriate or offensive. Please use this tool responsibly and be aware of the potential for unexpected output.
