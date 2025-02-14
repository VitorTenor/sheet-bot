# Sheet Bot

Sheet Bot is an application designed to automate the process of reading and writing data to Google Sheets based on messages received from a WhatsApp group. The bot uses Playwright for web automation and interacts with the Google Sheets API to manage data.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)

## Features

- **WhatsApp Crawler**: Automatically reads messages from a specified WhatsApp group.
- **Google Sheets Integration**: Reads and writes data to Google Sheets based on the messages received.
- **Message Processing**: Processes different types of messages such as income, outcome, daily expenses, and balance inquiries.
- **System Messages**: Handles system messages for errors and invalid inputs.

## Technologies Used

- **Go**: The primary programming language used for the application.
- **Playwright**: Used for web automation to interact with WhatsApp Web.
- **Google Sheets API**: Used to interact with Google Sheets for reading and writing data.
- **YAML**: Used for configuration management.
- **Labstack Gommon Log**: Used for logging.

## Project Structure

```
/
├── cmd/
│   ├── before_run.sh
│   └── main.go
├── internal/
│   ├── client/
│   │   └── google_sheets_client.go
│   ├── configs/
│   │   ├── setup_application_config.go
│   │   └── setup_google_client.go
│   ├── domain/
│   │   └── message.go
│   ├── services/
│   │   ├── google_sheets_service.go
│   │   ├── message_service.go
│   │   └── whatsapp_crawler_service.go
│   └── utils/
│       └── google_sheets_utils.go
├── application.yaml
├── go.mod
└── go.sum
```

## Configuration

The application requires a configuration file named `application.yaml` in the root directory. This file should contain the necessary configurations for Google Sheets API and WhatsApp Web.

Example `application.yaml`:

```yaml
google:
  client_email: your-client-email
  private_key: your-private-key
  api_url: https://www.googleapis.com/auth/spreadsheets
  sheet_id: your-sheet-id

whatsapp:
  web_url: https://web.whatsapp.com
  group_name: your-group-name

crawler:
  user_data_dir: /path/to/your/user/data/dir
```

## Running the Application

1. **Install Dependencies**: Ensure you have Go installed and run the following command to install Playwright:
   ```bash
   ./cmd/before_run.sh
   ```

2. **Run the Application**: Execute the main Go file to start the application:
   ```bash
   go run cmd/main.go
   ```

## License

This project is licensed under the MIT License.