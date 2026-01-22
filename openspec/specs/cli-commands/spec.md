# cli-commands Specification

## Purpose

The command-line interface provides user interaction for building static sites, serving content during development, and displaying version information.

## Requirements

### Requirement: Execute-Root-Command

The system SHALL provide a root command that serves as the entry point for all operations.

#### Scenario: No arguments provided

- **WHEN** the user executes the command without arguments
- **THEN** the system SHALL display usage information

#### Scenario: Help flag provided

- **WHEN** the user provides the `--help` flag
- **THEN** the system SHALL display help text describing the static site generator

#### Scenario: Version flag provided

- **WHEN** the user provides the `--version` flag
- **THEN** the system SHALL display the application name and version

#### Scenario: Execution failure

- **WHEN** a subcommand fails during execution
- **THEN** the system SHALL exit with a non-zero exit code

### Requirement: Execute-Build-Command

The system SHALL provide a build command that generates the static site from source content.

#### Scenario: Successful build

- **WHEN** the user executes the build command with valid configuration
- **THEN** the system SHALL process markdown files and generate HTML output
- **AND** the system SHALL display a success message

#### Scenario: Help for build command

- **WHEN** the user provides `build --help`
- **THEN** the system SHALL document the `--config`, `--output`, and `--content` flags

#### Scenario: Missing configuration file

- **WHEN** the user specifies a non-existent configuration file
- **THEN** the system SHALL report an error about loading configuration

### Requirement: Configure-Build-Options

The system SHALL support command-line flags to configure the build process.

#### Scenario: Config flag specified

- **WHEN** the user provides `--config` or `-c` with a path
- **THEN** the system SHALL use the specified configuration file
- **AND** the default configuration path SHALL be "ssg.yaml"

#### Scenario: Output flag overrides configuration

- **WHEN** the user provides `--output` or `-o` with a directory path
- **THEN** the system SHALL use the specified output directory instead of the configured value

#### Scenario: Content flag overrides configuration

- **WHEN** the user provides `--content` with a directory path
- **THEN** the system SHALL use the specified content directory instead of the configured value

#### Scenario: Assets flag overrides configuration

- **WHEN** the user provides `--assets` with a directory path
- **THEN** the system SHALL use the specified assets directory instead of the configured value
- **AND** the default SHALL be an empty string (no override)

### Requirement: Execute-Serve-Command

The system SHALL provide a serve command that starts a development server.

#### Scenario: Successful server start

- **WHEN** the user executes the serve command with valid configuration
- **THEN** the system SHALL start an HTTP server
- **AND** the system SHALL display the server address

#### Scenario: Help for serve command

- **WHEN** the user provides `serve --help`
- **THEN** the system SHALL document the `--port`, `--config`, `--dir`, and `--build` flags

#### Scenario: Missing configuration file

- **WHEN** the user specifies a non-existent configuration file for serve
- **THEN** the system SHALL report an error about loading configuration

#### Scenario: Serve directory does not exist

- **WHEN** the serve directory does not exist and build flag is not set
- **THEN** the system SHALL report an error indicating the directory does not exist

### Requirement: Configure-Serve-Options

The system SHALL support command-line flags to configure the development server.

#### Scenario: Port flag specified

- **WHEN** the user provides `--port` or `-p` with a number
- **THEN** the system SHALL listen on the specified port
- **AND** the default port SHALL be 8080

#### Scenario: Invalid port negative

- **WHEN** the user provides a port number less than 0
- **THEN** the system SHALL report an invalid port error

#### Scenario: Invalid port too high

- **WHEN** the user provides a port number greater than 65535
- **THEN** the system SHALL report an invalid port error

#### Scenario: Directory flag overrides output

- **WHEN** the user provides `--dir` or `-d` with a directory path
- **THEN** the system SHALL serve files from the specified directory instead of the configured output directory

#### Scenario: Build flag triggers build before serve

- **WHEN** the user provides `--build` or `-b` flag
- **THEN** the system SHALL build the site before starting the server
- **AND** the system SHALL copy assets to the output directory

### Requirement: Handle-Graceful-Shutdown

The system SHALL gracefully shut down the development server on termination signals.

#### Scenario: Interrupt signal received

- **WHEN** the server receives an interrupt signal
- **THEN** the system SHALL display a shutdown message
- **AND** the system SHALL stop accepting new connections
- **AND** the system SHALL complete in-progress requests within a timeout period

### Requirement: Execute-Version-Command

The system SHALL provide a version command that displays version information.

#### Scenario: Version command executed

- **WHEN** the user executes the version command
- **THEN** the system SHALL display the application name and version string
- **AND** the output format SHALL be "ssg version {version}"

### Requirement: Embed-Version-In-Output

The system SHALL embed version information in generated output.

#### Scenario: Build embeds version in footer

- **WHEN** the build command generates HTML output
- **THEN** the generated HTML SHALL contain the version string

#### Scenario: Serve with build embeds version

- **WHEN** the serve command runs with the build flag
- **THEN** the generated HTML SHALL contain the version string

### Requirement: Support-Build-Time-Variables

The system SHALL support version variables set at build time.

#### Scenario: Default version value

- **WHEN** version is not set via build-time flags
- **THEN** the version SHALL default to "dev"

#### Scenario: Default build date value

- **WHEN** build date is not set via build-time flags
- **THEN** the build date SHALL default to "unknown"
