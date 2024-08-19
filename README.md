<a href="https://opensource.newrelic.com/oss-category/#new-relic-experimental"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Experimental.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"><img alt="New Relic Open Source experimental project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"></picture></a>

![GitHub forks](https://img.shields.io/github/forks/newrelic-experimental/newrelic-metrics-selector?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic-experimental/newrelic-metrics-selector?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic-experimental/newrelic-metrics-selector?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic-experimental/newrelic-metrics-selector/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic-experimental/newrelic-metrics-selector)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic-experimental/newrelic-metrics-selector)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic-experimental/newrelic-metrics-selector)


![GitHub issues](https://img.shields.io/github/issues/newrelic-experimental/newrelic-metrics-selector)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic-experimental/newrelic-metrics-selector)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic-experimental/newrelic-metrics-selector)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic-experimental/newrelic-metrics-selector)


# New Relic Metrics Selector Utility

New Relic Metrics Selector Utility ( `nrms` ) is a CLI tool that identifies Prometheus metrics that are not being used in either alert definitions or dashboard definitions in New Relic. The tool is written in Go and uses New Relic's NerdGraph and NRQL APIs to fetch and analyze the metrics.

## What This Utility Does

This utility performs the following tasks:

1. **Fetch Prometheus Metrics**: Reads all Prometheus metrics from the Metric table using the NRQL GraphQL call.
2. **Analyze Dashboard Definitions**: Loads all dashboard definitions for the account and examines the queries for the presence of the metric names.
3. **Analyze Alert Definitions**: Loads all alert NRQL and examines the queries for the presence of the metric names.
4. **Identify Unused Metrics**: Outputs the metric names that were not found in either dashboard or alert queries, which you may consider dropping.


## Prerequisites

- Go 1.16 or later
- New Relic API key
- New Relic account ID

## Installation

### Option 1: Download Pre-built Binaries

You can  directly download the pre-built binaries from the [releases page](https://github.com/newrelic-experimental/newrelic-metrics-selector/releases) based on your platform needs.

1. Download the appropriate tarball for your platform (Linux or macOS).
2. Extract the tarball:

    ```sh
    tar -xzf nrms-<platform>-<version>.tar.gz
    ```

3. Move the binary to a directory in your PATH, for example:

    ```sh
    mv nrms /usr/local/bin/
    ```

### Option 2: Build from Source
    
1. Clone the repository:

    ```sh
    git clone https://github.com/newrelic-experimental/newrelic-metrics-selector.git
    cd newrelic-metrics-selector
    ```

2. Install dependencies:

    ```sh
    go get github.com/newrelic/newrelic-client-go/v2/newrelic
    go get github.com/sirupsen/logrus
    go get github.com/briandowns/spinner
    ```

3. Build the application using `make`:

    ```sh
    make build
    ```

   To build for specific platforms, use:

    ```sh
    make build-linux
    make build-mac
    ```

4. Package the binaries for release:

    ```sh
    make package-linux
    make package-mac
    ```

## Usage

1. Set the required environment variables:

    ```sh
    export NEW_RELIC_API_KEY=your_new_relic_api_key
    export NEW_RELIC_ACCOUNT_ID=your_new_relic_account_id
    ```

2. Optionally, set the NRQL query for fetching Prometheus metrics and the log level:

    ```sh
    export NRQL_PROMETHEUS_METRICS="YOUR_CUSTOM_NRQL_QUERY"
    export LOG_LEVEL=debug # Set to 'info', 'warn', 'error' as needed
    ```

3. Run the application:

    ```sh
    ./bin/nrms
    ```

4. To see the help message:

    ```sh
    ./bin/nrms --help
    ```

## Output

The application will generate two output files:

- `<accountID>_used_<timestamp>.txt`: Contains the list of used Prometheus metrics.
- `<accountID>_unused_<timestamp>.txt`: Contains the list of unused Prometheus metrics.

## Example

```sh
export NEW_RELIC_API_KEY=your_new_relic_api_key
export NEW_RELIC_ACCOUNT_ID=your_new_relic_account_id
export LOG_LEVEL=debug

./bin/nrms
```

You should see the processing indicator while the application fetches data and processes the metrics. Once complete, you will see a message indicating that processing is complete and the output files have been generated.

## Details

- **Step 1**: Fetch all Prometheus metrics using the NRQL query:
  ```sql
  SELECT uniques(metricName) FROM Metric WHERE (instrumentation.name = 'remote-write') AND (instrumentation.provider = 'prometheus') LIMIT MAX
  ```
  This ensures that only Prometheus metrics are fetched.

- **Step 2**: Load all dashboard definitions using NerdGraph and examine the queries for the presence of the metric names.

- **Step 3**: Load all alert NRQL using NerdGraph and examine the queries for the presence of the metric names.

- **Step 4**: Output the metric names that were not found in either dashboard or alert queries.

## Makefile Targets

- **all**: Build the project (default).
- **clean**: Clean the build directory.
- **build**: Build the `nrms` binary for the current platform.
- **build-linux**: Build the `nrms` binary for Linux.
- **build-mac**: Build the `nrms` binary for macOS.
- **package-linux**: Package the `nrms` binary for Linux.
- **package-mac**: Package the `nrms` binary for macOS.
- **lint**: Lint the code.
- **deps**: Install dependencies.

## Running the Utility

To run the utility, follow these steps:

1. **Set the required environment variables**:
    ```sh
    export NEW_RELIC_API_KEY=your_new_relic_api_key
    export NEW_RELIC_ACCOUNT_ID=your_new_relic_account_id
    ```

2. **Optionally, set the NRQL query for fetching Prometheus metrics and the log level**:
    ```sh
    export NRQL_PROMETHEUS_METRICS="YOUR_CUSTOM_NRQL_QUERY"
    export LOG_LEVEL=debug # Set to 'info', 'warn', 'error' as needed
    ```

3. **Run the application**:
    ```sh
    ./bin/nrms
    ```

4. **To see the help message**:
    ```sh
    ./bin/nrms --help
    ```

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR DEDICATED SUPPORT. Issues and contributions should be reported to the project here on GitHub.

>We encourage you to bring your experiences and questions to the [Explorers Hub](https://discuss.newrelic.com) where our community members collaborate on solutions and new ideas.

## Contributing

We encourage your contributions to improve Salesforce Commerce Cloud for New Relic Browser! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project. If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company, please drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

## License

New Relic Metrics Selector Utility is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.

>[If applicable: [Project Name] also uses source code from third-party libraries. You can find full details on which libraries are used and the terms under which they are licensed in the third-party notices document.]
