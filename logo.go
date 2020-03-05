package main

import (
	ctoai "github.com/cto-ai/sdk-go"
)

func printLogo(client ctoai.Client) {
	switch client.Sdk.GetInterfaceType() {
	case "terminal":
		client.Ux.Print(`
  [94m██████[39m[33m╗[39m [94m████████[39m[33m╗[39m  [94m██████[39m[33m╗ [39m      [94m█████[39m[33m╗[39m  [94m██[39m[33m╗[39m
 [94m██[39m[33m╔════╝[39m [33m╚══[39m[94m██[39m[33m╔══╝[39m [94m██[39m[33m╔═══[39m[94m██[39m[33m╗[39m     [94m██[39m[33m╔══[39m[94m██[39m[33m╗[39m [94m██[39m[33m║[39m
 [94m██[39m[33m║     [39m [94m   ██[39m[33m║   [39m [94m██[39m[33m║[39m[94m   ██[39m[33m║[39m     [94m███████[39m[33m║[39m [94m██[39m[33m║[39m
 [94m██[39m[33m║     [39m [94m   ██[39m[33m║   [39m [94m██[39m[33m║[39m[94m   ██[39m[33m║[39m     [94m██[39m[33m╔══[39m[94m██[39m[33m║[39m [94m██[39m[33m║[39m
 [33m╚[39m[94m██████[39m[33m╗[39m [94m   ██[39m[33m║   [39m [33m╚[39m[94m██████[39m[33m╔╝[39m [94m██[39m[33m╗[39m [94m██[39m[33m║[39m[94m  ██[39m[33m║[39m [94m██[39m[33m║[39m
 [33m ╚═════╝[39m [33m   ╚═╝   [39m [33m ╚═════╝ [39m [33m╚═╝[39m [33m╚═╝  ╚═╝[39m [33m╚═╝[39m

We’re building the world’s best developer experiences.
  `)
	default:
		client.Ux.Print(`:white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square:
:white_square::white_square::black_square::black_square::white_square::white_square::black_square::black_square::black_square::white_square::white_square::white_square::black_square::black_square::black_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::white_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::black_square::black_square::white_square::white_square::white_square::black_square::black_square::white_square::white_square:
:white_square::black_square::white_square::white_square::black_square::white_square::black_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::black_square::white_square:
:white_square::white_square::black_square::black_square::white_square::white_square::black_square::white_square::white_square::white_square::white_square::black_square::black_square::black_square::white_square::white_square:
:white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square::white_square:`)
	}

	client.Ux.Print(`👋  Hi there! Welcome to the CTO.ai SonarQube Op!
This Op allows you to trigger a SonarQube scan of a public Git repo.
For more information, see the README.
`)
}
