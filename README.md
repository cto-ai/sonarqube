![](https://raw.githubusercontent.com/cto-ai/sonarqube/master/assets/banner.png)

# SonarQube

An Ops port of the SonarScanner CLI tool.

![](https://raw.githubusercontent.com/cto-ai/sonarqube/master/assets/screenshot_cli.png)

## Requirements

To run this or any other Op, install the [Ops Platform](https://cto.ai/platform).

Find information about how to run and build Ops via the [Ops Platform Documentation](https://cto.ai/docs/overview).

This Op also requires:

- A SonarQube server
  - `Anyone` and `sonar-users` group should have all permissions removed
  - `Default visibility of new projects` should be set to private
    - Administration Tab -> Projects -> Managements -> Default visibility of new projects
  - Create a user for the ops, with permissions `Administer System`, `Execute Analysis`, and `Create Projects`
    - Security -> Global Permission -> update checkboxes
  - Get the token for that user and add it to secrets store via `ops secrets:set`. Key should ideally be `sonarToken`
  - Get the server address ready to input when it prompts
- A publicly accessible Git repository, with HTTP(S) protocol enabled
  - Currently this op does not work with SSH protocol

## Usage

To initiate the interactive SonarQube CLI prompt run:

```bash
ops run sonarqube [gitRepoURL]
```

## Local Development / Running from Source

**1. Clone the repo:**

```bash
git clone <git url>
```

**2. Navigate into the directory and install dependencies:**

```bash
cd sonarqube && go get ./...
```

**3. Run the Op from your current working directory with:**

```bash
ops run .
```

## Debugging Issues

Use the `DEBUG` flag in the terminal to see verbose Op output like so:

```bash
DEBUG=sonarqube:* ops run sonarqube
```

When submitting issues or requesting help, be sure to also include the version information. To get your ops version run:

```bash
ops -v
```

## Resources

### SonarQube Docs

- [SonarQube Documentation](https://docs.sonarqube.org/latest/)
- [SonarScanner Reference](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/)

## Contributing

See the [Contributing Docs](CONTRIBUTING.md) for more information.

## Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/jmariomejiap"><img src="https://github.com/jmariomejiap.png" width="100px;" alt="Mario Mejia"/><br /><sub><b>Mario Mejia</b></sub></a><br/></td>
    <td align="center"><a href="https://github.com/CalHoll"><img src="https://github.com/CalHoll.png" width="100px;" alt="Calvin Holloway"/><br /><sub><b>Calvin Holloway</b></sub></a><br/></td>
    <td align="center"><a href="https://github.com/aschereT"><img src="https://github.com/aschereT.png" width="100px;" alt="aschereT's face here"/><br /><sub><b>Vincent Tan</b></sub></a><br/></td>
    <td align="center"><a href="https://github.com/minsohng"><img src="https://github.com/minsohng.png" width="100px;" alt="Min Sohng"/><br /><sub><b>Min Sohng</b></sub></a><br/></td>
  </tr>
</table>

## LICENSE

[MIT](LICENSE)
