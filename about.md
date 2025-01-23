# About ThreshAI

ThreshAI is a versatile AI interaction platform that provides multiple ways to engage with large language models (LLMs). It offers both a command-line interface and a web server, supporting multiple LLM providers including Ollama (local) and DeepSeek (cloud-based).

## Core Features

### 1. Multi-Modal Interaction
- **CLI Interface**: Direct command-line access for quick prompts and responses
- **Web Server**: RESTful API endpoints for integration with other applications
- **Interactive Chat**: Persistent chat sessions with memory and context awareness

### 2. LLM Provider Integration
- **Ollama Integration**: Local LLM deployment using Ollama (default endpoint: http://localhost:11434)
- **DeepSeek Integration**: Cloud-based LLM access using DeepSeek's API
- **Extensible Architecture**: Modular design allowing easy addition of new LLM providers

### 3. Advanced Features
- **Memory Management**: Persistent conversation history with context retrieval
- **Caching System**: In-memory caching for improved response times
- **Context Awareness**: Maintains conversation context for more coherent interactions
- **Clarification System**: Automatic detection and handling of ambiguous queries

## Architecture

### Core Components
1. **Generation Engine**
   - Abstract generator interface for consistent interaction across providers
   - Provider-specific adapters for Ollama and DeepSeek
   - Centralized prompt handling and response generation

2. **Memory System**
   - Redis-backed persistent storage
   - Context retrieval for relevant conversation history
   - Interaction tracking and management

3. **Web Server**
   - RESTful endpoints for generation requests
   - Provider-specific routes (/generate/ollama and /generate/deepseek)
   - Simple HTTP interface for easy integration

4. **CLI Application**
   - Direct access to generation capabilities
   - Provider selection (ollama/deepseek)
   - Environment-based configuration

### Key Design Patterns
- **Adapter Pattern**: For LLM provider integration
- **Factory Pattern**: For generator instantiation
- **Repository Pattern**: For memory management
- **Strategy Pattern**: For provider selection

## Technical Stack
- **Language**: Go 1.21+
- **Dependencies**:
  - Redis (for memory management)
  - Ollama (for local LLM deployment)
  - Cobra (for CLI implementation)
  - Standard Go HTTP package (for web server)

## Use Cases

1. **Development Integration**
   - Code generation assistance
   - API integration testing
   - Documentation generation

2. **Interactive Applications**
   - Chatbot development
   - Customer service automation
   - Knowledge base querying

3. **Local Development**
   - Private LLM deployment using Ollama
   - Offline development capabilities
   - Custom model integration

## Performance Features

1. **Caching**
   - In-memory response caching
   - Reduced API calls for common queries
   - Configurable cache strategies

2. **Memory Optimization**
   - Context-aware memory management
   - Efficient storage and retrieval
   - Automatic cleanup of old conversations

3. **Error Handling**
   - Graceful degradation
   - Automatic retries for failed requests
   - Comprehensive error reporting

## Security Considerations
- API key management for cloud providers
- Local-first approach available through Ollama
- Environment-based configuration
- No persistent storage of sensitive information