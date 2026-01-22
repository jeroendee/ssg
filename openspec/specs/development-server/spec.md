# development-server Specification

## Purpose
Provide a local HTTP server for previewing static site output during development, serving files from a configured directory with support for multiple index file types and security protections.

## Requirements

### Requirement: Configure-Server
The system SHALL accept configuration specifying a network port and a directory path to serve.

#### Scenario: Default port configuration
- **WHEN** server is configured with port 8080 and directory "public"
- **THEN** the server SHALL listen on port 8080
- **AND** the server SHALL serve files from the "public" directory

#### Scenario: Custom port configuration
- **WHEN** server is configured with port 3000 and directory "/tmp/site"
- **THEN** the server SHALL listen on port 3000
- **AND** the server SHALL serve files from "/tmp/site"

#### Scenario: Automatic port assignment
- **WHEN** server is configured with port 0
- **THEN** the system SHALL assign an available port automatically
- **AND** the assigned port SHALL be retrievable after startup

### Requirement: Start-Server
The system SHALL start an HTTP server that listens for incoming requests.

#### Scenario: Server startup success
- **WHEN** server is started with a valid configuration
- **THEN** the server SHALL begin accepting HTTP requests
- **AND** the server address SHALL be available for retrieval

#### Scenario: Server address retrieval before start
- **WHEN** server address is queried before startup
- **THEN** the system SHALL return an empty address

### Requirement: Serve-Static-Files
The system SHALL serve static files from the configured directory.

#### Scenario: Serve HTML file
- **WHEN** client requests "/index.html"
- **AND** the file exists in the configured directory
- **THEN** the server SHALL return the file contents
- **AND** the response status SHALL be 200 OK

#### Scenario: File not found
- **WHEN** client requests "/missing.html"
- **AND** the file does not exist
- **THEN** the response status SHALL be 404 Not Found

### Requirement: Set-Content-Types
The system SHALL set appropriate content type headers based on file extensions.

#### Scenario: CSS file content type
- **WHEN** client requests a file with ".css" extension
- **THEN** the Content-Type header SHALL be "text/css; charset=utf-8"

#### Scenario: XML file content type
- **WHEN** client requests a file with ".xml" extension
- **THEN** the Content-Type header SHALL be "application/xml" or "text/xml; charset=utf-8"

### Requirement: Serve-Directory-Index
The system SHALL serve index files when a directory is requested.

#### Scenario: Directory with index.html
- **WHEN** client requests a directory path ending with "/"
- **AND** the directory contains "index.html"
- **THEN** the server SHALL serve the index.html file

#### Scenario: Directory with index.xml only
- **WHEN** client requests a directory path ending with "/"
- **AND** the directory does not contain "index.html"
- **AND** the directory contains "index.xml"
- **THEN** the server SHALL serve the index.xml file
- **AND** the Content-Type header SHALL be "text/xml; charset=utf-8"

#### Scenario: Sitemap directory with index.xml
- **WHEN** client requests "/sitemap/"
- **AND** the sitemap directory contains only "index.xml"
- **THEN** the server SHALL serve the sitemap index.xml file
- **AND** the response status SHALL be 200 OK

### Requirement: Prevent-Path-Traversal
The system SHALL reject requests that attempt to access files outside the configured directory.

#### Scenario: Obvious path traversal attack
- **WHEN** client requests "/../../../etc/passwd"
- **THEN** the response status SHALL be 403 Forbidden

#### Scenario: Traversal to temporary directory
- **WHEN** client requests "/../../../tmp/secret"
- **THEN** the response status SHALL be 403 Forbidden

#### Scenario: Hidden traversal with valid prefix
- **WHEN** client requests "/valid/../../../etc/passwd"
- **THEN** the response status SHALL be 403 Forbidden

#### Scenario: Encoded path traversal
- **WHEN** client requests a path containing encoded traversal sequences
- **THEN** the response status SHALL be 403 Forbidden

### Requirement: Allow-Legitimate-Paths
The system SHALL allow requests for paths that resolve within the configured directory.

#### Scenario: Root directory request
- **WHEN** client requests "/"
- **AND** the root directory contains index.html
- **THEN** the response status SHALL be 200 OK
- **AND** the response SHALL contain the index.html content

#### Scenario: Subdirectory with XML index
- **WHEN** client requests "/sitemap/"
- **AND** the subdirectory contains index.xml
- **THEN** the response status SHALL be 200 OK
- **AND** the response SHALL contain the index.xml content

### Requirement: Configure-Timeouts
The system SHALL configure connection timeouts to protect against denial-of-service attacks.

#### Scenario: Header read timeout
- **WHEN** server is started
- **THEN** the server SHALL limit time to read request headers to 10 seconds

#### Scenario: Idle connection timeout
- **WHEN** server is started
- **THEN** the server SHALL close idle keepalive connections after 120 seconds

### Requirement: Shutdown-Server
The system SHALL support graceful shutdown allowing in-flight requests to complete.

#### Scenario: Graceful shutdown
- **WHEN** shutdown is requested with a context timeout
- **THEN** the server SHALL stop accepting new connections
- **AND** the server SHALL allow in-flight requests to complete within the timeout

#### Scenario: Shutdown before start
- **WHEN** shutdown is requested before server has started
- **THEN** the shutdown operation SHALL complete without error
