# Sheet Bot

Sheet Bot is an application designed to automate the process of reading and writing data to Google Sheets based on messages received from a WhatsApp group. The bot uses Playwright for web automation and interacts with the Google Sheets API to manage data.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Logging](#logging)

## Features

- **WhatsApp Crawler**: Automatically reads messages from a specified WhatsApp group.
- **Google Sheets Integration**: Reads and writes data to Google Sheets based on the messages received.
- **Message Processing**: Processes different types of messages such as income, outcome, daily expenses, and balance inquiries.
- **System Messages**: Handles system messages for errors and invalid inputs.
- **Multi-User Support**: Handles multiple WhatsApp users with different configurations.
- **Advanced Logging**: Custom log formatting with file information for better debugging.

## Technologies Used

- **Go**: The primary programming language used for the application.
- **Playwright**: Used for web automation to interact with WhatsApp Web.
- **Google Sheets API**: Used to interact with Google Sheets for reading and writing data.
- **YAML**: Used for configuration management.
- **Labstack Gommon Log**: Used for logging.
- **Ollama AI**: Integration for message interpretation.

## Project Structure

```
/
├── cmd/
│   ├── before_run.sh
│   ├── main.go
│   └── run_test.sh
├── internal/
│   ├── client/
│   │   ├── google_sheets_client.go
│   │   └── ollama_ai_client.go
│   ├── configuration/
│   │   ├── setup_application_config.go
│   │   ├── setup_google_client.go
│   │   └── setup_logs_interceptor.go
│   ├── domain/
│   │   ├── message.go
│   │   ├── message_test.go
│   │   └── whatsapp_user.go
│   ├── services/
│   │   ├── google_sheets_service.go
│   │   ├── message_interpreter_service.go
│   │   ├── message_service.go
│   │   ├── ollama_ai_service.go
│   │   └── whatsapp_crawler_service.go
│   └── utils/
│       └── google_sheets_utils.go
├── user_data/
│   └── [browser data]
├── application.yaml
├── go.mod
├── go.sum
├── user.json
└── LICENSE
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

## WhatsApp User Configuration

The application supports both single and multiple WhatsApp users through the `WhatsappUser` structure. Users are configured by creating a `user.json` file in the project's root directory (this file should not be committed to version control as it contains sensitive information).

### user.json Structure

```json
[
  {
    "group_name": "Group Name 1",
    "sheet_id": "spreadsheet-id-1",
    "is_archived": false
  }
]
```

For a single user, you only need to include one object in the array. For multiple users, add additional objects:

```json
[
  {
    "group_name": "Group Name 1",
    "sheet_id": "spreadsheet-id-1",
    "is_archived": false
  },
  {
    "group_name": "Group Name 2",
    "sheet_id": "spreadsheet-id-2",
    "is_archived": true
  }
]
```

Each entry in the array represents a user with the following properties:

```go
type WhatsappUser struct {
  GroupName  string  // Name of the WhatsApp group
  SheetId    string  // Google Sheet ID for this user
  IsArchived bool    // Whether the chat is archived
}
```

**Important:** Add the `user.json` file to your `.gitignore` to avoid sharing sensitive data in public repositories.

The configured users will be loaded automatically during service initialization:

```go
wcs := services.NewWhatsAppCrawlerService(ctx, appConfig, ms, &[]domain.WhatsappUser{*user1, *user2})
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

3. **Run Tests**: You can run tests using:
   ```bash
   ./cmd/run_test.sh
   ```

## Logging

The application uses a custom log interceptor to format logs for better readability and debugging:

- Logs include timestamp, log level, message, and the source file
- Format: `YYYY/MM/DD HH:MM:SS LEVEL MESSAGE [FILE]`

The log interceptor can be found in `internal/configuration/setup_logs_interceptor.go`.

## License

This project is licensed under the MIT License.
