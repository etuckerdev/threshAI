# REST API

## Overview
The threshAI REST API provides endpoints for interacting with the platform, allowing users to manage prompts, monitor system health, and utilize various features.

## Base URL
The base URL for all API endpoints is:
```
https://api.threshai.com/v1
```

## Authentication
All API requests require authentication. Use your API key to authenticate requests.

### Headers
Include the following headers in your API requests:
- `Authorization: Bearer YOUR_API_KEY`

## Endpoints

### Get Prompt Status
Retrieve the status of a specific prompt.

**Endpoint**: `GET /prompts/{promptId}/status`

**Parameters**:
- `promptId` (path): The ID of the prompt.

**Response**:
```json
{
  "status": "processing",
  "progress": 50,
  "output": "Partial output so far..."
}
```

### Create Prompt
Create a new prompt.

**Endpoint**: `POST /prompts`

**Request Body**:
```json
{
  "input": "Your prompt input here",
  "config": {
    "optimizationLevel": "high",
    "plugins": ["plugin1", "plugin2"]
  }
}
```

**Response**:
```json
{
  "promptId": "12345",
  "status": "queued"
}
```

### List Prompts
Retrieve a list of all prompts.

**Endpoint**: `GET /prompts`

**Response**:
```json
[
  {
    "promptId": "12345",
    "input": "Your prompt input here",
    "status": "completed",
    "output": "Final output..."
  },
  {
    "promptId": "67890",
    "input": "Another prompt input",
    "status": "failed",
    "output": "Error message..."
  }
]
```

### Delete Prompt
Delete a specific prompt.

**Endpoint**: `DELETE /prompts/{promptId}`

**Parameters**:
- `promptId` (path): The ID of the prompt.

**Response**:
```json
{
  "message": "Prompt deleted successfully"
}
```

## Error Handling
The API returns standard HTTP status codes to indicate the success or failure of the API request. For errors, the API will return a JSON object with an error message.

**Example Error Response**:
```json
{
  "error": "Invalid API key"
}
```

## Rate Limiting
The API has rate limits to ensure fair usage. If you exceed the rate limit, you will receive a `429 Too Many Requests` status code.

## Conclusion
The threshAI REST API provides a robust interface for managing prompts and monitoring system health. Ensure you follow the authentication and rate limiting guidelines to ensure smooth API usage.