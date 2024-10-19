# Binance Orderbook Average

This project provides a real-time WebSocket server that calculates and broadcasts the average price of orders from the Binance order book.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the project, clone the repository and build the Go application:

```bash
git clone https://github.com/yourusername/binance-orderbook-average.git
cd binance-orderbook-average
go build
```

## Usage

To start the server, run the following command:

```bash
./binance-orderbook-average
```

## Configuration

The application can be configured using environment variables or a configuration file. The configuration includes settings for the Binance WebSocket connection and the server port.

## API Endpoints

- **Root Route**: `GET /`

  - Returns a simple "HelloWorld" message.

- **WebSocket Route**: `GET /average-price`
  - Establishes a WebSocket connection to receive real-time average price updates.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
